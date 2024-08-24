const config = @import("config.zig");
const gs = @import("game_state.zig");

pub fn update(state: *gs.GameState) void {
    state.frameCounter += 1;
    if (state.player.ammo < 255 and state.frameCounter % 25 == 0) {
        state.player.ammo += 1;
    }
    shiftBgLayers(state);
}

pub fn player_right(state: *gs.GameState) void {
    if (state.player.pos.x < config.SCREEN_WIDTH-state.player.sprite.width) { 
        state.player.pos.x += state.player.speed;
    }
}

pub fn player_left(state: *gs.GameState) void {
    if (state.player.pos.x > 0) { 
        state.player.pos.x -= state.player.speed;
    }
}

pub fn player_up(state: *gs.GameState) void {
    if (state.player.pos.y > 0) { 
        state.player.pos.y -= state.player.speed;
    }
}

pub fn player_down(state: *gs.GameState) void {
    if (state.player.pos.y < config.SCREEN_HEIGHT-state.player.sprite.height) { 
        state.player.pos.y += state.player.speed;
    }
}

pub fn player_fire(state: *gs.GameState) void {
    if (state.player.ammo > 0) {
        state.player.ammo -= 1;
    }
}


fn shiftBgLayers(state: *gs.GameState) void {
    // shift the layers layer 0 with the sun and clouds remains stationary
    state.l1[1] += -0.1;
    state.l1[2] += -0.2;
    state.l1[3] += -0.5;
    state.l1[4] += -0.9;
    state.l1[5] += -1.3;

    //check if layer needs to be reset
    for (0..state.layers.items.len) |layer_nr| {
        if (state.l1[layer_nr] < 0) {
            state.l1[layer_nr] = config.SCREEN_WIDTH;
        }
    }
}
