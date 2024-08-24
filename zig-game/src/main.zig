const config = @import("config.zig");
const rl = @import("raylib");
const std = @import("std");

const gs = @import("game_state.zig");
const gl = @import("game_logic.zig");
const gr = @import("game_render.zig");
const gi = @import("game_input.zig");

const pl = @import("player.zig");  

pub fn main() !void {
    const allocator = std.heap.page_allocator;

    rl.initWindow(config.SCREEN_WIDTH, config.SCREEN_HEIGHT, "Simple Test Scroller");
    defer rl.closeWindow();

    var game_state = gs.GameState.init(allocator) catch |err| {
        std.log.err("Failed to initialize game state: {s}", .{@errorName(err)});
        return;
    };
    defer game_state.deinit();


    rl.setTargetFPS(60);
    while (!rl.windowShouldClose()) {
        gi.handleInput(&game_state);
        
        gl.update(&game_state);
        
        gr.draw(game_state) catch |err| {
            std.log.err("Failed to draw game state: {s}", .{@errorName(err)});
            gs.GameState.deinit(&game_state);
            return;
        };
    }
}
