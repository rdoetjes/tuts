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

// Initialize game state
void init_game(GameState *game) {
    for (int i = 0; i < NUM_CHOICES; i++) {
        game->numbers[i] = i + 1;
        game->available[i] = X_NUMBER;
    }
    game->total = 0;
}

// Function to display available numbers
void show_available(const GameState *game) {
    printf("\033[2J\033[H");  // Clear screen
    printf("Welcome to the 31 Game! Inspired by Robin from 8Bit Show And Tell\n");
    printf("First to reach exactly 31 wins.\n\n");
    printf("Available numbers:\n");
    
    for (int j = X_NUMBER; j > 0; j--) {
        for (int i = 0; i < NUM_CHOICES; i++) {
            printf("%s", game->available[i] >= j ? "■ " : "□ ");
        }
        printf("\n");
    }
    
    for (int i = 0; i < NUM_CHOICES; i++) {
        printf("%d ", i + 1);
    }
    printf("\n\n");
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
    // Early termination for deep recursion (optimization)
    if (depth > 10) return false;
    
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
            printf("You win!\n");
        } else {
            printf("Computer wins!\n");
        }
        return true;
    }
    return false;
}

int main() {
    srand((unsigned int)time(NULL));
    GameState game;
    init_game(&game);
    
    bool player_turn = rand() % 2 == 0;
    printf("You %s first.\n", player_turn ? "go" : "go second");
    
    while (game.total < TARGET_SUM) {
        show_available(&game);
        printf("Current total: %d\n", game.total);

        bool valid_move;
        if (player_turn) {
            valid_move = player_turn_handler(&game);
            if (!valid_move) continue;
            
            if (check_game_over(&game, true)) break;
        } else {
            valid_move = computer_turn_handler(&game);
            if (!valid_move) {
                printf("You win!\n");
                break;
            }
            
            if (check_game_over(&game, false)) break;
        }
        
        player_turn = !player_turn;
    }

    return 0;
}
