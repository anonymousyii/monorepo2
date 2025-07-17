package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    _ "github.com/lib/pq"
    "github.com/redis/go-redis/v9"
    "context"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func createTable() {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users_go (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    )`)
    if err != nil {
        panic(err)
    }
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    cached, _ := rdb.Get(ctx, "users_go_cache").Result()
    if cached != "" {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cached))
        return
    }

    rows, err := db.Query("SELECT id, name FROM users_go")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name)
        users = append(users, u)
    }

    data, _ := json.Marshal(users)
    rdb.Set(ctx, "users_go_cache", string(data), 0)
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var u User
    json.NewDecoder(r.Body).Decode(&u)
    _, err := db.Exec("INSERT INTO users_go(name) VALUES($1)", u.Name)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    rdb.Del(ctx, "users_go_cache")
    w.WriteHeader(201)
}

func main() {
    pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
    )
    var err error
    db, err = sql.Open("postgres", pgURL)
    if err != nil {
        panic(err)
    }
    createTable()

    rdb = redis.NewClient(&redis.Options{
        Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
    })

    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            getUsers(w, r)
        } else if r.Method == http.MethodPost {
            createUser(w, r)
        }
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    fmt.Println("Go server running on port", port)
    http.ListenAndServe("0.0.0.0:"+port, nil)
}
