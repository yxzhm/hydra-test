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

func introspectToken(token Token) {
	data := url.Values{}
	data.Set("token", token.AccessToken)
	start := time.Now()
	r, err := http.Post(config.IntrospectUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	duration := time.Since(start)
	if err != nil {
		fmt.Println(err.Error())
	}
	content, _ := ioutil.ReadAll(r.Body)

	log.Println(duration, string(content))
	fileLogger.Print(fmt.Sprintf("%s %s", duration, string(content)))
}

func getToken(logEnable bool) Token {
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
	if logEnable {
		log.Println(duration, string(content))
		fileLogger.Println(fmt.Sprintf("%s %s", duration, string(content)))
	}

	var token Token

	json.Unmarshal(content, &token)
	return token
}

func main() {

	argsWithoutProg := os.Args[1:]
	// By default, test the token generating.
	mode := "token"
	if len(argsWithoutProg) > 0 {
		if argsWithoutProg[0] == "-i" {
			mode = "introspect"
		}
	}
	fmt.Println(mode)

	file, err := os.OpenFile(fmt.Sprintf("%s.log", mode),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	fileLogger = log.New(io.MultiWriter(file, os.Stderr),
		"",
		log.Lmicroseconds)

	wg.Add(config.ThreadNum)
	for n := 0; n < config.ThreadNum; n++ {
		if mode == "token" {
			go func() {
				for n := 0; n < config.LoopNum; n++ {
					getToken(true)
				}
				wg.Done()
			}()
		} else {
			go func() {
				token := getToken(false)
				for n := 0; n < config.LoopNum; n++ {
					introspectToken(token)
				}
				wg.Done()
			}()
		}

	}

	wg.Wait()
	log.Println("Complete")
}
