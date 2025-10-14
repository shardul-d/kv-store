package kvstore

import (
	"encoding/gob"
	"os"
)

// KeyDir represents an in-memory hashmap for faster lookups of the key.
type KeyDir map[string]Meta 

// Meta does not hold the value. 
// It holds the file, size of the value,
// and position of the value in the file. This enables
// lookup of the value in a single disk seek of the datafile.
type Meta struct {
	FileID int
	RecordSize int
	RecordPos int
	Timestamp int
}

// Encode encodes the map to a gob file.
// This is typically used to generate a hints file.
// Caller of this program should ensure to lock/unlock the map before calling.
func (k *KeyDir) Encode(fPath string) error {
	// Create a file for storing gob data.
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gob encoder.
	encoder := gob.NewEncoder(file)

	// Encode the map and save it to the file.
	err = encoder.Encode(k)
	if err != nil {
		return err
	}

	return nil
}

// Decode decodes the gob data in the map.
func (k *KeyDir) Decode(fPath string) error {
	// Open the file for decoding gob data.
	file, err := os.Open(fPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new gob decoder.
	decoder := gob.NewDecoder(file)

	// Decode the file to the map.
	err = decoder.Decode(k)
	if err != nil {
		return err
	}

	return nil
}
