package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func convert_c_to_f(c float64) float64 {
	return (c - 32) * 5.0 / 9.0
}

func convert_line(line string) (string, error) {
	columns := strings.Split(line, " ")
	if columns[1] != "F" {
		return line, nil
	}

	f, err := strconv.ParseFloat(columns[0], 64)
	if err != nil {
		return "", err
	}
	s_f := fmt.Sprintf("%.1f C", convert_c_to_f(f))

	return s_f, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
			break
		}

		line, err := convert_line(scanner.Text())
		if err != nil {
			println("Error converting line:", err)
			continue
		}

		println(line)
	}
}
