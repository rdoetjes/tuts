/*
 * serial.cpp
 *
 * Implementation of simple serial command/response helpers used by the
 * bluebox-alsa project.
 *
 * Behavior:
 *  - open_serial opens and configures the serial device (9600 8N1, no hw flow control)
 *  - send_and_wait_serial writes a single byte command and waits up to timeout_ms
 *    for a single-byte response. Interprets ASCII '1' as success, '0' as failure.
 *
 * Note: callers are responsible for closing the returned file descriptor.
 */

#include "serial.h"

#include <cerrno>
#include <cstring>
#include <fcntl.h>
#include <iostream>
#include <sys/select.h>
#include <sys/time.h>
#include <termios.h>
#include <unistd.h>

namespace bluebox {
namespace serial {

int open_serial(const char* path) {
    if (path == nullptr) path = DEFAULT_DEVICE;

    int fd = ::open(path, O_RDWR | O_NOCTTY | O_NONBLOCK);
    if (fd < 0) {
        std::cerr << "open_serial: cannot open device '" << path << "': " << std::strerror(errno) << "\n";
        return -1;
    }

    struct termios tty;
    if (tcgetattr(fd, &tty) != 0) {
        std::cerr << "open_serial: tcgetattr failed for '" << path << "': " << std::strerror(errno) << "\n";
        ::close(fd);
        return -1;
    }

    // Set baud rate to 9600
    cfsetospeed(&tty, B9600);
    cfsetispeed(&tty, B9600);

    // 8N1, disable parity, one stop bit, 8 data bits
    tty.c_cflag = (tty.c_cflag & ~CSIZE) | CS8;
    tty.c_cflag &= ~PARENB;
    tty.c_cflag &= ~CSTOPB;

    // No hardware flow control
    tty.c_cflag &= ~CRTSCTS;

    // Enable receiver, local mode
    tty.c_cflag |= (CLOCAL | CREAD);

    // Raw input / output mode
    tty.c_lflag &= ~(ICANON | ECHO | ECHOE | ISIG);
    tty.c_iflag &= ~(IXON | IXOFF | IXANY | ICRNL | INLCR);
    tty.c_oflag &= ~OPOST;

    // Non-blocking read using select() for timeouts: VMIN=0, VTIME=0
    tty.c_cc[VMIN] = 0;
    tty.c_cc[VTIME] = 0;

    if (tcsetattr(fd, TCSANOW, &tty) != 0) {
        std::cerr << "open_serial: tcsetattr failed for '" << path << "': " << std::strerror(errno) << "\n";
        ::close(fd);
        return -1;
    }

    // Clear any preexisting data in the input queue (best-effort)
    {
        // Set non-blocking read temporarily if needed (fd already O_NONBLOCK)
        char buf[256];
        while (true) {
            ssize_t r = ::read(fd, buf, sizeof(buf));
            if (r <= 0) break;
        }
    }

    return fd;
}

int wait_and_read(int fd, int timeout_ms){
    // Prepare select() to wait for a readable byte
    fd_set rfds;
    FD_ZERO(&rfds);
    FD_SET(fd, &rfds);

    struct timeval tv;
    tv.tv_sec = timeout_ms / 1000;
    tv.tv_usec = (timeout_ms % 1000) * 1000;

    int sel = ::select(fd + 1, &rfds, nullptr, nullptr, &tv);
    if (sel < 0) {
        std::cerr << "send_and_wait_serial: select failed: " << std::strerror(errno) << "\n";
        return -1;
    } else if (sel == 0) {
        // timeout
        return -1;
    }

    if (FD_ISSET(fd, &rfds)) {
        char resp = 0;
        ssize_t r = ::read(fd, &resp, 1);
        if (r == 1) {
            if (resp == 'H') return 1;
            if (resp == '0') return 0;
            // Some devices may append newline; trim common whitespace and re-evaluate.
            // If unexpected byte, attempt to read until newline or timeout quickly (best-effort)
            // but for simplicity we treat unexpected as error.
            std::cerr << "send_and_wait_serial: unexpected response byte: 0x" << std::hex
                      << (static_cast<int>(static_cast<unsigned char>(resp)) & 0xFF) << std::dec << "\n";
            return -1;
        } else if (r == 0) {
            // EOF-ish
            return -1;
        } else {
            std::cerr << "send_and_wait_serial: read error: " << std::strerror(errno) << "\n";
            return -1;
        }
    }

    return -1;
}

int send_and_wait_serial(int fd, char cmd, int timeout_ms) {
    if (fd < 0) {
        return -1;
    }

    // Write the single byte command
    ssize_t w = ::write(fd, &cmd, 1);
    if (w != 1) {
        std::cerr << "send_and_wait_serial: write failed: " << std::strerror(errno) << "\n";
        return -1;
    }

    return wait_and_read(fd, timeout_ms);

}

} // namespace serial
} // namespace bluebox
