# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

COPY middleware ./middleware
COPY repository ./repository
COPY routes ./routes

RUN go build -o /docker-dictionary

# Copy .env file, if it is present
COPY .env* ./

CMD ["/docker-dictionary"]
