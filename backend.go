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

	//Process Sentiment function to get map
	mapScores := processSentiment(userTweets, testSentiment)

	//Returns Score structs from map[string]interfaces{}
	scores, err := convertToScoreStruct(mapScores)
	if err != nil {
		log.Println("ERROR WHILE DOING convertToScoreStruct")
		return
	}

	//Gets the floating point score values from the Score struct
	var listOfScores []float64
	for _, v := range scores {
		floatScore := getScore(v)
		listOfScores = append(listOfScores, floatScore)
	}

	//Aggregates slice of scores using aggregateScores function
	finalScore := aggregateScores(listOfScores)
	endTime := time.Now()
	fmt.Println("Final sum is:", finalScore)
	if finalScore <= 0 {
		fmt.Println("You're probably a negative person")
	} else {
		fmt.Println("You're probably a positive person")
	}

	fmt.Println("Start time was ...", startTime)
	fmt.Println("Finish time was ...", endTime)
}

/*
Function that aggregates scores from a slice of floats
*/
func aggregateScores(input []float64) (output float64) {
	var out float64
	for _, v := range input {
		out += v
	}
	return out
}

/*
Function that takes a Score struct and returns the float score property
*/
func getScore(input Score) (output float64) {
	if input.DocSentiment["score"] == nil {
		return
	}
	out := input.DocSentiment["score"].(string)
	floatOut, err := strconv.ParseFloat(out, 64)
	if err != nil {
		log.Println("ERROR IN RUNNING getScore function")
		return
	}
	return floatOut
}

/*
Function that processes map[string]interface{} from processSentiment function
*/
func convertToScoreStruct(input []map[string]interface{}) (output []Score, err error) {
	var out []Score
	for _, v := range input {
		var m Score
		err := mapstructure.Decode(v, &m)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		out = append(out, m)
	}
	return out, nil
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

func processSentiment(tweets []tweet, testSentiment chan map[string]interface{}) (result []map[string]interface{}) {
	fmt.Println("STARTING processSentiment")
	//Variable to store temporary sentiment score for type assertion (for passing through channel)

	keyBytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
		log.Println("Couldn't find api key file correctly.")
		return
	}
	sentiment_doctor := alchemyAPI.NewAlchemist(strings.NewReader(string(keyBytes[:40])))
	var res []map[string]interface{}
	for k, v := range tweets {
		if k == 100 {
			break
		}

		score, err := sentiment_doctor.Sentiment("text", url.Values{}, v.Text)
		if err != nil {
			log.Println("Error while doing sentiment analysis on function")
			return
		}
		res = append(res, score)
		//testSentiment <- score
	}
	return res
}
