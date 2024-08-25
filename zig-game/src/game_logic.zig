const config = @import("config.zig");
const gs = @import("game_state.zig");
const gi = @import("game_input.zig");
const rl = @import("raylib");
const std = @import("std");

pub fn update(state: *gs.GameState) void {
    state.frame_counter += 1;
    state.player.rot = 0;
    gi.handleInput(state);

    reload_ammo(state);
    shiftBgLayers(state);
    moveEnemies(state);
    
    rl.updateMusicStream(state.snd_music);
}

fn moveEnemies(state: *gs.GameState) void {
    for (state.enemies.items) |*enemy| {
        enemy.pos.x -= enemy.speed;
        if (enemy.pos.x < -64) {

            if (state.score > 10000 and state.score < 20000){
                enemy.max_speed = 10;
            }
            else if (state.score > 20000 and state.score < 30000){
                enemy.max_speed = 15;
            }
            else if (state.score > 30000){
                enemy.max_speed = 20;
            }
            
            enemy.respawn();
        }
    }
}

pub fn player_right(state: *gs.GameState) void {
    if (state.player.pos.x < config.SCREEN_WIDTH-state.player.sprite.width) { 
        state.player.pos.x += state.player.speed;
        state.player.rot = 0;
    }
}

pub fn player_left(state: *gs.GameState) void {
    if (state.player.pos.x > 0) { 
        state.player.pos.x -= state.player.speed;
        state.player.rot = 0;
    }
}

pub fn player_up(state: *gs.GameState) void {
    if (state.player.pos.y > 0 ) {

        state.player.pos.y -= state.player.speed;
        if (state.player.pos.x > 0){
            state.player.pos.x -= state.player.speed/2;
        } 
        else {
            state.player.pos.x = 0;
        }
        state.player.rot = -20;
    }
}

pub fn player_up_release(state: *gs.GameState) void {
        state.player.rot = -10;
}

pub fn player_down(state: *gs.GameState) void {
    if (state.player.pos.y < config.SCREEN_HEIGHT-state.player.sprite.height and state.player.pos.x < config.SCREEN_WIDTH-state.player.sprite.width) { 
        state.player.pos.y += state.player.speed;
        state.player.pos.x += state.player.speed;
        state.player.rot = 20;
    }
}

pub fn player_down_release(state: *gs.GameState) void {
        state.player.rot = 10;
}

pub fn player_fire(state: *gs.GameState) void {
    if (state.player.ammo > 0 and state.frame_counter % 3 == 0) {
        state.player.ammo -= 1;
    }

    if (!rl.isSoundPlaying(state.snd_gun) and state.player.ammo > 0) {
        rl.playSound(state.snd_gun);
    }

    if (rl.isSoundPlaying(state.snd_gun) and state.player.ammo == 0) {
        rl.stopSound(state.snd_gun);
    }
}

pub fn player_fire_release(state: *gs.GameState) void {
    if (rl.isSoundPlaying(state.snd_gun)) {
         rl.stopSound(state.snd_gun);
    }
}

fn shiftBgLayers(state: *gs.GameState) void {
    // shift the layers layer 0 with the sun and clouds remains stationary
    state.l1[1] += -0.1;
    state.l1[2] += -0.2;
    state.l1[3] += -0.5;
    state.l1[4] += -0.9;
    state.l1[5] += -1.3;

    //check if layer needs to be reset
    for (0..state.layers.items.len) |layer_nr| {
        if (state.l1[layer_nr] < 0) {
            state.l1[layer_nr] = config.SCREEN_WIDTH;
        }
    }
}

fn reload_ammo(state: *gs.GameState) void {
    if (state.player.ammo < state.player.max_ammo and state.frame_counter % 120 == 0) {
        state.player.ammo += 5;
        if (state.player.ammo > state.player.max_ammo) {
            state.player.ammo = state.player.max_ammo;
        }
    }
}