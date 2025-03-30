# RSS Feed Aggregator

## Introduction
The RSS Feed Aggregator is a web server built with Go that allows users to:

- Add RSS feeds to be collected
- Follow and unfollow RSS feeds that other users have added
- Fetch all the latest posts from the RSS feeds they follow

RSS feeds are a way for websites to publish updates to their content. This project lets you stay updated with your favorite blogs, news sites, and podcasts all in one place!

## Features
- **User Authentication**: Users can sign up and log in.
- **Add Feeds**: Users can add RSS feeds to be aggregated.
- **Follow Feeds**: Users can follow or unfollow RSS feeds.
- **Retrieve Posts**: Users can fetch the latest posts from the feeds they follow.
- **Background Worker**: Periodically fetches new posts from subscribed feeds.
- **PostgreSQL Integration**: Stores user data, feeds, and posts.

## Prerequisites
To run this project, ensure you have the following installed:

1. [Go](https://golang.org/doc/install) (latest version recommended)
2. [PostgreSQL](https://www.postgresql.org/download/)
3. [VS Code](https://code.visualstudio.com/) (or any editor of your choice)
4. An HTTP client like [Thunder Client](https://www.thunderclient.io/) or [Postman](https://www.postman.com/)

## Setup
### 1. Clone the repository
```sh
 git clone https://github.com/tonge3199/go-RSS-project.git
 cd go-RSS-project
```

### 2. Set up environment variables
Create a `.env` file in the project root with the following:
```ini
DB_URL=postgres://youruser:yourpassword@localhost:5432/rss_aggregator?sslmode=disable
PORT=8080
```

### 3. Start PostgreSQL 

### 4. Run database migrations

​For detailed information, check out [SQLC](https://sqlc.dev/) and [Goose](https://pressly.github.io/goose/).
See the config-file in sqlc.yml 

### 5. Start the server
```sh
go run main.go
```

## Usage
### API Endpoints
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/users` | Create a new user |
| `POST` | `/feeds` | Add a new RSS feed |
| `GET` | `/feeds` | Get all available RSS feeds |
| `POST` | `/follow` | Follow a feed |
| `DELETE` | `/unfollow/:feed_id` | Unfollow a feed |
| `GET` | `/posts` | Get the latest posts from followed feeds |

### Example Request (Add a Feed)
```sh
curl -X POST http://localhost:8080/feeds -H "Content-Type: application/json" -d '{"url": "https://example.com/rss"}'
```

### Example Response
```json
{
  "id": "1234",
  "url": "https://example.com/rss",
  "created_at": "2025-03-30T12:00:00Z"
}
```

## Future Improvements
- Implement OAuth for authentication
- Enhance UI with a front-end client
- Add filtering and search capabilities

## License
This project is licensed under the MIT License. Feel free to use and contribute!

## Author
Developed by **Your Name**. Contributions welcome!

