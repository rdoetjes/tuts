const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");


pub const Bullet = struct {
    pos: pos.Position,
    direction: pos.Position,
    speed: i32,
    health: u32,
    collision_box: rl.Rectangle,
    
    pub fn init(x: i32, y: i32, xd: i32, xy: i32) Bullet {
        const collision_box = rl.Rectangle.init(@floatFromInt(x), @floatFromInt(y), 2, 2);
        return .{
            .pos = .{ .x = x, .y = y },
            .direction = .{ .x = xd, .y = xy },
            .health = 1,
            .speed = 3,
            .collision_box = collision_box,
        };
    }

    pub fn move(self: *Bullet) void {
        self.pos.x += self.direction.x * self.speed;
        self.pos.y += self.direction.y * self.speed;
        self.collision_box.x=@floatFromInt(self.pos.x); 
        self.collision_box.y=@floatFromInt(self.pos.y);
    }

    pub fn draw(self: Bullet) void {
        rl.drawRectangle(self.pos.x, self.pos.y, 2, 2, rl.Color.black);
    }
};
