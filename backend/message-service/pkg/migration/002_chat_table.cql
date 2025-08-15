CREATE TABLE IF NOT EXISTS cluster.chat_history (
    chat_id text,
    message text,
    created_at timestamp,
    PRIMARY KEY (chat_id, created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);