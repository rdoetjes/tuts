#ifndef CONFIG_H
#define CONFIG_H

typedef struct {
    int baud;
    int do_read;
    int timeout_ms;
    int pause_ms;
    int dry_run;
    int fd;
    int num_digits;
    const char *device;
    const char *phone;
} Config;

#endif
