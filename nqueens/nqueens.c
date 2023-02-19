#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>

static const size_t N = 8;

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

bool is_legal_move(int board[N][N], int row, int col){
    int i, j;
 
    /* Check this row on left side */
    for (i = 0; i < col; i++)
        if (board[row][i]>0)
            return false;
 
   /* Check this row on right side */
    for (i = row; i < N; i++)
        if (board[row][i]>0)
            return false;

    /* Check upper diagonal on left side */
    for (i = row, j = col; i >= 0 && j >= 0; i--, j--)
        if (board[i][j]>0)
            return false;
 
    /* Check lower diagonal on left side */
    for (i = row, j = col; j >= 0 && i < N; i++, j--)
        if (board[i][j]>0)
            return false;
 
    return true;
}

bool is_piece_in_col(int board[N][N], int col){
    for (int i=0; i<N; i++)
        if (board[i][col]>0) return true;
    return false;
}

bool solve(int board[N][N], int col)
{
    /* base case: If all queens are placed
      then return true */
    if (col >= N)
        return true;
 
    if (is_piece_in_col(board, col)) {
         if (solve(board, col + 1) )
            return true;
    };

    /* Consider this column and try placing
       this queen in all rows one by one */
    for (int i = 0; i < N; i++) {

        if (is_legal_move(board, i, col)) {
            board[i][col] = 1;

            if (solve(board, col + 1))
                return true;
            
            board[i][col] = 0; // BACKTRACK
        }   
    }
    return false;
}

int main(){
    int board[N][N];
    clear_board(board);

    //set the puzzle pieces (2 means puzzle piece that may not be moved)
    board[4][1] = 2;
    board[1][6] = 2;


    if ( solve(board, 0) ){
        show_board(board);
    }
    else{
        printf("Couldn't be solved\n");
    }
}
