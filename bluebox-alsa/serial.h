#ifndef BLUEBOX_ALSA_SERIAL_H
#define BLUEBOX_ALSA_SERIAL_H

#include <cstdint>
#include <string>

namespace bluebox {
namespace serial {

/*
 * Default serial device path used when none is supplied on the CLI.
 */
inline constexpr const char* DEFAULT_DEVICE = "/dev/ttyUSB0";

/*
 * Open and configure a serial device for simple command/response usage.
 *
 * Behavior:
 *  - Opens the device path for read/write.
 *  - Configures common settings (9600 baud, 8N1, no hardware flow control).
 *  - Non-blocking mode is used; reads are performed using select() with timeouts.
 *
 * Returns:
 *  - A non-negative file descriptor on success.
 *  - -1 on error (caller may consult errno to determine cause).
 *
 * Note: The implementation is responsible for calling ::close(fd) when done.
 */
int open_serial(const char* path);

/*
 * Send a single-byte command over the given serial file descriptor and wait
 * for a single-byte response with a timeout.
 *
 * Parameters:
 *  - fd:     Open serial file descriptor (as returned from open_serial).
 *  - cmd:    Single command byte to write (e.g. 'H').
 *  - timeout_ms: Maximum time to wait for a 1-byte response in milliseconds.
 *
 * Response interpretation:
 *  - Returns  1 if the device responded with ASCII '1' (success/OK).
 *  - Returns  0 if the device responded with ASCII '0' (failure).
 *  - Returns -1 on timeout, write/read error, or unexpected response.
 *
 * The function performs exactly one write of the command byte and attempts
 * to read a single byte as reply (using select/poll internally to implement
 * the timeout).
 */
int send_and_wait_serial(int fd, char cmd, int timeout_ms);

/*
 *  * Parameters:
 *  - fd:     Open serial file descriptor (as returned from open_serial).
 *  - cmd:    Single command byte to write (e.g. 'H').
 *  - timeout_ms: Maximum time to wait for a 1-byte response in milliseconds.
 *
 *   Response interpretation:
 *  - Returns  1 if the device responded with ASCII '1' (success/OK).
 *  - Returns  0 if the device responded with ASCII '0' (failure).
 *  - Returns -1 on timeout, write/read error, or unexpected response.
 */
int wait_and_read(int fd, int timeout_ms);
} // namespace serial
} // namespace bluebox

#endif // BLUEBOX_ALSA_SERIAL_H
