package tgbothandler

import (
	"context"
	"encoding/json"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/krizanauskas/winko/internal/service"
	"github.com/krizanauskas/winko/internal/store"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

type TgHandlerI interface {
	ProcessIncommingMessage(b *gotgbot.Bot, ctx *ext.Context) error
}

type TgHandler struct {
	userStore      store.UserStoreI
	assistantStore store.AssistantStoreI
	openAiService  service.OpenAiServiceI
	cryptoService  service.CryptoServiceI
}

func New(openAiService service.OpenAiServiceI, cryptoService service.CryptoServiceI) *TgHandler {
	return &TgHandler{
		userStore:      store.NewUserStore(),
		assistantStore: store.NewAssistantStore(),
		openAiService:  openAiService,
		cryptoService:  cryptoService,
	}
}

func (h TgHandler) ProcessIncommingMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	chatUser := ctx.EffectiveSender.User

	if chatUser == nil {
		return errors.New("chat user is nil")
	}

	if chatUser.IsBot {
		return errors.New("chat user is bot")
	}

	userId := chatUser.Id

	user, err := h.userStore.FindUserByTgId(userId)
	if err != nil {
		return errors.Wrap(err, "failed to find user by tg id")
	}

	assistant, err := h.assistantStore.GetAssistant()
	if err != nil {
		// TODO: create new assistant
		return errors.Wrap(err, "failed to get assistant")
	}

	var run openai.Run

	if user.LastThreadId == nil {
		// TODO create new thread and send message to openai
	} else {
		err := h.openAiService.SendMessage(*user.LastThreadId, ctx.EffectiveMessage.Text)
		if err != nil {
			if openapiErr, ok := err.(*openai.APIError); ok {
				if openapiErr.HTTPStatusCode == 400 {
					h.openAiService.CancelActiveRun(context.Background(), *user.LastThreadId)
				} else {
					return errors.Wrap(err, "failed to send message to openai")
				}
			}
		}

		run, err = h.openAiService.CreateRun(*user.LastThreadId, assistant.Id)
		if err != nil {
			return errors.Wrap(err, "failed to create run")
		}
	}

	var responseMessages []string

	messages, functionCalls, err := h.procesAssistantResponse(run)
	if err != nil {
		return errors.Wrap(err, "failed to process assistant response")
	}

	if len(functionCalls) > 0 {
		toolOutputs := make([]openai.ToolOutput, 0)

		for _, functionCall := range functionCalls {
			switch functionCall.Name {
			case service.GetBnbAllocationFunctionName:
				if user.CryptoWalletAddress == nil {
					toolOutputs = append(toolOutputs, openai.ToolOutput{
						ToolCallID: functionCall.Id,
						Output:     "bnb wallet address is not set",
					})
					continue
				}

				bnbAllocation, err := h.cryptoService.GetBnbAllocation(*user.CryptoWalletAddress)
				if err != nil {
					toolOutputs = append(toolOutputs, openai.ToolOutput{
						ToolCallID: functionCall.Id,
						Output:     "failed to get bnb allocation",
					})
					continue
				}

				toolOutputs = append(toolOutputs, openai.ToolOutput{
					ToolCallID: functionCall.Id,
					Output:     bnbAllocation,
				})
				continue
			case service.SendBnbToAddressFunctionName:
				params := service.SendBnbToAddressParams{}
				err = json.Unmarshal([]byte(functionCall.Arguments), &params)
				if err != nil {
					toolOutputs = append(toolOutputs, openai.ToolOutput{
						ToolCallID: functionCall.Id,
						Output:     "failed to parse function arguments",
					})
					continue
				}

				if user.CryptoWalletAddress != nil && params.Address == *user.CryptoWalletAddress {
					toolOutputs = append(toolOutputs, openai.ToolOutput{
						ToolCallID: functionCall.Id,
						Output:     "sending bnb to the same address is not allowed",
					})
					continue
				}

				response, err := h.cryptoService.SendBnbToAddress(params.Address, params.Amount)
				if err != nil {
					toolOutputs = append(toolOutputs, openai.ToolOutput{
						ToolCallID: functionCall.Id,
						Output:     "failed to send bnb to address",
					})
					continue
				}

				toolOutputs = append(toolOutputs, openai.ToolOutput{
					ToolCallID: functionCall.Id,
					Output:     response,
				})
				continue
			}
		}

		if len(toolOutputs) > 0 {
			run, err = h.openAiService.SubmitToolOutputs(context.Background(), *user.LastThreadId, run.ID, toolOutputs)
			if err != nil {
				return errors.Wrap(err, "failed to submit tool outputs")
			}

			messages, functionCalls, err = h.procesAssistantResponse(run)
			if err != nil {
				return errors.Wrap(err, "failed to process assistant response after submitting tool outputs")
			}
		}
	}

	for _, message := range messages {
		if message.Role == openai.ChatMessageRoleAssistant {
			if message.Content[0].Type == "text" {
				responseMessages = append(responseMessages, message.Content[0].Text.Value)
			}
		}

		if message.Role == openai.ChatMessageRoleUser {
			break
		}
	}

	if len(responseMessages) > 0 {
		for _, message := range responseMessages {
			_, err = ctx.EffectiveMessage.Reply(b, message, nil)
			if err != nil {
				return errors.Wrap(err, "failed to echo message")
			}
		}
	}

	return nil
}

func (h TgHandler) procesAssistantResponse(run openai.Run) ([]openai.Message, []service.OpenAiFunctionCall, error) {
	var err error
	functionCalls := make([]service.OpenAiFunctionCall, 0)
	// temp location for this code
	retries := 0
	for {
		if run.Status == openai.RunStatusRequiresAction {
			log.Printf("Thread status: %s", run.Status)

			if run.RequiredAction != nil && run.RequiredAction.SubmitToolOutputs != nil && len(run.RequiredAction.SubmitToolOutputs.ToolCalls) > 0 {
				for _, toolCall := range run.RequiredAction.SubmitToolOutputs.ToolCalls {
					if string(toolCall.Type) == string(openai.AssistantToolTypeFunction) {
						functionCalls = append(functionCalls, service.OpenAiFunctionCall{
							Id:        toolCall.ID,
							Name:      toolCall.Function.Name,
							Arguments: toolCall.Function.Arguments,
						})
					}
				}

				return nil, functionCalls, nil
			}
		}

		if run.Status == openai.RunStatusCompleted {
			messages, err := h.openAiService.ListMessage(context.Background(), run.ThreadID)
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to list messages")
			}

			return messages, nil, nil
		}

		if retries < 5 {
			log.Printf("Thread status: %s", run.Status)
			retries++
			time.Sleep(3 * time.Second)

			run, err = h.openAiService.RetrieveRun(run.ThreadID, run.ID)
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to retrieve run")
			}

		} else {
			break
		}
	}

	return nil, nil, errors.Wrap(err, "thread status is not completed after 5 retries")
}
