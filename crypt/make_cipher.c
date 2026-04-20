/*
 * crypt/make_cipher.c
 *
 * Small helper to produce a ciphertext file and a crib (known plaintext fragment)
 * using the same XOR+rotate obfuscation used by the repository's encrypt()
 * implementation. This is intended to generate demo data for the recovery tool.
 *
 * Usage:
 *   make_cipher <cipher_out> <crib_out> <crib_offset> [crib_len] [message_file]
 *
 * - cipher_out   : path to write ciphertext (binary)
 * - crib_out     : path to write the crib (binary/text)
 * - crib_offset  : 0-based offset inside the plaintext where the crib begins
 * - crib_len     : optional, number of bytes to extract for the crib (default 32)
 * - message_file : optional path to a file containing the plaintext message; if
 *                  omitted a built-in long demo string will be used.
 *
 * Example:
 *   ./make_cipher cipher.bin crib.txt 10 32
 *
 * After running successfully this prints a one-line token you can pass to
 * the recovery tool:
 *   cipher.bin:crib.txt:<offset>
 *
 * It also prints the key length you should try when running the recovery tool
 * (the length of the KEY used below).
 *
 * NOTE: This tool implements the same transform as 'encrypt' in the repo:
 *   start = message_len % key_len
 *   kidx  = (start + i) % key_len
 *   shift = (KEY[kidx] + i) % 8
 *   tmp   = plaintext[i] ^ KEY[kidx]
 *   cipher[i] = rol8(tmp, shift)
 *
 * This file is intentionally self-contained so it can be built independently:
 *   cc -std=c11 -O2 -o make_cipher make_cipher.c
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <errno.h>

static const char *KEY =
    "ThKUHIT678t2bkhabxI^&TI*gyhu2lh3jlhbiut987OYHLIUJBLWEUDG3278OYGHLIUQHBSQlHUTI&^tgJHBHJ6TGKUYGKJYiguykjgiyuti7^TIGyukhb";

static inline uint8_t rol8(uint8_t v, unsigned int s) {
    s &= 7u;
    if (s == 0) return v;
    return (uint8_t)((v << s) | (v >> (8 - s)));
}

static char *read_whole_file(const char *path, size_t *out_len) {
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
    char *buf = malloc((size_t)s + 1);
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
    buf[got] = '\0';
    if (out_len) *out_len = got;
    return buf;
}

static int write_file(const char *path, const void *data, size_t len) {
    FILE *f = fopen(path, "wb");
    if (!f) {
        fprintf(stderr, "fopen(%s): %s\n", path, strerror(errno));
        return -1;
    }
    size_t w = fwrite(data, 1, len, f);
    if (w != len) {
        fprintf(stderr, "fwrite(%s) failed\n", path);
        fclose(f);
        return -1;
    }
    fclose(f);
    return 0;
}

int main(int argc, char **argv) {
    if (argc < 4) {
        fprintf(stderr, "Usage: %s <cipher_out> <crib_out> <crib_offset> [crib_len] [message_file]\n", argv[0]);
        return 2;
    }
    const char *cipher_out = argv[1];
    const char *crib_out = argv[2];
    char *endptr = NULL;
    errno = 0;
    long off = strtol(argv[3], &endptr, 10);
    if (errno || *endptr != '\0' || off < 0) {
        fprintf(stderr, "Invalid crib_offset: %s\n", argv[3]);
        return 2;
    }
    size_t crib_offset = (size_t)off;

    size_t crib_len = 32;
    if (argc >= 5) {
        errno = 0;
        long cl = strtol(argv[4], &endptr, 10);
        if (!errno && *endptr == '\0' && cl >= 0) crib_len = (size_t)cl;
    }

    char *message = NULL;
    size_t message_len = 0;
    if (argc >= 6) {
        message = read_whole_file(argv[5], &message_len);
        if (!message) return 3;
    } else {
        /* default built-in message */
        const char *demo =
            "This is a demo plaintext used to generate ciphertext for the recovery tool. "
            "It contains ASCII, punctuation, and is long enough to wrap the key. "
            "Using a demo message helps illustrate the recovery process.";
        message_len = strlen(demo);
        message = malloc(message_len + 1);
        if (!message) {
            fprintf(stderr, "OOM\n");
            return 3;
        }
        memcpy(message, demo, message_len + 1);
    }

    if (crib_offset > message_len) {
        fprintf(stderr, "Error: crib_offset %zu beyond message length %zu\n", crib_offset, message_len);
        free(message);
        return 4;
    }
    if (crib_offset + crib_len > message_len) {
        crib_len = message_len - crib_offset; /* adjust */
    }

    size_t key_len = strlen(KEY);
    char *cipher = calloc(1, message_len);
    if (!cipher) {
        fprintf(stderr, "OOM allocating cipher buffer\n");
        free(message);
        return 5;
    }

    size_t start = message_len % key_len;
    for (size_t i = 0; i < message_len; ++i) {
        size_t kidx = (start + i) % key_len;
        uint8_t keybyte = (uint8_t)KEY[kidx];
        uint8_t plain = (uint8_t)message[i];
        unsigned int shift = (unsigned int)((keybyte + (uint8_t)i) % 8u);
        uint8_t tmp = (uint8_t)(plain ^ keybyte);
        cipher[i] = (char)rol8(tmp, shift);
    }

    /* write ciphertext and crib */
    if (write_file(cipher_out, cipher, message_len) != 0) {
        fprintf(stderr, "Failed to write ciphertext\n");
        free(cipher);
        free(message);
        return 6;
    }
    if (write_file(crib_out, message + crib_offset, crib_len) != 0) {
        fprintf(stderr, "Failed to write crib\n");
        free(cipher);
        free(message);
        return 7;
    }

    /* Print info for the user and how to call recovery tool */
    printf("Wrote ciphertext: %s (%zu bytes)\n", cipher_out, message_len);
    printf("Wrote crib (plaintext fragment): %s (offset=%zu, len=%zu)\n", crib_out, crib_offset, crib_len);
    printf("KEY length to try: %zu\n", key_len);
    printf("Recovery token (use with recover_key): %s:%s:%zu\n", cipher_out, crib_out, crib_offset);

    free(cipher);
    free(message);
    return 0;
}