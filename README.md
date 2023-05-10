# Dictionary Go backend

This is a backend for dictionary web application.

## Main features:

* Written purely in Golang
* MySQL as a database backend
* Quiz mode (you can test yourself and your knowledge)

## Setup

1. Install Go from the official [download page](https://golang.org/dl/).

2. Clone the repository to your local machine:

```bash
git clone https://github.com/RubyLegend/dictionary-backend.git
```

3. Install the required packages:

```go
go mod download
```

4. Set up a [MySQL](https://www.mysql.com/downloads/) database, and create a .env file in the root directory of the project with the following variables:

```makefile
DB_USER=username
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=database_name
JWT_SECRET=secret_key
```

Replace the `username`, `password`, and `database_name` values with your MySQL database username, password, and database name respectively. You can set the `JWT_SECRET` variable to any secret key of your choice.

5. Run the database migrations to create the required tables:

```go
go run cmd/migrate/main.go up
```

6. Start the server:

```go
go run cmd/server/main.go
```

The server will start running on http://localhost:8080.

## API Endpoints

Dictionary backend exposes the following endpoints:

```
POST /api/v1/word
DELETE /api/v1/word/:id
GET /api/v1/word
PATCH /api/v1/word/:id

POST /api/v1/translation
DELETE /api/v1/translation/:id
GET /api/v1/translation
PATCH /api/v1/translation/:id

POST /api/v1/user/login
POST /api/v1/user/signup
POST /api/v1/user/logout
GET /api/v1/user/status
POST /api/v1/user/restore-username
POST /api/v1/user/restore-password
DELETE /api/v1/user
PATCH /api/v1/user

GET /api/v1/dictionary
POST /api/v1/dictionary
PATCH /api/v1/dictionary/:id
DELETE /api/v1/dictionary/:id

GET /api/v1/history
DELETE /api/v1/history

GET /api/v1/quiz/new
POST /api/v1/quiz/:quizId
GET /api/v1/quiz/status
```

Refer to the API documentation for more information on the parameters and responses of each endpoint.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](https://github.com/RubyLegend/dictionary-backend/blob/main/LICENSE) file for details.
