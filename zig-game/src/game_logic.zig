const config = @import("config.zig");
const gs = @import("game_state.zig");
const gi = @import("game_input.zig");
const game_bullet = @import("bullet.zig");
const rl = @import("raylib");
const std = @import("std");

// this is called every frame and controls the game's behaviour
pub fn update(state: *gs.GameState) !void {
    state.frame_counter += 1;

    //always start with plane in straight flight position as release up or down key is used for an intermediate frame ;)
    state.player.move_to_xy(state.player.pos.x, state.player.pos.y, 0);

    if (state.screen == .playing) {
        rl.updateMusicStream(state.sound.music);
    }

    if (state.screen == .gameover) {
        if (rl.isSoundPlaying(state.sound.engine)) {
            rl.stopSound(state.sound.engine);
        }
        rl.updateMusicStream(state.sound.music);
    }

    if (state.screen == .playing) {
        gamePlay(state);
    }

    try gi.handle_input(state);
}

//this is called every frame and controls the game's behaviour
fn gamePlay(state: *gs.GameState) void {
    check_collisions(state);

    if (!rl.isSoundPlaying(state.sound.engine)) {
        rl.playSound(state.sound.engine);
    }

    if (state.player.health <= 0) {
        state.screen = .gameover;
        return;
    }

    if (@mod(state.frame_counter, 10) == 0) {
        state.player.score += 20;
    }

    delete_inactive_bullets(state);

    reload_ammo(state);
    shift_bg_layers(state);
    move_enemies(state);
    move_bullets(state);
}

// check the collision boxes between player and enemies. When collision happened deduct 1 point of health from enemy and player
fn check_collisions(state: *gs.GameState) void {
    for (state.enemies.items) |*enemy| {
        if (rl.checkCollisionRecs(state.player.collision_box, enemy.collision_box)) {
            state.player.health -= 3;
            enemy.health -= 1;
            if (!rl.isSoundPlaying(state.sound.hit)) {
                rl.playSound(state.sound.hit);
            }
        }

        if (enemy.health <= 0) {
            enemy.respawn();
        }
    }

    for (state.bullets.items) |*bullet| {
        for (state.enemies.items) |*enemy| {
            if (rl.checkCollisionRecs(bullet.collision_box, enemy.collision_box)) {
                bullet.health -= 1;
                enemy.health -= 2;
                state.player.score += 10; //for every hit add 10 points to the score
                if (enemy.health <= 0) {
                    state.player.score += 100; //for every downed enemy add 1an extra 00 points to the score
                }
            }
        }
    }
}

// An alarm is sounded as warning when we move a stage up
fn progress_stage(state: *gs.GameState, current_stage: u8) void {
    if (state.stage == current_stage) {
        rl.playSound(state.sound.alert);
        state.stage += 1;
    }
}

//move the bullets in the list
fn move_bullets(state: *gs.GameState) void {
    for (state.bullets.items) |*bullet| {
        bullet.move();
    }
}

fn delete_inactive_bullets(state: *gs.GameState) void {
    var i: usize = state.bullets.items.len;
    while (i > 0) {
        i -= 1;
        const bullet = &state.bullets.items[i];
        if (bullet.pos.x < 0 or bullet.pos.x > config.SCREEN_WIDTH or
            bullet.pos.y < 0 or bullet.pos.y > config.SCREEN_HEIGHT or
            bullet.health <= 0)
        {
            _ = state.bullets.swapRemove(i);
        }
    }
}

// progress the enemies forwards, at certain scores enemy behavour speeds up and a sin movement is added
fn move_enemies(state: *gs.GameState) void {
    const sin_offset_y: i32 = @intFromFloat(std.math.sin(state.frame_counter / 20) * 10);

    for (state.enemies.items) |*enemy| {
        switch (state.player.score) {
            5000...9999 => {
                progress_stage(state, 0);
                enemy.max_speed = 7;
                enemy.move_to_xy(enemy.pos.x - enemy.speed, enemy.pos.y + (@divFloor(sin_offset_y, 8)));
            },
            10000...19999 => {
                enemy.move_to_xy(enemy.pos.x - enemy.speed, enemy.pos.y + (@divFloor(sin_offset_y, 6)));
                progress_stage(state, 1);
                enemy.max_speed = 10;
            },
            20000...29999 => {
                progress_stage(state, 2);
                enemy.move_to_xy(enemy.pos.x - enemy.speed, enemy.pos.y + (@divFloor(sin_offset_y, 4)));
                enemy.max_speed = 15;
            },
            30000...std.math.maxInt(i32) => {
                progress_stage(state, 3);
                enemy.move_to_xy(enemy.pos.x - enemy.speed, enemy.pos.y + (@divFloor(sin_offset_y, 2)));
                enemy.max_speed = 20;
            },
            else => {
                enemy.move_to_xy(enemy.pos.x - enemy.speed, enemy.pos.y);
            },
        }

        if (enemy.pos.x < -64) { //when enemy s of screen respawn
            enemy.respawn();
        }
    }
}

// when button is pressed speed up plane (move to the right) but keep level
pub fn player_right(state: *gs.GameState) void {
    if (state.player.pos.x < config.SCREEN_WIDTH - state.player.sprite.width) {
        state.player.move_to_xy(state.player.pos.x + state.player.speed, state.player.pos.y, 0);
    }
}

// when button is pressed slow down plane but keep level
pub fn player_left(state: *gs.GameState) void {
    if (state.player.pos.x > 0) {
        state.player.move_to_xy(state.player.pos.x - state.player.speed, state.player.pos.y, 0);
    }
}

// when button is pressed tilt plane 20 degrees move up and slighly backwards (from air resistance)
pub fn player_up(state: *gs.GameState) void {
    if (state.player.pos.y > 0) {
        state.player.move_to_xy(state.player.pos.x, state.player.pos.y - state.player.speed, state.player.rot);
        if (state.player.pos.x > 0) {
            state.player.move_to_xy(state.player.pos.x - @divFloor(state.player.speed, 2), state.player.pos.y, state.player.rot);
        } else {
            state.player.move_to_xy(0, state.player.pos.y, state.player.rot);
        }
        state.player.move_to_xy(state.player.pos.x, state.player.pos.y, -20);
    }
}

// when up button is pressed tilt plane 10 degrees from 20 degrees as an intermediate animation frame
pub fn player_up_release(state: *gs.GameState) void {
    state.player.move_to_xy(state.player.pos.x, state.player.pos.y, -10);
}

// when down button is pressed tilt plane 20 degrees and move plane down and forwards (airspeed going down pushes plane forwards faster)
pub fn player_down(state: *gs.GameState) void {
    if (state.player.pos.y < config.SCREEN_HEIGHT - state.player.sprite.height and state.player.pos.x < config.SCREEN_WIDTH - state.player.sprite.width) {
        state.player.move_to_xy(state.player.pos.x + state.player.speed, state.player.pos.y + state.player.speed, 20);
    }
}

// when down button is released tilt plane 10 degrees from 20 degrees as an intermediate animation frame
pub fn player_down_release(state: *gs.GameState) void {
    state.player.move_to_xy(state.player.pos.x, state.player.pos.y, 10);
}

fn create_bullet(state: *gs.GameState, bullet_speed: i32) game_bullet.Bullet {
    const player = &state.player;
    if (player.rot == -20.0) {
        return game_bullet.Bullet.init(player.pos.x + 64, player.pos.y + 5, 1, -1, bullet_speed);
    } else if (player.rot == 20.0) {
        return game_bullet.Bullet.init(player.pos.x + 64, player.pos.y + 40, 2, 1, bullet_speed);
    } else {
        return game_bullet.Bullet.init(player.pos.x + 64, player.pos.y + 16, 1, 0, bullet_speed);
    }
}

// when fire butten is pressed sound bullets ever 3 frames and play the sound
pub fn player_fire(state: *gs.GameState) !void {
    if (state.player.ammo > 0 and @mod(state.frame_counter, 3) == 0) {
        state.player.ammo -= 1;
        const bullet_speed: i32 = @intCast(state.player.speed + 2);
        const b = create_bullet(state, bullet_speed);
        try state.bullets.append(b);
    }

    if (!rl.isSoundPlaying(state.sound.gun) and state.player.ammo > 0) {
        rl.playSound(state.sound.gun);
    }

    if (rl.isSoundPlaying(state.sound.gun) and state.player.ammo == 0) {
        rl.stopSound(state.sound.gun);
    }
}

// when fire butten is released stop playing machine gun sound
pub fn player_fire_release(state: *gs.GameState) void {
    if (rl.isSoundPlaying(state.sound.gun)) {
        rl.stopSound(state.sound.gun);
    }
}

//scroll the background layers at different speeds for parallex effect
fn shift_bg_layers(state: *gs.GameState) void {
    // shift the layers layer 0 with the sun and clouds remains stationary
    state.background_layer_speed[1] += -0.1;
    state.background_layer_speed[2] += -0.2;
    state.background_layer_speed[3] += -0.5;
    state.background_layer_speed[4] += -0.9;
    state.background_layer_speed[5] += -1.3;

    //check if layer needs to be reset
    for (0..state.layers.items.len) |layer_nr| {
        if (state.background_layer_speed[layer_nr] < 0) {
            state.background_layer_speed[layer_nr] = config.SCREEN_WIDTH;
        }
    }
}

//every 120 frames 5 bullets are reloaded
fn reload_ammo(state: *gs.GameState) void {
    if (state.player.ammo < state.player.max_ammo and @mod(state.frame_counter, 120) == 0) {
        state.player.ammo += 5;
        if (state.player.ammo > state.player.max_ammo) {
            state.player.ammo = state.player.max_ammo;
        }
    }
}
