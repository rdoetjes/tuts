//This is what we will be coding in assembler

#include <stdio.h>
#include <unistd.h>

int main(){
  write(1, "Hello World, ARM cpus rule!\n", 28);
  return 0;
}
