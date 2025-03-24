CREATE TABLE IF NOT EXISTS metadata
(
    install_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS data
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id TEXT,
    data       BLOB,
    data_type  TEXT NOT NULL,
    exec_time  DATETIME DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS logs
(
    message TEXT,
    level TEXT,
    log_time DATETIME DEFAULT CURRENT_TIMESTAMP
)