/*
 * crypt/recover_key.c
 *
 * Enhanced key recovery tool for the XOR+rotate obfuscation used in this repo.
 *
 * Features added:
 *  - Try multiple key lengths automatically (single values, comma lists, ranges).
 *  - Accept multiple ciphertext:crib[:offset] pairs on the command line and
 *    combine evidence across all pairs.
 *  - Write best-effort recovered key bytes to a binary output file and a
 *    readable text file (hex + ASCII).
 *
 * Usage:
 *   recover_key [options] cipher1:crib1[:offset1] [cipher2:crib2[:offset2] ...]
 *
 * Options:
 *   -k, --keylens <list>   Comma-separated key lengths or ranges. Examples:
 *                             64             (single)
 *                             64,128,256     (list)
 *                             16-128         (range inclusive)
 *                           If omitted, defaults to trying 1..256.
 *
 *   -o, --output <file>    Write best-effort recovered key (binary) to <file>.
 *   -t, --text <file>      Write human-readable recovered key to <file>.
 *   -h, --help             Show usage.
 *
 * Input arguments:
 *   Each input is "cipherfile:cribfile" or "cipherfile:cribfile:offset"
 *   where offset is a non-negative decimal integer (0-based index where the
 *   crib aligns in the ciphertext).
 *
 * Example:
 *   recover_key -k 64,128 -o recovered_key.bin cipher.bin:crib.txt:10 other.bin:known.txt
 *
 * Notes/Limitations:
 *  - The tool requires you to supply ciphertext(s) and at least one crib.
 *  - You must provide (or try) the correct key length(s). The tool will print
 *    candidate sets for each key length tried and the best-effort key for each.
 *  - This is an educational/testing tool to be used only on data you own / are
 *    authorized to test.
 *
 * Implementation notes:
 *  - The scheme assumed is the one used in this repository:
 *      start = message_len % key_len
 *      kidx = (start + i) % key_len
 *      shift = (k + i) % 8
 *      tmp = plaintext[i] ^ k
 *      cipher[i] = rol8(tmp, shift)
 *
 *  - Given cipher[i], plaintext[i], and index i, we brute-force k in 0..255:
 *      shift' = (k' + i) % 8
 *      tmp'   = ror8(cipher[i], shift')
 *      check  (tmp' ^ k') == plaintext[i]
 *
 *  - Candidate sets for each key index are represented as 256-bit bitmaps.
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <errno.h>
#include <inttypes.h>
#include <ctype.h>

#define MAX_AUTO_KEYLEN 256   /* default max to try when no -k provided */
#define CAND_BITMAP_BYTES 32  /* 256 bits / 8 */

static inline uint8_t ror8(uint8_t v, unsigned int s) {
    s &= 7u;
    if (s == 0) return v;
    return (uint8_t)((v >> s) | (v << (8 - s)));
}

/* Read entire file into a malloc'd buffer. Caller must free. */
static unsigned char *read_file(const char *path, size_t *out_len) {
    FILE *f = fopen(path, "rb");
    if (!f) {
        fprintf(stderr, "fopen(%s): %s\n", path, strerror(errno));
        return NULL;
    }
    if (fseek(f, 0, SEEK_END) != 0) {
        fprintf(stderr, "fseek(%s) failed: %s\n", path, strerror(errno));
        fclose(f);
        return NULL;
    }
    long s = ftell(f);
    if (s < 0) {
        fprintf(stderr, "ftell(%s) failed: %s\n", path, strerror(errno));
        fclose(f);
        return NULL;
    }
    rewind(f);
    unsigned char *buf = malloc((size_t)s + 1);
    if (!buf) {
        fprintf(stderr, "malloc(%ld) failed\n", (long)s + 1);
        fclose(f);
        return NULL;
    }
    size_t got = fread(buf, 1, (size_t)s, f);
    if (got != (size_t)s) {
        if (ferror(f)) {
            fprintf(stderr, "fread(%s) error: %s\n", path, strerror(errno));
            free(buf);
            fclose(f);
            return NULL;
        }
    }
    fclose(f);
    buf[got] = 0;
    if (out_len) *out_len = got;
    return buf;
}

/* Candidate bitmap container for key_len positions */
typedef struct {
    unsigned char *bitmap; /* CAND_BITMAP_BYTES bytes */
} cand_t;

/* Allocate candidate array of length key_len, initializing to "all ones" */
static cand_t *alloc_candidates(size_t key_len) {
    cand_t *arr = calloc(key_len, sizeof(cand_t));
    if (!arr) return NULL;
    for (size_t i = 0; i < key_len; ++i) {
        arr[i].bitmap = malloc(CAND_BITMAP_BYTES);
        if (!arr[i].bitmap) {
            for (size_t j = 0; j < i; ++j) free(arr[j].bitmap);
            free(arr);
            return NULL;
        }
        /* set all bits to 1 (all 0..255 possible) */
        memset(arr[i].bitmap, 0xFFu, CAND_BITMAP_BYTES);
    }
    return arr;
}

static void free_candidates(cand_t *arr, size_t key_len) {
    if (!arr) return;
    for (size_t i = 0; i < key_len; ++i) free(arr[i].bitmap);
    free(arr);
}

/* Intersect candidate bitmap at index 'pos' with newmap (CAND_BITMAP_BYTES bytes). */
static void intersect_bitmap(cand_t *arr, size_t pos, const unsigned char *newmap) {
    unsigned char *bm = arr[pos].bitmap;
    for (int i = 0; i < CAND_BITMAP_BYTES; ++i) bm[i] &= newmap[i];
}

/* Build a bitmap for a single (cipher_byte, plain_byte, index) triple */
static unsigned char *build_bitmap_for_triple(uint8_t cbyte, uint8_t pbyte, size_t index) {
    unsigned char *map = calloc(CAND_BITMAP_BYTES, 1);
    if (!map) return NULL;
    for (unsigned int k = 0; k < 256; ++k) {
        unsigned int shift = (unsigned int)((k + (uint8_t)index) % 8u);
        uint8_t tmp = ror8(cbyte, shift);
        if ((uint8_t)(tmp ^ (uint8_t)k) == pbyte) {
            map[k >> 3] |= (1u << (k & 7));
        }
    }
    return map;
}

/* Count bits set in a byte (simple popcount for 8-bit) */
static int popcount8(unsigned char b) {
    /* Kernighan's method compact */
    int c = 0;
    while (b) { b &= (b - 1); ++c; }
    return c;
}

/* Count candidates in bitmap */
static int count_candidates_bitmap(const unsigned char *bm) {
    int total = 0;
    for (int i = 0; i < CAND_BITMAP_BYTES; ++i) total += popcount8(bm[i]);
    return total;
}

/* If single candidate, return its value and set ok=1. Otherwise ok=0. */
static unsigned int single_candidate_value(const unsigned char *bm, int *ok) {
    unsigned int value = 0;
    int found = 0;
    for (unsigned int b = 0; b < 256; ++b) {
        if (bm[b >> 3] & (1u << (b & 7))) {
            if (!found) { found = 1; value = b; }
            else { *ok = 0; return 0; }
        }
    }
    if (found) { *ok = 1; return value; }
    *ok = 0; return 0;
}

/* Print readable candidate list for a bitmap (limited output) */
static void print_candidate_list(const unsigned char *bm) {
    int total = count_candidates_bitmap(bm);
    if (total == 0) { printf("(none)"); return; }
    if (total == 1) {
        int ok;
        unsigned int v = single_candidate_value(bm, &ok);
        if (ok) {
            if (v >= 32 && v <= 126) printf("0x%02X '%c'", v, (char)v);
            else printf("0x%02X", v);
            return;
        }
    }
    /* print up to first 32 candidates */
    int printed = 0;
    for (unsigned int b = 0; b < 256; ++b) {
        if (bm[b >> 3] & (1u << (b & 7))) {
            if (printed) printf(", ");
            printf("0x%02X", b);
            printed++;
            if (printed >= 32) { printf(", ..."); break; }
        }
    }
}

/* Parse key length specification string like "64,128,256,16-32" and append parsed
 * values into dynamically allocated array. Returns count, or -1 on error.
 *
 * Caller must free *out_arr.
 */
static int parse_keylens(const char *spec, size_t **out_arr) {
    if (!spec || !*spec) return -1;
    size_t *arr = NULL;
    size_t arr_len = 0;

    const char *p = spec;
    while (*p) {
        /* parse number */
        while (*p == ' ' || *p == '\t') ++p;
        if (!isdigit((unsigned char)*p)) { free(arr); return -1; }
        char *end = NULL;
        long a = strtol(p, &end, 10);
        if (end == p || a <= 0) { free(arr); return -1; }
        p = end;
        /* check if range */
        while (*p == ' ' || *p == '\t') ++p;
        if (*p == '-') {
            ++p;
            while (*p == ' ' || *p == '\t') ++p;
            if (!isdigit((unsigned char)*p)) { free(arr); return -1; }
            long b = strtol(p, &end, 10);
            if (end == p || b <= 0 || b < a) { free(arr); return -1; }
            p = end;
            for (long v = a; v <= b; ++v) {
                size_t *tmp = realloc(arr, (arr_len + 1) * sizeof(size_t));
                if (!tmp) { free(arr); return -1; }
                arr = tmp;
                arr[arr_len++] = (size_t)v;
            }
        } else {
            size_t *tmp = realloc(arr, (arr_len + 1) * sizeof(size_t));
            if (!tmp) { free(arr); return -1; }
            arr = tmp;
            arr[arr_len++] = (size_t)a;
        }
        while (*p == ' ' || *p == '\t') ++p;
        if (*p == ',') { ++p; continue; }
        else if (*p == '\0') break;
        else { free(arr); return -1; }
    }

    *out_arr = arr;
    return (int)arr_len;
}

/* Structure to hold one cipher+crib pair */
typedef struct {
    char *cipher_path;
    char *crib_path;
    size_t offset;        /* offset into ciphertext where crib aligns */
    unsigned char *cipher; size_t cipher_len;
    unsigned char *crib;   size_t crib_len;
} pair_t;

/* Parse an input token of form cipher:crib or cipher:crib:offset */
static int parse_pair_token(const char *token, pair_t *out) {
    char *tok = strdup(token);
    if (!tok) return -1;
    char *p1 = strchr(tok, ':');
    if (!p1) { free(tok); return -1; }
    *p1 = '\0';
    char *p2 = strchr(p1 + 1, ':');
    if (p2) {
        *p2 = '\0';
        /* offset specified */
        char *end = NULL;
        errno = 0;
        long off = strtol(p2 + 1, &end, 10);
        if (errno || *end != '\0' || off < 0) {
            free(tok); return -1;
        }
        out->offset = (size_t)off;
    } else {
        out->offset = 0;
    }
    out->cipher_path = strdup(tok);
    out->crib_path = strdup(p1 + 1);
    free(tok);
    if (!out->cipher_path || !out->crib_path) {
        free(out->cipher_path); free(out->crib_path);
        return -1;
    }
    out->cipher = NULL; out->crib = NULL;
    out->cipher_len = out->crib_len = 0;
    return 0;
}

static void free_pair(pair_t *p) {
    if (!p) return;
    free(p->cipher_path);
    free(p->crib_path);
    free(p->cipher);
    free(p->crib);
}

/* Process one key length and the provided pairs. Returns 0 on success. */
static int process_keylen(size_t key_len, pair_t *pairs, size_t pair_count,
                          int write_output, const char *out_bin, const char *out_text) {
    printf("=== Trying key length: %zu ===\n", key_len);
    cand_t *cands = alloc_candidates(key_len);
    if (!cands) {
        fprintf(stderr, "OOM allocating candidate arrays for key_len %zu\n", key_len);
        return -1;
    }

    /* Initialize candidate bitmaps to all ones (alloc_candidates already did) */

    /* For each pair, load files if not already loaded, then process crib bytes */
    for (size_t pi = 0; pi < pair_count; ++pi) {
        pair_t *pp = &pairs[pi];
        if (!pp->cipher) {
            pp->cipher = read_file(pp->cipher_path, &pp->cipher_len);
            if (!pp->cipher) { free_candidates(cands, key_len); return -1; }
        }
        if (!pp->crib) {
            pp->crib = read_file(pp->crib_path, &pp->crib_len);
            if (!pp->crib) { free_candidates(cands, key_len); return -1; }
        }
        if (pp->offset + pp->crib_len > pp->cipher_len) {
            fprintf(stderr, "Crib exceeds ciphertext length: %s offset %zu crib_len %zu cipher_len %zu\n",
                    pp->cipher_path, pp->offset, pp->crib_len, pp->cipher_len);
            free_candidates(cands, key_len); return -1;
        }
        size_t start = pp->cipher_len % key_len;
        for (size_t j = 0; j < pp->crib_len; ++j) {
            size_t msg_index = pp->offset + j;
            size_t kidx = (start + msg_index) % key_len;
            uint8_t cbyte = (uint8_t)pp->cipher[msg_index];
            uint8_t pbyte = (uint8_t)pp->crib[j];
            unsigned char *bm = build_bitmap_for_triple(cbyte, pbyte, msg_index);
            if (!bm) { free_candidates(cands, key_len); return -1; }
            intersect_bitmap(cands, kidx, bm);
            free(bm);
        }
    }

    /* Print summary and best-effort key */
    printf("Summary for key_len=%zu\n", key_len);
    int total_sure = 0;
    for (size_t i = 0; i < key_len; ++i) {
        int cnt = count_candidates_bitmap(cands[i].bitmap);
        if (cnt == 1) ++total_sure;
    }
    printf("Positions with single candidate: %d / %zu\n", total_sure, key_len);

    /* Print best-effort recovered key line (ASCII / ? for ambiguous) */
    printf("Best-effort key string (single-candidate bytes shown ASCII, others '?'):\n");
    for (size_t i = 0; i < key_len; ++i) {
        int ok; unsigned int v = single_candidate_value(cands[i].bitmap, &ok);
        if (ok) {
            unsigned char ch = (unsigned char)v;
            if (ch >= 32 && ch <= 126) putchar((char)ch);
            else printf("\\x%02X", (unsigned int)ch);
        } else putchar('?');
    }
    printf("\n\nDetailed per-position candidates (limited):\n");
    for (size_t i = 0; i < key_len; ++i) {
        printf("key[%4zu] (%3d candidates): ", i, count_candidates_bitmap(cands[i].bitmap));
        print_candidate_list(cands[i].bitmap);
        printf("\n");
    }

    /* Optionally write binary and text outputs */
    if (write_output && out_bin) {
        FILE *fb = fopen(out_bin, "wb");
        if (!fb) {
            fprintf(stderr, "fopen(%s) for writing failed: %s\n", out_bin, strerror(errno));
        } else {
            for (size_t i = 0; i < key_len; ++i) {
                int ok; unsigned int v = single_candidate_value(cands[i].bitmap, &ok);
                unsigned char outb = ok ? (unsigned char)v : 0x00;
                fwrite(&outb, 1, 1, fb);
            }
            fclose(fb);
            printf("Wrote best-effort binary key to %s (ambiguous positions written as 0x00)\n", out_bin);
        }
    }
    if (write_output && out_text) {
        FILE *ft = fopen(out_text, "w");
        if (!ft) {
            fprintf(stderr, "fopen(%s) for writing failed: %s\n", out_text, strerror(errno));
        } else {
            fprintf(ft, "# Best-effort recovered key for key_len=%zu\n", key_len);
            fprintf(ft, "# Single-candidate bytes shown as hex; ambiguous positions listed as '?'\n");
            for (size_t i = 0; i < key_len; ++i) {
                int ok; unsigned int v = single_candidate_value(cands[i].bitmap, &ok);
                if (ok) fprintf(ft, "%02X", v);
                else fprintf(ft, "??");
                if (i + 1 < key_len) fputc(' ', ft);
            }
            fputc('\n', ft);
            fprintf(ft, "\nDetailed per-position candidates:\n");
            for (size_t i = 0; i < key_len; ++i) {
                fprintf(ft, "key[%4zu] (%3d): ", i, count_candidates_bitmap(cands[i].bitmap));
                int cnt = count_candidates_bitmap(cands[i].bitmap);
                if (cnt == 0) fprintf(ft, "(none)\n");
                else if (cnt == 1) {
                    int ok2; unsigned int v2 = single_candidate_value(cands[i].bitmap, &ok2);
                    if (ok2) fprintf(ft, "0x%02X\n", v2); else fprintf(ft, "(err)\n");
                } else {
                    int printed = 0;
                    for (unsigned int b = 0; b < 256; ++b) {
                        if (cands[i].bitmap[b >> 3] & (1u << (b & 7))) {
                            if (printed) fputs(", ", ft);
                            fprintf(ft, "0x%02X", b);
                            printed++;
                            if (printed >= 64) { fputs(", ...", ft); break; }
                        }
                    }
                    fputc('\n', ft);
                }
            }
            fclose(ft);
            printf("Wrote readable key candidate output to %s\n", out_text);
        }
    }

    free_candidates(cands, key_len);
    return 0;
}

static void usage(const char *prog) {
    fprintf(stderr,
        "Usage: %s [options] cipher:crib[:offset] [cipher:crib[:offset] ...]\n"
        "Options:\n"
        "  -k, --keylens <spec>   comma-separated key lengths and ranges (e.g. 64,128,16-32)\n"
        "                         If omitted, tries key lengths 1..%d\n"
        "  -o, --output <file>    write best-effort recovered key (binary)\n"
        "  -t, --text <file>      write human-readable key candidates (text)\n"
        "  -h, --help             show this help\n",
        prog, MAX_AUTO_KEYLEN);
}

/* Simple helper to check if string starts with prefix */
static int starts_with(const char *s, const char *pref) {
    return strncmp(s, pref, strlen(pref)) == 0;
}

int main(int argc, char **argv) {
    if (argc < 2) { usage(argv[0]); return 2; }

    size_t *keylens = NULL;
    int keylens_count = 0;
    const char *out_bin = NULL;
    const char *out_text = NULL;
    int i = 1;
    /* parse options */
    while (i < argc && argv[i][0] == '-') {
        if (strcmp(argv[i], "-k") == 0 || strcmp(argv[i], "--keylens") == 0) {
            if (i + 1 >= argc) { fprintf(stderr, "%s requires an argument\n", argv[i]); return 2; }
            int cnt = parse_keylens(argv[i + 1], &keylens);
            if (cnt <= 0) { fprintf(stderr, "Invalid keylens spec: %s\n", argv[i + 1]); return 2; }
            keylens_count = cnt;
            i += 2;
        } else if (strcmp(argv[i], "-o") == 0 || strcmp(argv[i], "--output") == 0) {
            if (i + 1 >= argc) { fprintf(stderr, "%s requires an argument\n", argv[i]); return 2; }
            out_bin = argv[i + 1];
            i += 2;
        } else if (strcmp(argv[i], "-t") == 0 || strcmp(argv[i], "--text") == 0) {
            if (i + 1 >= argc) { fprintf(stderr, "%s requires an argument\n", argv[i]); return 2; }
            out_text = argv[i + 1];
            i += 2;
        } else if (strcmp(argv[i], "-h") == 0 || strcmp(argv[i], "--help") == 0) {
            usage(argv[0]); return 0;
        } else {
            fprintf(stderr, "Unknown option: %s\n", argv[i]); usage(argv[0]); return 2;
        }
    }

    /* remaining args are pair tokens */
    int remaining = argc - i;
    if (remaining < 1) { fprintf(stderr, "No cipher:crib pairs provided\n"); usage(argv[0]); return 2; }

    pair_t *pairs = calloc((size_t)remaining, sizeof(pair_t));
    if (!pairs) { fprintf(stderr, "OOM\n"); return 2; }

    int pair_count = 0;
    for (int j = 0; j < remaining; ++j) {
        const char *token = argv[i + j];
        if (parse_pair_token(token, &pairs[pair_count]) != 0) {
            fprintf(stderr, "Invalid pair token: %s (expected cipher:crib[:offset])\n", token);
            for (int k = 0; k < pair_count; ++k) free_pair(&pairs[k]);
            free(pairs);
            return 2;
        }
        pair_count++;
    }

    /* if no keylens specified, generate 1..MAX_AUTO_KEYLEN */
    if (keylens_count == 0) {
        keylens_count = (int)MAX_AUTO_KEYLEN;
        keylens = malloc((size_t)keylens_count * sizeof(size_t));
        if (!keylens) { fprintf(stderr, "OOM\n"); for (int k = 0; k < pair_count; ++k) free_pair(&pairs[k]); free(pairs); return 2; }
        for (int k = 0; k < keylens_count; ++k) keylens[k] = (size_t)(k + 1);
    }

    /* For each key length, process pairs and optionally write outputs.
     * If multiple key lengths are specified and an output filename is requested,
     * we will append the key_len to the filename when writing so results don't
     * clobber each other, unless the user provided only one keylen.
     */
    int exit_status = 0;
    for (int kl = 0; kl < keylens_count; ++kl) {
        size_t key_len = keylens[kl];
        /* build per-keylength output names if needed */
        char *bin_out = NULL;
        char *text_out = NULL;
        int write_output = (out_bin != NULL) || (out_text != NULL);
        if (write_output && keylens_count > 1) {
            /* append .<keylen>.bin or .<keylen>.txt */
            if (out_bin) {
                size_t need = strlen(out_bin) + 32;
                bin_out = malloc(need);
                if (!bin_out) { fprintf(stderr, "OOM\n"); exit_status = 1; break; }
                snprintf(bin_out, need, "%s.%zu.bin", out_bin, key_len);
            }
            if (out_text) {
                size_t need = strlen(out_text) + 32;
                text_out = malloc(need);
                if (!text_out) { fprintf(stderr, "OOM\n"); free(bin_out); exit_status = 1; break; }
                snprintf(text_out, need, "%s.%zu.txt", out_text, key_len);
            }
        } else {
            bin_out = out_bin ? strdup(out_bin) : NULL;
            text_out = out_text ? strdup(out_text) : NULL;
        }
        int r = process_keylen(key_len, pairs, (size_t)pair_count, write_output, bin_out, text_out);
        if (r != 0) exit_status = 1;
        free(bin_out);
        free(text_out);
    }

    /* cleanup */
    for (int k = 0; k < pair_count; ++k) free_pair(&pairs[k]);
    free(pairs);
    free(keylens);

    return exit_status;
}