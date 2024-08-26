const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");
const std = @import("std");

const col_y_offset = 20;
const col_x_offset = 13;
var prng = std.rand.DefaultPrng.init(666); //fixed seed so that every game is the same
const rand = prng.random();

pub const Enemy = struct {
    pos: pos.Position,
    health: u8,
    max_health: u8,
    speed: u8,
    max_speed: u8,
    sprite: *const rl.Texture2D,
    collision_box: rl.Rectangle,

     pub fn init(sprite: *const rl.Texture2D) !Enemy {
        const position: pos.Position = .{ .x = rand.intRangeAtMost(i32, config.SCREEN_HEIGHT, config.SCREEN_HEIGHT*2), .y = rand.intRangeAtMost(i32, 0, config.SCREEN_HEIGHT-64), };
        const max_speed = 3;
        const max_health = 5;
        const collision_box = rl.Rectangle.init(@floatFromInt(position.x+col_x_offset), @floatFromInt(position.y+col_y_offset), 44, 30);

        return .{
            .pos = position,
            .speed = rand.intRangeAtMost(u8, 3, max_speed),
            .max_speed = max_speed,
            .max_health = max_health,
            .health = max_health,
            .sprite = sprite,
            .collision_box = collision_box,
         };
     }

    pub fn respawn(self: *Enemy) void {
        self.speed = rand.intRangeAtMost(u8, 3, self.max_speed);
        self.pos.x = rand.intRangeAtMost(i32, config.SCREEN_WIDTH+64, config.SCREEN_WIDTH*2);
        self.pos.y = rand.intRangeAtMost(i32, 0, config.SCREEN_HEIGHT);
        self.health = self.max_health;
    }

    pub fn moveToXY(self: *Enemy, x: i32, y: i32)  void {
      self.pos.x = x;
      self.pos.y = y;
      self.collision_box.x = @floatFromInt(x+col_x_offset);
      self.collision_box.y = @floatFromInt(y+col_y_offset);
    }

    pub fn draw(self: Enemy) void {
        rl.drawTexture(self.sprite.*, self.pos.x, self.pos.y, rl.Color.white);
        //rl.drawRectangleRec(self.collision_box, rl.Color.red);
    }
};