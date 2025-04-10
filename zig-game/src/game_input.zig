const rl = @import("raylib");
const gs = @import("game_state.zig");
const gl = @import("game_logic.zig");
const config = @import("config.zig");
const std = @import("std");

// handle all input and map accordingly to the right method in game_logic.zig
pub fn handle_input(state: *gs.GameState) !void {
    if (rl.isKeyDown(rl.KeyboardKey.up)) {
        gl.player_up(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.down)) {
        gl.player_down(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.right)) {
        gl.player_right(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.left)) {
        gl.player_left(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.space)) {
        // avoids the game to fire when in splash screen or gameover screen
        if (state.screen == .playing) {
            try gl.player_fire(state);
        }
    }

    if (rl.isKeyReleased(rl.KeyboardKey.space)) {
        gl.player_fire_release(state);
    }

    if (rl.isKeyReleased(rl.KeyboardKey.s)) {
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

    if (rl.isKeyReleased(rl.KeyboardKey.down)) {
        gl.player_down_release(state);
    }

    if (rl.isKeyReleased(rl.KeyboardKey.up)) {
        gl.player_up_release(state);
    }
}
