// Main - correct spelling using "spwnn" via Lambda on AWS
// (Shallow, pair-wise, neural network)
// Based on outline here:  https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/
// Must deploy with knownWords.txt file
package main

import (
	"errors"
	//"fmt"
	"log"

	"github.com/above-the-garage/spwnn"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNoWordToCheck is thrown when a name is not provided
	ErrNoWordToCheck = errors.New("Query did not contain a value for '?word='")

	// Read dictionary once on startup; will be reused if lambda is called w/o shutting down
	dict *spwnn.SpwnnDictionary
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	wordToCorrect := request.QueryStringParameters["word"]
	log.Printf("Spwnn request '%s'\n", wordToCorrect)

	if len(wordToCorrect) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNoWordToCheck
	}

	correctedWords, _ := spwnn.CorrectSpelling(dict, wordToCorrect, false /* strictLen */)
	var resultText string
	for _, val := range correctedWords {
		resultText = resultText + val.Word + "\n"
	}

	log.Printf("Spwnn result: %v\n", resultText)

	return events.APIGatewayProxyResponse{
			Body:       resultText,
			StatusCode: 200},
		nil

}

// InitDict loads the dictionary from knownWords.txt
func InitDict() {
	dict = spwnn.ReadDictionary(false)
}

func main() {
	// Read dictionary once at startup
	// Will naturally be reused if this lambda is called multiple
	// times without getting shut down
	InitDict()
	lambda.Start(Handler)
}

// Editorial comment:  Lambdas are awesome.
