package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
)

func main() {
	storage, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Unable to initialize Storage client: %s", err)
	}

	dump(storage)

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	cron := time.NewTicker(7 * 24 * time.Hour).C

	for {
		select {
		case <-cron:
			dump(storage)
		case <-kill:
			return
		}
	}
}

func dump(storage *storage.Client) {
	bucketName := os.Getenv("GCS_BUCKET")
	if bucketName == "" {
		return
	}
	bucket := storage.Bucket(bucketName)
	orgs := strings.Split(os.Getenv("ORGANIZATIONS"), ",")
	for _, org := range orgs {
		dumpfile := fmt.Sprintf("%s.%d.dump", org, time.Now().Unix())
		filepath := fmt.Sprintf("/tmp/%s", dumpfile)
		output, err := os.Create(filepath)
		if err != nil {
			log.Fatalf("Unable to create dump file: %s", err)
		}
		defer output.Close()

		table := fmt.Sprintf("%s.messages", org)
		cmd := exec.Command("pg_dump", "-t", table, "historislack")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("Unable to execute pg_dump command: %s", err)
		}

		writer := bufio.NewWriter(output)

		if err := cmd.Start(); err != nil {
			log.Fatalf("Unable to start pg_dump command: %s", err)
		}

		go io.Copy(writer, stdout)
		cmd.Wait()
		writer.Flush()

		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatalf("Unable to read dump file: %s", err)
		}

		obj := bucket.Object(dumpfile)
		w := obj.NewWriter(context.Background())
		if _, err := w.Write(data); err != nil {
			log.Fatalf("Unable to write GCS file: %s", err)
		}

		if err := w.Close(); err != nil {
			log.Fatalf("Unable to close GCS file: %s", err)
		}

		if err := os.Remove(filepath); err != nil {
			log.Fatalf("Unable to remove dump file: %s", err)
		}
	}
}
