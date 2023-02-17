fn print_grid(buttons: &[u8]){
    let n = buttons.len() + 1;
    for i in 1..n{
        print!("{}", buttons[i-1]);
        if i % 10 == 0 { println!("") }
    }
    println!();
}

fn switch_logic(buttons: &mut [u8]){
    let n = buttons.len() + 1;
    for person in 1..n{
        println!("person: {}", person);
        for button in 1..n{
            if button % person == 0 { buttons[button-1] ^= 1 }
        }
        print_grid(buttons);
    }
}

fn main() {
    const N: usize = 100;
    let mut buttons: [u8;N] = [0;N];
    switch_logic(&mut buttons);
}