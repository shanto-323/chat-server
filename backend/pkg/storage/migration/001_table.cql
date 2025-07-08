CREATE TABLE IF NOT EXISTS clients (
  id TEXT PRIMARY KEY,
  conn TEXT,
  alive BOOLEAN,
  last_alive TIMESTAMP
);