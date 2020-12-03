package main

import (
	"github.com/aws/aws-lambda-go/events"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"
	"encoding/json"
	"errors"
	"os"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,error) {
	fmt.Println("Request body", request.Body)
	var postData map[string]string
	json.Unmarshal([]byte(request.Body),&postData)
	group, ok := postData["group"]
	if !ok {
		return events.APIGatewayProxyResponse{Body: respJSON("could not process request or group param missing","FAILED",""), StatusCode: 422},nil
	}
	nextINT,err := GetNext(group)
	if err != nil {
		errorMessage := fmt.Sprintf("%s,%s","Could not calculate next in sequence",err)
		return events.APIGatewayProxyResponse{Body: respJSON(errorMessage,"FAILED",""), StatusCode: 422},nil
	}
	// Pad single digit numbers with leading zero so that 1 turns into 01
	if len(nextINT) == 1 {
		nextINT = fmt.Sprintf("0%s",nextINT)
	}
	next := fmt.Sprintf("%s-%s",group,nextINT)
	return events.APIGatewayProxyResponse{Body: respJSON("Success","Success",next), StatusCode: 201}, nil
}

func respJSON(message, status, data string) string {
	r := make(map[string]string)
	r["Message"] = message
	r["Status"] = status
	r["Data"] = data
	j , err := json.Marshal(r)
	if err != nil {
		return `{"Message": "Failed to marshall response","Status": FAIL}`
	}
	return string(j)
}

func GetNext(identifier string) (nextSequence string, e error){
	//https://docs.aws.amazon.com/cli/latest/reference/dynamodb/update-item.html#options
	//https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.UpdateExpressions.html#Expressions.UpdateExpressions.ADD
	sessionInfo, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		return "0",errors.New("service could not create AWS session:" + err.Error())
	}
	svc := dynamodb.New(sessionInfo)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {
				N: aws.String("1"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"GroupName": {
				S: aws.String(identifier),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		TableName:        aws.String(os.Getenv("DYNAMO_TABLE")),
		UpdateExpression: aws.String("ADD SequenceNumber :incr"),
	}
	res, err := svc.UpdateItem(input)
	if err != nil {
		log.Println(err)
		if awsErr, ok := err.(awserr.Error); ok {
			return "0", awsErr
		}
	}
	log.Println(*res.Attributes["SequenceNumber"].N)
	return *res.Attributes["SequenceNumber"].N,nil
}


func main() {
	lambda.Start(Handler)
}
