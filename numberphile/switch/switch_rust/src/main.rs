fn print_grid(person: usize, buttons: &[usize]){
    let n = buttons.len() + 1;
    println!("person: {}", person);
    for i in 1..n{
        print!("{}", buttons[i-1]);
        if i % 10 == 0 { println!("") }
    }
    println!();
}

fn switch_logic(person: usize, buttons: &mut [usize]){
    let n = buttons.len() + 1;
    for button in 1..n{
        if button % person == 0 { buttons[button-1] ^= 1 }
    }
}

fn main() {
    const N: usize = 100;
    let mut buttons: [usize;N] = [0;N];    
    
    for person in 1..N + 1 {
        switch_logic(person, &mut buttons);
        print_grid(person, &buttons);
    }    
}
