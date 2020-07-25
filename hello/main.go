package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration

type Input struct {
	Text      string `message`
	ChannelID string ``
}

type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, input Input) (Response, error) {
	var buf bytes.Buffer
	Token := os.Getenv("BOT_KEY")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating bot reason: ", err)
	}

	fmt.Println(input.ChannelID)

	client := dg.Open()
	if client != nil {
		fmt.Println("Error opening client session. Reason: ", client)
	}

	random, err := dg.ChannelMessageSend(input.ChannelID, input.Text)
	if err != nil {
		fmt.Println("Message send failed, readin: ", err)
	}
	fmt.Println(random)
	body, err := json.Marshal(map[string]interface{}{
		"message": input.Text,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
