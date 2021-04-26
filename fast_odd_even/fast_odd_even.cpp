
/*
 * Don't be that fucking guy that does: ( i % 2 == 0) ? bgBlack: bgWhite
 * modulos are expensive as they rely on division, and processors can't divide.
 * a bit wise and (or or) on bit 0 will be faster.
 *
 * compile: g++ fast_odd_even.cpp -o mod -std=c++17
 */

#include <iostream>
#include <chrono>
#include <vector>
#include <numeric>
#include <algorithm>

using namespace std;
using namespace std::chrono;


unsigned long moduloMethod(int ITS){
  // Starting time for the clock
  auto start = high_resolution_clock::now();
  unsigned int temp = 0;

  for (int i=0; i<=ITS; ++i)
    if (i % 2 == 0) { temp+=i; }
	std::cout << temp << endl;
  auto stop = high_resolution_clock::now();

  return duration_cast<microseconds>(stop - start).count();
}

unsigned long andMethod(int ITS){
  // Starting time for the clock
  auto start = high_resolution_clock::now();
  unsigned int temp = 0;

  for (int i=0; i<=ITS; ++i)
    if ( !(i & 1) ) { temp+=i; }
	std::cout << temp << endl;

  auto stop = high_resolution_clock::now();
  return duration_cast<microseconds>(stop - start).count();
}

int main()
{
	const int ITS = 10000000;

	vector<long> aC;
	vector<long> bC;

	for(int i=0; i<100; i++){
		auto a = andMethod(ITS);
		aC.push_back(a);

		auto b = moduloMethod(ITS);
		bC.push_back(b);
	}

	float aCA = accumulate( aC.begin(), aC.end(), 0.0) / aC.size();
	float maxA = *max_element(aC.begin(), aC.end());
	float minA = *min_element(aC.begin(), aC.end());

	float bCA = accumulate( bC.begin(), bC.end(), 0.0) / bC.size();
	float maxB = *max_element(bC.begin(), bC.end());
	float minB = *min_element(bC.begin(), bC.end());
 
  cout << "and took on average: " << aCA << " uSec max: " << maxA << " uSec min: " << minA << " uSec" << endl;
  cout << "modulo took on average: " << bCA << " uSec max: " << maxB << " uSec min: " << minB << " uSec" << endl;

  return 0;
}
