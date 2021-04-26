
/*
 * Don't be that fucking guy that does: ( i % 2 == 0) ? bgBlack: bgWhite
 * modulos are expensive as they rely on division, and processors can't divide.
 * a bit wise and (or or) on bit 0 will be faster.
 *
 * compile: g++ fast_odd_even.cpp -o mod -std=c++17
 */

#include <algorithm>
#include <chrono>
#include <iostream>
#include <numeric>
#include <vector>

using namespace std;
using namespace std::chrono;

void moduloMethod(const int ITS) {
  // Starting time for the clock
  int temp = 0;

  for (int i = 0; i < ITS; ++i) {
    if (i % 2 == 0) {
      temp++;
    }
  }
  std::cout << temp << endl;
}

void andMethod(const int ITS) {
  int temp = 0;

  for (int i = 0; i < ITS; ++i) {
    if ((i & 1) == 0) {
      temp++;
    }
  }
  std::cout << temp << endl;
}

int main() {
  const int ITS = 10000000;

  vector<long> aC;
  vector<long> bC;

  std::chrono::time_point<std::chrono::system_clock> startTime;
  std::chrono::time_point<std::chrono::system_clock> endTime;

  for (int i = 0; i < 100; i++) {
    startTime = high_resolution_clock::now();
    andMethod(ITS);
    endTime = high_resolution_clock::now();
    aC.push_back(duration_cast<microseconds>(endTime - startTime).count());

    startTime = high_resolution_clock::now();
    moduloMethod(ITS);
    endTime = high_resolution_clock::now();
    bC.push_back(duration_cast<microseconds>(endTime - startTime).count());
  }

  float aCA = accumulate(aC.begin(), aC.end(), 0.0) / aC.size();

  float bCA = accumulate(bC.begin(), bC.end(), 0.0) / bC.size();

  cout << "and    took on average: " << aCA << endl;
  cout << "modulo took on average: " << bCA << endl;

  return 0;
}
