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
	src := flag.String("src", "", "from which table")
	dest := flag.String("dest", "", "to which table")
	awsRegion := flag.String("aws-region", "", "which AWS Region")
	flag.Parse()

	sess, err := session.NewSession(
		&aws.Config{
			Region: awsRegion,
		},
	)
	if err != nil {
		fmt.Println("Cannot connect to AWS, error was:")
		panic(err)
	}

	db := dynamodb.New(sess)
	err = db.ScanPages(&dynamodb.ScanInput{
		TableName:      src,
		ConsistentRead: aws.Bool(true),
	}, func(scanOutput *dynamodb.ScanOutput, lastPage bool) bool {
		for _, item := range scanOutput.Items {
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
