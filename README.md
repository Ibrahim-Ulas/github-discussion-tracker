# Github Discussion Tracker

A tracker to pull discussions from a designated repository, mark the unanswered ones and notify the user.

I wrote this for personal use and my boot.dev personal project course.

## Features

* **Category Filter** An optional feature where user can set a category name to filter with.
* **Comment Follow** Any discussion unanswered or new comments coming to an already replied discussion is being tracked.
* **Notification** With a given interval to get notified, program will notify you if there is a new comment or an unanswered discussion.

## Requirements

* Go (1.18+)
* Github Personal Access Token

## Setup And Configuration
1. Clone the project.
2. Copy the "env.example" and change the values to your own configuration.

```bash
# Bağımlılıkları yükle
go mod tidy
```

```bash
# Programı çalıştır
go run .

```

## Alternative Setup
1. Download the compiled file from releases for your OS.
2. Create a new .env file or copy .env.example and configure it with your values.
3. Start the program with the compiled file.