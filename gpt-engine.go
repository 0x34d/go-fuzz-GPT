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
You are a Golang fuzzing expert. Given a function name and number of arguments, determine whether the function is worth fuzzing. Answer Yes or No.

Here are some general guidelines:
	* Functions with names like Add, Remove, Get, Set, Update, Delete, Save, Call, and New are usually internal APIs and are not good targets for fuzzing.
	* Functions with a lot of arguments are also not good for fuzzing.
	* However, there are exceptions to these guidelines. For example, a function called ParseRequest() might be a good target for fuzzing, even though it has a generic name.
	* Ultimately, the best way to determine whether a function is worth fuzzing is to consider its specific implementation and the types of inputs that it accepts.
	* Function that call a network, or A database request consider them as a bad request, because for them i need to spin a network fuzzer.
	* Don't use function that's call offical golang library like : json,bas64,regexp

Example queries:
	* Is the function ParseRequest() worth fuzzing? (Answer: Yes)
	* Is the function Add() worth fuzzing? (Answer: No)
	* Is the function SaveUser() worth fuzzing? (Answer: Maybe, depending on its implementation)
	* Please answer the following query:
`

func GPTWork(funcName string, functions string, tests string) {

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
					Content: functions,
				},
			},
			Temperature: 0.5,
		},
	)

	if err != nil {
		fmt.Println(BlueColor + strings.Repeat("+", 80) + "" + ResetColor)
		fmt.Println(RedColor + "Function Name: " + funcName + ResetColor)
		fmt.Printf(CyanColor+"ChatCompletion error: %v\n"+ResetColor, err)
		fmt.Println(BlueColor + strings.Repeat("-", 80) + "\n" + ResetColor)
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
