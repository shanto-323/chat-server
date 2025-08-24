CREATE TABLE IF NOT EXISTS cluster.chat_history (
    chat_id UUID,
    conversation_id TEXT
    sender_id text,
    receiver_id text,
    message text,
    created_at timestamp,
    offline boolean,
    PRIMARY KEY (chat_id, created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);