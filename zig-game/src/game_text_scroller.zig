const config = @import("config.zig");
const rl = @import("raylib");
const position = @import("position.zig");

pub const Scroller = struct {
    pos: position.Position,
    msg: [*:0]const u8,
    start_x_pos: i32,
    end_x_pos: i32,
    speed: i32,

    pub fn init(msg: [*:0]const u8, start_x: i32, end_x: i32, start_y: i32, speed: i32) Scroller {
        return .{
            .msg = msg,
            .pos = .{
                .x = start_x,
                .y = start_y,
            },
            .start_x_pos = start_x,
            .end_x_pos = end_x,
            .speed = speed,
        };
    }

    pub fn update(self: *Scroller) void {
        self.pos.x += self.speed;
        if (self.speed < 0 and self.pos.x <= self.end_x_pos) {
            self.pos.x = self.start_x_pos;
        } else if (self.speed > 0 and self.pos.x >= self.end_x_pos) {
            self.pos.x = self.start_x_pos;
        }
    }

    pub fn draw(self: Scroller) void {
        rl.drawText(self.msg, self.pos.x, self.pos.y, 20, rl.Color.red);
    }
};
