# Build stage
FROM golang:1.19-alpine AS builder

WORKDIR /
COPY infrastructure .

# Set ARGs
ARG POSTGRES_DB
ARG POSTGRES_HOST
ARG POSTGRES_PORT
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG REDIS_HOST
ARG REDIS_PORT

ENV POSTGRES_DB=$POSTGRES_DB
ENV POSTGRES_HOST=$POSTGRES_HOST
ENV POSTGRES_PORT=$POSTGRES_PORT
ENV POSTGRES_USER=$POSTGRES_USER
ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV REDIS_HOST=$REDIS_HOST
ENV REDIS_PORT=$REDIS_PORT

ENV GO111MODULE=on
ENV TZ=Asia/Jakarta

# Set workdir
RUN mkdir -p /app
WORKDIR /app

# Copy all project code
COPY . .

RUN apk update && apk add git

# Download dependencies
RUN --mount=type=ssh go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o /tmp/app main.go

# Final stage
FROM alpine:latest AS production

# Copy output binary file from build stage
COPY --from=builder /tmp/app .
COPY --from=builder /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

#RUN go install github.com/pressly/goose/v3/cmd/goose@v3.15.0
#CMD "goose -dir /go/src/mini-wallet/infrastructure postgres \"host=localhost port=5432 user=postgres password=postgres dbname=mini-wallet sslmode=disable\" up"

CMD ["./app"]
