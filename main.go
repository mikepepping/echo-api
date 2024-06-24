package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

//go:embed config.json
var cfgFile []byte

type Config struct {
	Delay int               `json:"delay"`
	Codes map[string]string `json:"codes"`
	Port  int               `json:"port"`
}

func main() {
	var (
		cfg *Config = &Config{}
	)

	// Load server config
	if err := json.Unmarshal(cfgFile, cfg); err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/{code}", func(w http.ResponseWriter, r *http.Request) {
		// Start timing so we know how long to delay for before finising
		beggining := time.Now()

		code := r.PathValue("code")
		message := cfg.Codes[r.PathValue("code")]
		if message == "" {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "status code not found")
			return
		}

		status, err := strconv.Atoi(code)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "not a status code")
			return
		}

		// if we are supplied a delay param use it, if not use the default from the config
		delay := cfg.Delay
		delayParam := r.URL.Query().Get("delay")
		if delayParam != "" {
			if delay, err = strconv.Atoi(delayParam); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, "invalid delay time")
				return
			}
		}

		// TODO: find a better way to simulate delay
		delay = delay * int(time.Millisecond)
		time.Sleep(time.Duration(delay) - time.Since(beggining))

		w.WriteHeader(status)
		if echo := r.URL.Query().Get("echo"); echo != "" {
			io.WriteString(w, echo)
			return
		}

		io.WriteString(w, message)

	})

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}
