#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>

#define NUM_CHOICES 6
#define X_NUMBER 4
#define TARGET_SUM 31

// Game state
typedef struct {
    int numbers[NUM_CHOICES];
    int available[NUM_CHOICES];
    int total;
} GameState;

// just added for proving that starting player always wins (in a perfect game)
int player_wins = 0;
int computer_wins = 0;

// Initialize game state
void init_game(GameState *game) {
    for (int i = 0; i < NUM_CHOICES; i++) {
        game->numbers[i] = i + 1;
        game->available[i] = X_NUMBER;
    }
    game->total = 0;
}

// Function to display available numbers
void show_available_moves(const GameState *game) {
    for (int j = X_NUMBER; j > 0; j--) {
        for (int i = 0; i < NUM_CHOICES; i++) {
            printf("%s", game->available[i] >= j ? "■ " : "□ ");
        }
        printf("\n");
    }

    for (int i = 0; i < NUM_CHOICES; i++) {
        printf("%d ", i + 1);
    }
}

//Show game screem
void show_game_screen(const GameState *game) {
  printf("\033[2J\033[H");  // Clear screen
  printf("Welcome to the %d Game! Inspired by Robin from 8Bit Show And Tell\n", TARGET_SUM);
  printf("First to reach exactly %d wins.\n\n", TARGET_SUM);
  printf("Available numbers:\n");

  show_available_moves(game);

  printf("\n\n");
  printf("Current total: %d\n", game->total);
}

// Check if a number is available
bool is_available(const GameState *game, int n) {
    return (n >= 1 && n <= NUM_CHOICES && game->available[n-1] > 0);
}

// Use a number and update the total
void use_number(GameState *game, int n) {
    game->available[n-1]--;
    game->total += n;
}

// Evaluate recursively with alpha-beta pruning: returns true if current player can force a win
bool can_force_win(int current_total, int *available_pool, bool is_ai_turn, int depth) {
    // Early termination for deep recursion (optimization) -- removed for this brute force prove, to show that first player always wins
    //if (depth > 10) return false;

    if (current_total == TARGET_SUM) {
        return !is_ai_turn; // Whoever just moved wins
    }

    if (current_total > TARGET_SUM) {
        return is_ai_turn; // Current player loses (exceeded target)
    }

    // Check for winning move
    for (int i = 0; i < NUM_CHOICES; i++) {
        if (available_pool[i] > 0) {
            if (current_total + (i+1) == TARGET_SUM) {
                return is_ai_turn; // Current player can win immediately
            }
        }
    }

    bool can_win = !is_ai_turn; // Default: AI assumes loss, player assumes win

    for (int i = 0; i < NUM_CHOICES; i++) {
        if (available_pool[i] > 0) {
            available_pool[i]--;
            bool win = can_force_win(current_total + (i+1), available_pool, !is_ai_turn, depth+1);
            available_pool[i]++;

            if (is_ai_turn && win) {
                return true; // AI found a winning path
            }
            if (!is_ai_turn && !win) {
                return false; // Player has a forced win, bad for AI
            }

            // Update can_win based on current findings
            if (is_ai_turn) {
                can_win |= win;
            } else {
                can_win &= win;
            }
        }
    }

    return can_win;
}

// Find the best move for the computer
int computer_move(GameState *game) {
    int best_move = -1;
    bool found_winning_move = false;

    // First, check for immediate win
    for (int i = 0; i < NUM_CHOICES; i++) {
        if (game->available[i] > 0 && game->total + (i+1) == TARGET_SUM) {
            return i+1; // Immediate win
        }
    }

    // Then look for forced wins
    for (int i = 0; i < NUM_CHOICES; i++) {
        if (game->available[i] > 0 && game->total + (i+1) < TARGET_SUM) {
            game->available[i]--;
            bool win = can_force_win(game->total + (i+1), game->available, false, 0);
            game->available[i]++; // Backtrack

            if (win) {
                best_move = i+1;
                found_winning_move = true;
                break;
            }
        }
    }

    // If no winning strategy, pick a move that doesn't immediately lose
    if (!found_winning_move) {
        int valid_moves[NUM_CHOICES];
        int valid_count = 0;

        for (int i = 0; i < NUM_CHOICES; i++) {
            if (game->available[i] > 0 && game->total + (i+1) < TARGET_SUM) {
                valid_moves[valid_count++] = i+1;
            }
        }

        if (valid_count > 0) {
            best_move = valid_moves[rand() % valid_count];
        }
    }

    return best_move;
}

// Handle player's turn
bool player_turn_handler(GameState *game) {
    int move;
    printf("Your move (1-%d): ", NUM_CHOICES);

    if (scanf("%d", &move) != 1) {
        // Clear input buffer on invalid input
        while (getchar() != '\n');
        printf("Invalid input. Please enter a number.\n");
        return false;
    }

    if (!is_available(game, move)) {
        printf("Invalid move. Try again.\n");
        return false;
    }

    use_number(game, move);
    printf("You played: %d → Total: %d\n", move, game->total);
    if (game->total > TARGET_SUM) printf("YOU OVERSHOT AND LOST!\n");

    return true;
}

// Handle computer's turn
bool computer_turn_handler(GameState *game) {
    int comp = computer_move(game);
    if (comp == -1) {
        printf("Computer has no valid moves left.\n");
        return false;
    }

    use_number(game, comp);
    printf("Computer plays: %d → Total: %d\n", comp, game->total);

    return true;
}

// Check if game is over
bool check_game_over(const GameState *game, bool is_player_turn) {
    if (game->total == TARGET_SUM) {
        if (is_player_turn) {
            show_game_screen(game); // shows the final move of the player (winning move)
            printf("You win!\n");
            player_wins++;
        } else {
            show_game_screen(game);
            printf("Computer wins!\n"); // shows the final move of the computer (winning move)
            computer_wins++;
        }
        return true;
    }
    return false;
}

int main() {
    srand((unsigned int)time(NULL));
    GameState game;

    for(int i = 0; i < 10000; i++) {
      init_game(&game);

      bool player_turn = 1;
      while (game.total < TARGET_SUM) {
          show_game_screen(&game);

          //extra game dynamics that give player 2 a 15% win average
          if (game.total == 10 && game.available[0] < NUM_CHOICES){
               game.available[0]++;
          }

          if (game.total == 24 && game.available[0] > 0){
               game.available[0]--;
               if (game.available[NUM_CHOICES-1]>0) game.available[NUM_CHOICES-1]--;
          }

          bool valid_move;
          if (player_turn) {
              valid_move = computer_turn_handler(&game);
              if (!valid_move) {
                  show_game_screen(&game); //shows the final move of the computer (when it can no longer make a move)
                  printf("You win!\n");
                  computer_wins++;
                  break;
              }

              if (check_game_over(&game, true)) break;
          } else {
              valid_move = computer_turn_handler(&game);
              if (!valid_move) {
                  show_game_screen(&game); //shows the final move of the computer (when it can no longer make a move)
                  printf("You win!\n");
                  player_wins++;
                  break;
              }

              if (check_game_over(&game, false)) break;
          }

          player_turn = !player_turn;
      }
  }
  printf("Player wins: %d computer wins %d \n", player_wins, computer_wins);
}
