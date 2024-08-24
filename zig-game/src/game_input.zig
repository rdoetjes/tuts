const rl = @import("raylib");
const gs = @import("game_state.zig");
const gl = @import("game_logic.zig");
const config = @import("config.zig");

pub fn handleInput(state: *gs.GameState) void {
    if (rl.isKeyDown(rl.KeyboardKey.key_right)) {
        gl.player_right(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_left)) {
       gl.player_left(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_up)) {
       gl.player_up(state);
    }

    if (rl.isKeyDown(rl.KeyboardKey.key_down)) {
       gl.player_down(state);
    }
}
