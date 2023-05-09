package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var APIKEY = os.Getenv("OPENAI_API_KEY")

var SYSTEM_AI = `
As an expert in Golang Fuzz testing, I utilize the official fuzzing framework introduced in Go 1.18 to create fuzz tests.
Your primary goal is to determine whether a given function, based on its implementation and associated tests, along with the corresponding GitHub link, is worth fuzzing. Focus on identifying if the function handles complex input data, as such functions are generally more valuable for fuzzing. Provide your answer as either YES or NO.

Your secondary objective, if the answer is YES, is to generate a fuzz test function. When creating the fuzz test function, keep the following in mind:
- The output should consist solely of the fuzz test function code without package and import statements.
- If you lack accurate information about the f.Add seed corpus, do not include it.
- Examine the test function for potential seed corpus.
- Avoid performing error checking if it is not necessary.
- Utilize return values from the function under test to call other functions, based on related examples from the test function, in order to enhance fuzzing coverage.
- Use the recommended approach for fuzz testing as introduced in Go 1.18, which involves the "testing" package and the "f.Add" method for providing a seed corpus.
`

func gptWork(functions string, tests string) {

	GPTinput := RemoteGitURL + "\n" + functions + "\n" + tests

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
					Content: GPTinput,
				},
			},
			Temperature: 0.5,
			TopP:        0.5,
			MaxTokens:   1024,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
