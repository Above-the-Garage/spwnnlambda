module spwnnlambda

go 1.18

require (
	github.com/above-the-garage/spwnn v0.0.0-20220814165059-edf91933ae0c
	github.com/aws/aws-lambda-go v1.36.1
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/above-the-garage/spwnn => ../spwnn
