# Ghost Minion

This is a Linux backdoor tool designed to collect system data and securely transmit it to a remote server.   
It supports multiple types of data collection, including system information, files, and periodic commands.  
The tool ss intended for research, or red team exercises in controlled environments.

> **⚠️ Important:** This tool is for authorized security testing and research only. Unauthorized use on systems without
> permission is illegal and unethical.

## Configuration

On installation, provide a **YAML** config file in the same directory as the tool binary. It defines agent ID, database
paths, logging, communication, and app options.  
You can see the structure in the file structure in the [config.go](ghostminion/config/config.go).
> **⚠️ Important:** Config must stay in the same directory as the binary.

## Installation process

Place the tool binary and config file in the same directory.
The config file should follow the YAML structure described above.  
Run the tool and optionally set a custom --id for the agent-id:

```shell
./ghostminion --id agent123
```

the tool will:

- Load the configuration.
- Initialize logging.
- Hide itself on the system.
- Initialize the database with the configured path and password.
- Generate or apply the AgentID and save it in the config.
- Start default apps (Keylogger, Screenshot, Security Guard).
- The tool continuously listens for tasks from the communication routine and starts apps as needed.

> Ensure the config file is writable so the tool can update the AgentID.

## Apps

The agent fetches tasks from the server, which are instantiated as apps and run as goroutines. Each app performs a
specific function and can run concurrently.

App Types

- Keylogger – records keystrokes.
- Screenshot – captures screen at intervals.
- Periodic Command – executes predefined commands periodically.
- Periodic Get File – fetches files from the system.
- Connect Online – attempts to connect to a server (depends on environment).
- Security Guard – protects the backdoor; cannot be modified after initialization.

App Manager
The `AppManager` is a singleton that manages all apps. Key features:

- `StartApp(name, app)` – starts an app as a goroutine.
- `StopApp(name)` – stops and removes an app.
- `StartAll()` / `StopAll()` – manage all apps at once.
- `ListApps()` – returns currently running apps.
- `GetApp(name)` – retrieve a specific app instance.

App Factory
Tasks fetched from the server are converted into apps using `NewAppFactory(task)`. The factory automatically
instantiates the correct app type based on the task type.

> **⚠️ Security Guard is critical:** it monitors the backdoor state and may terminate the agent if necessary. It cannot
> be
> added or edited after the first initialization.

## Communication Protocol

The agent communicates with servers using HTTP by default, but the communication layer can be overridden to use custom
protocols.

#### How it Works

Authentication:

- The agent requests a challenge from the server using its AgentID.
- Computes an HMAC using the challenge and the server key.
- Send the HMAC back to verify authenticity.

Task Fetching:

- Periodically, the agent fetches tasks from a randomly selected server.
- Each task is converted into an app by the App Factory and started as a goroutine.

Data Ex-filtration:

- Old data rows (logs, collected data) are sent to the server in JSON format.
- Communication is designed to avoid sending too much data at once and can respect system constraints like CPU usage.

## Hider Capability

The tool includes a process-hiding feature to make the agent less visible on the system. On startup, the agent attempts
to:

- Hide the process: Renames its /proc/<pid> entry to a hidden path in /tmp.
- Overwrite the executable: Write garbage bytes to /proc/self/exe to obscure the binary.
- Delete itself: Removes the original executable file from disk.
  These steps help reduce detection on the host system.

> ⚠️ Hiding techniques are OS-specific and may not work on all Linux distributions. They are intended for controlled,
> authorized testing environments only.