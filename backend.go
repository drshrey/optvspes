package main

import (
	"encoding/json"

	"github.com/lineback/alchemyapi_go/alchemyAPI"

	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

type sentimentScore struct {
	DocSentiment interface{}
}

type tweet struct {
	Text string `json:"text"`
}

func main() {
	//Parse tweets into "tweet" struct format
	user_tweets := getTweets("user.json")

	//Get Document Sentiment scores from processSentiment function
	var testSentiment chan sentimentScore = make(chan sentimentScore)
	go processSentiment(user_tweets, testSentiment)
	go getScores(testSentiment)

	//Use channels to send information from processSentiment function value to getScores function value without any latency
}

/*
Through the testSentiment channel, this function takes the score and returns an array of ints (to eventually be aggregated)
*/
func getScores(testSentiment chan sentimentScore)

/*
Function that parses tweets from a json file, and returns an array of tweets
*/

func getTweets(filename string) (tweets []tweet) {
	configFile, err := os.Open(filename)
	if err != nil {
		log.Println("Error at opening file")
		return
	}
	var test []tweet
	jsonParser := json.NewDecoder(configFile)
	for {
		var partOf tweet
		if err = jsonParser.Decode(&partOf); err != nil {
			break
		}
		test = append(test, partOf)
	}
	return test

}

/*
Function that takes parsed tweets sends a document sentiment score of each tweet through a channel
*/

func processSentiment(tweets []tweet, testSentiment chan sentimentScore) {
	//Variable to store temporary sentiment score for type assertion (for passing through channel)
	var m sentimentScore

	keyBytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		log.Println("Couldn't find api key file correctly.")
		return
	}
	sentiment_doctor := alchemyAPI.NewAlchemist(strings.NewReader(string(keyBytes[:40])))
	for _, v := range tweets {
		score, err := sentiment_doctor.Sentiment("text", url.Values{}, v.Text)
		if err != nil {
			log.Println("Error while doing sentiment analysis on function")
			return
		}
		docSentiment := score["docSentiment"]
		m.DocSentiment = docSentiment
		do <- docSentiment
	}
}
