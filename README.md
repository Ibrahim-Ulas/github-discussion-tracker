# Github Discussion Tracker

A tracker to pull discussions from a designated repository, mark the unanswered ones and notify the user.

## Features

* **CategoryFilter** An optional feature where user can set a category name to filter with
* **Comment Follow** Any discussion unanswered or new comments coming to an already replied discussion is being tracked

## Requirements

* Go (1.18+)
* Github Personal Access Token

## Setup And Configuration
1. Clone to project.
2. Copy the "env.example" and change the values to your own configuration.

```bash
# Bağımlılıkları yükle
go mod tidy
```

```bash
# Programı çalıştır
go run main.go

```