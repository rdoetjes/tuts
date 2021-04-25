
/*
 * Don't be that fucking guy that does: ( i % 2 == 0) ? bgBlack: bgWhite
 * modulos are expensive as they rely on division, and processors can't divide.
 * a bit wise and (or or) on bit 0 will be faster.
 *
 * compile: g++ fast_odd_even.cpp -o mod -std=c++17
 */

#include <iostream>
#include <chrono>

using namespace std;
using namespace std::chrono;

const int ITS = 10000000;

unsigned long slow(){
  // Starting time for the clock
  auto start = high_resolution_clock::now();

  for (int i=0; i<=ITS; ++i)
    if (i % 2 == 0) { }

  auto stop = high_resolution_clock::now();
  return duration_cast<microseconds>(stop - start).count();
}

unsigned long fast(){
  // Starting time for the clock
  auto start = high_resolution_clock::now();

  for (int i=0; i<=ITS; ++i)
    if ( !(i & 1) ) { }

  auto stop = high_resolution_clock::now();
  return duration_cast<microseconds>(stop - start).count();
}

int main()
{
  auto a = slow();
  auto b = fast();

  cout << "modulo took " << a << " uSec " << endl;
  cout << "and took " << b << " uSec " << endl;
  return 0;
}
