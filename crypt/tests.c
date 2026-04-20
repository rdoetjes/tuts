#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <inttypes.h>

/*
 * Tests for encrypt/decrypt roundtrip.
 *
 * Assumes the following functions are available (from main.c):
 *
 *   char *encrypt(const char *message, size_t message_len);
 *   char *decrypt(const char *message, size_t message_len);
 *
 * - encrypt returns a newly-allocated ciphertext buffer of message_len bytes
 *   (binary). Caller must free().
 * - decrypt returns a newly-allocated NUL-terminated plaintext string of
 *   length message_len. Caller must free().
 *
 * These tests exercise:
 *  - empty messages
 *  - ASCII strings
 *  - strings with embedded NUL bytes
 *  - long messages (longer than the key)
 *  - binary data (0..255)
 *
 * The test runner treats an encrypt returning NULL for a zero-length input as
 * acceptable (some calloc implementations return NULL for size 0). For any
 * non-zero lengths, encrypt/decrypt must succeed.
 */

/* Prototypes from the implementation under test */
char *encrypt(const char *message, size_t message_len);
char *decrypt(const char *message, size_t message_len);

static int tests_passed = 0;
static int tests_failed = 0;

static void report_result(const char *name, int ok) {
    if (ok) {
        printf("[PASS] %s\n", name);
        ++tests_passed;
    } else {
        printf("[FAIL] %s\n", name);
        ++tests_failed;
    }
}

/* Generic test helper */
static void run_roundtrip_test(const char *name, const unsigned char *plain, size_t len) {
    int ok = 0;

    /* For len==0 we still call with a valid pointer (non-NULL) where possible. */
    const char *input_ptr = (const char *)plain;
    if (input_ptr == NULL && len == 0) input_ptr = ""; /* safe pointer for empty test */

    char *cipher = encrypt(input_ptr, len);
    if (cipher == NULL) {
        if (len == 0) {
            /* Some implementations may return NULL for 0-length allocation; accept that. */
            ok = 1;
            report_result(name, ok);
            return;
        } else {
            report_result(name, 0);
            return;
        }
    }

    /* decrypt should return message_len bytes (plus a NUL terminator) */
    char *dec = decrypt(cipher, len);
    if (dec == NULL) {
        /* decrypt failed */
        free(cipher);
        report_result(name, 0);
        return;
    }

    /* Compare binary memory (dec may be NUL-terminated) */
    if (len == 0) {
        /* expect an empty string */
        ok = (dec[0] == '\0');
    } else {
        ok = (memcmp(dec, plain, len) == 0);
    }

    if (!ok) {
        /* On failure, print diagnostics for debugging */
        printf("  Test '%s' diagnostics: len=%zu\n", name, len);
        printf("  Plain: ");
        for (size_t i = 0; i < len; ++i) {
            printf("%02X ", (unsigned char)plain[i]);
            if (i >= 64) { printf("..."); break; }
        }
        printf("\n  Cipher: ");
        for (size_t i = 0; i < (len < 64 ? len : 64); ++i) {
            printf("%02X ", (unsigned char)cipher[i]);
        }
        if (len > 64) printf("...");
        printf("\n  Dec: ");
        for (size_t i = 0; i < len; ++i) {
            printf("%02X ", (unsigned char)dec[i]);
            if (i >= 64) { printf("..."); break; }
        }
        printf("\n");
    }

    free(cipher);
    free(dec);

    report_result(name, ok);
}


int main(void) {
    /* 1) Empty string */
    {
        const unsigned char empty[] = "";
        run_roundtrip_test("empty-string", empty, 0);
    }

    /* 2) Short ASCII */
    {
        const unsigned char s[] = "hello world";
        run_roundtrip_test("short-ascii", s, strlen((const char*)s));
    }

    /* 3) Embedded NUL */
    {
        const unsigned char s[] = { 'a', 'b', 0x00, 'c', 'd' };
        run_roundtrip_test("embedded-nul", s, sizeof(s));
    }

    /* 4) Long message (repeat a sentence to exceed key length) */
    {
        const char *base =
            "This is a long test message used to exceed the key length. "
            "It contains multiple repetitions and should force the key to wrap. ";
        /* make ~1024 bytes */
        size_t repeats = 1024 / strlen(base) + 2;
        size_t len = repeats * strlen(base);
        unsigned char *buf = malloc(len);
        if (!buf) {
            fprintf(stderr, "OOM during test setup\n");
            return 2;
        }
        for (size_t i = 0, off = 0; i < repeats; ++i) {
            memcpy(buf + off, base, strlen(base));
            off += strlen(base);
        }
        run_roundtrip_test("long-repeated", buf, len);
        free(buf);
    }

    /* 5) Binary 0..255 (truncated to 256 bytes) */
    {
        size_t len = 256;
        unsigned char *buf = malloc(len);
        if (!buf) {
            fprintf(stderr, "OOM during test setup\n");
            return 2;
        }
        for (size_t i = 0; i < len; ++i) buf[i] = (unsigned char)i;
        run_roundtrip_test("binary-0-255", buf, len);
        free(buf);
    }

    /* 6) Deterministic pseudo-random bytes */
    {
        size_t len = 15000;
        unsigned char *buf = malloc(len);
        if (!buf) {
            fprintf(stderr, "OOM during test setup\n");
            return 2;
        }
        unsigned int seed = 12345;
        for (size_t i = 0; i < len; ++i) {
            /* simple LCG for deterministic bytes */
            seed = (1103515245u * seed + 12345u) & 0x7fffffff;
            buf[i] = (unsigned char)(seed & 0xFF);
        }
        run_roundtrip_test("deterministic-rand-15000", buf, len);
        free(buf);
    }

    /* 7) Many small messages in a batch (key reuse patterns) */
    {
        const unsigned char *msgs[] = {
            (const unsigned char *)"a",
            (const unsigned char *)"ab",
            (const unsigned char *)"abc",
            (const unsigned char *)"abcd",
            (const unsigned char *)"abcde",
            (const unsigned char *)"abcdef",
            (const unsigned char *)"abcdefg"
        };
        for (size_t i = 0; i < sizeof(msgs)/sizeof(msgs[0]); ++i) {
            run_roundtrip_test("batch-small", msgs[i], strlen((const char*)msgs[i]));
        }
    }

    printf("\nTest summary: %d passed, %d failed\n", tests_passed, tests_failed);
    return tests_failed == 0 ? 0 : 1;
}
