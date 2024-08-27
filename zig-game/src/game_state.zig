const config = @import("config.zig");
const rl = @import("raylib");
const std = @import("std");
const text_scroller = @import("game_text_scroller.zig");
const game_logic = @import("game_logic.zig");
const game_player = @import("player.zig");
const game_enemy = @import("enemy.zig");
const ArrayList = std.ArrayList;
const screen = enum { splash, playing, gameover };
var sprite_enemy_1: rl.Texture2D = undefined;
var sprite_player_1: rl.Texture2D = undefined;

pub const Sound = struct {
    gun: rl.Sound,
    hit: rl.Sound,
    alert: rl.Sound,
    music: rl.Music,

    pub fn deinit(self: *Sound) void {
        rl.unloadSound(self.gun);
        rl.unloadSound(self.hit);
        rl.unloadSound(self.alert);
        rl.unloadMusicStream(self.music);
    }
};

pub const GameState = struct {
    layers: ArrayList(rl.Texture2D),
    gameover_image: rl.Texture2D,
    splash_image: rl.Texture2D,
    enemies: ArrayList(game_enemy.Enemy),
    sound: Sound,
    score: u32,
    allocator: std.mem.Allocator,
    l1: [config.NR_BG_LAYERS]f32,
    player: game_player.Player,
    frame_counter: f32,
  
    stage: u32,
    screen: screen,
    font: rl.Font,

    pub fn init(allocator: std.mem.Allocator) !GameState {
        sprite_player_1 = rl.loadTexture("resources/sprites/player.png");
        const player = try game_player.Player.init(&sprite_player_1);

        sprite_enemy_1 = rl.loadTexture("resources/sprites/blimp.png");
        var enemies = ArrayList(game_enemy.Enemy).init(allocator);
        for (0..config.NR_ENEMIES) |_| {
            try enemies.append(try game_enemy.Enemy.init(&sprite_enemy_1));
        }

        const sound = Sound{.gun = rl.loadSound("resources/sounds/gun.wav"),
                            .hit = rl.loadSound("resources/sounds/hit.wav"),
                            .alert = rl.loadSound("resources/sounds/alert.wav"),
                            .music = rl.loadMusicStream("resources/sounds/music.wav"),
        };

        var layers = ArrayList(rl.Texture2D).init(allocator);
        var l1: [6]f32 = undefined;
        for (0..6) |i| {
            const layer_name = std.fmt.allocPrintZ(allocator, "resources/layers/l{}.png", .{i + 1}) catch return error.OutOfMemory;
            defer allocator.free(layer_name);
            try layers.append(rl.loadTexture(layer_name));
            l1[i] = 0.0;
        }
        const gameover_image = rl.loadTexture("resources/layers/gameover.png");
        const splash_image= rl.loadTexture("resources/layers/splash.png");
        const font = rl.loadFontEx("resources/fonts/Blankenburg.ttf", 20, null);
        rl.playMusicStream(sound.music);

        return GameState{
            .layers = layers,
            .l1 = l1,
            .allocator = allocator,
            .score = 0,
            .player = player,
            .enemies = enemies,
            .frame_counter = 0,
            .sound = sound,
            .font = font,
            .stage = 0,
            .screen = screen.splash,
            .gameover_image = gameover_image,
            .splash_image = splash_image,
        };
    }

    pub fn reset(self: *GameState) void {
        self.frame_counter = 0;
        self.score = 0;
        self.stage = 0;
        rl.seekMusicStream(self.sound.music, 0.0);
        self.player = try game_player.Player.init(&sprite_player_1);
        for (self.enemies.items) |*enemy| {
            enemy.* = try game_enemy.Enemy.init(&sprite_enemy_1);
        }
        self.screen = screen.playing;
    }

    pub fn deinit(self: *GameState) void {
        for (self.layers.items) |layer| {
            rl.unloadTexture(layer);
        }
        rl.unloadTexture(self.gameover_image);
        rl.unloadTexture(self.splash_image);
        self.layers.deinit();
        self.enemies.deinit();
        rl.unloadTexture(sprite_player_1);
        rl.unloadTexture(sprite_enemy_1);
        self.sound.deinit();
    }
};
