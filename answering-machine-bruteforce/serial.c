#include "serial.h"

/* Map common integer baud to termios speed_t (0 => unsupported) */
speed_t baud_to_speed(int baud) {
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

/* Open serial device and return fd or -1 */
int open_serial(const char *device) {
    int fd = open(device, O_RDWR | O_NOCTTY | O_SYNC);
    if (fd < 0) {
        fprintf(stderr, "open(%s): %s\n", device, strerror(errno));
        return -1;
    }
    return fd;
}

/* Configure serial port for raw 8N1 no-flow-control. If speed == 0, do not set speeds. */
int configure_serial(int fd, speed_t speed) {
    struct termios tty;
    if (tcgetattr(fd, &tty) != 0) { perror("tcgetattr"); return -1; }

    if (speed != (speed_t)0) {
        if (cfsetispeed(&tty, speed) != 0 || cfsetospeed(&tty, speed) != 0) {
            perror("cfset[io]speed");  // note: you might want to return -1 here
        }
    }

    // --- cflag: 8N1, no hw flow, ignore modem lines ---
    tty.c_cflag &= ~PARENB;
    tty.c_cflag &= ~CSTOPB;
    tty.c_cflag &= ~CSIZE;
    tty.c_cflag |= CS8;
    tty.c_cflag &= ~CRTSCTS;
    tty.c_cflag |= (CLOCAL | CREAD);

    // --- iflag: DISABLE all input processing ICRNL WAS REQUIRED FOR LOOP BACK READ! IN CASE OF ACTUAL MODEM YOU MAY NEED TO REMOVE ICRNL ---
    tty.c_iflag &= ~(IGNBRK | BRKINT | PARMRK | ISTRIP | INLCR | IGNCR | ICRNL | IXON | IXOFF);

    // --- lflag: already good, but make sure fully raw ---
    tty.c_lflag &= ~(ICANON | ECHO | ECHOE | ECHONL | ISIG | IEXTEN);

    // --- oflag: raw output ---
    tty.c_oflag &= ~OPOST;

    // --- timeouts: your current setting is ok for select() ---
    tty.c_cc[VMIN]  = 0;
    tty.c_cc[VTIME] = 10;   // 1 second inter-char timeout

    if (tcsetattr(fd, TCSANOW, &tty) != 0) {
        perror("tcsetattr");
        return -1;
    }

    // Bonus: flush any garbage before starting
    tcflush(fd, TCIOFLUSH);

    return 0;
}
/* Write all bytes robustly */
ssize_t write_all(int fd, const char *buf, size_t len) {
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

        fprintf(stderr, "write_all: len: %zd wrote: %zd bytes total: %zd\n", len, w, total);
    }
    //if (tcflush(fd, TCIOFLUSH) != 0) perror("tcflush");
    return (ssize_t)total;
}

/* Read until inactivity timeout_ms elapses. Caller must free *out_buf if bytes>0.
 * Returns bytes read (>0), 0 if none, -1 on error.
 */
ssize_t read_response(int fd, char **out_buf, int timeout_ms) {
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
