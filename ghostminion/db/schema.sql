CREATE TABLE IF NOT EXISTS metadata
(
    install_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS data
(
    request_id TEXT,
    data       BLOB,
    data_type  TEXT,
    exec_time  TIMESTAMP

);
