fn print_grid(switches: &[u8]){
    let n = switches.len() + 1;
    for i in 1..n{
        print!("{}", switches[i-1]);
        if i % 10 == 0 { println!("") }
    }
    println!();
}

fn switch_logic(switches: &mut [u8]){
    let n = switches.len() + 1;
    for person in 1..n{
        println!("person: {}", person);
        for button in 1..n{
            if button % person == 0 { switches[button-1] ^= 1 }
        }
        print_grid(switches);
    }
}

fn main() {
    const N: usize = 100;
    let mut switches: [u8;N] = [0;N];
    switch_logic(&mut switches);
}
