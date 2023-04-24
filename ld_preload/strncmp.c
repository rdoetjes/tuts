#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>

extern int strncmp(const char *s1, const char *s2, size_t n){
   printf("s1: %s s2: %s len %d\n", s1, s2, n);
   unsetenv("LD_PRELOAD");
   return 0;
}
