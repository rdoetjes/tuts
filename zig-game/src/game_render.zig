const config = @import("config.zig");
const std = @import("std");
const rl = @import("raylib");
const gs = @import("game_state.zig");

pub fn draw(state: gs.GameState) !void {
    rl.beginDrawing();
    rl.clearBackground(rl.Color.white);


    if (state.screen == .gameover) {
        rl.drawTexture(state.gameover_image, 0, 0, rl.Color.white);
        const score = try std.fmt.allocPrintZ(state.allocator, "SCORE: {d:0>6}", .{state.player.score});
        defer state.allocator.free(score);
        rl.drawTextEx(state.font, score,  (rl.Vector2.init(@floatFromInt( (config.SCREEN_WIDTH/2)-96 ), @floatFromInt( (config.SCREEN_HEIGHT-40)))), 50, 2, rl.Color.white);
    }

    if (state.screen == .splash) {
        rl.drawTexture(state.splash_image, 0, 0, rl.Color.white);
    }

    if (state.screen == .playing) {
        for (0..state.layers.items.len) |layer_nr| {
            rl.drawTextureEx(state.layers.items[layer_nr], rl.Vector2.init(state.background_layer_speed[layer_nr], 0), 0.0, config.SCREEN_WIDTH / config.BG_IMAGE_WIDTH, rl.Color.white);
            rl.drawTextureEx(state.layers.items[layer_nr], rl.Vector2.init(state.background_layer_speed[layer_nr] - config.SCREEN_WIDTH, 0), 0.0, config.SCREEN_WIDTH / config.BG_IMAGE_WIDTH, rl.Color.white);

            // this is the layer with the action
            if (layer_nr == config.PLAYFIELD_LAYER) {
                drawGameItems(state);
                try drawHud(state);
            }
        }
    }
    rl.endDrawing();
}

fn drawGameItems(state: gs.GameState) void {
    state.player.draw();
    for(state.enemies.items) |enemy| {
        enemy.draw();
    }
}

fn drawHud(state: gs.GameState) !void {
    //const fps = rl.getFPS();
    
    const score = try std.fmt.allocPrintZ(state.allocator, "SCORE: {d:0>6}", .{
        state.player.score,
    });
    defer state.allocator.free(score);

    const ammo = try std.fmt.allocPrintZ(state.allocator, "AMMO: {d:0>3}", .{
        state.player.ammo,
    });
    defer state.allocator.free(ammo);


    const health = try std.fmt.allocPrintZ(state.allocator, "HEALTH: {d:0>2}", .{
        state.player.health,
    });
    defer state.allocator.free(health);


    rl.drawTextEx(state.font, score, (rl.Vector2.init(10, 10)), 30, 2, rl.Color.black);
    rl.drawTextEx(state.font, ammo,  (rl.Vector2.init(@floatFromInt((config.SCREEN_WIDTH/2)-(ammo.len/2*11)), 10)), 30, 2, rl.Color.black);
    rl.drawTextEx(state.font, health,  (rl.Vector2.init(@floatFromInt( (config.SCREEN_WIDTH-10)-(health.len*11)), 10)), 30, 2, rl.Color.black);

    rl.drawFPS(10, 30);
}
