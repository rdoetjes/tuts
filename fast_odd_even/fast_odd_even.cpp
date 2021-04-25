#include <iostream>

/*
 * Don't be that fucking guy that does: ( i % 2 == 0) ? bgBlack: bgWhite
 * modulos are expensive as they rely on division, and processors can't divide.
 * a bit wise and (or or) on bit 0 will be faster.
 */
int main()
{
  for (int i=0; i<=1000000; ++i){
    if ( (i & 1) )
      std::cout << i << " odd" << std::endl;
    else
      std::cout << i << " even" << std::endl;
  }
  return 0;
}
