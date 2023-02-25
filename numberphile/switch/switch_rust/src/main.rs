fn print_grid(person: usize, buttons: &[usize]){
    let n = buttons.len() + 1;  //derive length from array size; this will prevent us from having to also pass it
    println!("person: {}", person);
    for i in 1..n{
        print!("{}", buttons[i-1]);
        if i % 10 == 0 { println!("") }
    }
    println!();
}

fn switch_logic(person: usize, buttons: &mut [usize]){
    let n = buttons.len() + 1;  //derive length from array size; this will prevent us from having to also pass it
    for button in 1..n{
        if button % person == 0 { buttons[button-1] ^= 1 }
    }
}

fn fondle_my_knobs(buttons: &mut [usize]){
    let n = buttons.len() + 1;  //derive length from array size; this will prevent us from having to also pass it
    for person in 1..n {
        switch_logic(person, buttons);
        print_grid(person, &buttons);
    }    
}

fn main() {
    const N: usize = 100;
    let mut buttons: [usize;N] = [0;N];    
    
    fondle_my_knobs(&mut buttons);
}
