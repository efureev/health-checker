package checkers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/efureev/health-checker"
)

type ElasticsearchConfig struct {
	Host string
	Port int
}

type EsVersion struct {
	Number        string `json:"number"`
	BuildType     string `json:"build_type"`
	LuceneVersion string `json:"lucene_version"`
}
type EsResponse struct {
	Name    string    `json:"name"`
	Version EsVersion `json:"version"`
}

func (r EsResponse) String() string {
	return fmt.Sprintf(`%s (v%s), for %s, lucene: v%s`, r.Name, r.Version.Number, r.Version.BuildType, r.Version.LuceneVersion)
}

func CheckingElasticsearchFn(config ElasticsearchConfig) checker.CmdFn {
	return func(result *checker.Result, log checker.ILogger) {
		url := fmt.Sprintf(`%s:%d`, config.Host, config.Port)
		resp, err := http.Get(url)
		step := result.AddStep(fmt.Sprintf(`Try to connect to %s`, url))
		if err != nil {
			step.Failed(err)
			result.AddError(err)
			return
		}
		defer resp.Body.Close()
		step.Success()

		step = result.AddStep(`Try to decode a response`)
		var cResp EsResponse
		if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
			step.Failed(err)
			result.AddError(err)
			return
		}
		step.Success()
		result.Info.Version = cResp.String()

		return
	}
}
