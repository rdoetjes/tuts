/*
 * main.c
 *
 * Dial a phone number and brute-force DTMF sequences in batches of N iterations per call.
 *
 * Features:
 *  - getopt-based CLI:
 *      -d DEVICE   serial device (default: /dev/ttyS1)
 *      -b BAUD     baud rate (default: 57600)
 *      -r          read modem response after each send
 *      -t MS       read inactivity timeout in milliseconds (default: 1000)
 *      -p MS       pause between calls in milliseconds (default: 1000)
 *      -T          dry-run: print commands to screen but do not send to serial
 *
 * Build:
 *   run make
 *     or manualy
 *   gcc -std=gnu99 -Wall -Wextra -o main main.c
 *
 * Use responsibly.phone
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <unistd.h>
#include <fcntl.h>
#include <termios.h>
#include <errno.h>
#include <stdint.h>
#include <sys/select.h>
#include <sys/time.h>
#include <sys/types.h>
#include <time.h>

#include "config.h"
#include "serial.h"

#define DEFAULT_DEVICE "/dev/ttyS1"
#define DEFAULT_BAUD 57600
#define DEFAULT_TIMEOUT_MS 1000
#define DEFAULT_PAUSE_MS 1000

/* allocation heuristics */
#define APPROX_BASE 64
#define APPROX_PER_2DIG 32
#define APPROX_PER_3DIG 48

void sleep_us(unsigned long microseconds) {
    struct timespec ts;
    ts.tv_sec = microseconds / 1000000UL;
    ts.tv_nsec = (microseconds % 1000000UL) * 1000;
    nanosleep(&ts, NULL);
}

static void usage(const char *prog) {
    fprintf(stderr,
            "Usage: %s [options] <digits:2|3> <phonenumber>\n"
            "Options:\n"
            "  -d DEVICE   serial device (default: %s)\n"
            "  -b BAUD     baud rate (default: %d)\n"
            "  -r          read modem response after sending\n"
            "  -t MS       read inactivity timeout in milliseconds (default: %d)\n"
            "  -p MS       pause between calls in milliseconds (default: %d)\n"
            "  -T          dry-run: print commands to screen but do not send to serial\n",
            prog, DEFAULT_DEVICE, DEFAULT_BAUD, DEFAULT_TIMEOUT_MS, DEFAULT_PAUSE_MS);
}


/* Ensure the buffer has at least `need` bytes capacity. Resize by doubling. */
static int ensure_capacity(char **bufp, size_t *capp, size_t need) {
    if (*capp >= need) return 0;
    size_t newcap = (*capp == 0) ? 256 : *capp * 2;
    while (newcap < need) newcap *= 2;
    char *nb = realloc(*bufp, newcap);
    if (!nb) return -1;
    *bufp = nb;
    *capp = newcap;
    return 0;
}

/* Append formatted data into buf at position *pos, growing the buffer as needed.
 * Returns number of bytes appended (excluding terminating NUL), or -1 on error.
 */
static int appendf(char **bufp, size_t *capp, size_t *pos, const char *fmt, ...) {
    va_list ap;
    va_start(ap, fmt);
    va_list ap2;
    va_copy(ap2, ap);
    int needed = vsnprintf(NULL, 0, fmt, ap2);
    va_end(ap2);
    if (needed < 0) { va_end(ap); return -1; }
    size_t need_total = *pos + (size_t)needed + 1;
    if (ensure_capacity(bufp, capp, need_total) != 0) { va_end(ap); return -1; }
    int written = vsnprintf(*bufp + *pos, *capp - *pos, fmt, ap);
    va_end(ap);
    if (written < 0) return -1;
    *pos += (size_t)written;
    return written;
}

/* Build a single call's AT command that dials phone and sends combos[start..end-1].
 * Returns a malloc'd string (caller must free) or NULL on error.
 *
 * Format: ATDT<phone>;AT+VTD=1; (for each combo) AT+VTS=x;... ,; \r
 */
static char *build_batch_command(int num_digits, const char *phone, int start, int end) {
    printf("Building batch command for phone: %s, range: %d-%d\n", phone, start, end);
    if (!phone || start >= end) {
        fprintf(stderr, "Invalid input for build_batch_command\n");
        return NULL;
    }
    int combos = end - start;
    size_t phone_len = strlen(phone);
    size_t approx_per = (num_digits == 3) ? APPROX_PER_3DIG : APPROX_PER_2DIG;
    size_t cap = APPROX_BASE + phone_len + (size_t)combos * approx_per + 8;

    char *buf = malloc(cap);
    if (!buf) {
        fprintf(stderr, "Failed to allocate memory for batch command\n");
        return NULL;
    }

    size_t pos = 0;
    int r = appendf(&buf, &cap, &pos, "ATDT%s;AT+VTD=1;", phone);
    if (r < 0) { fprintf(stderr, "append error\n"); free(buf); return NULL; }

    for (int idx = start; idx < end; ++idx) {
        if (num_digits == 3) {
            if (appendf(&buf, &cap, &pos, "AT+VTS=%d;AT+VTS=%d;AT+VTS=%d,;",
                        idx / 100, (idx / 10) % 10, idx % 10) < 0) {
                free(buf);
                return NULL;
            }
        } else {
            if (appendf(&buf, &cap, &pos, "AT+VTS=%d;AT+VTS=%d,;", idx / 10, idx % 10) < 0) {
                free(buf);
                return NULL;
            }
        }
    }

    if (ensure_capacity(&buf, &cap, pos + 2) != 0) { free(buf); return NULL; }
    buf[pos++] = '\r';
    buf[pos] = '\0';

    return buf;
}

/* Sleep for ms milliseconds */
static void msleep(int ms) {
    if (ms <= 0) return;
    usleep((useconds_t)ms * 1000);
}

/*
 * Dial and brute force DTMF
 */
void dtmf_out(Config *c){
    const int batch_size = 3;
    const int max_combos = (c->num_digits == 3) ? 1000 : 100;

    for (int i = 0; i < max_combos; i += batch_size) {
        int end = i + batch_size;
        if (end > max_combos) end = max_combos;

        char *cmd = build_batch_command(c->num_digits, c->phone, i, end);
        if (!cmd) { fprintf(stderr, "Failed to build batch command for combos %d..%d\n", i, end-1); break; }

        if (c->dry_run) {
            printf("DRY-RUN: batch %d..%d command:\n%s\n", i, end-1, cmd);
            free(cmd);
            continue;
        }

        printf("Sending batch %d..%d command:\n%s\n", i, end-1, cmd);
        ssize_t wrote = write_all(c->fd, cmd, strlen(cmd));
        printf("Send done\n\n");

        if (wrote < 0) { fprintf(stderr, "Failed to send command for combos %d..%d\n", i, end-1); free(cmd); break; }
        printf("Sent call for combos %d..%d (%zd bytes)\n", i, end-1, wrote);

        if (!c->do_read) {
            char *resp = NULL;
            ssize_t r = read_response(c->fd, &resp, c->timeout_ms);
            if (r < 0) {
                fprintf(stderr, "Error reading response after sending batch %d..%d\n", i, end-1);
            } else if (r == 0) {
                printf("No modem response for batch %d..%d within %d ms\n", i, end-1, c->timeout_ms);
            } else {
                printf("Modem response (%zd bytes):\n%s\n", r, resp);
                free(resp);
            }
        }
        free(cmd);
    }
}

/* Allow -T anywhere by pre-scanning argv for "-T" entries and removing them
 * before getopt runs. This makes -T usable even after positional args.
 */
int main(int argc, char *argv[]) {
    int opt;
    Config c;
    c.device = DEFAULT_DEVICE;
    c.baud = DEFAULT_BAUD;
    c.do_read = 0;
    c.timeout_ms = DEFAULT_TIMEOUT_MS;
    c.pause_ms = DEFAULT_PAUSE_MS;
    c.dry_run = 0;

    /* Pre-scan argv and remove any "-T" occurrences so getopt will parse options
     * reliably even if -T appears after non-option arguments. */
    for (int i = 1; i < argc; ++i) {
        if (strcmp(argv[i], "-T") == 0) {
           c. dry_run = 1;
            /* shift left */
            for (int j = i; j < argc - 1; ++j) argv[j] = argv[j + 1];
            argc--;
            argv[argc] = NULL;
            i--; /* re-check this index after shift */
            if (argc <= 1) break;
        }
    }

    /* Now parse remaining options normally */
    while ((opt = getopt(argc, argv, "d:b:rt:p:T")) != -1) {
        switch (opt) {
            case 'd': c.device = optarg; break;
            case 'r': c.do_read = 1; break;

            case 'b': c.baud = atoi(optarg);
            if (c.baud <= 0) { fprintf(stderr, "Invalid baud\n"); return 1; } break;

            case 't': c.timeout_ms = atoi(optarg);
            if (c.timeout_ms < 0) { fprintf(stderr, "Invalid timeout\n"); return 1; } break;

            case 'p': c.pause_ms = atoi(optarg);
            if (c.pause_ms < 0) { fprintf(stderr, "Invalid pause\n"); return 1; } break;

            case 'T': c.dry_run = 1; break; /* also handle if user specified it normally */
            default: usage(argv[0]); return 1;
        }
    }

    if (optind + 2 != argc) { usage(argv[0]); return 1; }

    c.num_digits = atoi(argv[optind]);
    if (c.num_digits != 2 && c.num_digits != 3) { fprintf(stderr, "Digits must be 2 or 3\n"); return 1; }
    c.phone = argv[optind + 1];

    c.fd = -1;
    speed_t sp = baud_to_speed(c.baud);
    if (sp == (speed_t)0) fprintf(stderr, "Warning: baud %d not mapped to termios constant; attempting to configure anyway\n", c.baud);
    if (!c.dry_run) {
        c.fd = open_serial(c.device);
        if (c.fd < 0) return 1;
        if (configure_serial(c.fd, sp) != 0) { close(c.fd); return 1; }
    }

    dtmf_out(&c);

    /* Hang up the call */
    printf("Hang up the call\n");
    const char *hang = "ATH\r";
    ssize_t hw = write_all(c.fd, hang, strlen(hang));
    if (hw < 0) { fprintf(stderr, "Failed to send hangup\n"); }
    if (tcflush(c.fd, TCIOFLUSH) != 0) perror("tcflush");
    msleep(c.pause_ms);

    if (c.fd >= 0) close(c.fd);
    return 0;
}
