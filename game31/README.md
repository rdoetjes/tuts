# Game31

This was inspired by the video made by Robin from 8Bit Show and Tell

[![Game31 by 8Bit Show and Tell](https://img.youtube.com/vi/hyHeQgrvu4w/0.jpg)](https://www.youtube.com/watch?v=hyHeQgrvu4w)

## Compilation of User game
This compiles the user game, it randomly selects who goes first (the cpu or the user) and then plays the game.
```
gcc -O3 -o game31 game.c
```

## Compilation of Computer game
This is "quick batardized" version of the game.c where the computer plays itself.

All the prints are removed and the game is played automatically, making sure to reset the game after each round and start with "player 1" (the first player).

This proved my suspicion, that this game is onbalanced, and that the person who starts first and plays a perfect game will always win!
```
gcc -O3 -o game_autoplay game_autoplay.c
```

