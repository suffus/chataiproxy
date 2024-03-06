package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

var OpenAIClient *openai.Client

func ChatOne(openAIKey string, message string) string {
	//openAIKey := "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"
	client := openai.NewClient(openAIKey)
	OpenAIClient = client
	// Create a new chat completion
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		Messages:  getExampleSMSChatCompletions(message),
		MaxTokens: 75,
	}
	completion, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return "Error"
	}
	fmt.Println(completion.Choices[0].Message.Content)
	return completion.Choices[0].Message.Content
}

func getExampleSMSChatCompletions(message string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{Role: "system", Content: "You are helpful assistant for a recruiter who is communicating with candidates."},
		{Role: "system", Content: "You are working for Apex Systems and are communicating with a candidate who is applying for a truck driver role.  The candidate does not have a resume and does not have a computer."},
		{Role: "system", Content: "The candidate is applying for a truck driver role. The pay is $45 per hour plus a comprehensive benefits package. The job is full-time and requires a Class A CDL. The candidate must be able to drive a semi and have a clean driving record. The candidate must be able to pass a drug test.  The job starts in May 2024.  "},
		{Role: "system", Content: "Your client for this role is Walmart, and the job involves driving a semi from New York to Richmond, VA."},
		{Role: "system", Content: "Apex Systems has a mobile-friendly truck driver portal which includes a resume builder and an interview scheduler."},
		{Role: "system", Content: "The messages you return should be suitable for sending as SMS messages."},
		{Role: "system", Content: "If you do not know how to respond to a message, or if you have not been prompted with sufficient informatin to be sure of an answer, you can respond with a message that says 'I do not know how to respond to this message.'"},
		{Role: "user", Content: "Recruiter: Hello Martin I see you applied for our truck driver role.  Please send me your resume and I will get back to you."},
		{Role: "user", Content: "Candidate: OK I will send you my resume.  Where should I send it?"},
		{Role: "user", Content: "Recruiter: You can send it to me at this email address: mgf@apex.com"},
		{Role: "user", Content: "Candidate:" + message},
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

func makeWebserver() {

	// Set up a backend webserver that takes JSON requests and returns JSON responses
	// The JSON request should contain a message and the JSON response should contain the response message
	// we shall use the net/http package
	http.HandleFunc("/suggestions", func(w http.ResponseWriter, r *http.Request) {
		// Parse the Body into a ChatRequest
		fmt.Println("In suggestions")
		req := ChatRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Call open AI like we do in ChatOne()

		chatMessages := make([]openai.ChatCompletionMessage, len(req.Messages))
		for i, m := range req.Messages {
			chatMessages[i] = openai.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			}
		}

		// log all the messages
		for _, m := range chatMessages {
			fmt.Println(m.Role, m.Content)
		}

		// Create a new chat completion
		req2 := openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			Messages:  chatMessages,
			MaxTokens: 60,
		}
		completion, err := OpenAIClient.CreateChatCompletion(context.Background(), req2)
		if err != nil {
			fmt.Println(err)
		}

		// Return the response as JSON
		resp := ChatRequest{}
		for _, m := range completion.Choices {
			resp.Messages = append(resp.Messages, ChatMessage{
				Role:    m.Message.Role,
				Content: m.Message.Content,
			})
		}
		json.NewEncoder(w).Encode(resp)

	})

	http.HandleFunc("/dummy", func(w http.ResponseWriter, r *http.Request) {
		// Parse the Body into a ChatRequest
		fmt.Println("In dummy")
		req := ChatRequest{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Call open AI like we do in ChatOne()

		chatMessages := make([]openai.ChatCompletionMessage, len(req.Messages))
		for i, m := range req.Messages {
			chatMessages[i] = openai.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			}
		}

	})

}

func main() {
	openAIKey := flag.String("openAIKey", "", "OpenAI key")
	flag.Parse()
	OpenAIClient = openai.NewClient(*openAIKey)
	makeWebserver()
	http.ListenAndServe(":8080", nil)
}
