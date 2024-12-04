package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type PostStore struct {
	db *sql.DB
}

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comment_count"`
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
  INSERT INTO posts (content, title, user_id, tags)
  VALUES ($1, $2, $3, $4) RETURNING id,
  created_at,
  updated_at
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	query := `
  SELECT id,
       content,
       title,
       user_id,
       tags,
       created_at,
       updated_at,
       VERSION
  FROM posts
  WHERE id = $1
`
	var post Post

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
  DELETE FROM posts WHERE id = $1
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrorNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
  UPDATE posts
  SET title = $1 , content = $2, version = version + 1
  WHERE id = $3 AND version = $4
  RETURNING version
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userId int64) ([]*PostWithMetaData, error) {
	return nil, nil
}
