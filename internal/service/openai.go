package service

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAiFunctionCall struct {
	Id        string
	Name      string
	Arguments string
}

type OpenAiServiceI interface {
	SendMessage(threadId string, message string) error
	CreateRun(threadId string, assistantId string) (openai.Run, error)
	RetrieveRun(threadId string, runId string) (openai.Run, error)
	ListMessage(ctx context.Context, threadId string) ([]openai.Message, error)
	SubmitToolOutputs(ctx context.Context, threadId string, runId string, outputs []openai.ToolOutput) (openai.Run, error)
	CancelActiveRun(ctx context.Context, threadId string) error
}

type OpenAiService struct {
	client *openai.Client
}

func NewOpenAiService(client *openai.Client) *OpenAiService {
	return &OpenAiService{
		client: client,
	}
}

func (s *OpenAiService) SendMessage(threadId string, message string) error {
	_, err := s.client.CreateMessage(context.Background(), threadId, openai.MessageRequest{
		Content: message,
		Role:    openai.ChatMessageRoleUser,
	})

	return err
}

func (s *OpenAiService) CreateRun(threadId string, assistantId string) (openai.Run, error) {
	run, err := s.client.CreateRun(context.Background(), threadId, openai.RunRequest{
		AssistantID: assistantId,
	})
	if err != nil {
		return openai.Run{}, err
	}

	return run, nil
}

func (s *OpenAiService) RetrieveRun(threadId string, runId string) (openai.Run, error) {
	run, err := s.client.RetrieveRun(context.Background(), threadId, runId)
	if err != nil {
		return openai.Run{}, err
	}

	return run, nil
}

func (s *OpenAiService) ListMessage(ctx context.Context, threadId string) ([]openai.Message, error) {
	messages, err := s.client.ListMessage(ctx, threadId, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return messages.Messages, nil
}

func (s *OpenAiService) CancelRun(ctx context.Context, threadId string, runId string) error {
	_, err := s.client.CancelRun(ctx, threadId, runId)
	if err != nil {
		return err
	}

	return nil
}

func (s *OpenAiService) SubmitToolOutputs(ctx context.Context, threadId string, runId string, outputs []openai.ToolOutput) (openai.Run, error) {
	run, err := s.client.SubmitToolOutputs(ctx, threadId, runId, openai.SubmitToolOutputsRequest{
		ToolOutputs: outputs,
	})
	if err != nil {
		return openai.Run{}, err
	}

	return run, nil
}

func (s *OpenAiService) CancelActiveRun(ctx context.Context, threadId string) error {
	run, err := s.client.ListRuns(ctx, threadId, openai.Pagination{})
	if err != nil {
		return err
	}

	if len(run.Runs) > 0 {
		for _, run := range run.Runs {
			if run.Status == openai.RunStatusRequiresAction {
				err = s.CancelRun(ctx, threadId, run.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
