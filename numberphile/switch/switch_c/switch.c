#include <stdlib.h>
#include <stdio.h>

void print_grid(size_t person, char *buttons, size_t size){
    size_t n = size + 1;
    printf("person: %ld\n", person);
    for(int i=1; i<n; i++){
        printf("%d", buttons[i-1]);
        if (i % 10 == 0) printf("\n");
    }
    printf("\n");
}

void switch_logic(char *buttons, size_t size){
    size_t n = size + 1;
    for(int person=1; person<n; person++){
        for (int button=1; button<n; button++){
            if (button % person == 0) buttons[button-1] ^= 1;
        }
        print_grid(person, buttons, size);
    }
}

int main(){
    size_t n = 100; 
    
    char *buttons = calloc(n,1);
    switch_logic(buttons, n);
    
    free(buttons);
    return 0;
}
