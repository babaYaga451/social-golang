package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/babaYaga451/social/internal/store"
	"github.com/bxcodec/faker/v3"
)

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user", err)
			return
		}
	}

	tx.Commit()

	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating posts", err)
			return
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comment.Create(ctx, comment); err != nil {
			log.Println("Error creating comment", err)
			return
		}
	}

	log.Println("Seeding complete")
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]

		comments[i] = &store.Comment{
			UserID:  user.ID,
			PostID:  post.ID,
			Content: faker.Sentence(),
		}
	}
	return comments
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   faker.Sentence(),
			Content: faker.Paragraph(),
			Tags:    generateTags(),
		}
	}
	return posts
}

func generateTags() []string {
	tags := []string{"golang", "aws", "cloud", "coding", "backend", "frontend"}
	selectedTags := []string{}
	for i := 0; i < rand.Intn(4)+1; i++ {
		selectedTags = append(selectedTags, tags[rand.Intn(len(tags))])
	}
	return selectedTags
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			UserName: faker.FirstName() + fmt.Sprintf("%d", i),
			Email:    faker.Email(),
			RoleID:   1,
		}
	}
	return users
}
