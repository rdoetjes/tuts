
#ifndef SERIAL_H
#define SERIAL_H

#include <termios.h>
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
#include <stdio.h>

/*Converts baud rate to termios speed_t */
speed_t baud_to_speed(int baud);

/* Open serial device and return fd or -1 */
int open_serial(const char *device);

/* Configure serial port for raw 8N1 no-flow-control. If speed == 0, do not set speeds. */
int configure_serial(int fd, speed_t speed);

/* Write all bytes robustly */
ssize_t write_all(int fd, const char *buf, size_t len);

/* Read until inactivity timeout_ms elapses. Caller must free *out_buf if bytes>0.
 * Returns bytes read (>0), 0 if none, -1 on error.
 */
ssize_t read_response(int fd, char **out_buf, int timeout_ms);

#endif
