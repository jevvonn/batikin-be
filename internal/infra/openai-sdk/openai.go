package openaisdk

import (
	"batikin-be/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewOpenAIClient() openai.Client {
	conf := config.Load()

	return openai.NewClient(
		option.WithAPIKey(conf.OPENAIAPIKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)
}
