const config = @import("config.zig");
const std = @import("std");
const rl = @import("raylib");
const gs = @import("game_state.zig");

// draw the game every frame
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

// draw the in game items
fn drawGameItems(state: gs.GameState) void {
    state.player.draw();

    // bullets go behind the enemies as to suggest a hit
    for(state.bullets.items) |bullet| {
        bullet.draw();
    }

    for(state.enemies.items) |enemy| {
        enemy.draw();
    }
}

//draw the text HUD
fn drawHud(state: gs.GameState) !void {
    const score = try std.fmt.allocPrintZ(state.allocator, "SCORE: {d:0>6}", .{
        state.player.score,
    });
    defer state.allocator.free(score);

    state.ammo_bar.draw(@floatFromInt(state.player.ammo));
    state.health_bar.draw(state.player.health);
    rl.drawTextEx(state.font, score, (rl.Vector2.init(10, 10)), 30, 2, rl.Color.black);
    rl.drawTextEx(state.font, "  AMMO  ",  (rl.Vector2.init(@floatFromInt((config.SCREEN_WIDTH/2)-50), 10)), 30, 2, rl.Color.black);
    rl.drawTextEx(state.font, " HEALTH ",  (rl.Vector2.init(@floatFromInt( (config.SCREEN_WIDTH/2)+185), 10)), 30, 2, rl.Color.black);

    //rl.drawFPS(10, 30);
}
