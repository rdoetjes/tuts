#include <unistd.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>

char* to_roman(const uint32_t year){
    typedef struct numeral{
        char* roman;
        int dec;
    } numeral;

    static numeral numerals[] = {
        { "M", 1000 },
        { "CM", 900 },
        { "D", 500 },
        { "CD", 400 },
        { "C", 100 },
        { "XC", 90 },
        { "L", 50 },
        { "XL", 40 },
        { "X", 10 },
        { "IX", 9 },
        { "V", 5 },
        { "IV", 4 },
        { "I", 1 },
    };

    int remainder = year;
    char *result = (char *)calloc(1, 1);
    size_t i = 0;
    while(remainder>0){
        if (remainder < numerals[i].dec){
            i++;
        }else{
            int grow_to_nr_bytes = strlen(result) + strlen(numerals[i].roman);
            result = realloc(result, grow_to_nr_bytes);
            if (result == NULL) {
                fprintf(stderr, "Memory allocation failed\n");
                free(result);
                return NULL;
            }
            strcat(result, numerals[i].roman);
            remainder -= numerals[i].dec;
        }
    }
    return result;
}

int read_stdin(char *buffer, size_t size)
{
    size_t cnt = 0;
    char c;

    if(buffer == NULL || size <= 0)
        return 0;

    while(read(STDIN_FILENO, &c, 1) == 1 && cnt <= size){
        if(c == '\n') {
            buffer[cnt] = 0;
            return 1;
        }
        buffer[cnt++] = c;
    }
    buffer[cnt] = 0;
    return 0;
}

#define MAX_LENGTH 6

void write_error_and_exit(const char* error){
    write(STDERR_FILENO, error, strlen(error));
    exit(1);
}

int main(){
    char input[MAX_LENGTH]="";
    if (read_stdin((char *)&input, MAX_LENGTH-1) == 0) write_error_and_exit("ABORTED! Input to long\n");

    char *roman = to_roman(atoi(input));
    printf("%s\n", roman);
    free(roman);
    return 0;
}
