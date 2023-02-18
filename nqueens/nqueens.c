#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>

static const size_t N = 8;

void clear_board(int board[N][N]) {
    for(int row=0; row < N; row++){
        for(int col=0; col < N; col++){
            board[row][col] = 0;
        }    
    }
}

void show_board(int board[N][N]) {
    for(int row=0; row < N; row++){
        for(int col=0; col < N; col++){
            if (board[row][col] == 0) 
                printf(".");
            else
                printf("Q");
        }
        printf("\n");    
    }
}

void solve(int board[N][N]){

}

bool is_col_free(int board[N][N], int x){
    for(int y=0; y < N; y++){
        if(board[y][x] != 0) 
            return false;
    }
    return true;
}

bool is_row_free(int board[N][N], int y){
    for(int x=0; x < N; x++){
        if(board[y][x] != 0) 
            return false;
    }
    return true;
}

void set_queens(int board[N][N], int x, int y){
    if (x==N || y==N)
        return;
    if(!is_row_free(board, y))
        set_queens(board, x, y+1);
    else if(!is_col_free(board,x))
        set_queens(board, x+1, y);
    else
        board[y][x]=1;
}

int main(){
    int board[N][N];
    clear_board(board);

    board[0][0] = 2;
    board[2][1] = 2;
    board[7][6] = 2;
    board[4][3] = 2;


    for(int x=0; x < N; x++)
        set_queens(board, x, 0);
    
    solve(board);
    show_board(board);
}
