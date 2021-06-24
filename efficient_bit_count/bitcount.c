#include <stdio.h>

static int bitcount(int n) {
  int c = 0;
  while (n != 0) {
    c++;
    n &= n - 1;
  }
  return c;
}

int main() { printf("%d bits or set\n", bitcount(254)); }
