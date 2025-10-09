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
    char *result = (char *)calloc(1, sizeof(char));
    if (result == NULL) {
        fprintf(stderr, "Memory allocation failed\n");
        return NULL;
    }

    size_t i = 0;
    while(remainder>0){
        if (remainder < numerals[i].dec){
            i++;
        }else{
            int grow_to_nr_bytes = strlen(result) + strlen(numerals[i].roman) + 1;
            char *temp_result = realloc(result, grow_to_nr_bytes);
            if (temp_result == NULL) {
                fprintf(stderr, "Memory allocation failed\n");
                free(result);
                return NULL;
            }
            result = temp_result;
            if (result == NULL) {
                fprintf(stderr, "Memory allocation failed\n");
                free(result);
                return NULL;
            }
            strncat(result, numerals[i].roman, strlen(numerals[i].roman));
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

#define MAX_LENGTH 5

void write_error_and_exit(const char* error){
    write(STDERR_FILENO, error, strlen(error));
    exit(1);
}

int main(){
    char input[MAX_LENGTH+1] = "";
    if (read_stdin(input, MAX_LENGTH) == 0)
        write_error_and_exit("ABORTED! Input too long\n");

    char *endptr;
    long year = strtol(input, &endptr, 10);
    if (*endptr != '\0' || year <= 0 || year > 99999) {
        write_error_and_exit("Invalid input! Please enter a valid year.\n");
    }

    char *roman = to_roman((uint32_t)year);
    if (roman != NULL)
        printf("%s\n", roman);
    else
        write_error_and_exit("Conversion failed!");

    free(roman);
    return 0;
}
