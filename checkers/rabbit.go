package checkers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/efureev/health-checker"
)

type RabbitResponse struct {
	Name    string `json:"product_name"`
	Version string `json:"product_version"`
}

type RabbitConfig struct {
	Host     string
	Username string
	Password string
	Port     int
}

func CheckingRabbitFn(config RabbitConfig) checker.CmdFn {
	return func(result *checker.Result, log checker.ILogger) {
		url := fmt.Sprintf(`%s:%d/api/overview`, config.Host, config.Port)
		step := result.AddStep(`Try to connect to ` + url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			step.Failed(err)
			result.AddError(err)
			return
		}
		step.Success()
		req.Header.Add("Authorization", "Basic "+basicAuth(config.Username, config.Password))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			result.AddError(err)
			return
		}
		defer resp.Body.Close()

		step = result.AddStep(`Try to decode a response`)
		var cResp RabbitResponse
		if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
			step.Failed(err)
			result.AddError(err)
			return
		}
		step.Success()
		result.Info.Version = fmt.Sprintf(`%s v%s`, cResp.Name, cResp.Version)

		return
	}
}
