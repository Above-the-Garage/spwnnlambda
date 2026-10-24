// Main - correct spelling using "spwnn" via Lambda on AWS
// (Shallow, pair-wise, neural network)
// Based on outline here:  https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/
// Must deploy with knownWords.txt file
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/above-the-garage/spwnn"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// ErrNoWordToCheck is thrown when a name is not provided
	ErrNoWordToCheck = errors.New("Query did not contain a value for '?word='")

	// Read dictionary once on startup; will be reused if lambda is called w/o shutting down
	dict *spwnn.SpwnnDictionary

	// Startup message shows time to load dictionary
	startTimeMsg string
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var resultText string

	resultText = startTimeMsg
	startTimeMsg = "Lambda reused from previous execution\n"
	resultText = resultText + "\n"

	start := time.Now()

	// Optional "?length=any" (default) or "?length=same";
	// "?strict" is kept as a shorthand for "?length=same"
	lengthMatch := spwnn.AnyLength
	if _, strict := request.QueryStringParameters["strict"]; strict {
		lengthMatch = spwnn.SameLength
	}
	if name, ok := request.QueryStringParameters["length"]; ok {
		var err error
		lengthMatch, err = parseLength(name)
		if err != nil {
			return events.APIGatewayProxyResponse{
					Body:       err.Error() + "\n",
					StatusCode: 400},
				nil
		}
	}
	wordToCorrect := request.QueryStringParameters["word"]
	log.Printf("Spwnn request '%s'\n", wordToCorrect)

	if len(wordToCorrect) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNoWordToCheck
	}

	// Optional "?scoring=capped" (default) or "?scoring=original"
	scoring := spwnn.CappedPairs
	if name, ok := request.QueryStringParameters["scoring"]; ok {
		var err error
		scoring, err = parseScoring(name)
		if err != nil {
			return events.APIGatewayProxyResponse{
					Body:       err.Error() + "\n",
					StatusCode: 400},
				nil
		}
	}

	correctedWords, _ := spwnn.CorrectSpelling(dict, wordToCorrect, lengthMatch, scoring)
	for _, val := range correctedWords {
		resultText = resultText + val.Word + "\n"
	}

	duration := time.Since(start)
	resultText = resultText + fmt.Sprintln("\nExecution ", duration)

	log.Printf("Spwnn result: %v\n", resultText)

	return events.APIGatewayProxyResponse{
			Body:       resultText,
			StatusCode: 200},
		nil

}

// InitDict loads the dictionary from knownWords.txt
func InitDict() {
	dict = spwnn.ReadDictionary("knownWords.txt")
}

func main() {
	// Read dictionary once at startup
	// Will naturally be reused if this lambda is called multiple
	// times without getting shut down
	start := time.Now()
	InitDict()
	duration := time.Since(start)
	startTimeMsg = fmt.Sprintln("Lambda start and dictionary load", duration)
	lambda.Start(Handler)
}

// parseScoring converts the scoring query text to a spwnn.PairScoring
func parseScoring(name string) (spwnn.PairScoring, error) {
	switch name {
	case "capped":
		return spwnn.CappedPairs, nil
	case "original":
		return spwnn.OriginalPairs, nil
	}
	return spwnn.CappedPairs, fmt.Errorf("unknown scoring %q (want capped or original)", name)
}

// parseLength converts the length query text to a spwnn.LengthMatch
func parseLength(name string) (spwnn.LengthMatch, error) {
	switch name {
	case "any":
		return spwnn.AnyLength, nil
	case "same":
		return spwnn.SameLength, nil
	}
	return spwnn.AnyLength, fmt.Errorf("unknown length %q (want any or same)", name)
}

// Editorial comment:  Lambdas are awesome.
