#include <string.h>
#include <stdio.h>

int main(void){
	char a[20];
	for(int i=0; i<100; i++){
		sprintf(a, "%02d", i);
		for(int j=0; j<2; j++){
			printf("D\t%c\t50\t50\n", a[j]);
		}
		printf("~\t1000\n");
	}
	return 0;
}
