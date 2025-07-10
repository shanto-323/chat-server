CREATE KEYSPACE IF NOT EXISTS cluster WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': 1
};
