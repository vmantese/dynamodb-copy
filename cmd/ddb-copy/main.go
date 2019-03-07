package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

func main() {
	//flags from command line
	fromTable := flag.String("src", "", "from which table")
	toTable := flag.String("dest", "", "to which table")
	awsRegion := flag.String("aws-region", "", "which AWS Region")
	flag.Parse()

	sess, err := session.NewSession(
		&aws.Config{
			Region: awsRegion,
		},
	)
	if err != nill {
		fmt.Println("Cannot connect to AWS, error was:")
		panic(err)
	}

	db := dynamodb.New(sess)
	err = db.ScanPages(&dynamodb.ScanInput{
		TableName:      fromTable,
		ConsistentRead: aws.Bool(true),
	}, func(scanOutput *dynamodb.ScanOutput, lastPage bool) bool {
		for _, item := range scanOutput.Items {
			input := &dynamodb.PutItemInput{
				Item:      item,
				TableName: toTable,
			}

			_, err = db.PutItem(input)

			if err != nill {
				fmt.Println("Error with inserting item:")
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
		return !lastPage
	})
	if err != nill {
		fmt.Println(err)
	}
}
