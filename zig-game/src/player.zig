const config = @import("config.zig");
const rl = @import("raylib");
const pos = @import("position.zig");
const col_x_offset = 0;
const col_y_offset = 22;

pub const Player = struct {
   pos: pos.Position,
   health: i32,
   speed: u8,
   ammo: u8,
   max_ammo: u8,
   rot: f32,
   score: u32,
   sprite: *const rl.Texture2D,
   collision_box: rl.Rectangle,

   pub fn init(sprite: *const rl.Texture2D) !Player {
      const position: pos.Position = .{ .x = 10, .y = config.SCREEN_HEIGHT/2,};
      const start_ammo = 100;
      const collision_box = rl.Rectangle.init(position.x+col_x_offset, position.y+col_y_offset, 62, 18);

      return .{
         .pos = position,
         .speed = 3,
         .health = 100,
         .max_ammo = start_ammo,
         .ammo = start_ammo,
         .sprite = sprite,
         .rot = 0,
         .score = 0,
         .collision_box = collision_box,
      };
   }

   pub fn moveToXY(self: *Player, x: i32, y: i32, rot: f32)  void {
      self.pos.x = x;
      self.pos.y = y;
      self.rot = rot;
      if (self.rot >= 0) {
         self.collision_box.x = @floatFromInt(x + col_x_offset);
         self.collision_box.y = @floatFromInt(y + col_y_offset);
      } else {
         self.collision_box.x = @floatFromInt(x + col_x_offset + 9); // this compensates from some weird rotation offset because of sprite
      }
   }

   pub fn draw(self: Player) void {
      //rl.drawRectanglePro(self.collision_box, (rl.Vector2.init(@floatFromInt(col_x_offset/2), @floatFromInt(col_y_offset/2))), self.rot, rl.Color.blue);
      rl.drawTextureEx(self.sprite.*, (rl.Vector2.init(@floatFromInt(self.pos.x), @floatFromInt(self.pos.y))), self.rot, 1, rl.Color.white);
   }
};