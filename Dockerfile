FROM golang:1.27.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY --parents web pkg main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main

ENV TODO_PORT=7540
ENV TODO_PASSWORD=test1234

# For persistent storage
VOLUME ["/app/db"]
ENV TODO_DBFILE=/app/db/scheduler.db

EXPOSE 7540

CMD ["./main"]