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
  std::map<std::string, std::string (*)()> functions;

  //fill the map (you'd do that in your own function normally or based on configuration file for example)
  functions["one"] = &one;
  functions["1"] = &one;
  functions["two"] = &two;
  functions["2"] = &two;
  functions["three"] = &three;
  functions["3"] = &three;
 
  std::cout << "enter: one, two or three" << std::endl;
  std::cin >> choice;

  std::cout << ( (functions.find(choice)!=functions.end()) ? functions[choice]() : "Whoops wrong input motherfucker" );  

  //or when you are not a fan of ternary statement, you can do with a single readable if else
  if (functions.find(choice)!=functions.end())
    std::cout << functions[choice]() << "\n\n";
  else
    std::cout << "Whoops, wrong input motherfucker!";
}
