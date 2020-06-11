package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var config *Config
var wg sync.WaitGroup
var fileLogger *log.Logger

func init() {
	configFile, err := os.Open("./config.json")
	if err != nil {
		log.Fatalln("Can't open config file")

	}
	defer configFile.Close()
	if err = json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatalln("Can't parse config file")
	}
}

func worker() {
	for n := 0; n < config.LoopNum; n++ {
		getToken()

	}
	wg.Done()
}

func getToken() {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("scope", config.Scope)

	start := time.Now()
	r, err := http.Post(config.TokenUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	duration := time.Since(start)
	if err != nil {
		fmt.Println(err.Error())
	}
	content, _ := ioutil.ReadAll(r.Body)

	log.Println(duration, string(content))
	fileLogger.Println(fmt.Sprintf("%s %s", duration, string(content)))
}

func main() {

	file, err := os.OpenFile("result.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	fileLogger = log.New(io.MultiWriter(file, os.Stderr),
		"",
		log.Lmicroseconds)

	wg.Add(config.ThreadNum)
	for n := 0; n < config.ThreadNum; n++ {
		go worker()
	}

	// Interrupt by Ctrl + C

	wg.Wait()
	log.Println("Complete")
}
