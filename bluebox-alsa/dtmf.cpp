/*
 * dtmf.cpp
 *
 * Implementation of DTMF frequency lookup for bluebox-alsa.
 *
 * Provides get_dtmf_freqs(char digit, int& f1, int& f2) which maps a single
 * DTMF key to its two constituent frequencies (in Hz).
 *
 * Supported keys:
 *   '0'..'9', '*', '#', 'A', 'B', 'C', 'D'  (letters are case-insensitive)
 */

#include "dtmf.h"
#include <cctype>

namespace bluebox {
namespace dtmf {

bool get_dtmf_freqs(char digit, int& f1, int& f2) {
    switch (std::toupper(static_cast<unsigned char>(digit))) {
        case '1': f1 = 697;  f2 = 1209; return true;
        case '2': f1 = 697;  f2 = 1336; return true;
        case '3': f1 = 697;  f2 = 1477; return true;
        case '4': f1 = 770;  f2 = 1209; return true;
        case '5': f1 = 770;  f2 = 1336; return true;
        case '6': f1 = 770;  f2 = 1477; return true;
        case '7': f1 = 852;  f2 = 1209; return true;
        case '8': f1 = 852;  f2 = 1336; return true;
        case '9': f1 = 852;  f2 = 1477; return true;
        case '0': f1 = 941;  f2 = 1336; return true;
        case '*': f1 = 941;  f2 = 1209; return true;
        case '#': f1 = 941;  f2 = 1477; return true;
        case 'A': f1 = 697;  f2 = 1633; return true;
        case 'B': f1 = 770;  f2 = 1633; return true;
        case 'C': f1 = 852;  f2 = 1633; return true;
        case 'D': f1 = 941;  f2 = 1633; return true;
        default:  return false;
    }
}

} // namespace dtmf
} // namespace bluebox