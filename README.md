# Slack Message Logger

## Via RTM Bot

Run `make docker-build`, edit [schema](./schema/01_schema.sql) and create tables per organization, run `make docker-run`.

## Via Outgoing Webhooks

### Usage

`DB_INFO="user=postgres dbname=slack-history sslmode=disable" PORT=":8080" TOKEN="asdf1234" go run main.go`
