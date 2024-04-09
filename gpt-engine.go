package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

var APIKEY = os.Getenv("OPENAI_API_KEY")
var SYSTEM_AI = `
Determine if a given Golang function should undergo fuzz testing. Answer with 'Yes' or 'No'.
Internal API functions typically do not benefit from fuzzing.
Functions that unmarshal, decode, or parse are good candidates; those that marshal or encode are not.
Functions such as get, set, add, remove, delete, create, destroy, and init are unsuitable for fuzz testing.
Network setup functions are also poor choices for fuzz testing.
Assess a function's implementation and the nature of its inputs to decide if it is appropriate for fuzz testing.
`

func GPTWork(funcName string, functions string, tests string) {

	gptinput := functions + "\n" + tests

	client := openai.NewClient(APIKEY)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4-0125-preview",
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
			Temperature: 0.2,
		},
	)

	if err != nil {
		fmt.Printf(RedColor+"Function Name: "+funcName+CyanColor+"\nChatCompletion error: %v\n"+ResetColor, err)
		return
	}

	//	if resp.Choices[0].Message.Content == "No" {
	//		return
	//	}

	fmt.Println(BlueColor + strings.Repeat("+", 80) + "" + ResetColor)
	fmt.Println(RedColor + "Function Name: " + funcName + ResetColor)
	fmt.Println(CyanColor + resp.Choices[0].Message.Content + ResetColor)
	fmt.Println(BlueColor + strings.Repeat("-", 80) + "\n" + ResetColor)

	time.Sleep(2 * time.Second)
}
