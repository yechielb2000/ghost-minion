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
    save_time  DATETIME DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE IF NOT EXISTS logs
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    message  TEXT,
    level    TEXT,
    save_time DATETIME DEFAULT CURRENT_TIMESTAMP
)