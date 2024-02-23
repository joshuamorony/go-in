package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	openai "github.com/sashabaranov/go-openai"
)

type Note struct {
	Title string
	Body  string
}

func main() {
	var input string

	if len(os.Args) < 2 {
		fmt.Print("Enter a note: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input = scanner.Text()
	} else {
		input = os.Args[1]
	}

	title, err := generateTitle(input)
	if err != nil {
		fmt.Printf("Error generating title: %v\n", err)
		return
	}

	createNote(Note{
		Title: title,
		Body:  input,
	})
}

func generateTitle(content string) (string, error) {
	prompt := "Summarize the following note in a title of at most 8 words." +
		"This title is intended for the creator of the note and should focus on being clear" +
		"and consise it does not need to be editorialised. It should be in kebab-case: '" + content + "'"

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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
	}

	filePath := filepath.Join(homeDir, "Documents", "notes", "$in", note.Title+".md")

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(note.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Created note: ", note.Title)
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
