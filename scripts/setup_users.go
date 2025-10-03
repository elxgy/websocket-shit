// This is a standalone utility to create default users
// Run with: cd scripts && go run setup_users.go ../models.go ../database.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from parent directory
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get MongoDB URI
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is required")
	}

	fmt.Printf("Connecting to MongoDB at: %s\n", mongoURI)

	// Connect to database
	db, err := NewDatabase(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Connected to MongoDB successfully!")
	fmt.Println("Setting up default users...")

	// Create default users
	if err := db.CreateDefaultUsers(); err != nil {
		log.Fatalf("Error creating default users: %v", err)
	}

	fmt.Println("\n✅ Default users created successfully!")
	fmt.Println("\nLogin credentials:")
	fmt.Println("- Username: alice    | Password: password123")
	fmt.Println("- Username: bob      | Password: password123")
	fmt.Println("- Username: charlie  | Password: password123")
	fmt.Println("- Username: diana    | Password: password123")
	fmt.Println("\nYou can now use these credentials to test the chat application.")
}