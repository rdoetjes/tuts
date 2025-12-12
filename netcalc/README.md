# Teaching Bitwise Operations and Network Calculations

This Go program is designed as an educational tool to help students understand **bitwise operations** and **network calculation concepts** through practical, hands-on examples.  

The project demonstrates how bit manipulation can be used to perform common networking tasks, such as calculating subnet masks, determining network and broadcast addresses, and understanding how IP addressing works behind the scenes.

## What This Program Teaches

### 🔹 Bitwise Operations  
Students will learn how operators like **AND**, **OR**, **NOT**, and bit shifting are used to manipulate binary data. These examples show how bitwise logic forms the foundation of low-level computing and networking.

### 🔹 Network Calculations  
The program includes demonstrations of:  
- Converting IP addresses to binary  
- Applying subnet masks  
- Determining network and broadcast addresses  
- Calculating the number of hosts in a subnet  
- Understanding CIDR notation and its practical use  

By exploring these topics in Go, students gain both theoretical understanding and real coding experience.

## Why Go?

Go is well-suited for teaching these concepts because of its:  
- Clean and readable syntax  
- Powerful standard library  
- Strong emphasis on performance and low-level operations  
- Built-in support for network-related features  

## How to Use

Clone the repository and run the program:

```bash
go run main.go -cidr 192.168.1.1/24
or
go run main.go -ip 192.168.1.1 -netmask 255.255.255.0
```
(After compiling with go build you don't need the go run main.go anymore)
