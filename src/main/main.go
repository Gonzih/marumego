package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const gifsURL = "https://raw.githubusercontent.com/paulhenrich/marume-server/master/resources/gifs.txt"

func getURLs() string {
	resp, _ := http.Get(gifsURL)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

func getRandomUrl(urls string) string {
	urlsArray := strings.Split(urls, "\n")
	randomIndex := rand.Intn(len(urlsArray))
	randomUrl := urlsArray[randomIndex]

	return randomUrl
}

func main() {
	urls := getURLs()
	urlsMutex := &sync.Mutex{}

	rand.Seed(time.Now().UTC().UnixNano())

	go func() {
		for {
			time.Sleep(time.Minute * 10)
			fmt.Printf("Updating urls\n")
			urlsMutex.Lock()
			defer urlsMutex.Unlock()
			urls = getURLs()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<html><body><img src='/random.gif'/></body></html>")

	})

	http.HandleFunc("/random.gif", func(w http.ResponseWriter, r *http.Request) {
		urlsMutex.Lock()
		defer urlsMutex.Unlock()
		randomUrl := getRandomUrl(urls)
		fmt.Printf("Redirecting to %s\n", randomUrl)
		http.Redirect(w, r, randomUrl, 301)
	})

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
