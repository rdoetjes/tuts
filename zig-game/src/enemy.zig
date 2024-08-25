const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");
const std = @import("std");

var prng = std.rand.DefaultPrng.init(121221);
const rand = prng.random();

pub const Enemy = struct {
     pos: pos.Position,
     health: u8,
     max_health: u8,
     speed: u8,
     max_speed: u8,
     sprite: rl.Texture2D,

     pub fn init() !Enemy {
        const max_speed = 5;
        const max_health = 5;

        return .{
            .pos = .{ .x = rand.intRangeAtMost(i32, config.SCREEN_HEIGHT, config.SCREEN_HEIGHT*2), .y = rand.intRangeAtMost(i32, 0, config.SCREEN_HEIGHT-64), },
            .speed = rand.intRangeAtMost(u8, 3, max_speed),
            .max_speed = max_speed,
            .max_health = max_health,
            .health = 10,
            .sprite = rl.loadTexture("resources/sprites/blimp.png"),
         };
     }

    pub fn respawn(self: *Enemy) void {
        self.speed = rand.intRangeAtMost(u8, 3, self.max_speed);
        self.pos.x = rand.intRangeAtMost(i32, config.SCREEN_WIDTH+64, config.SCREEN_WIDTH*2);
        self.pos.y = rand.intRangeAtMost(i32, 0, config.SCREEN_HEIGHT);
        self.health = self.max_health;
    }

    pub fn deinit(self: Enemy) void {
        rl.unloadTexture(self.sprite);
     }

     pub fn draw(self: Enemy) void {
        rl.drawTexture(self.sprite, self.pos.x, self.pos.y, rl.Color.white);
     }
};