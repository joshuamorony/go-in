package main

import (
	"context"
	"fmt"
	"os/exec"

	openai "github.com/sashabaranov/go-openai"
)

type Note struct {
	Title string
	Body  string
}

func main() {
	// TODO allow passing input through CLI
	example := "use go to make a program that can create a summary title given a full note and create a file"

	title, err := generateTitle(example)
	if err != nil {
		fmt.Printf("Error generating title: %v\n", err)
		return
	}

	fmt.Println("generate title:", title)
	// TODO pass to createNote to create the actual file
}

func generateTitle(content string) (string, error) {
	prompt := "Summarize the following note in a title of at most 8 words." +
		"This title is intended for the creator of the note and should focus on being clear" +
		"and consise it does not need to be editorialised: '" + content + "'"

	token := getToken()
	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func createNote(note Note) {
	// use title/body to create note in $in folder
}

func getToken() string {
	cmd := exec.Command("op", "read", "op://private/OpenAI/credential", "--no-newline")

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		return ""
	}

	return string(output)
}
