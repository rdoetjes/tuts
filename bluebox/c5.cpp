// c5.cpp
//
// Implementation of CCITT / C5 tone lookup helpers for bluebox-alsa.
//
// Maps C5-style codes (single digit or named codes like "KP1", "KP2", "ST",
// "SEIZE", "PROCEED", etc.) to frequency pairs (Hz).
//
// See c5.h for the public API.

#include "c5.h"

#include <algorithm>
#include <cctype>
#include <string>

namespace bluebox {
namespace c5 {

// Helper: uppercase a string in-place
static std::string to_upper(std::string s) {
    std::transform(s.begin(), s.end(), s.begin(),
                   [](unsigned char c){ return static_cast<char>(std::toupper(c)); });
    return s;
}

bool get_c5_freqs(const std::string& code_in, int& f1, int& f2) {
    if (code_in.empty()) return false;

    std::string code = to_upper(code_in);

    // If the code is a single digit, handle the digit -> (f1,f2) mapping.
    if (code.size() == 1 && std::isdigit(static_cast<unsigned char>(code[0]))) {
        int d = code[0] - '0';
        switch (d) {
            case 1: f1 =  700; f2 =  900; return true;
            case 2: f1 =  700; f2 = 1100; return true;
            case 3: f1 =  900; f2 = 1100; return true;
            case 4: f1 =  700; f2 = 1300; return true;
            case 5: f1 =  900; f2 = 1300; return true;
            case 6: f1 = 1100; f2 = 1300; return true;
            case 7: f1 =  700; f2 = 1500; return true;
            case 8: f1 =  900; f2 = 1500; return true;
            case 9: f1 = 1100; f2 = 1500; return true;
            case 0: f1 = 1300; f2 = 1500; return true;
            default: return false;
        }
    }

    // Named codes (case-insensitive)
    if (code == "KP1")                             { f1 = 1100; f2 = 1700; return true; }
    if (code == "KP2")                             { f1 = 1300; f2 = 1700; return true; }
    if (code == "ST")                              { f1 = 1500; f2 = 1700; return true; }
    if (code == "CODE11")                          { f1 =  700; f2 = 1700; return true; }
    if (code == "CODE12")                          { f1 =  900; f2 = 1700; return true; }
    if (code == "SEIZE")                           { f1 = 2400; f2 =    0; return true; }
    if (code == "PROCEED" || code == "PK")         { f1 = 2600; f2 =    0; return true; }
    if (code == "ANSWER")                          { f1 = 2400; f2 =    0; return true; }
    if (code == "BUSY" || code == "BUSYFLASH")     { f1 = 2600; f2 =    0; return true; }
    if (code == "CLEARBACK")                       { f1 = 2600; f2 =    0; return true; }
    if (code == "CLEARFWD" || code == "RELGUARD")  { f1 = 2400; f2 = 2600; return true; }

    // Unknown code
    return false;
}

bool is_valid_c5(const std::string& code) {
    int f1 = 0, f2 = 0;
    return get_c5_freqs(code, f1, f2);
}

} // namespace c5
} // namespace bluebox