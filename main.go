package main

import (
  "os"
  "log"
  "fmt"
  "time"

  // Webserver dependencies
  "net/http"
  "github.com/gorilla/mux"

  // Database dependencies
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
  "io"
  "io/ioutil"
  "net/url"
)

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", PingHandler)
  r.HandleFunc("/message", MessageHandler)

  log.Fatal(http.ListenAndServe(os.Getenv("PORT"), r))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello world")
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err)
  }

  parsed, err := url.ParseQuery(string(body))
  if err != nil {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  if parsed.Get("token") != os.Getenv("TOKEN") {
    w.WriteHeader(http.StatusBadRequest)
  } else {
    db, err := sqlx.Connect("postgres", os.Getenv("DB_INFO"))
    if err != nil {
      panic(err)
    }
    db.MustExec(
      `INSERT INTO messages(token,team_id,team_domain,channel_id,channel_name,timestamp,user_id,user_name,text,trigger_word,service_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);`,
      parsed.Get("token"), parsed.Get("team_id"), parsed.Get("team_domain"), parsed.Get("channel_id"), parsed.Get("channel_name"), parsed.Get("timestamp"), parsed.Get("user_id"), parsed.Get("user_name"), parsed.Get("text"), parsed.Get("trigger_word"), parsed.Get("service_id"), time.Now(), time.Now())
    defer db.Close()
    w.WriteHeader(http.StatusOK)
  }
}
