const rl = @import("raylib");
const std = @import("std");
const ArrayList = std.ArrayList;

const SCREEN_WIDTH = 640;
const SCREEN_HEIGHT = 480;

const Position = struct {
    x: i32,
    y: i32,
};

const Scroller = struct {
    pos: Position,
    msg: [*:0]const u8,
    start_x_pos: i32, 
    end_x_pos: i32, 
    speed: i32,

    pub fn init(msg: [*:0]const u8, start_x: i32, end_x: i32, start_y: i32, speed: i32) Scroller {
        return .{
            .msg = msg,
            .pos = .{.x = start_x, .y = start_y,},
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

const GameState = struct {
    scrollers: ArrayList(Scroller),
    layers: ArrayList(rl.Texture2D),
    l1: [6]f32,

    pub fn init(allocator: std.mem.Allocator) !GameState {
        var scrollers = ArrayList(Scroller).init(allocator);
        try scrollers.append(Scroller.init("Your first scroller in ZIG", SCREEN_WIDTH, -450, SCREEN_HEIGHT*0.25, -5));
        try scrollers.append(Scroller.init("This one scrolls slower", -300, SCREEN_WIDTH, SCREEN_HEIGHT*0.75, 2));
        
        var layers = ArrayList(rl.Texture2D).init(allocator);
        var l1: [6]f32 = undefined;
        for (0..6) |l| {
            const layer_name = std.fmt.allocPrintZ(allocator, "layers/l{}.png", .{l+1}) catch return error.OutOfMemory;
            defer allocator.free(layer_name);
            try layers.append(rl.loadTexture(layer_name));

            l1[l] = 0.0;
        }
        return GameState{ .scrollers = scrollers,  .layers = layers, .l1=l1};
    }

    pub fn deinit(self: *GameState) void {
        self.scrollers.deinit();
    }

    pub fn update(self: *GameState) void {
        for (self.scrollers.items) |*scroller| {
            scroller.update();
        }
        
        self.l1[0] += -0.0;
        self.l1[1] += -0.1;
        self.l1[2] += -0.4;
        self.l1[3] += -0.6;
        self.l1[4] += -0.8;
        self.l1[5] += -1.0;

        var i: usize = 0;
        for (self.layers.items) |_| {
            if (self.l1[i] < 0) {
                self.l1[i] = 640.0;
            }
            i += 1;
        }
    }

    pub fn draw(self: GameState) void {
        rl.clearBackground(rl.Color.white);

        var i: usize = 0;
        for (self.layers.items) |layer| {
            rl.drawTextureEx(layer, rl.Vector2.init(self.l1[i], 0), 0.0, 2.0, rl.Color.white);
            rl.drawTextureEx(layer, rl.Vector2.init(self.l1[i] - 640, 0), 0.0, 2.0, rl.Color.white);
            i += 1;
        }

        for (self.scrollers.items) |scroller| {
            scroller.draw();
        }
    }
};

pub fn main() !void {
    const allocator = std.heap.page_allocator;
 
    rl.initWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Simple Test Scroller");
    defer rl.closeWindow();

    var game_state = try GameState.init(allocator);
    defer game_state.deinit();

    rl.setTargetFPS(60);
        
    while (!rl.windowShouldClose()) {
        game_state.update();

        rl.beginDrawing();
        game_state.draw();
        rl.endDrawing();
    }
}