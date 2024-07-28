# Preface 🚀
Hola! I created this utility CLI to streamline my daily tasks. It's grown to include backend support, with an API on the horizon! 🌐

# Pre-installation Requirements 📋

### Software Needed
- 🐳 [Docker](https://docs.docker.com/engine/install/)
- 🐹 [Golang](https://go.dev/doc/install)
    - 🔧 [sqlboiler CLI](https://github.com/volatiletech/sqlboiler)
    - 🔧 [golang-migrate](https://github.com/golang-migrate/migrate)
- 🪝 [pre-commit](https://pre-commit.com/)

### Environment Variables
Specify your environment variables:
```sh
export THEOVERWATCHTOOLS_APP_DIR=/Users/dembygenesisabella/apps/utils
export THEOVERWATCHTOOLS_DB_HOST="localhost"
export THEOVERWATCHTOOLS_DB_PORT="3308"
export THEOVERWATCHTOOLS_DB_USER="root"
export THEOVERWATCHTOOLS_DB_PASS="secret"
export THEOVERWATCHTOOLS_DB_DATABASE="example"
export THEOVERWATCHTOOLS_DB_USE_EXISTING_MARIADB="0"
export THEOVERWATCHTOOLS_DB_QUERY_TIMEOUT="10s"
export THEOVERWATCHTOOLS_DB_EXEC_TIMEOUT="10s"
```
- 🌍 `MIGRATION_DIR`: Specify your app directory (e.g., `./internal/database/migrations`).

# Installation Steps 🔧

1. Execute migration and install pre-commit hooks:
   ```sh
   sh ./scripts/dev_setup.sh
   ```

2. Generate docs:
   ```sh
   Generate your docs here
   ```

# Script Commands 🛠️ (Context: Main Directory)
- `sh ./scripts/build-cli.sh`: Compiles the CLI.
- `sh ./scripts/build-di.sh`: Compiles the container.
- `sh ./scripts/build-sqlboiler.sh`: Generates sqlboiler ORM files.
- `sh ./scripts/migrate.sh`: Performs database migration.

# Convenience Commands
- `clear_all && go test ./... -parallel=100 -count=1`: Runs tests with parallel execution and no-cache.

## Features 🌟

### Copy Root Path to Clipboard ✅
- **Command**: `clip-file-contents`
- Copies the specified root path's contents to the clipboard, excluding `.GIT`, IDE configurations, and non-essential files. Each copy includes a header for file identification.

### Clip GPT Code Standards Preface ✅
- **Command**: `clip-gpt-preface`
- Enhances ChatGPT code quality by incorporating a preface that focuses on defensive programming, testability, readability, and modularity.

### Copy One Folder to Another ✅
- Facilitates folder content transfer with options for exclusions and pre-transfer cleanup, preserving essential metadata.

### Todo Roadmap 🗺️
- Implement a `Makefile` for rapid development setup in a Docker environment, including binary compilation and CLI integration into shell configurations.
- Enhance CLI documentation with detailed command descriptions.
- Introduce automated testing in CI workflows for pull requests and merges - ✅
- Expand Docker configurations for backend development.
- Refactor utility package to be more idiomatic by creating sub-packages.

---