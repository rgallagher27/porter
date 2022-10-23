package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	if err := Run(); err != nil {
		log.Fatal(err)
	}

}

func Run() error {
	f, err := os.Open("testdata/ports.json")
	if err != nil {
		return err
	}
	defer f.Close()

	if err := decodeStream(f); err != nil {
		return err
	}

	return nil
}

func decodeStream(r io.Reader) error {
	dec := json.NewDecoder(r)

	// Get first token
	t, err := dec.Token()
	if err != nil {
		return err
	}

	// Expect to be at the start of a JSON object
	if t != json.Delim('{') {
		return fmt.Errorf("expected { but got %v", t)
	}

	for dec.More() {
		// Read the key.
		t, err := dec.Token()
		if err != nil {
			return err
		}
		key := t.(string)

		// Decode the value.
		var port map[string]any
		if err := dec.Decode(&port); err != nil {
			return err
		}

		// Add your code to process the key and value here.
		fmt.Printf("key %q, value %#v\n", key, port)
	}
	return nil
}
