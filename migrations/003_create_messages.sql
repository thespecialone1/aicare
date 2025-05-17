-- migrations/003_create_messages.sql 
CREATE TABLE IF NOT EXISTS messages (
  id            SERIAL      PRIMARY KEY, 
  user_id       INT         NOT NULL,
  role          VARCHAR(10) NOT NULL,  -- "USER" OR "assistent"
  content       TEXT        NOT NULL,
  created_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);
