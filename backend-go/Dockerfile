FROM golang:1.22-alpine

WORKDIR /app

# Copy everything first
COPY . .

# Force fresh module resolution (ignore cache)
RUN rm -f go.sum && go mod tidy && go mod download

RUN go build -o main .

EXPOSE 8082
CMD ["./main"]
