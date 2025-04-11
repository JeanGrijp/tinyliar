CREATE TABLE links
(
  id SERIAL PRIMARY KEY,
  short_url VARCHAR(255) UNIQUE NOT NULL,
  original_url TEXT NOT NULL,
  clicks BIGINT DEFAULT 0,
  owner_id BIGINT,
  expired_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
