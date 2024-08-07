use std::io;
use std::io::prelude::*;
use std::num::ParseFloatError;

/// Converts a Fahrenheit temperature to Celsius
///
/// # Arguments
///
/// * `fahrenheit` - The input temperature in Fahrenheit
///
/// # Returns
///
/// The temperature converted to Celsius
fn fahrenheit_to_celsius(fahrenheit: f64) -> f64 {
    (fahrenheit - 32.0) * 5.0 / 9.0
}

/// Converts a Fahrenheit temperature to Celsius
/// Takes in a record in b that looks like "36.6 F" and returns a record in c that looks like "2.6 C"
/// # Arguments
///
/// * `b` - The record that looks like "36.6 F"
///
/// # Returns
/// Record in Celsius that looks like "2.6 C"
///
/// The temperature converted to Celsius
fn read_record(b: String) -> Result<String, ParseFloatError> {
    let list = b.split(' ').collect::<Vec<&str>>();
    if list.len() < 2 || list[1] != "F" {
        return Ok(b);
    }

    let fahrenheit = list[0].parse::<f64>();
    match fahrenheit {
        Ok(value) => Ok(format!("{:.1} C", fahrenheit_to_celsius(value))),
        Err(err) => Err(err),
    }
}

fn main() {
    for line in io::stdin().lock().lines() {
        match line {
            Ok(line) => {
                println!("{}", read_record(line).expect("Error parsing records!"));
            }
            Err(err) => {
                println!("Error: {}", err);
                continue;
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_read_record_valid() {
        let input = "36.6 F".to_string();
        let expected = "2.6 C".to_string();
        let actual = read_record(input);
        assert_eq!(expected, actual.unwrap());
    }

    #[test]
    fn test_read_record_valid_c() {
        let input = "36.6 C".to_string();
        let expected = "36.6 C".to_string();
        let actual = read_record(input);
        assert_eq!(expected, actual.unwrap());
    }

    #[test]
    fn test_read_record_invalid_temp() {
        let input = "invalid F".to_string();
        let actual = read_record(input);
        assert!(actual.is_err());
    }

    #[test]
    fn test_read_record_missing_unit() {
        let input = "36.6".to_string();
        let expected = input.clone();
        let actual = read_record(input);
        assert_eq!(expected, actual.unwrap());
    }
}
