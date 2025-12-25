#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#define N 8 // Chessboard size

// Directions in which a knight can move
int knight_moves[8][2] = {
    {2, 1}, {1, 2}, {-1, 2}, {-2, 1},
    {-2, -1}, {-1, -2}, {1, -2}, {2, -1}
};

// Check if a move is valid
bool is_valid(int x, int y, int board[N][N]) {
    return (x >= 0 && y >= 0 && x < N && y < N && board[x][y] == -1);
}

// Count possible moves from a position
int get_degree(int x, int y, int board[N][N]) {
    int count = 0;
    for (int i = 0; i < 8; i++) {
        int new_x = x + knight_moves[i][0];
        int new_y = y + knight_moves[i][1];
        if (is_valid(new_x, new_y, board)) {
            count++;
        }
    }
    return count;
}

// Helper to sort moves by minimum degree heuristic
void sort_moves(int x, int y, int board[N][N], int moves[8][2]) {
    int degree[8];
    for (int i = 0; i < 8; i++) {
        degree[i] = get_degree(x + knight_moves[i][0], y + knight_moves[i][1], board);
    }
    for (int i = 0; i < 8 - 1; i++) {
        for (int j = 0; j < 8 - i - 1; j++) {
            if (degree[j] > degree[j + 1]) {
                int temp_deg = degree[j];
                int temp_move[2] = { moves[j][0], moves[j][1] };

                degree[j] = degree[j + 1];
                moves[j][0] = moves[j + 1][0];
                moves[j][1] = moves[j + 1][1];

                degree[j + 1] = temp_deg;
                moves[j + 1][0] = temp_move[0];
                moves[j + 1][1] = temp_move[1];
            }
        }
    }
}

// Backtracking algorithm to find the Knight's Tour
bool knight_tour(int x, int y, int movei, int board[N][N]) {
    if (movei == N * N) {
        return true;
    }

    int next_x, next_y;
    int moves[8][2];
    for (int i = 0; i < 8; i++) {
        moves[i][0] = knight_moves[i][0];
        moves[i][1] = knight_moves[i][1];
    }
    sort_moves(x, y, board, moves);

    for (int i = 0; i < 8; i++) {
        next_x = x + moves[i][0];
        next_y = y + moves[i][1];
        if (is_valid(next_x, next_y, board)) {
            board[next_x][next_y] = movei;
            if (knight_tour(next_x, next_y, movei + 1, board)) {
                return true;
            } else {
                board[next_x][next_y] = -1; // Backtracking
            }
        }
    }
    return false;
}

// Convert board indices to chess notation
void to_chess_notation(int x, int y, char* notation) {
    notation[0] = 'A' + y;
    notation[1] = '1' + (7 - x);
    notation[2] = '\0';
}

// set all the squares to -1 (start position)
void reset_board(int board[N][N]){
    for (int i = 0; i < N; i++)
        for (int j = 0; j < N; j++)
            board[i][j] = -1;
}

// Start the tour and print the result
bool solve_knight_tour(int start_x, int start_y) {
    int board[N][N];
    reset_board(board);

    char chess_notation[3];

    board[start_x][start_y] = 0; // Start at the initial position

    if (knight_tour(start_x, start_y, 1, board)) {
        int sequence[N * N][2];

        for (int move = 0; move < N * N; move++) {
            for (int x = 0; x < N; x++) {
                for (int y = 0; y < N; y++) {
                    if (board[x][y] == move) {
                        sequence[move][0] = x;
                        sequence[move][1] = y;
                    }
                }
            }
        }

        for (int move = 0; move < N * N; move++) {
            to_chess_notation(sequence[move][0], sequence[move][1], chess_notation);
            printf("Move %d: %s\n", move, chess_notation);
        }
        return true;
    } else {
        printf("No solution exists from starting position (%d, %d)\n", start_x, start_y);
        return false;
    }
}

// Test to try the tour from every square on the board
// algorithm needs to work guaranteed for every position (it does)
void test_knight_tour(void) {
    int successful_tours = 0;
    char chess_notation[3];
    for (int x = 0; x < N; x++) {
        for (int y = 0; y < N; y++) {
            to_chess_notation(x, y, chess_notation);
            printf("Starting at %s:\n", chess_notation);
            if (solve_knight_tour(x, y)) {
                successful_tours++;
            }
            printf("\n");
        }
    }
    printf("Total successful tours: %d out of %d\n", successful_tours, N * N);
}

int main(int argc, char **argv) {
    if (argc!=3) {
        printf("Usage %s start_x start_y\nf.i. %s 0 0 to start at A8\n\n or %s >7 >7 to test\n",argv[0], argv[0], argv[0]);
        exit(1);
    }

    // test the algorithm for all 64 positions
    if (atoi(argv[1]) > 7 || atoi(argv[2]) > 7 ) {
        test_knight_tour();
        exit(1);
    }

    solve_knight_tour(atoi(argv[1]), atoi(argv[2]));
    return 0;
}
