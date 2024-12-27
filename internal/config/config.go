package config

import (
	"errors"
	"github.com/goloop/env"
	"log"
	"os"
)

// MustLoad loading parameters from the environment file
func MustLoad(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("file %s does not exist", filename)
		} else {
			log.Fatal(err)
		}
	}

	err = env.Load(filename)
	if err != nil {
		log.Fatalf("Incorrect data in the configuration file: %s", err)
	}

	log.Printf("the configuration file '%s' has been uploaded successfully", filename)
}
