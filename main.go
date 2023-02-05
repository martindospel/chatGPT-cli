package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(100),
		Temperature: gpt3.Float32Ptr(0),
	}, func(res *gpt3.CompletionResponse) {
		fmt.Print(res.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

func main() {
	//set config file
	viper.SetConfigFile(".env")
	//read the config file
	viper.ReadInConfig()
	//get the api key from config file
	apiKey := viper.GetString("API_KEY")
	//validate api key
	if apiKey == "" {
		panic("Missing API key")
	}

	//set up new context
	ctx := context.Background()
	//set up new client
	client := gpt3.NewClient(apiKey)
	//cobra library for command line arguments
	rootCommand := &cobra.Command{
		Use:   "chatpgt",
		Short: "Chat with ChatGPT in console.",
		//function that scans the input, set the exit string, assign question according to user input
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Print("Say something or write 'quit' to end: ")
				if !scanner.Scan() {
					break
				}
				question := scanner.Text()
				switch question {
				//this user input ends the communication
				case "quit":
					quit = true
				default:
					//send user input to get response function
					GetResponse(client, ctx, question)
				}
			}
		},
	}
	rootCommand.Execute()
}
