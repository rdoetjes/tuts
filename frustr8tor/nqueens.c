#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <string.h>


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
A prepopuated "puzzle" Queen is index 2 and is shown as a RED Q
A "player" queen that the player should put down is index 1 and is shown as a YELLOW Q 
An empty field by a .
*/
void show_board(const int board[N][N]) {
    printf("Board with solution\n");
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

/*
Usually with an NQUEENS problem you only have to check the diagonals
up and down left of you, since you will put down the pieces
in order for each column.
In this case since we already have puzzle pieces populating the board,
we need to check all directions.

TODO: find an algorithm that determines the start and end vectors of
an arbitrary vector it's diagonals. That way I just need two loops
instead of 4 to check all the directions.

I initially started with just an array of 8 queens and I could quickly
iterate over each of the queens in relation to each other and calculate
a valid move by just comparing the two pieces' location.
But backtacking and solving became a hassle so I reverted to the old way
of using a board.
*/
bool is_legal_move(const int board[N][N], int row, int col){
    int i,j;

    //Check horizontal
    for (i=0; i < 8; i++){
        if (board[row][i]>0)
            return false;
    }

    //Check vertical
    for (i=0; i < 8; i++){
        if (board[i][col]>0)
            return false;
    }

    //TODO: FIND OUT ABSOLUTE START AND END POINT SO WE ONLY HAVE TWO DIAGONALS
    //diagonal down left
    for (i = row, j = col; i < N && j >= 0; i++, j--){
        if(board[i][j]>0) return false;
    }
 
    //diagonal up left
    for (i = row, j = col; i >=0 && j >= 0; i--, j--){
        if(board[i][j]>0) return false;
    }
 
    //diagonal up right
    for (i = row, j = col; i >=0 && j < N; i--, j++){
        if(board[i][j]>0) return false;
    }
 
    //diagonal down right
    for (i = row, j = col; i < N && j < N; i++, j++){
        if(board[i][j]>0) return false;
    }
    
    return true;
}

/*
Check to see if there's already a puzzle piece in the column
This way we will find an unpopulated column to put our piece into
*/
bool is_piece_in_col(const int board[N][N], int col){
    for (int i=0; i<N; i++)
        if (board[i][col]>0) return true;
    return false;
}

/*
Standard nqueens backtrack algorithm.
With the addition of checking to see if there's already a puzzle piece
in the column, if so we skip to the next column
*/
bool solve(int board[N][N], int col)
{
    if (col >= N)        
        return true;

    if (is_piece_in_col(board, col)) {
         if (solve(board, col + 1) )
            return true;
    };

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

/*
Very crude file reading for the puzzle pieces setup
Nothing robust or clever about it. 
But good enough ;)
*/
void parse_file(const char *filename, int board[N][N]){
    FILE *ptr;
    ptr = fopen(filename, "r");
    size_t l = 255;
    char *line = calloc(l, 1);

    if (NULL == ptr) {
        printf("file can't be opened \n");
        exit(1);
    }

    printf("Putting down the puzzle pieces\n");
    while (!feof(ptr)) {
        getline(&line, &l, ptr);
        printf("%s", line);
        if (strlen(line) == 4){
            int x = atoi(&line[2]) - 1;
            int y = atoi(&line[0]) - 1;
            if (x>=0 && x<=7 && y>=0 && y<=7)
                board[x][y] = 2;
        }
    }

    fclose(ptr);
    free(line);
}

void print_solution(const int board[N][N]){
    printf("\nPrinting solution vectors\n");
    for(int i=0; i<N; i++){
        for(int j=0; j<N; j++){
            if (board[i][j] == 1) 
                printf("%d,%d\n", j+1, i+1);
        }
    }
}

int main(int argc, char **argv){
    int board[N][N];
    clear_board(board);

    if (argc == 2)
        parse_file(argv[1], board);

    solve(board, 0);

    show_board(board);

    print_solution(board);
}
