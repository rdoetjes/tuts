/*
 * parsing.cpp
 *
 * Implementation of sequence file parsing for bluebox-alsa.
 *
 * See parsing.h for API and documentation.
 */

#include "parsing.h"

#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <cctype>

namespace bluebox {
namespace parsing {

namespace {

inline std::string rstrip(std::string s) {
    while (!s.empty() && (s.back() == '\r' || s.back() == '\n')) s.pop_back();
    return s;
}

inline bool is_comment_or_empty(const std::string& s) {
    // skip leading whitespace to detect comment chars after optional indentation
    size_t i = 0;
    while (i < s.size() && std::isspace(static_cast<unsigned char>(s[i]))) ++i;
    if (i >= s.size()) return true;
    char c = s[i];
    return (c == '#' || c == ';');
}

} // anonymous


static void parse_wait_step(std::istringstream& iss, SequenceStep& step) {
    step.type = '~';
    std::string token;
    if (std::getline(iss, token, '\t')) {
        try {
            size_t pos = token.find_first_not_of(" \t\r\n");
            if (pos != std::string::npos) {
                step.duration_ms = std::stoi(token);
            } else {
                step.duration_ms = 0;
            }
        } catch (...) {
            step.duration_ms = 0;
        }
    } else {
        step.duration_ms = 0;
    }
}

static bool parse_tone_or_h_step(std::istringstream& iss, SequenceStep& step, char t, int lineno) {
    if (t != 'D' && t != 'C' && t != 'H') {
        std::cerr << "Line " << lineno << ": unsupported type '" << t << "'\n";
        return false;
    }
    step.type = t;
    std::string token;

    if (!std::getline(iss, token, '\t')) {
        step.tone.clear();
    } else {
        step.tone = token;
    }

    if (std::getline(iss, token, '\t')) {
        try {
            size_t first = token.find_first_not_of(" \t\r\n");
            size_t last  = token.find_last_not_of(" \t\r\n");
            if (first != std::string::npos && last != std::string::npos) {
                std::string sub = token.substr(first, last - first + 1);
                step.duration_ms = std::stoi(sub);
            } else {
                step.duration_ms = 0;
            }
        } catch (...) {
            step.duration_ms = 0;
        }
    } else {
        step.duration_ms = 0;
    }

    if (std::getline(iss, token, '\t')) {
        try {
            size_t first = token.find_first_not_of(" \t\r\n");
            size_t last  = token.find_last_not_of(" \t\r\n");
            if (first != std::string::npos && last != std::string::npos) {
                std::string sub = token.substr(first, last - first + 1);
                step.pause_ms = std::stoi(sub);
            } else {
                step.pause_ms = 0;
            }
        } catch (...) {
            step.pause_ms = 0;
        }
    } else {
        step.pause_ms = 0;
    }
    return true;
}

bool parse_sequence_stream(std::istream& in, std::vector<SequenceStep>& out_steps) {
    out_steps.clear();

    std::string line;
    int lineno = 0;

    while (std::getline(in, line)) {
        ++lineno;
        line = rstrip(line);

        if (line.empty() || is_comment_or_empty(line)) {
            continue;
        }

        std::istringstream iss(line);
        std::string token;

        SequenceStep step{};
        if (!std::getline(iss, token, '\t')) {
            continue;
        }

        if (token == "~") {
            parse_wait_step(iss, step);
            out_steps.push_back(std::move(step));
            continue;
        }

        if (token.length() != 1) {
            std::cerr << "Line " << lineno << ": invalid type token '" << token << "'\n";
            continue;
        }

        if (parse_tone_or_h_step(iss, step, token[0], lineno)) {
            out_steps.push_back(std::move(step));
        }
    }

    if (in.bad()) {
        std::cerr << "parse_sequence_stream: I/O error while reading stream\n";
        return false;
    }

    return true;
}

bool parse_sequence_file(const char* filename, std::vector<SequenceStep>& out_steps) {
    std::ifstream ifs(filename);
    if (!ifs.is_open()) {
        std::cerr << "Cannot open file: " << filename << "\n";
        return false;
    }

    bool ok = parse_sequence_stream(ifs, out_steps);
    if (!ok) {
        std::cerr << "Error parsing sequence file: " << filename << "\n";
        return false;
    }

    return true;
}

} // namespace parsing
} // namespace bluebox
