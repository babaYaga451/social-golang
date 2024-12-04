package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
   SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.id, u.username
   FROM comments c
   JOIN users u ON u.id = c.user_id
   WHERE c.post_id = $1
   ORDER BY c.created_at DESC
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.ID, &c.User.UserName)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
  INSERT INTO comments (user_id, post_id, content)
  VALUES ($1, $2, $3) RETURNING id, created_at
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.UserID,
		comment.PostID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	return err
}
