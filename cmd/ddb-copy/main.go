package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func main() {
	//flags from command line
	src := flag.String("src", "", "source table")
	dest := flag.String("dest", "", "destination table")
	awsRegion := flag.String("aws-region", "", "which AWS Region")
	transformerOpt := flag.String("transformer-type", "", "transformer interface")
	flag.Parse()

	if *src == "" || *dest == "" {
		fmt.Println("Source and destination table names are required")
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region: awsRegion,
		},
	)
	if err != nil {
		fmt.Println("Cannot connect to AWS, error was:")
		panic(err)
	}

	transformer, err := makeTransformer(*transformerOpt)
	if *transformerOpt != "" && err != nil {
		fmt.Println("There was a problem creating the specified transformer:")
		fmt.Println(err.Error())
	}

	db := dynamodb.New(sess)
	err = db.ScanPages(&dynamodb.ScanInput{
		TableName:      src,
		ConsistentRead: aws.Bool(true),
	}, func(scanOutput *dynamodb.ScanOutput, lastPage bool) bool {
		for _, item := range scanOutput.Items {

			item = transformer.Transform(item)
			input := &dynamodb.PutItemInput{
				Item:      item,
				TableName: dest,
			}
			_, err = db.PutItem(input)

			if err != nil {
				fmt.Println("Error with inserting item:")
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
		return !lastPage
	})
	if err != nil {
		fmt.Println(err)
	}
}

type Transformer interface {
	Transform(map[string]*dynamodb.AttributeValue) map[string]*dynamodb.AttributeValue
}

type DefaultTransformer struct{}

func (*DefaultTransformer) Transform(av map[string]*dynamodb.AttributeValue) map[string]*dynamodb.AttributeValue {
	return av
}

func makeTransformer(name string) (Transformer, error) {
	switch name {
	case "":
		return new(DefaultTransformer{}),nil
	default:
		return new(DefaultTransformer{}), errors.New("unable to find selected transformer")
	}
}
