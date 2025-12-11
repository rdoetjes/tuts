/*
Quick and dirty bootstrap tool for initializing the database schema and creating the first user(s) with password.
Without this tool, there is no way to be able to login and create/update/delete users via the REST API.
Usage:

	go run tools/setup.go init-db       - Initialize database schema with password column
	go run tools/setup.go create-user   - Create a new user with email and password
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
	"phonax.com/db/auth"
	"phonax.com/db/db"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	switch command {
	case "init-db":
		initDB()
	case "create-user":
		createUser()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
Usage:
  go run tools/setup.go init-db       - Initialize database schema with password column
  go run tools/setup.go create-user   - Create a new user with email and password
`)
}

// initDB adds the password column to users table if it doesn't exist
func initDB() {
	fmt.Println("Connecting to database...")
	conn := db.Connect()
	defer conn.Close()

	// Check if password column exists
	var exists bool
	err := conn.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'password'
		)
	`).Scan(&exists)

	if err != nil {
		log.Fatalf("Failed to check schema: %v", err)
	}

	if exists {
		fmt.Println("✓ Password column already exists")
		return
	}

	// Add password column
	_, err = conn.Exec(`
		ALTER TABLE users ADD COLUMN password VARCHAR(255) NOT NULL DEFAULT '';
	`)

	if err != nil {
		log.Fatalf("Failed to add password column: %v", err)
	}

	fmt.Println("✓ Password column added successfully")
}

// createUser prompts for email and password, then creates a user
func createUser() {
	fmt.Println("Connecting to database...")
	conn := db.Connect()
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if email == "" {
		log.Fatal("Email cannot be empty")
	}

	fmt.Print("Enter password: ")
	passwordBytes, _ := term.ReadPassword(int(syscall.Stdin))
	password := strings.TrimSpace(string(passwordBytes))
	fmt.Println() // newline after password input

	if password == "" {
		log.Fatal("Password cannot be empty")
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if user already exists
	var userID int
	err = conn.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err == nil {
		// User exists, update password
		_, err = conn.Exec(
			"UPDATE users SET password = $1 WHERE email = $2",
			hashedPassword, email,
		)
		if err != nil {
			log.Fatalf("Failed to update user password: %v", err)
		}
		fmt.Printf("✓ User %s password updated\n", email)
		return
	}

	// Create new user
	var newID int
	err = conn.QueryRow(
		`INSERT INTO users (firstname, lastname, email, password, dob) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id`,
		"", "", email, hashedPassword, "2000-01-01",
	).Scan(&newID)

	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("✓ User created with ID %d, email: %s\n", newID, email)
}
