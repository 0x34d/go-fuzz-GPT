package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var APIKEY = os.Getenv("OPENAI_API_KEY")
var SYSTEM_AI = `You are a Golang fuzzing expert, check if this function is worth fuzzing.`

func GPTWork(funcName string, functions string, tests string) {

	fmt.Println(BlueColor + strings.Repeat("+", 80) + "" + ResetColor)
	fmt.Println(RedColor + "Function Name:" + funcName + ResetColor)

	client := openai.NewClient(APIKEY)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "ft:gpt-3.5-turbo-0613:0x34d::7yZLjRMZ",
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
		fmt.Printf(CyanColor+"ChatCompletion error: %v\n"+ResetColor, err)
	} else {
		fmt.Println(CyanColor + resp.Choices[0].Message.Content + ResetColor)
	}

	fmt.Println(BlueColor + strings.Repeat("-", 80) + "\n" + ResetColor)
}

/*
	gptinput := functions + "\n\n\n" + tests
	jsonString, _ := json.Marshal(gptinput)
	fmt.Println(BlueColor + "GPT-3 Input: " + ResetColor + string(gptinput))
	fmt.Println(BlueColor + "GPT-3 Input: " + ResetColor + string(jsonString))
*/

/*
	if counter == 2 {
		return true
	}
	if len(funcDecls[0].Type.Params.List) > 2 {
		return
	}
	Blacklist := []string{"Add", "Remove", "Get", "Set", "Update", "Delete", "Save", "Call", "New"}
	for _, name := range Blacklist {
		if strings.Contains(strings.ToLower(funcName), strings.ToLower(name)) {
			fmt.Println(RedColor + "Blacklist name: `" + funcName + "`\n" + ResetColor)
			return
		}
	}
*/
