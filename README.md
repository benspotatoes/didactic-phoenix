# Slack Message Logger

## Via RTM Bot

See [Dockerfile](./Dockerfile)

## Via Outgoing Webhooks

### Usage

`DB_INFO="user=postgres dbname=slack-history sslmode=disable" PORT=":8080" TOKEN="asdf1234" go run cmd/http/main.go`

## Export

See [Dockerfile](./export.Dockerfile)
