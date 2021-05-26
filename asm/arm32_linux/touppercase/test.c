//This is what we will be coding in assembler

#include <stdio.h>
#include <unistd.h>

int main(){
  char *string = "Hello World, ARM cpus rule!\n";
  int len = 28;
  char uppercase[len];
  
  for(int i=0; i<len; i++){
    if (string[i] >= 'a' && string[i]<='z')
      uppercase[i]=string[i]-32;
    else
      uppercase[i]=string[i];
  }
 
  for(int i=0; i<0xFFFF; i++){
    write(1, uppercase, len);
  }
  return 0;
}
