#include <iostream>
#include <map>


std::string one(){  return "I am one "; }
std::string two() { return "Now this was option 2"; }
std::string three() { return "You chose option 3"; }

/*
* example of using function pointers in arrays and maps (hashes) to clean up code and avoid
* large switch or if blocks.
*/

int main(){
  std::string choice="";

  //declare a map (hash) of 3 function pointers that take int as argument. C++ as std::function which would be a cleaner more c++ approach for prototyping def.
  std::map<std::string, std::string (*)()> functions = { {"one", &one}, {"1", &one}, {"two", &two}, {"2", &two}, {"three", &three}, {"3", &three} };

  std::cout << "enter: one, two or three" << std::endl;
  std::cin >> choice;

  std::cout << ( (functions.find(choice)!=functions.end()) ? functions[choice]() : "Whoops wrong input motherfucker" );  
}
