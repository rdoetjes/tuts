#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>

static const size_t N = 8;

typedef struct QUEEN {
    int x;
    int y;
    unsigned int t;
} QUEEN;

/*
Set the whole board to 0 (no queens on the board)
*/
void clear_board(int board[N][N]) {
    for(int row=0; row < N; row++){
        for(int col=0; col < N; col++){
            board[row][col] = 0;
        }    
    }
}

/*
print the board to the screen.
A queen is denoted by a Q
An empty field by a .
*/
void show_board(int board[N][N]) {
    for(int row=0; row < N; row++){
        for(int col=0; col < N; col++){
            switch (board[row][col])
            {
                case 0:
                    printf(" . ");
                    break;

                case 1:
                    printf("\033[33m Q \033[0m");
                    break;

                case 2:
                    printf("\033[31m Q \033[0m");
                    break;
                
                default:
                    break;
            }
        }
        printf("\n");    
    }
}

bool is_x_empty(int board[N][N], int y){
    for(int x=0; x < N; x++){
        if(board[y][x]>0) return false;
    }
    return true;
}

bool is_y_empty(int board[N][N], int x){
    for(int y=0; y < N; y++){
        if(board[y][x]>0) return false;
    }
    return true;
}

void solve(int board[N][N], QUEEN queens[N]){
    printf("placing queen #%d\n", queen_idx+1);
}

/*
Set up the other pieces so that they are never on the same row or column
with any of the other already placed pieces
*/
void set_other_queens_on_board(int board[N][N], QUEEN queens[N], int queen_idx, int x, int y){
    if (x>=N || y>=N || queen_idx>=N)
        return;
    
    if (!is_y_empty(board, x)) 
        set_other_queens_on_board(board, queens, queen_idx, x + 1 ,y);
    else if (!is_x_empty(board, y)) 
        set_other_queens_on_board(board, queens, queen_idx, x ,y + 1);
    else {
        board[y][x] = 1;
        queens[queen_idx].x = x; queens[queen_idx].y = y;  queens[queen_idx].t = 1;
        set_other_queens_on_board(board, queens, queen_idx + 1, x, y);
    }
}

int get_next_queen_idx(QUEEN queens[N]){
    for(int i=0; i < N; i++)
        if (queens[i].t == 0) return i;
    return 0;
}

void set_start_queens_on_board(int board[N][N], QUEEN queens[N]){
    for(int i=0; i < N; i++) 
        board[queens[i].y][queens[i].x] = queens[i].t;

    set_other_queens_on_board(board, queens, get_next_queen_idx(queens), 0 ,0);
}

void init_queens(QUEEN queens[N]){
    for(int x=0; x < N; x++){
        queens[x].x = -1;
        queens[x].y = -1;
        queens[x].t = 0;
    }
}

int main(){
    int board[N][N];
    QUEEN queens[N];

    init_queens(queens);
    clear_board(board);

    //set the puzzle pieces (2 means puzzle piece that may not be moved)
    queens[0].x = 4; queens[0].y = 4; queens[0].t = 2;
    queens[1].x = 7; queens[1].y = 0; queens[1].t = 2;

    set_start_queens_on_board(board, queens);

    solve(board, queens);

    show_board(board);
}
