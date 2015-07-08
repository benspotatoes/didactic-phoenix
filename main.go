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
  "io"
  "io/ioutil"
  "net/url"
)

type QueryRequest struct {
  Field string `json:"field"`
  Query string `json:"query"`
}

type QueryResp struct {
  Messages []Message `json:"messages"`
}

type Message struct {
  ID int `db:"id"`
  ChannelName string `db:"channel_name"`
  Timestamp string `db:"timestamp"`
  UserName string `db:"user_name"`
  Text string `db:"text"`
  TeamDomain string `db:"team_domain"`
}

type Result struct {
  UserName string `db:"user_name"`
  Text string `db:"text"`
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/", PingHandler)
  r.HandleFunc("/message", MessageHandler)
  r.HandleFunc("/query", QueryHandler)
  r.HandleFunc("/search", SearchHandler)

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
    db, err := sqlx.Connect("postgres", os.Getenv("DB_INFO") + fmt.Sprintf(" dbname=%s", parsed.Get("channel_name")))
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

func QueryHandler(w http.ResponseWriter, r *http.Request) {
  var queryReq QueryRequest

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&queryReq)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, err.Error())
    return
  }

  db, err := sqlx.Connect("postgres", os.Getenv("DB_INFO"))
  if err != nil {
    panic(err)
  }
  // db = db.Unsafe()

  messages := []Message{}
  query := fmt.Sprintf("SELECT id,channel_name,timestamp,user_name,text,team_domain FROM messages WHERE %s LIKE '%s';", queryReq.Field, "%%" + queryReq.Query + "%%")
  if os.Getenv("DEBUG") == "true" {
    fmt.Printf(query)
  }

  err = db.Select(&messages, query)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, err.Error())
    return
  }

  messageResp := QueryResp{Messages: messages}

  resp, err := json.Marshal(messageResp)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, err.Error())
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write(resp)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err)
  }

  parsed, err := url.ParseQuery(string(body))
  if err != nil {
    panic(err)
  }

  if parsed.Get("token") != os.Getenv("SEARCH_TOKEN") {
    w.WriteHeader(http.StatusBadRequest)
  } else {
    db, err := sqlx.Connect("postgres", os.Getenv("DB_INFO") + fmt.Sprintf(" dbname=%s", parsed.Get("channel_name")))
    if err != nil {
      panic(err)
    }
    
    results := []Result{}
    query := fmt.Sprintf("SELECT user_name,text FROM messages WHERE text LIKE '%s' ORDER BY id desc LIMIT 10;", "%%" + parsed.Get("text") + "%%")
    if os.Getenv("DEBUG") == "true" {
      fmt.Printf(query)
    }

    err = db.Select(&results, query)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, err.Error())
      return
    }

    resp := ""
    for _, res := range results {
      resp = resp + fmt.Sprintf("%s: %s;\n", res.UserName, res.Text)
    }
    
    w.Write([]byte(resp))
  }
}
