package main

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"unicode/utf8"

	"cloud.google.com/go/pubsub"
	"github.com/coreos/go-systemd/sdjournal"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "vidston"

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new topic.
	topicName := "test-fluidor-logs"

	// Creates the new topic.
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}

	fmt.Printf("Topic %v created.\n", topic)

	j, err := sdjournal.NewJournal()
	check(err)
	defer j.Close()

	err = j.SeekTail()
	check(err)

	for {
		for {
			n, err := j.Next()
			check(err)
			if n <= 0 {
				break
			}
			entry, err := j.GetEntry()
			check(err)
			m := map[string]interface{}{}
			for k, v := range entry.Fields {
				if utf8.ValidString(v) {
					m[k] = v
				} else {
					sl := make([]int, len(v))
					for i := range v {
						sl[i] = int(v[i])
					}
					m[k] = sl
				}
			}
			m["__CURSOR"] = entry.Cursor
			m["__REALTIME_TIMESTAMP"] = entry.RealtimeTimestamp
			m["__MONOTONIC_TIMESTAMP"] = entry.MonotonicTimestamp
			data, err := json.Marshal(m)
			check(err)
			fmt.Println(string(data))
			fmt.Println(len(data))
			buf := &bytes.Buffer{}
			w := zlib.NewWriter(buf)
			io.Copy(w, bytes.NewBuffer(data))
			w.Flush()
			fmt.Println(len(buf.Bytes()))

		}

		j.Wait(sdjournal.IndefiniteWait)
	}
}
