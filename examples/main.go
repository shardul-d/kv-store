package main

import (
	"log"
	"os"
	"time"

	kvstore "github.com/shardul-d/kv-store"
)

var (
	lo = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	// Initialise.
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		lo.Fatalf("error creating data dir: %v", err)
	} // Creating data dir.

	kvstore, err := kvstore.Init(kvstore.WithDir("data/"), kvstore.WithAutoSync())
	if err != nil {
		lo.Fatalf("error initialising kvstore: %v", err)
	}

	var (
		key = "hello"
		val = []byte("world")
	)

	// Set a key.
	if err := kvstore.Put(key, val); err != nil {
		lo.Fatalf("error setting a key: %v", err)
	}

	// Fetch the key.
	v, err := kvstore.Get(key)
	if err != nil {
		lo.Fatalf("error fetching key %s: %v", key, err)
	}
	lo.Printf("fetched val: %s\n", string(v))

	// Set a new key with an expiry.
	key = "fruit"
	val = []byte("apple")
	ex := time.Second * 2
	if err := kvstore.PutEx(key, val, ex); err != nil {
		lo.Fatalf("error setting a key with ex: %v", err)
	}

	// Wait for 3 seconds for expiry.
	wait := time.Second * 3
	lo.Printf("waiting for %s for the key to get expired", wait.String())
	time.Sleep(wait)

	// Try fetching the expired key.
	_, err = kvstore.Get(key)
	if err != nil {
		lo.Printf("error fetching key %s: %v\n", key, err)
	}

	// Delete the key.
	if err := kvstore.Delete(key); err != nil {
		lo.Fatalf("error deleting key %s: %v", key, err)
	}

	// Fetch list of keys.
	keys := kvstore.List()
	for i, k := range keys {
		lo.Printf("key %d is %s\n", i, k)
	}

	if err := kvstore.Shutdown(); err != nil {
		lo.Fatalf("error closing kvstore: %v", err)
	}
}
