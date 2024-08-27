const rl = @import("raylib");
const config = @import("config.zig");
const pos = @import("position.zig");

pub const BarCauge = struct {
    pos: pos.Position,
    health: f32,
    max_health: f32,
    width: f32,
    height: f32,
    
    pub fn init(x: i32, y: i32, health: f32, max_health: f32, width: f32, height: f32) BarCauge {
        return .{
            .pos = .{ .x = x, .y = y },
            .health = health,
            .max_health = max_health,
            .width = width,
            .height = height,
        };
    }

    pub fn moveToXY(self: *BarCauge, x: i32, y: i32) void {
        self.pos.x = x;
        self.pos.y = y;
    }

    pub fn draw(self: *const BarCauge, health: f32) void {
        rl.drawRectangle(self.pos.x, self.pos.y, @intFromFloat(self.width), @intFromFloat(self.height), rl.Color.red);
        const bar_width = self.width * (health/self.max_health);
        rl.drawRectangle(self.pos.x, self.pos.y, @intFromFloat(bar_width), @intFromFloat(self.height), rl.Color.green);
    }
};