CREATE TABLE IF NOT EXISTS cluster.chat_offline (
    sender_id text,    
    receiver_id text,
    message text,
    created_at timestamp,
    PRIMARY KEY (receiver_id, created_at)
) WITH CLUSTERING ORDER BY (created_at ASC);