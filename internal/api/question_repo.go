package api

import "database/sql"

type QuestionRepo struct {
	db *sql.DB
}

func NewQuestionRepo(db *sql.DB) *QuestionRepo {
	return &QuestionRepo{db: db}
}

func (r *QuestionRepo) Save(userID int, question, answer string) error {
	_, err := r.db.Exec(
		`INSERT INTO questions (user_id, question, answer) VALUES  ($1, $2, $3)`, userID, question, answer,
	)
	return err
}
