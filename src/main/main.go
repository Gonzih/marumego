package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

var giphyAPIKey = os.Getenv("GIPHY_API_KEY")

const searchURL = "https://api.giphy.com/v1/gifs/search"

type GiphyResponse struct {
	Data []struct {
		Images map[string]struct {
			URL string `json:"url"`
		} `json:"images"`
	} `json:"data"`
}

func getRandomURL() string {
	url, _ := url.Parse(searchURL)
	q := url.Query()
	q.Set("api_key", giphyAPIKey)
	q.Set("q", "maru cat")
	q.Set("limit", "100")
	url.RawQuery = q.Encode()

	resp, _ := http.Get(url.String())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response GiphyResponse
	json.Unmarshal(body, &response)

	l := len(response.Data)

	if l == 0 {
		return ""
		log.Println("Nothing found")
	}

	n := rand.Intn(l)
	return response.Data[n].Images["original"].URL
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<html><head><meta http-equiv='refresh' content='30'><title>Maru!</title></head><body><img src='/random.gif'/></body></html>")
	})

	http.HandleFunc("/random.gif", func(w http.ResponseWriter, r *http.Request) {
		randomUrl := getRandomURL()
		fmt.Printf("Redirecting to %s\n", randomUrl)
		http.Redirect(w, r, randomUrl, 302)
	})

	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
