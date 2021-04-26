
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


unsigned long long moduloMethod(int ITS){
  // Starting time for the clock
  auto startTime = high_resolution_clock::now();
  unsigned int temp = 0;

  for (int i=0; i<=ITS; ++i)
    if (i % 2 == 0) { temp+=i; }
	std::cout << temp << endl;
	auto endTime = high_resolution_clock::now();

	return duration_cast<microseconds>(endTime - startTime).count();
}

unsigned long long andMethod(int ITS){
  // Starting time for the clock
  auto startTime = high_resolution_clock::now();
  unsigned int temp = 0;

  for (int i=0; i<=ITS; ++i)
    if ( !( i & 1 ) ) { temp+=i; }
	std::cout << temp << endl;

	auto endTime = high_resolution_clock::now();
	return duration_cast<microseconds>(endTime - startTime).count();
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

	float bCA = accumulate( bC.begin(), bC.end(), 0.0) / bC.size();
 
  cout << "and    took on average: " << aCA << endl;
  cout << "modulo took on average: " << bCA << endl;

  return 0;
}
