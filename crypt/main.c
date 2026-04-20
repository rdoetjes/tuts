#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

/*
 * Simple XOR + per-byte 8-bit rotate obfuscator.
 *
 * For each byte i:
 *   kidx  = (start + i) % key_len
 *   shift = ((unsigned char)KEY[kidx] + (unsigned char)i) % 8
 *   tmp   = (unsigned char)plaintext[i] ^ (unsigned char)KEY[kidx]
 *   cipher[i] = rol8(tmp, shift)
 *
 * Decryption reverses the steps:
 *   tmp = ror8(cipher[i], shift)
 *   plaintext[i] = tmp ^ KEY[kidx]
 *
 * Notes:
 * - The rotation amount must be derivable by both sides without knowing the
 *   plaintext; here we use key byte + index. Do NOT derive the rotation from
 *   the plaintext itself.
 * - Ciphertext is binary and may contain NUL bytes: always track lengths
 *   explicitly and never rely on strlen() for ciphertext.
 * - This is only an obfuscation layer on top of XOR. It does NOT make the
 *   construction cryptographically secure. Prefer AEAD ciphers for real use.
 */

char *KEY =
    "ThKUHIT678t2bkhabxI^&TI*gyhu2lh3jlhbiut987OYHLIUJBLWEUDG3278OYGHLIUQHBSQlHUTI&^tgJHBHJ6TGKUYGKJYiguykjgiyuti7^TIGyukhb";

/* 8-bit rotate left */
static inline uint8_t rol8(uint8_t v, unsigned int s) {
    s &= 7u; /* modulo 8 */
    if (s == 0) return v;
    return (uint8_t)((v << s) | (v >> (8 - s)));
}

/* 8-bit rotate right */
static inline uint8_t ror8(uint8_t v, unsigned int s) {
    s &= 7u;
    if (s == 0) return v;
    return (uint8_t)((v >> s) | (v << (8 - s)));
}

/* Encrypt a buffer of explicit length. Returns newly-allocated ciphertext buffer
 * of exactly message_len bytes. Caller must free(). */
char *encrypt(const char *message, size_t message_len) {
    if (!message && message_len != 0) return NULL;

    size_t key_len = strlen(KEY);
    if (key_len == 0) return NULL;

    char *enc = calloc(1, message_len);
    if (!enc) return NULL;

    size_t start = message_len % key_len;
    for (size_t i = 0; i < message_len; ++i) {
        size_t kidx = (start + i) % key_len;
        uint8_t keybyte = (uint8_t)KEY[kidx];
        uint8_t plain = (uint8_t)message[i];
        unsigned int shift = (unsigned int)((keybyte + (uint8_t)i) % 8u);
        uint8_t tmp = (uint8_t)(plain ^ keybyte);
        enc[i] = (char)rol8(tmp, shift);
    }
    return enc;
}

/* Decrypt a ciphertext buffer of explicit length. Returns a newly-allocated
 * NUL-terminated plaintext string (length message_len). Caller must free(). */
char *decrypt(const char *message, size_t message_len) {
    if (!message && message_len != 0) return NULL;

    size_t key_len = strlen(KEY);
    if (key_len == 0) return NULL;

    char *dec = malloc(message_len + 1);
    if (!dec) return NULL;

    size_t start = message_len % key_len;
    for (size_t i = 0; i < message_len; ++i) {
        size_t kidx = (start + i) % key_len;
        uint8_t keybyte = (uint8_t)KEY[kidx];
        unsigned int shift = (unsigned int)((keybyte + (uint8_t)i) % 8u);
        uint8_t cipher = (uint8_t)message[i];
        uint8_t tmp = ror8(cipher, shift);
        dec[i] = (char)(tmp ^ keybyte);
    }
    dec[message_len] = '\0';
    return dec;
}

#ifndef TESTING
int main(void) {
    const char *message_in =
        "Okay, how long can this message get before it starts to act all weird? "
        "Probably longer than I had imagined this seems to work okay";

    size_t message_len = strlen(message_in);

    char *enc = encrypt(message_in, message_len);
    if (!enc) {
        perror("encrypt");
        return 1;
    }

    /* Print ciphertext as hex (explicit length) */
    for (size_t i = 0; i < message_len; ++i) {
        printf("0x%02X ", (unsigned char)enc[i]);
    }
    printf("\n");

    /* Decrypt using the explicit length */
    char *dec = decrypt(enc, message_len);
    if (!dec) {
        perror("decrypt");
        free(enc);
        return 1;
    }

    printf("%s\n", dec);

    free(enc);
    free(dec);
    return 0;
}
#endif
