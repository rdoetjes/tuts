const config = @import("config.zig");
const gs = @import("game_state.zig");

pub fn update(state: *gs.GameState) void {
    scrollText(state);
    shiftBgLayers(state);
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

fn scrollText(state: *gs.GameState) void {
    for (state.scrollers.items) |*scroller| {
        scroller.update();
    }
}
