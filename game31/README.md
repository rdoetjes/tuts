# Game31: A Command-Line Game Implementation

This project is a C implementation of a game, referred to here as "Game31." The concept and inspiration for this implementation came from a video by Robin of 8Bit Show and Tell.

[![Game31 by 8Bit Show and Tell](https://img.youtube.com/vi/hyHeQgrvu4w/0.jpg)](https://www.youtube.com/watch?v=hyHeQgrvu4w)

*(**Developer Note:** You may want to add a brief 1-2 sentence description of the actual rules/objective of "Game31" here for clarity, e.g., "Players take turns adding a number, from the pool of available numners, to a running total, aiming to reach 31, without going over.")*

## Versions Included

This repository contains two C programs related to Game31:

### 1. Interactive Mode: Player vs. CPU (`game.c`)

This program allows a human user to play Game31 against a CPU opponent.

*   **Gameplay:**
    *   The game randomly determines whether the CPU or the user makes the first move.
    *   Players then take turns according to the rules of Game31.
*   **Compilation:**
    To compile the interactive game, use the following command:
    ```bash
    gcc -O3 -o game31 game.c
    ```
*   **Usage:**
    After successful compilation, run the game with:
    ```bash
    ./game31
    ```

### 2. Automated Mode: CPU vs. CPU Simulation (`game_autoplay.c`)

This program is a "quickly bastardized" version of `game.c` designed for automated gameplay where the computer plays against itself. It serves as a tool to prove that the game is imbalanced.

*   **Key Modifications:**
    *   User interaction (prompts, input) and detailed print statements have been removed to facilitate rapid, automated simulations.
    *   The game automatically resets after each round.
    *   "Player 1" (the first simulated player to move) consistently starts each new game.
*   **Compilation:**
    To compile the automated simulation, use the following command:
    ```bash
    gcc -O3 -o game_autoplay game_autoplay.c
    ```
*   **Usage:**
    After successful compilation, run the simulation with:
    ```bash
    ./game_autoplay
    ```
    *(Note: Output will likely be minimal or summarized, as per its design for analysis rather than interactive play.)*

## Key Finding: Game Balance

The development and use of the `game_autoplay.c` simulation led to a significant observation regarding the imbalance of Game31:

**The game appears to be deterministic with an optimal strategy. Specifically, the player who starts first and employs a perfect game strategy will consistently win.**

This finding highlights a potential imbalance in Game31, where the starting position, combined with optimal play, dictates the outcome.

## Prerequisites

*   A C compiler (e.g., GCC) is required to compile the programs.