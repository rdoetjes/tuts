#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>

#define NUM_CHOICES 6
#define X_NUMBER 4

int numbers[] = {1, 2, 3, 4, 5, 6};
int available[NUM_CHOICES] = {X_NUMBER, X_NUMBER, X_NUMBER, X_NUMBER, X_NUMBER, X_NUMBER}; // 4 of each number
int total = 0;

// Function to display available numbers
void show_available()
{
  printf("\033[2J\033[H");
  printf("Available numbers:\n");
  for (int j = X_NUMBER; j != 0; j--)
  {
    printf("%s %s %s %s %s %s\n", 
      available[0] >= j ? "1" : " ", 
      available[1] >= j ? "2" : " ", 
      available[2] >= j ? "3" : " ", 
      available[3] >= j ? "4" : " ", 
      available[4] >= j ? "5" : " ", 
      available[5] >= j ? "6" : " ");
  }
  printf("\n\n");
}

// Check if a number is available
int is_available(int n)
{
  for (int i = 0; i < NUM_CHOICES; i++)
  {
    if (numbers[i] == n && available[i] > 0)
      return 1;
  }
  return 0;
}

// Decrease availability
void use_number(int n)
{
  for (int i = 0; i < NUM_CHOICES; i++)
  {
    if (numbers[i] == n)
    {
      available[i]--;
      total += n;
      break;
    }
  }
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
      available[i]++;

      if (win)
      {
        best_move = numbers[i];
        break; // Take the first move that guarantees a win
      }
    }
  }

  // If no forced-win path, pick first valid move
  if (best_move == -1)
  {
    for (int i = 0; i < NUM_CHOICES; i++)
    {
      if (available[i] > 0 && total + numbers[i] <= 31)
      {
        best_move = numbers[i];
        break;
      }
    }
  }

  return best_move;
}

int main()
{
  srand(time(NULL));
  int player_turn = rand() & 1;
  int move;

  printf("Welcome to the 31 Game!\nFirst to reach exactly 31 wins.\n");

  while (total < 31)
  {
    show_available();
    printf("Current total: %d\n", total);

    if (player_turn)
    {
      printf("Your move: ");
      scanf("%d", &move);

      if (!is_available(move))
      {
        printf("Invalid move. Try again.\n");
        continue;
      }

      use_number(move);
      printf("You played: %d → Total: %d\n", move, total);

      if (total == 31)
      {
        printf("You win!\n");
        break;
      }
    }
    else
    {
      int comp = computer_move();
      if (comp == -1)
      {
        printf("You win!\n");
        break;
      }
      use_number(comp);
      printf("Computer plays: %d → Total: %d\n", comp, total);

      if (total == 31)
      {
        printf("Computer wins!\n");
        break;
      }
    }
    player_turn = !player_turn;
  }

  if (total < 31)
  {
    printf("I have no valid moves left.\n");
  }

  return 0;
}
