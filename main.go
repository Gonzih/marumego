package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
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

	go func() {
		for {
			time.Sleep(time.Minute * 10)
			fmt.Printf("Updating urls\n")
			urls = getURLs()
		}
	}()

	http.HandleFunc("/random.gif", func(w http.ResponseWriter, r *http.Request) {
		randomUrl := getRandomUrl(urls)
		fmt.Printf("Redirecting to %s\n", randomUrl)
		http.Redirect(w, r, randomUrl, 301)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
