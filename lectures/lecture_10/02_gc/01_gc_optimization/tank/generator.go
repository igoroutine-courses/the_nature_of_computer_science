package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func generateRequest(method, url, tag, body string) string {
	return strings.Join([]string{method, url, tag, body}, "||") + "\n"
}

func main() {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(wd, "/ammo_template.txt"))

	if err != nil {
		panic(err)
	}

	defer func() {
		err = f.Close()

		if err != nil {
			log.Println(err)
		}
	}()

	for range 100_000 {
		_, err = f.WriteString(generateRequest(http.MethodPost, "/v1/action", "action", ""))

		if err != nil {
			panic(err)
		}
	}
}
