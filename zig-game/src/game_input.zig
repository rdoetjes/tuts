const rl = @import("raylib");
const gs = @import("game_state.zig");
const gl = @import("game_logic.zig");
const config = @import("config.zig");
const std = @import("std");

// handle all input and map accordingly to the right method in game_logic.zig
pub fn handle_input(state: *gs.GameState) !void {
    if (rl.isKeyDown(rl.KeyboardKey.key_space)) {
        // avoids the game to fire when in splash screen or gameover screen
        if (state.screen == .playing) {
            try gl.player_fire(state);
        }
    }

    if (rl.isKeyReleased(rl.KeyboardKey.key_space)) {
        gl.player_fire_release(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_up)) {
        gl.player_up(state);
        return;
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_down)) {
        gl.player_down(state);
        return;
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_right)) {
        gl.player_right(state);
        return;
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_left)) {
        gl.player_left(state);
        return;
    }

    if (rl.isKeyReleased(rl.KeyboardKey.key_s)) {
        if (state.screen == .splash) {
            state.reset();
            state.screen = .playing;
            return;
        }

        if (state.screen == .gameover) {
            state.screen = .splash;
            return;
        }
    }

    if (rl.isKeyReleased(rl.KeyboardKey.key_down)) {
        gl.player_down_release(state);
    }

    if (rl.isKeyReleased(rl.KeyboardKey.key_up)) {
        gl.player_up_release(state);
    }
}
