# aggregator_program

Gator is a small multi-user CLI app I built for aggregating RSS feeds and browsing posts directly from the terminal.  
It’s designed to be simple, fast, and a good excuse to work with Go + Postgres.

---

## Installation

Before installing, make sure you have both Go and PostgreSQL set up.

P.S. Note that while you *can* use the program with Windows I highly recommend that you download WSL as that is a far better user experience for this program.

### 1. Install Go

Download and install Go from the official site:  
https://golang.org/dl/

Verify the installation:

```bash
go --version
```

---

### 2. Install PostgreSQL (v15 or later)

#### macOS (Homebrew)

```bash
brew install postgresql@15
```

#### Linux / WSL (Debian-based)

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Verify the installation using the `psql` CLI:

```bash
psql --version
```

Make sure it reports version 15 or later.

---

## Configuration

Gator expects a config file in your home directory:

`~/.gatorconfig.json`

Example:

```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```

Replace the connection string with your local database credentials.

### Custom Config Location (Optional)

By default, the config file is read from your home directory using:

```go
home, err := os.UserHomeDir()
```

If you want to store the config somewhere else, you can modify the  
`getConfigFilePath` function in `config/config.go`:

```go
home := new_folder_location
```

And then delete the following which is found right after:

```go
if err != nil {
   return "", err
}
```

⚠️ Note: After making this change, your existing compiled binary will be outdated.  
To continue using `aggregator_program`, rebuild it with:

```bash
go build
```

---

## Running the Program

While you *can* run commands using:

```bash
go run .
```

the project is intended to be used via the compiled binary:

```bash
./aggregator_program
```

---

## Getting Started

Create a user:

```bash
aggregator_program register <name>
```

Log in:

```bash
aggregator_program login <name>
```

Add a feed:

```bash
aggregator_program addfeed <url>
```

Start the aggregator (polls feeds on an interval):

```bash
aggregator_program agg 30s
```

> Keep this running in a separate terminal, or stop it temporarily if needed, as it runs indefinitely.

Browse posts:

```bash
aggregator_program browse [limit]
```

---

## Other Commands

- `aggregator_program users` – list all users  
- `aggregator_program feeds` – list all feeds  
- `aggregator_program follow <url>` – follow an existing feed  
- `aggregator_program unfollow <url>` – unfollow a feed  

---

## Help Command

There’s also a built-in `help` command that prints all available commands and their descriptions.

It dynamically lists every registered command (sorted alphabetically) along with a short description, so you don’t have to remember everything.

---

## Notes

- Feeds are stored and managed in Postgres
- The aggregator runs continuously and fetches updates on an interval you choose
- Designed as a lightweight terminal-first experience (no UI, no distractions)

---

## Why I Built This

This was created as part of the Boot.dev "Build a Blog Aggregator in Go" course.
