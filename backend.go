package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lineback/alchemyapi_go/alchemyAPI"
	"github.com/mitchellh/mapstructure"
)

type Score struct {
	Status       string
	Usage        string
	Url          string
	Language     string
	DocSentiment map[string]interface{}
}

type tweet struct {
	Text string `json:"text"`
}

func main() {
	startTime := time.Now()
	//Parse tweets into "tweet" struct format
	userTweets := getTweets("user.json")

	//Get Document Sentiment scores from processSentiment function
	testSentiment := make(chan map[string]interface{})

	//Use channels to send information from processSentiment function value to getScores function value without any latency
	go processSentiment(userTweets, testSentiment)

	var agg []float64
	//Run channel
	for s := range testSentiment {
		var m Score
		err := mapstructure.Decode(s, &m)
		if err != nil {
			log.Println("ERROR WHILE DECODING MAP INTERFACE STRING THING")
			return
		}
		if m.DocSentiment["score"] == nil {
			continue
		}
		score, err := strconv.ParseFloat(m.DocSentiment["score"].(string), 64)
		if err != nil {
			log.Println("ERROR WHILE RUNNING STRING CONVERSION IN MAIN")
			return
		}
		agg = append(agg, score)
		testSum := sumFloats(agg)
		fmt.Println("Current Sum is:", testSum)
	}
	//Sum everything up in the slice
	sum := sumFloats(agg)
	fmt.Println("Final sum is:", sum)
	endTime := time.Now()
	fmt.Println("Start time is,", startTime)
	fmt.Println("End time is", endTime)
	if sum > 0 {
		fmt.Println("You're probably an optimist.")
	} else {
		fmt.Println("You're probably a pessimist.")
	}
}

/*
Function that sums everything in the slice
*/
func sumFloats(toAdd []float64) (sum float64) {
	var result float64
	for _, v := range toAdd {
		result += v
	}
	return result
}

/*
Function that parses tweets from a json file, and returns an array of tweets
*/

func getTweets(filename string) (tweets []tweet) {
	fmt.Println("STARTING getTweets")
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

func processSentiment(tweets []tweet, testSentiment chan map[string]interface{}) {
	fmt.Println("STARTING processSentiment")
	//Variable to store temporary sentiment score for type assertion (for passing through channel)

	keyBytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		log.Println("Couldn't find api key file correctly.")
		return
	}
	sentiment_doctor := alchemyAPI.NewAlchemist(strings.NewReader(string(keyBytes[:40])))
	for k, v := range tweets {
		if k == 201 {
			break
		}
		score, err := sentiment_doctor.Sentiment("text", url.Values{}, v.Text)
		fmt.Println(score)
		if err != nil {
			log.Println("Error while doing sentiment analysis on function")
			return
		}
		testSentiment <- score
	}
}
