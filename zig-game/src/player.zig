const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");

pub const Player = struct {
     pos: pos.Position,
     health: u8,
     speed: u8,
     ammo: u8,
     max_ammo: u8,
     rot: f32,
     sprite: rl.Texture2D,

     pub fn init() !Player {
         const start_ammo = 100;
         return .{
            .pos = .{ .x = 10, .y = config.SCREEN_HEIGHT/2, },
            .speed = 3,
            .health = 10,
            .max_ammo = start_ammo,
            .ammo = start_ammo,
            .sprite = rl.loadTexture("resources/sprites/player.png"),
            .rot = 0,
         };
     }

     pub fn deinit(self: Player) void {
        rl.unloadTexture(self.sprite);
     }

     pub fn draw(self: Player) void {
        rl.drawTextureEx(self.sprite, (rl.Vector2.init(@floatFromInt(self.pos.x), @floatFromInt(self.pos.y))), self.rot, 1, rl.Color.white);
     }
};