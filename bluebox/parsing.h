#ifndef BLUEBOX_ALSA_PARSING_H
#define BLUEBOX_ALSA_PARSING_H

/*
 * parsing.h
 *
 * Sequence file parsing for bluebox-alsa.
 *
 * Sequence files are tab-delimited text files. Each non-comment line describes one step:
 *
 *   Type<TAB>Tone<TAB>Duration_ms<TAB>Pause_ms
 *
 * - Fields are separated by a single tab character ('\t').
 * - Lines beginning with '#' or ';' are treated as comments and ignored.
 * - Supported Type values:
 *     'D' - DTMF single-character tone (Tone must be a single character, e.g. '5' or '#')
 *     'C' - CCITT/C5 code (Tone may be numeric or a code string like "KP1", "SEIZE", etc.)
 *     'H' - Serial Host command (no tone output; Tone field is ignored)
 *     '~' - Wait/interactive. If Duration_ms > 0 the player sleeps that many ms; if 0 or missing, it waits for an Enter keypress.
 *
 * Example lines:
 *   D<TAB>5<TAB>120<TAB>80
 *   C<TAB>KP1<TAB>100<TAB>50
 *   ~<TAB><TAB><TAB>   (interactive wait for Enter)
 *
 * The parser is permissive: if a line cannot be parsed it is skipped and parsing continues.
 *
 * The original single-file program used an internal `SequenceStep` type; this header
 * exposes that type and parsing entry points so other translation units may use them.
 */

#include <string>
#include <vector>
#include <istream>
#include <cstddef>
#include <cstdint>

namespace bluebox {
namespace parsing {

/*
 * SequenceStep
 *
 * Represents a single step from a sequence file.
 *
 * Fields:
 *  - type: 'D', 'C', 'H', or '~'
 *  - duration_ms: how long to play/send the tone/command (milliseconds). May be 0.
 *  - pause_ms: how long to silence/wait after the step (milliseconds). May be 0.
 *  - tone: textual tone description (single character for D, code for C). May be empty.
 *
 * alignas(16) is preserved to mirror the original layout choices for alignment.
 */
struct alignas(16) SequenceStep {
    char        type = 0;
    int         duration_ms = 0;
    int         pause_ms = 0;
    std::string tone;
};

/*
 * parse_sequence_file
 *
 * Parse a sequence file from disk.
 *
 * Parameters:
 *  - filename: path to the sequence file.
 *  - out_steps: vector that will be appended with parsed SequenceStep entries.
 *
 * Returns:
 *  - true on success (file opened and parsed, even if some lines were skipped).
 *  - false if the file could not be opened or a fatal IO error occurred.
 */
bool parse_sequence_file(const char* filename, std::vector<SequenceStep>& out_steps);

/*
 * Convenience overload using std::string filename.
 */
inline bool parse_sequence_file(const std::string& filename, std::vector<SequenceStep>& out_steps) {
    return parse_sequence_file(filename.c_str(), out_steps);
}

/*
 * parse_sequence_stream
 *
 * Parse from an already-open input stream. This is useful for unit tests or
 * piping data from other sources. The function will read until EOF or error.
 *
 * Returns:
 *  - true if parsing completed without a fatal IO error (individual bad lines
 *    are ignored and parsing continues).
 *  - false on stream errors prior to EOF.
 */
bool parse_sequence_stream(std::istream& in, std::vector<SequenceStep>& out_steps);

} // namespace parsing
} // namespace bluebox

#endif // BLUEBOX_ALSA_PARSING_H