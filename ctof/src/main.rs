use std::io;
use std::io::prelude::*;

fn fahrenheit_to_celsius(fahrenheit: f64) -> f64 {
    return (fahrenheit - 32.0) * 5.0 / 9.0;
}

fn read_record(b: String) -> String {
    let list = b.split(" ").collect::<Vec<&str>>();
    if list.len() == 2 && list[1] == "F" {
        let fahrenheit = list[0].parse::<f64>();
        match fahrenheit {
            Ok(value) => {
                return format!("{:.1} C", fahrenheit_to_celsius(value));
            },
            Err(err) => {
                return format!("Error: {}", err);
            }
        }
    } else {
        return b;
    }
}

fn main() {
    for line in io::stdin().lock().lines() {
        match line {
            Ok(line) => {
                print!("{}\n", read_record(line));
            },
            Err(err) => {
                println!("Error: {}", err);
            }
        }
    }
}
