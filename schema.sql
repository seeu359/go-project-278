CREATE TABLE links (
  id   BIGSERIAL PRIMARY KEY,
  url TEXT  NOT NULL,
  short_name TEXT NOT NULL UNIQUE
);

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
