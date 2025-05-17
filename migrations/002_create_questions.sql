CREATE TABLE questions (
  id          SERIAL PRIMARY KEY,
  user_id     INT   NOT NULL,
  question    TEXT  NOT NULL,
  answer      TEXT  NOT NULL,
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

