package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type payload map[string]interface{}

type sequence []payload

func (s sequence) payloads() []payload {
	return s
}

func main() {
	var (
		targetHost, payloadDir string
		interval               time.Duration
	)

	flag.StringVar(&targetHost, "host", "", "The target host")
	flag.StringVar(&payloadDir, "data-dir", "data", "The payload data directory")
	flag.DurationVar(&interval, "interval", time.Second, "The interval between requests")
	flag.Parse()

	if targetHost == "" {
		flag.Usage()

		os.Exit(1)
	}

	ss, err := loadSequences(payloadDir)
	if err != nil {
		log.Fatalf("enable to load event from data directory, %s", err.Error())
	}

	var wg sync.WaitGroup

	for _, s := range ss {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := processSequence(targetHost, s, interval); err != nil {
				log.Print("error: ", err)

				return
			}
		}()
	}

	wg.Wait()
}

func loadSequences(dir string) ([]sequence, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var ss []sequence
	for _, file := range files {
		c, err := ioutil.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		var s sequence
		if err := json.Unmarshal(c, &s); err != nil {
			return nil, fmt.Errorf("file %s : %w", file.Name(), err)
		}

		ss = append(ss, s)
	}

	return ss, nil
}

func sendPayload(url string, p payload) (err error) {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {

		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return fmt.Errorf("error: wrong status code expected 2xx got %d  ", resp.StatusCode)
	}

	return nil
}

func processSequence(host string, s sequence, interval time.Duration) error {
	for _, p := range s.payloads() {
		if err := sendPayload(host, p); err != nil {
			return err
		}

		time.Sleep(interval)
	}

	return nil
}
