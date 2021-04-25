#include <iostream>
#include <map>

void one(int v){
  std::cout << "I am one " << v << std::endl;
}

void two(int v){
  std::cout << "I am two " << v << std::endl;
}

void three(int v){
  std::cout << "I am three " << v << std::endl;
}

/*
* example of using function pointers in arrays and maps (hashes) to clean up code and avoid
* large switch or if blocks.
*/

int main(){

  void (*p[4]) (int x);
  p[0] = &one;
  p[1] = &two;
  p[2] = &three;

  std::cout << "an example of using fuction pointers in an array, as order of execution" << std::endl;
  int ins[] = {1, 1, 2, 0, 1};
  for(int i=0; i<sizeof(ins)/sizeof(ins[0]); i++)
    (*p[ ins[i] ])(i);

  std::map<std::string, void (*)(int)> t;
  t["one"] = &one;
  t["two"] = &two;
  t["three"] = &three;
 
  std::cout << "an example of using fuction pointers with a map (hash) as a switch statement" << std::endl;
  std::string inp;
  while(1){

    std::cout << "enter a command: one, two or three" << std::endl;
    std::cin >> inp;

    if (t.find(inp)==t.end()){
      std::cout << "Not a valid command" << std::endl;
      continue;
    }

    t[inp](12);
  }
}
