#include "playing.h"
#include "dtmf.h"
#include "c5.h"
#include "serial.h"

#include <alsa/asoundlib.h>

#include <cmath>
#include <algorithm>
#include <iostream>
#include <thread>
#include <chrono>
#include <vector>
#include <cstring>

namespace bluebox {
namespace playing {

// Setup ALSA for playback using the parameters defined in the header.
bool setup_alsa(snd_pcm_t*& handle) {
    int err = snd_pcm_open(&handle, "default", SND_PCM_STREAM_PLAYBACK, 0);
    if (err < 0) {
        std::cerr << "setup_alsa: cannot open PCM device: " << snd_strerror(err) << "\n";
        return false;
    }

    snd_pcm_hw_params_t* hw_params = nullptr;
    snd_pcm_hw_params_alloca(&hw_params);

    unsigned int rate = SAMPLE_RATE;
    snd_pcm_uframes_t period = PERIOD_SIZE;
    snd_pcm_uframes_t bufsz  = static_cast<snd_pcm_uframes_t>(PERIOD_SIZE * BUFFER_MULTIPLIER);

    if ((err = snd_pcm_hw_params_any(handle, hw_params)) < 0 ||
        (err = snd_pcm_hw_params_set_access(handle, hw_params, SND_PCM_ACCESS_RW_INTERLEAVED)) < 0 ||
        (err = snd_pcm_hw_params_set_format(handle, hw_params, FORMAT)) < 0 ||
        (err = snd_pcm_hw_params_set_rate_near(handle, hw_params, &rate, 0)) < 0 ||
        (err = snd_pcm_hw_params_set_channels(handle, hw_params, CHANNELS)) < 0 ||
        (err = snd_pcm_hw_params_set_period_size_near(handle, hw_params, &period, 0)) < 0 ||
        (err = snd_pcm_hw_params_set_buffer_size_near(handle, hw_params, &bufsz)) < 0 ||
        (err = snd_pcm_hw_params(handle, hw_params)) < 0 ||
        (err = snd_pcm_prepare(handle)) < 0) {
        std::cerr << "setup_alsa: ALSA setup failed: " << snd_strerror(err) << "\n";
        snd_pcm_close(handle);
        handle = nullptr;
        return false;
    }

    // Small settle to avoid partial beginning artefacts when device is opened
    std::this_thread::sleep_for(std::chrono::milliseconds(300));
    return true;
}

// Write frames robustly with underrun recovery. Returns true on success.
bool write_frames(snd_pcm_t* handle, const int16_t* data, size_t frame_count) {
    if (handle == nullptr) return false;

    size_t remaining = frame_count;
    const int16_t* ptr = data;

    while (remaining > 0) {
        snd_pcm_sframes_t written = snd_pcm_writei(handle, ptr, remaining);

        if (written == -EPIPE) {
            // underrun - try to recover
            int err = snd_pcm_prepare(handle);
            if (err < 0) {
                std::cerr << "write_frames: snd_pcm_prepare after underrun failed: " << snd_strerror(err) << "\n";
                return false;
            }
            continue;
        }
        if (written < 0) {
            std::cerr << "write_frames: writei failed: " << snd_strerror(static_cast<int>(written)) << "\n";
            return false;
        }

        ptr += written;
        remaining -= static_cast<size_t>(written);
    }

    return true;
}

// Generate tone buffer summing up to two sinusoids, clamped to int16 range.
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

    const size_t ramp_samples = std::min((size_t)40, num_samples / 4);  // ~5 ms attack/decay at 8000 Hz

    for (size_t i = 0; i < num_samples; ++i) {
        double env = 1.0;

        // Attack ramp (fade-in at start)
        if (i < ramp_samples) {
            env = (double)i / ramp_samples;
        }
        // Decay ramp (fade-out at end)
        else if (i > num_samples - ramp_samples) {
            env = (double)(num_samples - i) / ramp_samples;
        }
        // Middle part stays at env = 1.0 (full volume)
        double sample = 0.0;
        if (f1 > 0) sample += AMPLITUDE * std::sin(phase1) * env;
        if (f2 > 0) sample += AMPLITUDE * std::sin(phase2) * env;
        // clamp to int16_t range
        double clamped = std::max(-32768.0, std::min(32767.0, sample));
        buffer[i] = static_cast<int16_t>(clamped);

        phase1 += incr1;
        phase2 += incr2;
    }
}

static void play_wait_step(const parsing::SequenceStep& step, size_t index) {
    if (step.duration_ms == 0) {
        std::cout << "Waiting for Enter.. (Press Enter)\n";
        std::string tmp;
        std::cout << "step: " << (index+1) << " " << step.type << std::endl;
        std::getline(std::cin, tmp);
    } else {
        std::cout << "step: " << (index+1) << " " << step.type << " " << step.duration_ms << " " << std::endl;
        std::this_thread::sleep_for(std::chrono::milliseconds(step.duration_ms));
    }
}

static void play_h_step(snd_pcm_t* handle, const parsing::SequenceStep& step, size_t index, int serial_fd, const std::vector<int16_t>& silence_chunk) {
    if (serial_fd < 0) {
        std::cerr << "play_sequence: H command requested but no serial device open; skipping\n";
    } else {
        int resp = bluebox::serial::send_and_wait_serial(serial_fd, 'H', 1000);
        std::cout << "step: " << (index+1) << " " << step.type << " " << step.duration_ms << " " << step.pause_ms << std::endl;
        if (resp == 0) {
            std::cout << "Step " << (index+1) << ": Serial response: FAILED\n";
        } else if (resp == -1) {
            std::cout << "Step " << (index+1) << ": Serial response: TIMEOUT/ERROR\n";
        }
    }

    int silence_ms = std::max(step.pause_ms, MIN_SILENCE_MS);
    size_t silence_samples = static_cast<size_t>(silence_ms * SAMPLE_RATE / 1000.0 + 0.5);
    size_t remaining = silence_samples;
    while (remaining > 0) {
        size_t chunk = std::min(remaining, silence_chunk.size());
        write_frames(handle, silence_chunk.data(), chunk);
        remaining -= chunk;
    }
}

static void play_tone_step(snd_pcm_t* handle, const parsing::SequenceStep& step, size_t index, const std::vector<int16_t>& silence_chunk) {
    int f1 = 0, f2 = 0;
    bool valid = false;

    if (step.type == 'D' && !step.tone.empty()) {
        char d = step.tone[0];
        valid = bluebox::dtmf::get_dtmf_freqs(d, f1, f2);
    } else if (step.type == 'C' && !step.tone.empty()) {
        valid = bluebox::c5::get_c5_freqs(step.tone, f1, f2);
    } else {
        valid = false;
    }

    if (!valid) {
        std::cerr << "play_sequence: Step " << (index+1) << ": Invalid or unsupported tone for type "
                  << step.type << " (\"" << step.tone << "\"); skipping tone\n";
    } else if (step.duration_ms > 0) {
        std::vector<int16_t> tone_buffer;
        generate_tone_buffer(f1, f2, step.duration_ms, tone_buffer);
        if (!tone_buffer.empty()) {
            if (!write_frames(handle, tone_buffer.data(), tone_buffer.size())) {
                std::cerr << "play_sequence: failed to write tone frames for step " << (index+1) << "\n";
            }
        }
    }

    int silence_ms = std::max(step.pause_ms, MIN_SILENCE_MS);
    size_t silence_samples = static_cast<size_t>(silence_ms * SAMPLE_RATE / 1000.0 + 0.5);
    size_t remaining = silence_samples;
    while (remaining > 0) {
        size_t chunk = std::min(remaining, silence_chunk.size());
        write_frames(handle, silence_chunk.data(), chunk);
        remaining -= chunk;
    }
    std::cout << "step: " << (index+1) << " " << step.type << " " << " " << step.tone << " " << step.duration_ms << " " << step.pause_ms << std::endl;
}

// Execute a parsed sequence. Uses dtmf, c5 and serial helpers as needed.
//
// serial_fd may be -1 to indicate no serial device is available; H steps will
// be skipped in that case (with a warning).
void play_sequence(snd_pcm_t* handle, const std::vector<parsing::SequenceStep>& sequence, int serial_fd) {
    // Pre-create a silence chunk (frames) to write in chunks.
    const size_t silence_chunk_frames = static_cast<size_t>(PERIOD_SIZE * 2);
    std::vector<int16_t> silence_chunk(silence_chunk_frames, 0);

    for (size_t i = 0; i < sequence.size(); ++i) {
        const auto& step = sequence[i];

        if (step.type == '~') {
            play_wait_step(step, i);
        } else if (step.type == 'H') {
            play_h_step(handle, step, i, serial_fd, silence_chunk);
        } else {
            play_tone_step(handle, step, i, silence_chunk);
        }
    }
}

} // namespace playing
} // namespace bluebox
