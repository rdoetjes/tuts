#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#define NUM_CHOICES 6
#define X_NUMBER 4

int numbers[] = {1, 2, 3, 4, 5, 6};
int available[NUM_CHOICES] = {X_NUMBER, X_NUMBER, X_NUMBER, X_NUMBER, X_NUMBER, X_NUMBER}; // 4 of each number
int total = 0;

// Function to display available numbers
void show_available()
{
  printf("\033[2J\033[H");
  printf("Welcome to the 31 Game! Inspried by Robin from 8Bit Show And Tell\nFirst to reach exactly 31 wins.\n\n");
  printf("Available numbers:\n");
  for (int j = X_NUMBER; j != 0; j--)
  {
    for (int i = 0; i < NUM_CHOICES; i++){
      available[i] >= j ? printf("%d ", i+1) : printf("  ");
    }
    printf("\n");
  }
  printf("\n\n");
}

// Check if a number is available
int is_available(int n)
{
  if (numbers[n-1] == n && available[n-1] > 0) return 1;
  return 0;
}

// Decrease availability
void use_number(int n)
{
    available[n-1]--;
    total += n;
}

// Evaluate recursively: returns 1 if current player can force a win
int can_force_win(int current_total, int *available_pool, int is_ai_turn)
{
  if (current_total == 31)
  {
    return !is_ai_turn; // Whoever just moved wins
  }

  for (int i = 0; i < NUM_CHOICES; i++)
  {
    if (available_pool[i] > 0 && current_total + numbers[i] <= 31)
    {
      available_pool[i]--;
      int win = can_force_win(current_total + numbers[i], available_pool,
                              !is_ai_turn);
      available_pool[i]++;

      if (is_ai_turn && win)
        return 1; // AI has a winning path
      if (!is_ai_turn && !win)
        return 0; // Player has a forced win, bad for AI
    }
  }

  return is_ai_turn ? 0 : 1; // No moves left: bad for current player
}

int computer_move()
{
  int best_move = -1;

  for (int i = 0; i < NUM_CHOICES; i++)
  {
    if (available[i] > 0 && total + numbers[i] <= 31)
    {
      available[i]--;
      int win =
          can_force_win(total + numbers[i], available, 0); // Player's turn next
      available[i]++; //backtrack

      if (win)
      {
        best_move = numbers[i];
        break; // Take the first move that guarantees a win
      }
    }
  }

    // If no forced-win path, pick a random valid move
    if (best_move == -1)
    {
      int valid_moves[NUM_CHOICES];
      int valid_count = 0;
      
      // Collect all valid moves
      for (int i = 0; i < NUM_CHOICES; i++)
      {
        if (available[i] > 0 && total + numbers[i] < 31)
        {
          valid_moves[valid_count] = numbers[i];
          valid_count++;
        }
      }
      
      // Pick a random move from valid moves
      if (valid_count > 0)
      {
        best_move = valid_moves[rand() % valid_count];
      }
    }  

  return best_move;
}

// New function to handle player's turn
int player_turn_handler() 
{
  int move;
  printf("Your move: ");
  scanf("%d", &move);

  if (!is_available(move)) {
    printf("Invalid move. Try again.\n");
    return 0; // Invalid move
  }

  use_number(move);
  printf("You played: %d → Total: %d\n", move, total);
  
  return 1; // Valid move
}

// New function to handle computer's turn
int computer_turn_handler() 
{
  int comp = computer_move();
  if (comp == -1) {
    printf("You win!\n");
    return 0; // Computer has no valid moves
  }
  
  use_number(comp);
  printf("Computer plays: %d → Total: %d\n", comp, total);
  
  return 1; // Valid move
}

// New function to check if game is over
int check_game_over(int is_player_turn) 
{
  if (total == 31) {
    if (is_player_turn) {
      printf("You win!\n");
    } else {
      printf("Computer wins!\n");
    }
    return 1; // Game over
  }
  return 0; // Game continues
}

int main()
{
  srand(time(NULL));
  int player_turn = rand() & 1;

  while (total < 31)
  {
    show_available();
    printf("Current total: %d\n", total);

    int valid_move;
    if (player_turn) {
      valid_move = player_turn_handler();
      if (valid_move==0) continue; //invlaide move
      
      if (check_game_over(1)) break;
    }
    else {
      valid_move = computer_turn_handler();
      if (valid_move==0) break;
      
      if (check_game_over(0)) break;
    }
    
    player_turn = !player_turn;
  }

  if (total < 31) {
    printf("I have no valid moves left.\n");
  }
  return 0;
}
