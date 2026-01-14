#include <stdio.h>
#include <stdlib.h>
#include <limits.h>

#define ROWS 6
#define COLS 7
#define DEPTH 6 // Minimax lookahead depth

char board[ROWS][COLS];

// Function prototypes
void initialize_board();
void print_board();
int make_move(int col, char player);
int undo_move(int col);
int check_winner(char player);
int minimax(int depth, int is_maximizing);
int ai_move();
int is_full();
int evaluate_board();
int evaluate_window(char window[4], char player);

int main() {
    initialize_board();
    int current_column, game_over = 0;

    printf("Welcome to Connect 4! You are 'O' and the computer is 'X'.\n");
    print_board();

    while (!game_over) {
        // Player's turn
        printf("Enter column (1-7) to drop your disc: ");
        scanf("%d", &current_column);
        if (current_column < 1 || current_column > 7 || !make_move(current_column - 1, 'O')) {
            printf("Invalid move. Try again.\n");
            continue;
        }
        print_board();
        if (check_winner('O')) {
            printf("Congratulations! You won!\n");
            break;
        }
        if (is_full()) {
            printf("It's a draw!\n");
            break;
        }

        // AI's turn
        printf("Computer's turn...\n");
        ai_move();
        print_board();
        if (check_winner('X')) {
            printf("The computer wins. Better luck next time!\n");
            break;
        }
        if (is_full()) {
            printf("It's a draw!\n");
            break;
        }
    }

    return 0;
}

// Initialize the game board
void initialize_board() {
    for (int i = 0; i < ROWS; i++)
        for (int j = 0; j < COLS; j++)
            board[i][j] = ' ';
}

// Print the game board
void print_board() {
    for (int i = 0; i < ROWS; i++) {
        for (int j = 0; j < COLS; j++)
            printf("| %c ", board[i][j]);
        printf("|\n");
    }
    printf(" -----------------------------\n");
    printf("  1   2   3   4   5   6   7  \n\n");
}

// Make a move on the board
int make_move(int col, char player) {
    for (int i = ROWS - 1; i >= 0; i--) {
        if (board[i][col] == ' ') {
            board[i][col] = player;
            return 1;
        }
    }
    return 0;
}

// Undo a move on the board
int undo_move(int col) {
    for (int i = 0; i < ROWS; i++) {
        if (board[i][col] != ' ') {
            board[i][col] = ' ';
            return 1;
        }
    }
    return 0;
}

// Check if there is a winner
int check_winner(char player) {
    for (int r = 0; r < ROWS; r++) {
        for (int c = 0; c < COLS; c++) {
            if (board[r][c] == player) {
                if (c <= COLS - 4 && board[r][c + 1] == player && board[r][c + 2] == player && board[r][c + 3] == player) return 1; // Horizontal
                if (r <= ROWS - 4 && board[r + 1][c] == player && board[r + 2][c] == player && board[r + 3][c] == player) return 1; // Vertical
                if (r <= ROWS - 4 && c <= COLS - 4 && board[r + 1][c + 1] == player && board[r + 2][c + 2] == player && board[r + 3][c + 3] == player) return 1; // Diagonal (down-right)
                if (r <= ROWS - 4 && c >= 3 && board[r + 1][c - 1] == player && board[r + 2][c - 2] == player && board[r + 3][c - 3] == player) return 1; // Diagonal (down-left)
            }
        }
    }
    return 0;
}

// Check if the board is full
int is_full() {
    for (int c = 0; c < COLS; c++)
        if (board[0][c] == ' ') return 0;
    return 1;
}

// Evaluate a window of 4 cells
int evaluate_window(char window[4], char player) {
    int score = 0;
    char opponent = (player == 'X') ? 'O' : 'X';
    int count_player = 0, count_opponent = 0, count_empty = 0;

    for (int i = 0; i < 4; i++) {
        if (window[i] == player) count_player++;
        else if (window[i] == opponent) count_opponent++;
        else count_empty++;
    }

    if (count_player == 4) score += 100;
    else if (count_player == 3 && count_empty == 1) score += 10;
    else if (count_player == 2 && count_empty == 2) score += 5;

    if (count_opponent == 3 && count_empty == 1) score -= 80; // block threats

    return score;
}

// Evaluate the entire board for the AI
int evaluate_board() {
    int score = 0;

    // Horizontal
    for (int r = 0; r < ROWS; r++) {
        for (int c = 0; c <= COLS - 4; c++) {
            char window[4] = {board[r][c], board[r][c + 1], board[r][c + 2], board[r][c + 3]};
            score += evaluate_window(window, 'X');
        }
    }

    // Vertical
    for (int c = 0; c < COLS; c++) {
        for (int r = 0; r <= ROWS - 4; r++) {
            char window[4] = {board[r][c], board[r + 1][c], board[r + 2][c], board[r + 3][c]};
            score += evaluate_window(window, 'X');
        }
    }

    // Diagonal down-right
    for (int r = 0; r <= ROWS - 4; r++) {
        for (int c = 0; c <= COLS - 4; c++) {
            char window[4] = {board[r][c], board[r + 1][c + 1], board[r + 2][c + 2], board[r + 3][c + 3]};
            score += evaluate_window(window, 'X');
        }
    }

    // Diagonal down-left
    for (int r = 0; r <= ROWS - 4; r++) {
        for (int c = 3; c < COLS; c++) {
            char window[4] = {board[r][c], board[r + 1][c - 1], board[r + 2][c - 2], board[r + 3][c - 3]};
            score += evaluate_window(window, 'X');
        }
    }

    return score;
}

// Minimax algorithm with heuristic evaluation
int minimax(int depth, int is_maximizing) {
    if (check_winner('X')) return 1000;
    if (check_winner('O')) return -1000;
    if (is_full() || depth == 0) return evaluate_board();

    if (is_maximizing) {
        int max_eval = INT_MIN;
        for (int col = 0; col < COLS; col++) {
            if (make_move(col, 'X')) {
                int eval = minimax(depth - 1, 0);
                undo_move(col);
                if (eval > max_eval) max_eval = eval;
            }
        }
        return max_eval;
    } else {
        int min_eval = INT_MAX;
        for (int col = 0; col < COLS; col++) {
            if (make_move(col, 'O')) {
                int eval = minimax(depth - 1, 1);
                undo_move(col);
                if (eval < min_eval) min_eval = eval;
            }
        }
        return min_eval;
    }
}

// AI move using minimax
int ai_move() {
    int best_score = INT_MIN;
    int best_col = 0;

    for (int col = 0; col < COLS; col++) {
        if (make_move(col, 'X')) {
            int score = minimax(DEPTH - 1, 0);
            undo_move(col);

            if (score > best_score) {
                best_score = score;
                best_col = col;
            }
        }
    }
    make_move(best_col, 'X');
    return best_col;
}
