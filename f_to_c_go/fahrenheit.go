package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func convert_c_to_f(c float64) float64 {
	return (c - 32) * 5.0 / 9.0
}

func conform_line(line string) (string, error) {
	r := regexp.MustCompile(`^(\d{1,}\.?\d?) F$`)
	fahrenheit_s := r.FindAllStringSubmatch(line, -1)

	if len(fahrenheit_s) == 0 {
		return line, nil
	}

	fahrenheit, err := strconv.ParseFloat(fahrenheit_s[0][1], 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.1f C", convert_c_to_f(fahrenheit)), nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "opening standard input:", err)
	}

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
			break
		}

		celcius_conformed_record, err := conform_line(scanner.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting line:", err)
			continue
		}

		fmt.Println(celcius_conformed_record)
	}
}
