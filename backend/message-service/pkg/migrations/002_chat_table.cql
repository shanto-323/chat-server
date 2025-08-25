CREATE TABLE IF NOT EXISTS cluster.chat_history (
    conversation_id text,
    created_at timestamp,
    chat_id UUID,
    payload blob,
    PRIMARY KEY (conversation_id, created_at, chat_id)
) WITH CLUSTERING ORDER BY (created_at DESC);