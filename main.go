package main

import (
  "os"
  "log"
  "fmt"
  "time"

  // Webserver dependencies
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"

  // Database dependencies
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
)

type Message struct {
  Token string `json:"token"`
  TeamID string `json:"team_id"`
  TeamDomain string `json:"team_domain"`
  ChannelID string `json:"channel_id"`
  ChannelName string `json:"channel_name"`
  Timestamp string `json:"timestamp"`
  UserID string `json:"user_id"`
  UserName string `json:"user_name"`
  Text string `json:"text"`
  TriggerWord string `json:"trigger_word"`
  ServiceID string `json:"service_id"`
}

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
  decoder := json.NewDecoder(r.Body)

  var message Message
  err := decoder.Decode(&message)
  if err != nil {
    panic(err)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")

  if message.Token != os.Getenv("TOKEN") {
    w.WriteHeader(http.StatusBadRequest)
  } else {
    db, err := sqlx.Connect("postgres", os.Getenv("DB_INFO"))
    if err != nil {
      panic(err)
    }
    db.MustExec(
      `INSERT INTO messages(token,team_id,team_domain,channel_id,channel_name,timestamp,user_id,user_name,text,trigger_word,service_id,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);`,
      message.Token, message.TeamID, message.TeamDomain, message.ChannelID, message.ChannelName, message.Timestamp, message.UserID, message.UserName, message.Text, message.TriggerWord, message.ServiceID, time.Now(), time.Now())
    defer db.Close()

    w.WriteHeader(http.StatusOK)
  }
}
