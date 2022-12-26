// Main - correct spelling using "spwnn" via Lambda on AWS
// (Shallow, pair-wise, neural network)
// Based on outline here:  https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/
// Must deploy with knownWords.txt file
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/above-the-garage/spwnn"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNoWordToCheck is thrown when a name is not provided
	ErrNoWordToCheck = errors.New("PathParameter did not contain a value for 'check'")

	// Read dictionary once on startup; reuse if lambda is called w/o shutting down
	dict *spwnn.SpwnnDictionary
)

// Handler returns an array of words that match the request string
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// See:  https://docs.aws.amazon.com/lambda/latest/dg/go-programming-model-handler-types.html
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	wordToCheck := request.PathParameters["check"]
	log.Printf("Spwnn request '%s'\n", wordToCheck)

	if len(wordToCheck) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNoWordToCheck
	}

	correctedWords, _ := spwnn.CorrectSpelling(dict, wordToCheck, false /* strictLen */)
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
	fmt.Println("main")
	// Read dictionary once at startup
	// Will naturally be reused if this lambda is called multiple
	// times without getting shut down
	InitDict()
	lambda.Start(Handler)
}

// Editorial comment:  Lambdas are awesome.