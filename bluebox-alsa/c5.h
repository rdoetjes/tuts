#ifndef BLUEBOX_ALSA_C5_H
#define BLUEBOX_ALSA_C5_H

/*
 * c5.h
 *
 * CCITT / C5 tone lookup helpers for bluebox-alsa.
 *
 * This header declares functions to map C5-style codes (either single
 * digit symbols or named codes such as "KP1", "KP2", "ST", "SEIZE", etc.)
 * to their corresponding frequency pairs (in Hz).
 *
 * The implementation is expected to:
 *  - Return two frequencies (f1, f2). A frequency value of 0 may be used
 *    for single-frequency tones (e.g. seizure/answer tones if modeled that way).
 *  - Return true when the provided code is recognized and f1/f2 are set.
 *  - Return false when the code is not recognized.
 *
 * Supported lookup forms:
 *  - Single digits '0'..'9' map to predefined frequency pairs for certain C5 sets.
 *  - Named codes (case-insensitive), for example:
 *      KP1, KP2, ST, CODE11, CODE12, SEIZE, PROCEED (or PK), ANSWER,
 *      BUSY (or BUSYFLASH), CLEARBACK, CLEARFWD (or RELGUARD)
 *
 * Example usage:
 *   int f1 = 0, f2 = 0;
 *   if (bluebox::c5::get_c5_freqs("KP1", f1, f2)) {
 *       // use f1/f2
 *   }
 */

#include <string>

namespace bluebox {
namespace c5 {

/*
 * Map a C5 code string to two frequencies (Hz).
 *
 * Parameters:
 *  - code: string containing the code (case-insensitive). For single digit
 *          codes the string may be a single character like "3".
 *  - f1: output parameter set to the first (typically lower) frequency in Hz.
 *  - f2: output parameter set to the second (typically higher) frequency in Hz.
 *
 * Returns:
 *  - true if the code was recognized and f1/f2 were set.
 *  - false if the code is unknown.
 */
bool get_c5_freqs(const std::string& code, int& f1, int& f2);

/*
 * Helper: returns true if the provided code string is a supported C5 code.
 * This is implemented as a utility predicate (may be implemented inline by
 * the .cpp or here as a simple wrapper around get_c5_freqs).
 */
bool is_valid_c5(const std::string& code);

} // namespace c5
} // namespace bluebox

#endif // BLUEBOX_ALSA_C5_H