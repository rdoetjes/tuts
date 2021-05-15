const readline = require('readline-sync');

const one = () => "This is option 1";
const two = () => "Now this was option 2";
const three = () => "You chose option 3";

const choice =  readline.question('Enter: one, two or three: \n');

switch(choice){
    case "1":
        console.log(one());
        break;
    case "one":
        console.log(one());    
        break;
    case "2":
        console.log(two());
        break;
    case "two":
        console.log(two());
        break;
    case "3":
        console.log(three());
        break;
    case "three":
        console.log(three());
        break;
    default:
        console.log("whoops motherfucker!");
}

//this is shorter and cleaner
functions = { "1": one, "one": one, "2": two, "two": two, "3": three, "three": three}
console.log( (choice in functions) ? functions[choice]() : "whoops motherfucker!");