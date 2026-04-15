#include <iostream>
#include <string>
#include <vector>
#include <cstdlib>
#include <unistd.h> // for close()

#include "parsing.h"
#include "playing.h"
#include "serial.h"

int main(int argc, char* argv[]) {
    if (argc < 2 || argc > 3) {
        std::cerr << "Usage: " << argv[0] << " sequence.txt [serial_device]\n\n"
                  << "  sequence.txt     Path to tab-delimited sequence file\n"
                  << "  serial_device    Optional serial device (default: " << bluebox::serial::DEFAULT_DEVICE << ")\n";
        return 1;
    }

    const char* sequence_file = argv[1];
    const char* serial_device = (argc >= 3) ? argv[2] : bluebox::serial::DEFAULT_DEVICE;

    std::vector<bluebox::parsing::SequenceStep> sequence;
    if (!bluebox::parsing::parse_sequence_file(sequence_file, sequence)) {
        std::cerr << "Failed to parse sequence file: " << sequence_file << "\n";
        return 1;
    }

    std::cout << "Loaded " << sequence.size() << " steps from '" << sequence_file << "'.\n";

    // Setup ALSA
    snd_pcm_t* handle = nullptr;
    if (!bluebox::playing::setup_alsa(handle)) {
        std::cerr << "ALSA setup failed; cannot play sequence.\n";
        return 1;
    }

    // Try to open serial device. If it fails we continue but H commands will be skipped.
    int serial_fd = bluebox::serial::open_serial(serial_device);
    if (serial_fd < 0) {
        std::cerr << "Warning: serial device '" << serial_device << "' not available; 'H' commands will be skipped.\n";
    } else {
        std::cout << "Serial device '" << serial_device << "' opened.\n";
    }

    // Play the sequence (this will block until finished)
    bluebox::playing::play_sequence(handle, sequence, serial_fd);

    // Cleanup
    if (serial_fd >= 0) {
        ::close(serial_fd);
    }

    if (handle) {
        snd_pcm_drain(handle);
        snd_pcm_close(handle);
    }

    std::cout << "Sequence complete.\n";
    return 0;
}