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

bool legal_move(QUEEN queens[N]){
    for (int q=0; q<N; q++){
        if (queens[q].t == -1) continue;

        for(int qc=0; qc<N; qc++){
            if (qc==q || queens[qc].t<0 || queens[qc].x<0 || queens[qc].y<0) continue;

            if(queens[q].x == queens[qc].x) return false;

            if(queens[q].y == queens[qc].y) return false;
            
            if( abs(queens[q].y - queens[qc].y) == abs(queens[q].x - queens[qc].x) ) return false;
        }
    }
    return true;
}

int get_next_new_queen_idx(QUEEN queens[N]){
    for(int i=0; i < N; i++)
        if (queens[i].t == -1) return i;
    return 0;
}

int find_empty_col(QUEEN queens[N]){
    int used[N];
    for(int i=0; i<N; i++)
        used[i] = 0;

    for(int i=0; i<N; i++)
        used[queens[i].x ]=1;

    for(int i=0; i<N; i++)
        if ( used[i] == 0) return i;

    return -1;
}

void put_piece(QUEEN queens[N], int queen_idx, int x, int y){
    if (queen_idx>=N)
        return;
    
     if (y>=N){
        queen_idx-=2;
        queens[queen_idx].y ++;
        return;
     }

    queens[queen_idx].x = x;
    queens[queen_idx].y = y;
    queens[queen_idx].t = 1;

    if (legal_move(queens)){
        queen_idx++;
        put_piece(queens, queen_idx, find_empty_col(queens), 0);
    }
    else{
        put_piece(queens, queen_idx, x, y+1);
    }
}

void solve(QUEEN queens[N], int start_idx){
    put_piece(queens, start_idx, find_empty_col(queens), 0);
}

void reflect_pieces_to_board(int board[N][N], QUEEN queens[N]){
    for(int i=0; i < N; i++) 
        board[queens[i].y][queens[i].x] = queens[i].t;
}

void init_queens(QUEEN queens[N]){
    for(int x=0; x < N; x++){
        queens[x].x = -1;
        queens[x].y = -1;
        queens[x].t = -1;
    }
}

int main(){
    int board[N][N];
    QUEEN queens[N];

    init_queens(queens);
    clear_board(board);

    //set the puzzle pieces (2 means puzzle piece that may not be moved)
    queens[0].x = 0; queens[0].y = 3; queens[0].t = 2;
    queens[1].x = 2; queens[1].y = 7; queens[1].t = 2;

    solve(queens, get_next_new_queen_idx(queens));

    reflect_pieces_to_board(board, queens);

    show_board(board);
}
