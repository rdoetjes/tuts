const rl = @import("raylib");
const gs = @import("game_state.zig");
const gl = @import("game_logic.zig");
const config = @import("config.zig");
const std = @import("std");

pub fn handleInput(state: *gs.GameState) void {

     if (rl.isKeyDown(rl.KeyboardKey.key_up)) {
          gl.player_up(state);
     }

     if (rl.isKeyDown(rl.KeyboardKey.key_down)) {
          gl.player_down(state);
     }

     if (rl.isKeyDown(rl.KeyboardKey.key_right)) {
          gl.player_right(state);
     }

     if (rl.isKeyDown(rl.KeyboardKey.key_left)) {
          gl.player_left(state);
     }

     if (rl.isKeyDown(rl.KeyboardKey.key_space)) {
          if (state.screen == .playing) {
               gl.player_fire(state);
          }
     }

     if (rl.isKeyReleased(rl.KeyboardKey.key_space)) {

          if (state.screen == .splash) {
               state.reset();
               state.screen = .playing;
               return;
          }

          if (state.screen == .gameover) {
               state.screen = .splash;
               return;
          }


          gl.player_fire_release(state);
     }

     if (rl.isKeyReleased(rl.KeyboardKey.key_down)) {
          gl.player_down_release(state);
     }


     if (rl.isKeyReleased(rl.KeyboardKey.key_up)) {
          gl.player_up_release(state);
     }
}
