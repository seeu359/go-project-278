-- +goose Up
CREATE TABLE visits (
  id   BIGSERIAL PRIMARY KEY,
  link_id INT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  ip TEXT NOT NULL,
  user_agent TEXT,
  status INT,

  CONSTRAINT fk_Link
  FOREIGN KEY (link_id)
  REFERENCES links(id)
);

-- +goose Down
DROP TABLE visits;
