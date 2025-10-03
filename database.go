package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	client   *mongo.Client
	db       *mongo.Database
	users    *mongo.Collection
	messages *mongo.Collection
}

func NewDatabase(uri string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database("chatapp")
	return &Database{
		client:   client,
		db:       db,
		users:    db.Collection("users"),
		messages: db.Collection("messages"),
	}, nil
}

func (d *Database) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d.client.Disconnect(ctx)
}

func (d *Database) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Username: username,
		Password: string(hashedPassword),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = d.users.InsertOne(ctx, user)
	return err
}

func (d *Database) AuthenticateUser(username, password string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err := d.users.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *Database) SaveMessage(username, content string, timestamp time.Time, messageType string) error {
	message := Message{
		Username:  username,
		Content:   content,
		Timestamp: timestamp,
		Type:      messageType,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.messages.InsertOne(ctx, message)
	return err
}

func (d *Database) GetRecentMessages() ([]Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"type": bson.M{"$in": []string{"message", ""}}}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(50)
	cursor, err := d.messages.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []Message
	err = cursor.All(ctx, &messages)
	if err != nil {
		return nil, err
	}

	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}

	return messages, nil
}

func (d *Database) CreateDefaultUsers() error {
	defaultUsers := []struct {
		username string
		password string
	}{
		{"yuri", "yuricbtt"},
		{"bernardo", "bernardocbtt"},
		{"pedro", "pedrocbtt"},
		{"marcelo", "marcelocbtt"},
		{"giggio", "giggio123"},
		{"ramos", "ramosgay"},
		{"markin", "markinviado"},
	}

	for _, user := range defaultUsers {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var existingUser User
		err := d.users.FindOne(ctx, bson.M{"username": user.username}).Decode(&existingUser)
		cancel()

		if err == mongo.ErrNoDocuments {
			err = d.CreateUser(user.username, user.password)
			if err != nil {
				log.Printf("Error creating default user %s: %v", user.username, err)
			} else {
				log.Printf("Created default user: %s", user.username)
			}
		}
	}

	return nil
}
