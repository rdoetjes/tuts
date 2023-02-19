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

bool are_diagonals_clear(int board[N][N], QUEEN queens[N]){
    for (int q=0; q<N; q++){
        
        int x = queens[q].x;
        int y = queens[q].y;

        for(int tx=x; tx<N; tx++){
            for(int ty=y; ty<N; ty++){
                if(board[y][x] >0 ) return false;
            }
        }

        for(int tx=x; tx<=0; tx--){
            for(int ty=y; ty<N; ty--){
                if(board[y][x] >0 ) return false;
            }
        }

        for(int tx=x; tx<=0; tx--){
            for(int ty=y; ty<N; ty++){
                if(board[y][x] >0 ) return false;
            }
        }

        for(int tx=x; tx<N; tx++){
            for(int ty=y; ty>=0; ty--){
                if(board[y][x] >0 ) return false;
            }
        }
    }
    return true;
}

int get_next_movable_queen_idx(QUEEN queens[N], int queen_idx){
    if (queen_idx==N-1) queen_idx = -1;

    for(int i=queen_idx+1; i < N; i++)
        if (queens[i].t == 1) return i;

    for(int i=0; i < queen_idx; i++)
        if (queens[i].t == 1) return i;

    return -1;
}

int get_next_new_queen_idx(QUEEN queens[N]){
    for(int i=0; i < N; i++)
        if (queens[i].t == 0) return i;
    return 0;
}

void swap_queen_with_next_one_over(int board[N][N], QUEEN queens[N], int i){
    if (are_diagonals_clear(board, queens))
        return;

    int q0 = get_next_movable_queen_idx(queens, i);
    int q1 = get_next_movable_queen_idx(queens, q0);

    printf("%d %d\n", q0, q1);
    //clear the old spots of these queens from the board
    board[queens[q1].y][queens[q1].x] = 0;
    board[queens[q0].y][queens[q0].x] = 0;

    printf("q0 %d %d,%d, q1 %d %d,%d\n", q0, queens[q0].y ,queens[q0].x, q1, queens[q1].y, queens[q1].x );

    //safe copy of location for swapping
    int y = queens[q0].y;
    //set queen on position q to queen on position i
    queens[q0].y = queens[q1].y;
    //set queen on position i to queen on position q
    queens[q1].y = y; 

    //update the board if the new queen positions
    board[queens[q0].y][queens[q0].x] = 1;
    board[queens[q1].y][queens[q1].x] = 1;
    
    show_board(board);
    i = get_next_movable_queen_idx(queens,q1);
    swap_queen_with_next_one_over(board, queens, i);
}

void solve(int board[N][N], QUEEN queens[N], int start_idx){
    swap_queen_with_next_one_over(board, queens, start_idx);
}

/*
Set up the other pieces so that they are never on the same row or column
as any of the other already placed pieces
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

void set_start_queens_on_board(int board[N][N], QUEEN queens[N], int start_idx){
    for(int i=0; i < N; i++) 
        board[queens[i].y][queens[i].x] = queens[i].t;

    set_other_queens_on_board(board, queens, start_idx, 0 ,0);
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
    queens[0].x = 0; queens[0].y = 0; queens[0].t = 2;
    
    int idx_first_movable_queen = get_next_new_queen_idx(queens);
    set_start_queens_on_board(board, queens, idx_first_movable_queen);

    show_board(board);

    solve(board, queens, idx_first_movable_queen);
}
