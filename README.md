# GhostMinion

/ghostminion
├── cmd/
│   ├── main.go                # Main entry point
│
├── config/
│   ├── config.yaml            # Configuration file (store encrypted data)
│
├── communication/
│   ├── https_server.go        # HTTPS server code for handling encrypted requests
│   ├── commands.go            # Command processing logic with encryption
│
├── executor/
│   ├── executor.go            # Command execution logic (encrypt results)
│
├── persistence/
│   ├── persistence.go         # Ensures persistence and encrypts data
│
├── db/
│   ├── db_handler.go          # SQLite interaction (encrypted database)
│   ├── schema.sql             # Database schema
│
├── encryption/
│   ├── encryption.go          # Encryption/decryption utilities (AES, etc.)
│
├── utils/
│   ├── logger.go              # Custom logger (with encrypted log files)
│
├── README.md                  # Documentation
└── .gitignore                 # Files to ignore