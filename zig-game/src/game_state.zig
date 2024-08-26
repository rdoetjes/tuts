const config = @import("config.zig");
const rl = @import("raylib");
const std = @import("std");
const text_scroller = @import("game_text_scroller.zig");
const game_logic = @import("game_logic.zig");
const game_player = @import("player.zig");
const game_enemy = @import("enemy.zig");
const ArrayList = std.ArrayList;

var sprite_enemy_1: rl.Texture2D = undefined;
var sprite_player_1: rl.Texture2D = undefined;

pub const GameState = struct {
    layers: ArrayList(rl.Texture2D),
    enemies: ArrayList(game_enemy.Enemy),
    score: u32,
    allocator: std.mem.Allocator,
    l1: [config.NR_BG_LAYERS]f32,
    player: game_player.Player,
    frame_counter: f32,
    snd_gun: rl.Sound,
    snd_music: rl.Music,
    font: rl.Font,

    pub fn init(allocator: std.mem.Allocator) !GameState {
        sprite_player_1 = rl.loadTexture("resources/sprites/player.png");
        const player = try game_player.Player.init(&sprite_player_1);

        sprite_enemy_1 = rl.loadTexture("resources/sprites/blimp.png");
        var enemies = ArrayList(game_enemy.Enemy).init(allocator);
        for (0..config.NR_ENEMIES) |_| {
            try enemies.append(try game_enemy.Enemy.init(&sprite_enemy_1));
        }

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

        const font = rl.loadFontEx("resources/fonts/Blankenburg.ttf", 20, null);
        rl.playMusicStream(snd_music);
        
        return GameState{
            .layers = layers,
            .l1 = l1,
            .allocator = allocator,
            .score = 0,
            .player = player,
            .enemies = enemies,
            .frame_counter = 0,
            .snd_gun = snd_gun,
            .snd_music = snd_music,
            .font = font,
        };
    }

    pub fn deinit(self: *GameState) void {
        self.layers.deinit();
        rl.unloadTexture(sprite_player_1);
        rl.unloadTexture(sprite_enemy_1);
        rl.unloadSound(self.snd_gun);
        rl.unloadMusicStream(self.snd_music);
    }
};
