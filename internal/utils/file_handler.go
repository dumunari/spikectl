package utils

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dumunari/spikectl/internal/config"
)

func CreateTmpFile(ca []byte) string {
	err := os.MkdirAll(config.KUBE_CA_PATH, os.ModePerm)
	if err != nil {
		log.Fatal("[🐶] Error creating temp file path: ", err)
	}

	err = ioutil.WriteFile(config.KUBE_CA_FILE_PATH, ca, 0644)
	if err != nil {
		log.Fatal("[🐶] Error saving tmp CA File: ", err)
	}

	return config.KUBE_CA_FILE_PATH
}
