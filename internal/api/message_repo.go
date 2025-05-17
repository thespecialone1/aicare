package api

import (
	"database/sql"
	"log"
	"time"
)

// Message represents a single turn in chat
type Message struct {
	ID        int
	UserID    int
	Role      string
	Content   string
	CreatedAt time.Time
}

// MessageRepo handles db ops for messages
type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

// Save inserts a message row
func (r *MessageRepo) Save(userID int, role, content string) error {
	_, err := r.db.Exec(
		`INSERT INTO messages (user_id, role, content) VALUES($1, $2, $3)`,
		userID, role, content,
	)
	return err
}

// ListMessages returns the entire history for a user.
func (r *MessageRepo) ListMessages(userID int) ([]Message, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, role, content, created_at
     FROM messages
     WHERE user_id = $1
     ORDER BY id ASC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows.Close() error: %v", cerr)
		}
	}()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.UserID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
