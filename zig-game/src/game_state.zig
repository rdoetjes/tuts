const config = @import("config.zig");
const rl = @import("raylib");
const std = @import("std");
const text_scroller = @import("game_text_scroller.zig");
const game_logic = @import("game_logic.zig");
const ArrayList = std.ArrayList;

pub const GameState = struct {
    scrollers: ArrayList(text_scroller.Scroller),
    layers: ArrayList(rl.Texture2D),
    score: u32,
    allocator: std.mem.Allocator,
    l1: [config.NR_BG_LAYERS]f32,

    pub fn init(allocator: std.mem.Allocator) !GameState {
        var scrollers = ArrayList(text_scroller.Scroller).init(allocator);
        try scrollers.append(text_scroller.Scroller.init("Your first scroller in ZIG", config.SCREEN_WIDTH, -450, config.SCREEN_HEIGHT * 0.25, -5));
        try scrollers.append(text_scroller.Scroller.init("This one scrolls slower", -300, config.SCREEN_WIDTH, config.SCREEN_HEIGHT * 0.75, 2));

        var layers = ArrayList(rl.Texture2D).init(allocator);
        var l1: [6]f32 = undefined;
        for (0..6) |i| {
            const layer_name = std.fmt.allocPrintZ(allocator, "resources/layers/l{}.png", .{i + 1}) catch return error.OutOfMemory;
            defer allocator.free(layer_name);
            try layers.append(rl.loadTexture(layer_name));

            l1[i] = 0.0;
        }

        return GameState{
            .scrollers = scrollers,
            .layers = layers,
            .l1 = l1,
            .allocator = allocator,
            .score = 0,
        };
    }

    pub fn deinit(self: *GameState) void {
        self.scrollers.deinit();
        self.layers.deinit();
    }
};
