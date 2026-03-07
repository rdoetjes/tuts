#include <iostream>
#include <fstream>
#include <vector>
#include <string>
#include <sstream>
#include <algorithm>
#include <cctype>
#include <cmath>
#include <chrono>
#include <thread>
#include <climits>
#include <alsa/asoundlib.h>

// ------------------------------------------------------------------------------------------------
// FEATURE IMPLEMENTATION INJECTED AT TOP
// ------------------------------------------------------------------------------------------------

// Constants needed for the new implementation
constexpr unsigned SAMPLE_RATE         = 8000;
constexpr unsigned CHANNELS            = 1;
constexpr snd_pcm_format_t FORMAT      = SND_PCM_FORMAT_S16_LE;
constexpr double   AMPLITUDE           = INT16_MAX * 0.4;
constexpr unsigned PERIOD_SIZE         = 128;               // ~16 ms @ 8 kHz
constexpr unsigned BUFFER_MULTIPLIER   = 4;
constexpr int      MIN_SILENCE_MS      = 5;                 // force at least this much zero after each tone

// Redefined struct to ensure we can use it in our new parser
struct alignas(16) SequenceStep {
    char        type;          // 'D', 'C', or '~'
    int         duration_ms;
    int         pause_ms;
    std::string tone;
};

// Forward declarations of helper functions from the preserved block that we will reuse
bool get_dtmf_freqs(char digit, int& f1, int& f2);
bool get_c5_freqs(const std::string& code, int& f1, int& f2);
void generate_tone_buffer(int f1, int f2, int duration_ms, std::vector<int16_t>& buffer);
bool setup_alsa(snd_pcm_t*& handle);
bool write_frames(snd_pcm_t* handle, const int16_t* data, size_t frame_count);
void play_sequence(snd_pcm_t* handle, const std::vector<SequenceStep>& sequence);

// New parsing logic supporting '~'
bool parse_sequence_file(const char* filename, std::vector<SequenceStep>& sequence) {
    std::ifstream file(filename);
    if (!file) {
        std::cerr << "Cannot open file: " << filename << "\n";
        return false;
    }

    std::string line;
    int lineno = 0;

    while (std::getline(file, line)) {
        ++lineno;
        if (line.empty() || line[0] == '#' || line[0] == ';') continue;

        std::istringstream iss(line);
        std::string token;
        SequenceStep step{};

        // Check first token
        if (!std::getline(iss, token, '\t')) continue;

        // Handle wait command
        if (token == "~") {
            step.type = '~';
            sequence.push_back(std::move(step));
            continue;
        }

        // Type
        if (token.length() != 1 || (token[0] != 'D' && token[0] != 'C')) {
            std::cerr << "Line " << lineno << ": Invalid type >" << token[0] << "< \n";
            continue;
        }
        step.type = token[0];

        // Tone
        if (!std::getline(iss, token, '\t')) continue;
        step.tone = std::move(token);

        // Duration
        if (!std::getline(iss, token, '\t')) continue;
        try { step.duration_ms = std::stoi(token); } catch (...) { continue; }

        // Pause
        if (!std::getline(iss, token, '\t')) continue;
        try { step.pause_ms = std::stoi(token); } catch (...) { continue; }

        sequence.push_back(std::move(step));
    }
    return true;
}

// New playback loop supporting '~'
void play_sequence(snd_pcm_t* handle, const std::vector<SequenceStep>& sequence) {
    std::vector<int16_t> silence_chunk(PERIOD_SIZE * 2, 0);

    for (size_t i = 0; i < sequence.size(); ++i) {
        const auto& step = sequence[i];

        if (step.type == '~') {
            std::cout << "Waiting for Enter.. (Press Enter)\n";
            std::cin.get();
            continue;
        }

        int f1 = 0, f2 = 0;
        bool valid = false;

        if (step.type == 'D') {
            if (step.tone.length() == 1) {
                valid = get_dtmf_freqs(step.tone[0], f1, f2);
            }
        } else if (step.type == 'C') {
            valid = get_c5_freqs(step.tone, f1, f2);
        }

        if (!valid) {
            std::cerr << "Step " << (i+1) << ": Invalid tone '" << step.tone
                      << "' for type " << step.type << " → skipping tone\n";
        } else if (step.duration_ms > 0) {
            std::vector<int16_t> tone_buffer;
            generate_tone_buffer(f1, f2, step.duration_ms, tone_buffer);
            if (!tone_buffer.empty()) {
                write_frames(handle, tone_buffer.data(), tone_buffer.size());
            }
        }

        int silence_ms = std::max(step.pause_ms, MIN_SILENCE_MS);
        size_t silence_samples = static_cast<size_t>(silence_ms * SAMPLE_RATE / 1000.0 + 0.5);

        size_t remaining = silence_samples;
        while (remaining > 0) {
            size_t chunk_frames = std::min(remaining, silence_chunk.size());
            write_frames(handle, silence_chunk.data(), chunk_frames);
            remaining -= chunk_frames;
        }
    }
}

// New main
int main(int argc, char* argv[]) {
    if (argc != 2) {
        std::cerr << "Usage: " << argv[0] << " sequence.txt\n\n"
                  << "Tab-delimited format (one line per step):\n"
                  << "Type\tTone\tDuration_ms\tPause_ms\n"
                  << "D\t5\t120\t80\n"
                  << "~\t\t\t\t(Wait for keypress)\n";
        return 1;
    }

    std::vector<SequenceStep> sequence;
    if (!parse_sequence_file(argv[1], sequence)) {
        return 1;
    }

    std::cout << "Loaded " << sequence.size() << " steps.\n";

    snd_pcm_t* handle = nullptr;
    if (!setup_alsa(handle)) {
        return 1;
    }

    play_sequence(handle, sequence);

    snd_pcm_drain(handle);
    snd_pcm_close(handle);

    return 0;
}


// ────────────────────────────────────────────────────────────────────────────────
// Tone frequency lookup
// ────────────────────────────────────────────────────────────────────────────────

bool get_dtmf_freqs(char digit, int& f1, int& f2) {
    digit = std::toupper(digit);
    switch (digit) {
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

bool get_c5_freqs(const std::string& code, int& f1, int& f2) {
    if (code.length() == 1 && std::isdigit(code[0])) {
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
        }
        return false;
    }

    std::string upper = code;
    std::transform(upper.begin(), upper.end(), upper.begin(), ::toupper);

    if (upper == "KP1" || upper == "PK1")   { f1 = 1100; f2 = 1700; return true; }
    if (upper == "KP2" || upper == "PK2")   { f1 = 1300; f2 = 1700; return true; }
    if (upper == "ST")                      { f1 = 1500; f2 = 1700; return true; }
    if (upper == "CODE11")                  { f1 =  700; f2 = 1700; return true; }
    if (upper == "CODE12")                  { f1 =  900; f2 = 1700; return true; }
    if (upper == "SEIZE")                   { f1 = 2400; f2 =    0; return true; }
    if (upper == "PROCEED" || upper == "PK") { f1 = 2600; f2 =    0; return true; }
    if (upper == "ANSWER")                  { f1 = 2400; f2 =    0; return true; }
    if (upper == "BUSY" || upper == "BUSYFLASH") { f1 = 2600; f2 = 0; return true; }
    if (upper == "CLEARBACK")               { f1 = 2600; f2 =    0; return true; }
    if (upper == "CLEARFWD" || upper == "RELGUARD") { f1 = 2400; f2 = 2600; return true; }

    return false;
}

// ────────────────────────────────────────────────────────────────────────────────
// Audio generation – plain sinusoids, no envelope
// ────────────────────────────────────────────────────────────────────────────────

void generate_tone_buffer(int f1, int f2, int duration_ms, std::vector<int16_t>& buffer) {
    if (duration_ms <= 0) {
        buffer.clear();
        return;
    }

    size_t num_samples = static_cast<size_t>(duration_ms * SAMPLE_RATE / 1000.0 + 0.5);
    buffer.resize(num_samples);

    const double tau = 2.0 * M_PI;
    double phase1 = 0.0, phase2 = 0.0;
    double incr1 = (f1 > 0) ? tau * f1 / SAMPLE_RATE : 0.0;
    double incr2 = (f2 > 0) ? tau * f2 / SAMPLE_RATE : 0.0;

    for (size_t i = 0; i < num_samples; ++i) {
        double sample = 0.0;
        if (f1 > 0) sample += AMPLITUDE * std::sin(phase1);
        if (f2 > 0) sample += AMPLITUDE * std::sin(phase2);
        buffer[i] = static_cast<int16_t>(std::clamp(sample, -32768.0, 32767.0));

        phase1 += incr1;
        phase2 += incr2;
    }
}

// ────────────────────────────────────────────────────────────────────────────────
// ALSA setup
// ────────────────────────────────────────────────────────────────────────────────

bool setup_alsa(snd_pcm_t*& handle) {
    int err = snd_pcm_open(&handle, "default", SND_PCM_STREAM_PLAYBACK, 0);
    if (err < 0) {
        std::cerr << "Cannot open PCM device: " << snd_strerror(err) << "\n";
        return false;
    }

    snd_pcm_hw_params_t* hw_params = nullptr;
    snd_pcm_hw_params_alloca(&hw_params);

    unsigned int rate = SAMPLE_RATE;
    snd_pcm_uframes_t period = PERIOD_SIZE;
    snd_pcm_uframes_t bufsz  = PERIOD_SIZE * BUFFER_MULTIPLIER;

    if ((err = snd_pcm_hw_params_any(handle, hw_params)) < 0 ||
        (err = snd_pcm_hw_params_set_access(handle, hw_params, SND_PCM_ACCESS_RW_INTERLEAVED)) < 0 ||
        (err = snd_pcm_hw_params_set_format(handle, hw_params, FORMAT)) < 0 ||
        (err = snd_pcm_hw_params_set_rate_near(handle, hw_params, &rate, 0)) < 0 ||
        (err = snd_pcm_hw_params_set_channels(handle, hw_params, CHANNELS)) < 0 ||
        (err = snd_pcm_hw_params_set_period_size_near(handle, hw_params, &period, 0)) < 0 ||
        (err = snd_pcm_hw_params_set_buffer_size_near(handle, hw_params, &bufsz)) < 0 ||
        (err = snd_pcm_hw_params(handle, hw_params)) < 0 ||
        (err = snd_pcm_prepare(handle)) < 0) {
        std::cerr << "ALSA setup failed: " << snd_strerror(err) << "\n";
        snd_pcm_close(handle);
        return false;
    }

    //have also settle, to prevent partial begin note
    std::this_thread::sleep_for(std::chrono::milliseconds(600));
    return true;
}

// ────────────────────────────────────────────────────────────────────────────────
// Write audio frames (tone or silence) with underrun recovery
// ────────────────────────────────────────────────────────────────────────────────

bool write_frames(snd_pcm_t* handle, const int16_t* data, size_t frame_count) {
    size_t remaining = frame_count;
    const int16_t* ptr = data;

    while (remaining > 0) {
        snd_pcm_sframes_t written = snd_pcm_writei(handle, ptr, remaining);

        if (written == -EPIPE) {
            snd_pcm_prepare(handle);
            continue;
        }
        if (written < 0) {
            std::cerr << "writei failed: " << snd_strerror(static_cast<int>(written)) << "\n";
            return false;
        }

        ptr += written;
        remaining -= static_cast<size_t>(written);
    }

    return true;
}
