package main

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
	"github.com/krizanauskas/winko/internal/handlers/tgbothandler"
	"github.com/krizanauskas/winko/internal/service"
	"github.com/krizanauskas/winko/pkg/bscclient"
	"github.com/krizanauskas/winko/pkg/config"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"net/http"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := gotgbot.NewBot(cfg.TgBotApiKEy, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout, // Customise the default request timeout here
				APIURL:  gotgbot.DefaultAPIURL,  // As well as the Default API URL here (in case of using local bot API servers)
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	openAiClient := openai.NewClient(cfg.OpenAiApiKey)
	openAiService := service.NewOpenAiService(openAiClient)
	bscClient := bscclient.NewClient(bscclient.Config{
		BaseURL: cfg.BscApiUrl,
		Timeout: cfg.Timeout,
		ApiKey:  cfg.BscApiKey,
	})

	cryptoService := service.NewCryptoService(bscClient)

	botHandler := tgbothandler.New(
		openAiService,
		cryptoService,
	)

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)
	dispatcher.AddHandler(handlers.NewMessage(message.Text, botHandler.ProcessIncommingMessage))

	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	//log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()

	//const AssistantId = "asst_SlrV2qbUDOWtNESqpi2uqyOD"

	//bscClient := bscclient.NewClient(bscclient.Config{
	//	BaseURL: cfg.BscApiUrl,
	//	Timeout: cfg.Timeout,
	//	ApiKey: cfg.BscApiKey,
	//})
	//
	//bnbHoldings, err := bscClient.GetBnbAllocation("0xAf0476C27A15b2A6C7b9BDFe410fe0E59Ef7bEAA");
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("BNB holdings: %s", bnbHoldings)

	//assistantName := "Winko"
	//instructions := "you are assistant which helps user with crypto related questions"
	//assistant, err := openaiClient.CreateAssistant(context.Background(), openai.AssistantRequest{
	//	Model:        "gpt-3.5-turbo-1106",
	//	Name:         &assistantName,
	//	Instructions: &instructions,
	//	Tools: []openai.AssistantTool{
	//		{
	//			Type: openai.AssistantToolTypeFunction,
	//			Function: &openai.FunctionDefinition{
	//				Name:        "get_bnb_allocation",
	//				Description: "get bnb allocation of user address",
	//			},
	//		},
	//		{
	//			Type: openai.AssistantToolTypeFunction,
	//			Function: &openai.FunctionDefinition{
	//				Name:        "send_bnb_to_address",
	//				Description: "send bnb token to recipient address",
	//				Parameters: jsonschema.Definition{
	//					Type: jsonschema.Object,
	//					Properties: map[string]jsonschema.Definition{
	//						"amount": {
	//							Type:        jsonschema.Number,
	//							Description: "Amount of bnb to send to recipient",
	//						},
	//						"recipient_address": {
	//							Type:        jsonschema.String,
	//							Description: "blockchain address of recipient to send bnb to",
	//						},
	//					},
	//					Required: []string{"amount", "recipient_address"},
	//				},
	//			},
	//		},
	//	},
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Printf("Assistant created: %s", assistant.ID)

	//run, err := openaiClient.CreateThreadAndRun(context.Background(), openai.CreateThreadAndRunRequest{
	//	RunRequest: openai.RunRequest{
	//		AssistantID: AssistantId,
	//	},
	//	Thread: openai.ThreadRequest{
	//		Messages: []openai.ThreadMessage{
	//			{
	//				Role:    openai.ThreadMessageRoleUser,
	//				Content: "Hello, whats up?",
	//			},
	//		},
	//	},
	//})
	//
	//runId := run.ID
	//threadId := run.ThreadID
	//
	//log.Printf("Thread created: %s", threadId)
	//log.Printf("Run created: %s", runId)

	//retries := 0
	//
	//for run.Status != openai.RunStatusCompleted {
	//	if retries < 5 {
	//		log.Printf("Thread status: %s", run.Status)
	//		retries++
	//		time.Sleep(1 * time.Second)
	//
	//		run, err = openaiClient.RetrieveRun(context.Background(), threadId, runId)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//	} else {
	//		log.Fatal("Thread status is not completed after 5 retries")
	//	}
	//}
	//
	//if run.Status == openai.RunStatusCompleted {

	//threadId := "thread_dh7cXcwMJ42fQEzJyUNT6UXU"

	//messages, err := openaiClient.ListMessage(context.Background(), threadId, nil, nil, nil, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, message := range messages.Messages {
	//	if message.Role == openai.ChatMessageRoleAssistant {
	//		if message.Content[0].Type == "text" {
	//			response := message.Content[0].Text.Value
	//
	//			log.Printf("Response: %s", response)
	//		}
	//	}
	//
	//	if message.Role == openai.ChatMessageRoleUser {
	//		break
	//	}
	//}
}

//}
