package controller

import (
	"context"
	"fmt"
	"main/internal/model"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	baseApiUrl   = "https://openrouter.ai/api/v1"
	currentModel = openai.ChatModel("openai/gpt-oss-20b:free")
)

type LlmController struct {
	client *openai.Client
}

func NewLlmController() (*LlmController, error) {
	apiKey, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return nil, fmt.Errorf("environment variable OPENAI_API_KEY required")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseApiUrl),
	)

	return &LlmController{client: &client}, nil
}

func (lc *LlmController) SummarizeChat(ctx context.Context, chat model.Chat) (string, error) {
	data, err := os.ReadFile(chat.Filepath)
	if err != nil {
		return "", err
	}

	chatCompletion, err := lc.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: currentModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
				Ты — аналитик переписок. 
				Твоя задача — кратко и чётко выделять обсуждения и темы.
				Пиши лаконично, без воды, без специальных символов, эмодзи и лишних фраз.
				Стиль деловой, структурированный, с акцентом на ключевые моменты.
				Отвечай в формате текста, если хочешь как-то структурировать вывод, то используй нумерацию.
			`),
			openai.UserMessage(fmt.Sprintf("Вот переписка для анализа:\n\n%s", string(data))),
			openai.UserMessage(`
				Составь краткое и чёткое резюме всех тем обсуждений,
				которые происходили в чате.

				Требования:
				- представь ответ в виде списка;
				- каждый пункт должен описывать одну ключевую тему обсуждения;
				- пиши коротко, ёмко, без воды;
				- выделяй наиболее важные моменты и выводы;
				- не используй вводные фразы типа "Понял" или "Вот анализ".

				Пиши только анализ.`,
			),
		},
	})

	if err != nil || len(chatCompletion.Choices) == 0 {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (lc *LlmController) DescribePersonality(ctx context.Context, chat model.Chat) (string, error) {
	data, err := os.ReadFile(chat.Filepath)
	if err != nil {
		return "", err
	}

	chatCompletion, err := lc.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: currentModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
				Ты – аналитик переписок. 
				Твоя задача – объективно анализировать поведение людей на основе чата. 
				Пиши кратко, строго по делу, без лишних фраз вроде “Понял”, “Конечно” и т.п.
				Не используй спецсимволы, эмодзи и украшения.
				Отвечай в формате текста, если хочешь как-то структурировать вывод, то используй нумерацию.
			`),
			openai.UserMessage(fmt.Sprintf("Вот переписка для анализа:\n\n%s", string(data))),
			openai.UserMessage(`
				На основе этой переписки опиши личности всех участников. 
				Для каждого участника укажи:

				- Имя участника чата
				- Предполагаемый возраст
				- Предполагаемый характер
				- Предполагаемые привычки и предпочтения
				- Особенности поведения

				Каждое предположение обязательно подкрепляй конкретной цитатой или фактом из чата.
				Пиши кратко, по существу, без воды.
			`),
		},
	})

	if err != nil || len(chatCompletion.Choices) == 0 {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (lc *LlmController) MeetingSearch(ctx context.Context, chat model.Chat) (string, error) {
	data, err := os.ReadFile(chat.Filepath)
	if err != nil {
		return "", err
	}

	chatCompletion, err := lc.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: currentModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
				Ты — аналитик переписок.
				Твоя задача — находить факты о встречах и событиях на основе контекста.
				Отвечай кратко, структурировано и строго по делу.
				Не используй спецсимволы, эмодзи, вводные фразы или воду.
				Отвечай в формате текста, если хочешь как-то структурировать вывод, то используй нумерацию.
			`),
			openai.UserMessage(fmt.Sprintf("Вот переписка для анализа:\n\n%s", string(data))),
			openai.UserMessage(`
				Найди в переписке все упоминания встреч — как прошедших, так и запланированных.
				Представь ответ в хронологическом порядке.

				Для каждой встречи укажи:
				- Название встречи (если нет явного названия, придумай краткое и точное по смыслу)
				- Дата встречи (ориентируйся на явные либо косвенные временные указания)
				- Контекст, который подтверждает, что встреча состоялась или планируется
				(краткая цитата или пересказ из чата)

				Пиши только анализ, без лишних слов и формальных фраз.
			`),
		},
	})

	if err != nil || len(chatCompletion.Choices) == 0 {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (lc *LlmController) ContextSearch(ctx context.Context, chat model.Chat, description string) (string, error) {
	data, err := os.ReadFile(chat.Filepath)
	if err != nil {
		return "", err
	}

	chatCompletion, err := lc.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: currentModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`
				Ты — аналитик переписок.
				Твоя задача — выявлять фрагменты чата, относящиеся к описанному событию.
				Пиши строго по делу, кратко, без воды и без специальных символов, эмодзи.
				Не используй вводные фразы вроде "Понял" или "Я считаю".
				Только анализ и факты из переписки.
				Отвечай в формате текста, если хочешь как-то структурировать вывод, то используй нумерацию.
			`),
			openai.UserMessage(fmt.Sprintf("Вот переписка для анализа:\n\n%s", string(data))),
			openai.UserMessage(fmt.Sprintf(`
				Найди в переписке часть диалога, которая соответствует следующему описанию:
				"%s"
				Правила:
				- если контекст полностью отсутствует — напиши, что подходящего фрагмента переписки не найдено, ничего не придумывай;
				- если событие или похожая ситуация есть — кратко перескажи соответствующий фрагмент общения;
				- обязательно укажи дату переписки, если она присутствует;
				- пиши лаконично, только анализ.

				Формат ответа при наличии подходящего контекста:
				- Дата переписки
				- Краткий пересказ подходящей части чата
			`, description)),
		},
	})

	if err != nil || len(chatCompletion.Choices) == 0 {
		return "", err
	}

	return chatCompletion.Choices[0].Message.Content, nil
}
