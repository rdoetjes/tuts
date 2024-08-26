const rl = @import("raylib");
const gs = @import("game_state.zig");
const gl = @import("game_logic.zig");
const config = @import("config.zig");
const std = @import("std");

pub fn handleInput(state: *gs.GameState) void {
     if (state.screen == .splash) {
          if (rl.isKeyPressed(rl.KeyboardKey.key_space)) {
               state.reset();
          }
          return;
     }

     if (state.screen == .gameover) {
          if (rl.isKeyPressed(rl.KeyboardKey.key_space)){
               state.screen = .splash;
          }
          return;
     }

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
          gl.player_fire(state);
     }

     if (rl.isKeyReleased(rl.KeyboardKey.key_space)) {
          gl.player_fire_release(state);
     }

     if (rl.isKeyReleased(rl.KeyboardKey.key_down)) {
          gl.player_down_release(state);
     }

     if (rl.isKeyReleased(rl.KeyboardKey.key_up)) {
          gl.player_up_release(state);
     }
     }
