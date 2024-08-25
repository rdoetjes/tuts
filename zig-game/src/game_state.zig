const config = @import("config.zig");
const rl = @import("raylib");
const std = @import("std");
const text_scroller = @import("game_text_scroller.zig");
const game_logic = @import("game_logic.zig");
const game_player = @import("player.zig");
const ArrayList = std.ArrayList;

pub const GameState = struct {
    layers: ArrayList(rl.Texture2D),
    score: u32,
    allocator: std.mem.Allocator,
    l1: [config.NR_BG_LAYERS]f32,
    player: game_player.Player,
    frame_counter: u32,
    snd_gun: rl.Sound,
    snd_music: rl.Music,

    pub fn init(allocator: std.mem.Allocator) !GameState {
        const player = try game_player.Player.init();
        const snd_gun = rl.loadSound("resources/sounds/gun.wav");
        const snd_music = rl.loadMusicStream("resources/sounds/music.wav");
        
        var layers = ArrayList(rl.Texture2D).init(allocator);
        var l1: [6]f32 = undefined;
        for (0..6) |i| {
            const layer_name = std.fmt.allocPrintZ(allocator, "resources/layers/l{}.png", .{i + 1}) catch return error.OutOfMemory;
            defer allocator.free(layer_name);
            try layers.append(rl.loadTexture(layer_name));
            l1[i] = 0.0;
        }

        rl.playMusicStream(snd_music);
        
        return GameState{
            .layers = layers,
            .l1 = l1,
            .allocator = allocator,
            .score = 0,
            .player = player,
            .frame_counter = 0,
            .snd_gun = snd_gun,
            .snd_music = snd_music,
        };
    }

    pub fn deinit(self: *GameState) void {
        self.layers.deinit();
        self.player.deinit();
        rl.unloadSound(self.snd_gun);
        rl.unloadMusicStream(self.snd_music);
    }
};
