package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var APIKEY = os.Getenv("OPENAI_API_KEY")
var SYSTEM_AI = `
As an expert in fuzzing with Golang, your task is to evaluate whether a given function, along with its number of arguments, is a suitable candidate for fuzzing. You should respond with either 'Yes' or 'No'.

Here are some general considerations to guide your decision:
	* Internal API functions are typically not optimal targets for fuzzing.
	* The best way to determine a function's suitability for fuzzing is to examine its specific implementation and the nature of inputs it accepts.
	* Functions that require network setup, call a network socket, or make a database request are generally not ideal for fuzzing.
	* Avoid functions that call official Golang libraries such as json, base64, regexp, etc.

Given these guidelines, evaluate the following functions for their suitability for fuzzing:
`
//If Yes, Generate a Golang fuzz test for given function, use the Go 1.18 recommended fuzz testing approach;
//- If you lack accurate information about the f.Add seed corpus, do not include it.
//- Examine the test function for potential seed corpus.
//- Avoid performing error checking if it is not necessary.
//- Utilize return values from the function under test to call other functions, based on related examples from the test function, in order to enhance fuzzing coverage.
//- Use the recommended approach for fuzz testing as introduced in Go 1.18, which involves the "testing" package and the "f.Add" method for providing a seed corpus.

func GPTWork(funcName string, functions string, tests string) {
	gptinput := functions + "\n" + tests

	client := openai.NewClient(APIKEY)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4",
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
			Temperature: 0.5,
		},
	)

	if err != nil {
		fmt.Printf(RedColor+"Function Name: "+funcName+CyanColor+"\nChatCompletion error: %v\n"+ResetColor, err)
		return
	}

	if resp.Choices[0].Message.Content == "No" {
		return
	}

	fmt.Println(BlueColor + strings.Repeat("+", 80) + "" + ResetColor)
	fmt.Println(RedColor + "Function Name: " + funcName + ResetColor)
	fmt.Println(CyanColor + resp.Choices[0].Message.Content + ResetColor)
	fmt.Println(BlueColor + strings.Repeat("-", 80) + "\n" + ResetColor)
}
