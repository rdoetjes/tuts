#ifndef BLUEBOX_ALSA_DTMF_H
#define BLUEBOX_ALSA_DTMF_H

/*
 * dtmf.h
 *
 * DTMF frequency lookup helpers for the bluebox-alsa project.
 *
 * The DTMF keypad is arranged as 4 rows × 4 columns. Each key is the sum
 * of one low-frequency and one high-frequency tone. This header exposes a
 * simple function to map a DTMF character to the two frequencies (in Hz).
 *
 * Supported characters:
 *   '0'..'9', '*', '#', 'A', 'B', 'C', 'D'  (case-insensitive for letters)
 *
 * Frequency table (Hz):
 *   Low  : 697, 770, 852, 941
 *   High : 1209, 1336, 1477, 1633
 *
 * Mapping:
 *   '1' -> 697  + 1209
 *   '2' -> 697  + 1336
 *   '3' -> 697  + 1477
 *   'A' -> 697  + 1633
 *   '4' -> 770  + 1209
 *   '5' -> 770  + 1336
 *   '6' -> 770  + 1477
 *   'B' -> 770  + 1633
 *   '7' -> 852  + 1209
 *   '8' -> 852  + 1336
 *   '9' -> 852  + 1477
 *   'C' -> 852  + 1633
 *   '*' -> 941  + 1209
 *   '0' -> 941  + 1336
 *   '#' -> 941  + 1477
 *   'D' -> 941  + 1633
 *
 * Usage:
 *   int f1 = 0, f2 = 0;
 *   if (bluebox::dtmf::get_dtmf_freqs('5', f1, f2)) {
 *       // f1 and f2 contain the frequencies in Hz
 *   }
 */

#include <cctype>

namespace bluebox {
namespace dtmf {

/*
 * Lookup the two frequencies (Hz) for a single DTMF character.
 *
 * Parameters:
 *  - digit: single character representing the DTMF key. For letters,
 *           case-insensitive ('a'..'d' allowed).
 *  - f1: output parameter set to the low frequency (Hz) on success.
 *  - f2: output parameter set to the high frequency (Hz) on success.
 *
 * Returns:
 *  - true if the digit is a valid DTMF key and f1/f2 were set.
 *  - false if the digit is not recognized.
 */
bool get_dtmf_freqs(char digit, int& f1, int& f2);

bool is_valid_dtmf(char digit);

} // namespace dtmf
} // namespace bluebox

#endif // BLUEBOX_ALSA_DTMF_H
