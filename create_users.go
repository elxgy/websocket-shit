// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("🚀 Setting up default users for chat application...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get MongoDB URI
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGODB_URI environment variable is required")
	}

	fmt.Printf("📡 Connecting to MongoDB...\n")

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer client.Disconnect(ctx)

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	db := client.Database("chatapp")
	usersCollection := db.Collection("users")

	fmt.Println("✅ Connected to MongoDB successfully!")
	fmt.Println("👥 Creating default users...")

	// Define default users
	defaultUsers := []struct {
		username string
		password string
	}{
		{"alice", "password123"},
		{"bob", "password123"},
		{"charlie", "password123"},
		{"diana", "password123"},
	}

	usersCreated := 0
	usersSkipped := 0

	for _, user := range defaultUsers {
		// Check if user already exists
		var existingUser bson.M
		err := usersCollection.FindOne(ctx, bson.M{"username": user.username}).Decode(&existingUser)

		if err == mongo.ErrNoDocuments {
			// User doesn't exist, create them
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("❌ Error hashing password for %s: %v", user.username, err)
				continue
			}

			newUser := bson.M{
				"username": user.username,
				"password": string(hashedPassword),
			}

			_, err = usersCollection.InsertOne(ctx, newUser)
			if err != nil {
				log.Printf("❌ Error creating user %s: %v", user.username, err)
			} else {
				fmt.Printf("✅ Created user: %s\n", user.username)
				usersCreated++
			}
		} else if err != nil {
			log.Printf("❌ Error checking user %s: %v", user.username, err)
		} else {
			fmt.Printf("⚠️  User %s already exists, skipping\n", user.username)
			usersSkipped++
		}
	}

	fmt.Println("\n🎉 User setup complete!")
	fmt.Printf("✅ Users created: %d\n", usersCreated)
	fmt.Printf("⚠️  Users skipped: %d\n", usersSkipped)

	fmt.Println("\n📋 Available login credentials:")
	fmt.Println("┌──────────┬─────────────┐")
	fmt.Println("│ Username │ Password    │")
	fmt.Println("├──────────┼─────────────┤")
	fmt.Println("│ alice    │ password123 │")
	fmt.Println("│ bob      │ password123 │")
	fmt.Println("│ charlie  │ password123 │")
	fmt.Println("│ diana    │ password123 │")
	fmt.Println("└──────────┴─────────────┘")

	fmt.Println("\n🚀 You can now start the chat server and test with these credentials!")
}