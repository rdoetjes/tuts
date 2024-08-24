const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");

pub const Player = struct {
     pos: pos.Position,
     health: u8,
     speed: u8,
     ammo: u8,
     max_ammo: u8,
     sprite: rl.Texture2D,

     pub fn init() !Player {
         const start_ammo = 100;
         return .{
            .pos = .{ .x = config.SCREEN_HEIGHT/2, .y = config.SCREEN_WIDTH/2, },
            .speed = 3,
            .health = 10,
            .max_ammo = start_ammo,
            .ammo = start_ammo,
            .sprite = rl.loadTexture("resources/sprites/player.png"),
         };
     }

     pub fn deinit(self: Player) void {
        rl.unloadTexture(self.sprite);
     }

     pub fn draw(self: Player) void {
        rl.drawTexture(self.sprite, self.pos.x, self.pos.y, rl.Color.white);
     }
};