#!/usr/bin/env python3
import sys

def to_roman(year: int) -> str:
    from collections import namedtuple

    Numeral = namedtuple('Numeral', ['roman', 'dec'])

    numerals = [
        Numeral("M", 1000),
        Numeral("CM", 900),
        Numeral("D", 500),
        Numeral("CD", 400),
        Numeral("C", 100),
        Numeral("XC", 90),
        Numeral("L", 50),
        Numeral("XL", 40),
        Numeral("X", 10),
        Numeral("IX", 9),
        Numeral("V", 5),
        Numeral("IV", 4),
        Numeral("I", 1),
    ]

    remainder = year
    result = ""
    i = 0

    while remainder > 0:
        if remainder < numerals[i].dec:
            i += 1
        else:
            result += numerals[i].roman
            remainder -= numerals[i].dec

    return result

def write_error_and_exit(msg: str):
    sys.stderr.write(msg);
    sys.exit(1)

def read_stdin(size: int) -> str:
    data = sys.stdin.read(size)
    if '\n' in data:
        data = data[:data.index('\n')]
    return data.strip()

def main():
    MAX_LENGTH = 6
    input_str = read_stdin(MAX_LENGTH)

    if len(input_str) >= MAX_LENGTH:
        write_error_and_exit("Input string too long\n")

    year = 0
    try:
        year = int(input_str)
    except ValueError:
        write_error_and_exit("Invalid input\n")

    print(to_roman(year))

if __name__ == "__main__":
    main()
