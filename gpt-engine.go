package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var APIKEY = os.Getenv("OPENAI_API_KEY")
var SYSTEM_AI = "You are a Golang Fuzz testing expart."

func GPTWork(funcName string, functions string, tests string) {

	Blacklist := []string{"Add", "Remove", "Get", "Set", "Update", "Delete", "Save"}
	for _, name := range Blacklist {
		if strings.Contains(strings.ToLower(funcName), strings.ToLower(name)) {
			fmt.Println(RedColor + "Sorry, this function is in the blacklist. You can manually check it. : `" + funcName + "`\n" + ResetColor)
			return
		}
	}

	gptinput := functions + "\n" + tests

	fmt.Println(RedColor + funcName + ResetColor)
	fmt.Println(BlueColor + strings.Repeat("+", 80) + ResetColor)

	client := openai.NewClient(APIKEY)
	resp, err := client.CreateChatCompletion(

		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: SYSTEM_AI,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: gptinput,
				},
			},
			Temperature: 0.1,
		},
	)

	if err != nil {
		fmt.Printf(CyanColor+"ChatCompletion error: %v\n"+ResetColor, err)
	} else {
		fmt.Println(CyanColor + resp.Choices[0].Message.Content + ResetColor)
	}

	fmt.Println(BlueColor + strings.Repeat("-", 80) + "\n" + ResetColor)
}
