const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");
const std = @import("std");
const bcauge = @import("bar_cauge.zig");
const col_y_offset = 20;
const col_x_offset = 13;
var prng = std.rand.DefaultPrng.init(666);
const rand = prng.random();

//The enemy has a collision box that is rectangle and is smaller than the sprite to make collision detection fairer
//using a box is the easieast way to do collision detection, not very accurate but good enough
pub const EnemyBlimp = struct {
    pos: pos.Position,
    health: f32,
    max_health: f32,
    speed: u8,
    max_speed: u8,
    sprite: *const rl.Texture2D,
    collision_box: rl.Rectangle,
    health_bar: bcauge.BarCauge,

    pub fn init(sprite: *const rl.Texture2D) !EnemyBlimp {
        const position: pos.Position = .{
            .x = rand.intRangeAtMost(i32, config.SCREEN_HEIGHT, config.SCREEN_HEIGHT * 2),
            .y = rand.intRangeAtMost(i32, 0, config.SCREEN_HEIGHT - 64),
        };
        const max_speed = 3;
        const max_health = 5;
        const collision_box = rl.Rectangle.init(@floatFromInt(position.x + col_x_offset), @floatFromInt(position.y + col_y_offset), 44, 30);
        const health_bar = bcauge.BarCauge.init(position.x + 40, position.y + 20, max_health, max_health, 20, 4);

        return .{
            .pos = position,
            .speed = rand.intRangeAtMost(u8, 3, max_speed),
            .max_speed = max_speed,
            .max_health = max_health,
            .health = max_health,
            .sprite = sprite,
            .collision_box = collision_box,
            .health_bar = health_bar,
        };
    }

    //when a enemy is dead or off screen, it is moved offscreen to a new random height
    //this way we don't need to waste time instatiating new enemies
    pub fn respawn(self: *EnemyBlimp) void {
        self.speed = rand.intRangeAtMost(u8, 3, self.max_speed);
        self.pos.x = rand.intRangeAtMost(i32, config.SCREEN_WIDTH + 64, config.SCREEN_WIDTH * 2);
        self.pos.y = rand.intRangeAtMost(i32, 0, config.SCREEN_HEIGHT);
        self.health = self.max_health;
    }

    //position enemey and it's collision(with offset) box on XY
    pub fn move_to_xy(self: *EnemyBlimp, x: i32, y: i32) void {
        const x_offset = 20;
        const y_offset = 40;
        self.pos.x = x;
        self.pos.y = y;
        self.collision_box.x = @floatFromInt(x + col_x_offset);
        self.collision_box.y = @floatFromInt(y + col_y_offset);
        self.health_bar.move_to_xy(x + x_offset, y + y_offset);
    }

    //draw the enemy sprite (not the collision box)
    pub fn draw(self: EnemyBlimp) void {
        rl.drawTexture(self.sprite.*, self.pos.x, self.pos.y, rl.Color.white);
        self.health_bar.draw(self.health);
        //rl.drawRectangleRec(self.collision_box, rl.Color.red);
    }
};
