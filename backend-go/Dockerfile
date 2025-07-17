FROM golang:1.22-alpine

WORKDIR /app

# Copy only go.mod first
COPY go.mod ./

# Force download of missing modules (this will also create go.sum)
RUN go mod tidy && go mod download

# Copy everything else
COPY . .

# Build binary
RUN go build -o main .

EXPOSE 8082
CMD ["./main"]
