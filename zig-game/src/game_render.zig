const config = @import("config.zig");
const std = @import("std");
const rl = @import("raylib");
const gs = @import("game_state.zig");

pub fn draw(state: gs.GameState) !void {
    rl.beginDrawing();
    rl.clearBackground(rl.Color.white);

    for (0..state.layers.items.len) |layer_nr| {
        rl.drawTextureEx(state.layers.items[layer_nr], rl.Vector2.init(state.l1[layer_nr], 0), 0.0, config.SCREEN_WIDTH / config.BG_IMAGE_WIDTH, rl.Color.white);
        rl.drawTextureEx(state.layers.items[layer_nr], rl.Vector2.init(state.l1[layer_nr] - config.SCREEN_WIDTH, 0), 0.0, config.SCREEN_WIDTH / config.BG_IMAGE_WIDTH, rl.Color.white);

        // this is the layer with the action
        if (layer_nr == config.PLAYFIELD_LAYER) {
            drawGameItems(state);
        }
    }
    try drawHud(state);
    rl.endDrawing();
}

fn drawGameItems(state: gs.GameState) void {
    state.player.draw();
    for (state.scrollers.items) |scroller| {
        scroller.draw();
    }
}

fn drawHud(state: gs.GameState) !void {
    const fps = rl.getFPS();
    const hud = std.fmt.allocPrintZ(state.allocator, "SCORE: {d:0<6}  FPS: {d}", .{
        state.score,
        fps,
    }) catch return error.OutOfMemory;
    defer state.allocator.free(hud);
    rl.drawText(hud, 10, 10, 20, rl.Color.black);
}
