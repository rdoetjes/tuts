#include <stdlib.h>
#include <stdio.h>

void print_grid(char *switches, size_t size){
    size_t n = size + 1;
    for(int i=1; i<n; i++){
        printf("%d", switches[i-1]);
        if (i % 10 == 0) printf("\n");
    }
    printf("\n");
}

void switch_logic(char *switches, size_t size){
    size_t n = size + 1;
    for(int person=1; person<n; person++){
        printf("person: %d\n", person);
        for (int button=1; button<n; button++){
            if (button % person == 0) switches[button-1] ^= 1;
        }
        print_grid(switches, size);
    }
}

int main(){
    u_int8_t n = 100; 
    char *switches = malloc(n);
    switch_logic(switches, n);
    free(switches);
    return 0;
}
