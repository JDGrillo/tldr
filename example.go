package main

import (
	"fmt"
	"encoding/json"
	"strings"
	"github.com/JesusIslam/tldr"
	"log"
	"net/http"
	"math"
)

type test_struct struct {
	Article_Body string
}

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "GET":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "get called"}`))
    case "POST":
		w.WriteHeader(http.StatusCreated)
		var t test_struct
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			panic(err)
		}
		bodyResponse := t.Article_Body
		originalWordNum := numberOfWords(bodyResponse)
		intoSentences := originalWordNum / 100
		bag := tldr.New()
		result, _ := bag.Summarize(bodyResponse, intoSentences)
		concatenatedString := ""
		for _, value := range result {
			concatenatedString = concatenatedString + " " + value
		}
		newWordNum := numberOfWords(concatenatedString)
		wordDifference := originalWordNum - newWordNum
		percentageRemoved := math.Round((float64(wordDifference) / float64(originalWordNum)) * 100.0)
		formattedString := fmt.Sprintf(`{"message": "%s","words_removed": %d, "percentage_removed": %f}`, concatenatedString, wordDifference, percentageRemoved)
		w.Write([]byte(formattedString))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
    }
}

func numberOfWords(sentence string) int {
	words := strings.Fields(sentence)
	return len(words)
}

func main() {
	http.HandleFunc("/tldr", home)
	fmt.Println("Server Started!")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
