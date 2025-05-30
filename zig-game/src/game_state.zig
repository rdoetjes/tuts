const config = @import("config.zig");
const rl = @import("raylib");
const std = @import("std");
const game_logic = @import("game_logic.zig");
const game_player = @import("player.zig");
const game_enemy = @import("enemy_blimp.zig");
const gb = @import("bullet.zig");
const bar_graph = @import("bar_cauge.zig");
const ArrayList = std.ArrayList;
const screen = enum { splash, playing, gameover };
var sprite_enemy_1: rl.Texture2D = undefined;
var sprite_player_1: rl.Texture2D = undefined;

//Structure that holds all sounds (just convenience)
pub const Sound = struct {
    gun: rl.Sound,
    hit: rl.Sound,
    alert: rl.Sound,
    music: rl.Music,
    engine: rl.Sound,

    pub fn deinit(self: *Sound) void {
        rl.unloadSound(self.gun);
        rl.unloadSound(self.hit);
        rl.unloadSound(self.alert);
        rl.unloadSound(self.engine);
        rl.unloadMusicStream(self.music);
    }
};

//GameState holds all the variables that are needed to run the game
//It initializes and resets them
pub const GameState = struct {
    layers: ArrayList(rl.Texture2D),
    gameover_image: rl.Texture2D,
    splash_image: rl.Texture2D,
    enemies: ArrayList(game_enemy.EnemyBlimp),
    bullets: ArrayList(gb.Bullet),
    sound: Sound,
    allocator: std.mem.Allocator,
    background_layer_speed: [config.NR_BG_LAYERS]f32,
    player: game_player.Player,
    frame_counter: f32,
    stage: u32,
    screen: screen,
    font: rl.Font,
    ammo_bar: bar_graph.BarCauge,
    health_bar: bar_graph.BarCauge,

    pub fn init(allocator: std.mem.Allocator) !GameState {
        sprite_player_1 = try rl.loadTexture("resources/sprites/player.png");
        const player = try game_player.Player.init(&sprite_player_1);

        sprite_enemy_1 = try rl.loadTexture("resources/sprites/blimp.png");
        var enemies = ArrayList(game_enemy.EnemyBlimp).init(allocator);
        for (0..config.NR_ENEMIES) |_| {
            try enemies.append(try game_enemy.EnemyBlimp.init(&sprite_enemy_1));
        }

        const bullets = ArrayList(gb.Bullet).init(allocator);

        const sound = Sound{
            .gun = try rl.loadSound("resources/sounds/gun.wav"),
            .hit = try rl.loadSound("resources/sounds/hit.wav"),
            .alert = try rl.loadSound("resources/sounds/alert.wav"),
            .music = try rl.loadMusicStream("resources/sounds/music.wav"),
            .engine = try rl.loadSound("resources/sounds/engine.wav"),
        };

        var layers = ArrayList(rl.Texture2D).init(allocator);
        var background_layer_speed: [6]f32 = undefined;
        for (0..6) |i| {
            const layer_name = std.fmt.allocPrintZ(allocator, "resources/layers/l{}.png", .{i + 1}) catch return error.OutOfMemory;
            defer allocator.free(layer_name);
            try layers.append(try rl.loadTexture(layer_name));
            background_layer_speed[i] = 0.0;
        }
        const gameover_image = try rl.loadTexture("resources/layers/gameover.png");
        const splash_image = try rl.loadTexture("resources/layers/splash.png");

        const ammo_bar = bar_graph.BarCauge.init(270, 10, 100, 100, 100, 30);
        const health_bar = bar_graph.BarCauge.init(500, 10, 100, 100, 100, 30);

        const font = try rl.loadFontEx("resources/fonts/Blankenburg.ttf", 20, null);
        rl.playMusicStream(sound.music);

        return GameState{
            .layers = layers,
            .background_layer_speed = background_layer_speed,
            .allocator = allocator,
            .player = player,
            .enemies = enemies,
            .frame_counter = 0,
            .sound = sound,
            .font = font,
            .stage = 0,
            .screen = screen.splash,
            .gameover_image = gameover_image,
            .splash_image = splash_image,
            .bullets = bullets,
            .ammo_bar = ammo_bar,
            .health_bar = health_bar,
        };
    }

    pub fn reset(self: *GameState) void {
        self.frame_counter = 0;
        self.player.score = 0;
        self.stage = 0;
        rl.seekMusicStream(self.sound.music, 0.0);
        self.player = try game_player.Player.init(&sprite_player_1);
        for (self.enemies.items) |*enemy| {
            enemy.* = try game_enemy.EnemyBlimp.init(&sprite_enemy_1);
        }
        self.bullets.clearRetainingCapacity();
        self.screen = screen.playing;
    }

    pub fn deinit(self: *GameState) void {
        for (self.layers.items) |layer| {
            rl.unloadTexture(layer);
        }
        self.bullets.deinit();
        rl.unloadTexture(self.gameover_image);
        rl.unloadTexture(self.splash_image);
        self.layers.deinit();
        self.enemies.deinit();
        rl.unloadTexture(sprite_player_1);
        rl.unloadTexture(sprite_enemy_1);
        self.sound.deinit();
    }
};
