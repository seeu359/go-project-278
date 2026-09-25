-- +goose Up
CREATE TABLE links (
  id   BIGSERIAL PRIMARY KEY,
  url TEXT  NOT NULL,
  short_name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE links;
