#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static const char *the_key = "just_a_demo";

void main(int argc, char **argv){
	if ( strncmp(argv[1], the_key, strlen(argv[1])) == 0){
	  printf("OK!\n");
	}
	else{
	  printf("NO WAY!\n");
	}
}
