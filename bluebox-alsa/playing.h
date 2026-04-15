#ifndef BLUEBOX_ALSA_PLAYING_H
#define BLUEBOX_ALSA_PLAYING_H

/*
 * playing.h
 *
 * Playback and sequence execution API for bluebox-alsa.
 *
 * This header declares the functions and constants used to:
 *  - configure and open ALSA for playback
 *  - generate tone audio buffers (two-frequency tones)
 *  - write frames reliably to ALSA (with underrun recovery)
 *  - execute a parsed sequence of steps, including DTMF, C5 and serial 'H' commands
 *
 * The library is intentionally minimal: implementations should live in
 * playing.cpp and link against other translation units (dtmf, c5, parsing, serial).
 *
 * Notes:
 *  - SequenceStep is defined in parsing.h and used here by reference.
 *  - Serial-related commands are performed by the serial module; play sequence
 *    accepts a serial file descriptor (or -1) to allow optional serial usage.
 */

#include <cstdint>
#include <cstddef>
#include <vector>
#include <string>

#include "parsing.h"   // for parsing::SequenceStep
#include <alsa/asoundlib.h> // for snd_pcm_t, snd_pcm_format_t

namespace bluebox {
namespace playing {

/* Audio configuration constants (compatible with existing implementation) */
inline constexpr unsigned SAMPLE_RATE       = 8000;
inline constexpr unsigned CHANNELS          = 1;
inline constexpr snd_pcm_format_t FORMAT    = SND_PCM_FORMAT_S16_LE;
inline constexpr double AMPLITUDE           = static_cast<double>(INT16_MAX) * 0.4;
inline constexpr unsigned PERIOD_SIZE       = 128;              // approx ~16 ms @ 8 kHz
inline constexpr unsigned BUFFER_MULTIPLIER = 4;
inline constexpr int MIN_SILENCE_MS         = 5;                // minimum enforced silence after a tone

/*
 * setup_alsa
 *
 * Open and configure the default ALSA PCM device for playback using the
 * parameters defined above. On success, `handle` will contain an opened and
 * prepared snd_pcm_t* (caller is responsible for calling snd_pcm_close).
 *
 * Returns true on success, false on failure (caller may log errors from ALSA).
 */
bool setup_alsa(snd_pcm_t*& handle);

/*
 * write_frames
 *
 * Write `frame_count` frames (int16_t samples for mono) to the ALSA device.
 * This function should recover from underruns and attempt to write all frames
 * before returning. Returns true on success, false on unrecoverable error.
 *
 * Parameters:
 *  - handle: opened snd_pcm_t* (playback)
 *  - data: pointer to interleaved frames (int16_t)
 *  - frame_count: number of frames to write
 */
bool write_frames(snd_pcm_t* handle, const int16_t* data, size_t frame_count);

/*
 * generate_tone_buffer
 *
 * Generate a mono int16_t buffer containing the sum of up to two pure
 * sinusoids at frequencies f1 and f2 (Hz) for `duration_ms` milliseconds.
 * If a frequency value is <= 0 it will be ignored (allowing single-frequency
 * tones).
 *
 * The output vector `buffer` will be resized to contain `num_samples` samples.
 * The function clamps to the int16 range and uses the AMPLITUDE constant for scaling.
 *
 * Parameters:
 *  - f1, f2: frequencies in Hz (0 to ignore second/first)
 *  - duration_ms: duration in milliseconds
 *  - buffer: output vector of int16_t samples (mono)
 */
void generate_tone_buffer(int f1, int f2, int duration_ms, std::vector<int16_t>& buffer);

/*
 * play_sequence
 *
 * Execute the parsed sequence steps. Tone-producing steps ('D' and 'C') will
 * be converted into audio using the DTMF / C5 lookup helpers and played via
 * the provided ALSA handle. Wait steps ('~') will sleep or prompt for Enter as
 * specified by the SequenceStep. Serial 'H' steps will use `serial_fd` to send
 * a single-byte 'H' command and wait for the response; if `serial_fd < 0`,
 * H steps will be skipped (with a warning).
 *
 * Parameters:
 *  - handle: opened snd_pcm_t* returned from setup_alsa
 *  - sequence: vector of parsing::SequenceStep (parsed from file)
 *  - serial_fd: file descriptor for serial device, or -1 when not available
 *
 * This function will block until the sequence finishes. It should perform
 * conservative silence padding between steps (at least MIN_SILENCE_MS).
 */
void play_sequence(snd_pcm_t* handle, const std::vector<parsing::SequenceStep>& sequence, int serial_fd);

} // namespace playing
} // namespace bluebox

#endif // BLUEBOX_ALSA_PLAYING_H