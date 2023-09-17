package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"vassaneer-tg-bot/internal/app"
)

func (n *Notion) Add(cmd app.ReturnCommand) error {
	props := make(map[string]interface{})
	for k, v := range cmd.Fields {
		switch v.DataType {
		case app.Number:
			{
				props[k] = NewNumberProperty(v.Value.(float64))
			}
		case app.Select:
			{
				props[k] = NewSelectProperty(v.Value.(string))
			}
		case app.Title:
			{
				title := v.Value.(string)
				if title != "" {
					props[k] = NewTitleProperty(title)
				} else {
					props[k] = NewTitleProperty(fmt.Sprintf("B spend with %s", title))
				}
			}
		}
	}

	//
	payload := NotionPostPayload{
		Parent: Parent{
			Type:       "database_id",
			DatabaseID: n.DatabaseID,
		},
		Properties: props,
	}
	payloadBytes, _ := json.Marshal(payload)
	payloadReader := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", "https://api.notion.com/v1/pages", payloadReader)
	req.Header = http.Header{
		"Authorization":  []string{"Bearer " + n.Token},
		"Content-Type":   []string{"application/json"},
		"Notion-Version": []string{"2022-06-28"},
	}
	resp, err := n.http.Do(req)
	if err != nil {
		n.logger.Error(
			"Request failed",
			slog.String("errorMessage", err.Error()),
			slog.String("requestUrl", req.URL.String()),
			slog.String("requestMethod", req.Method),
		)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		responseError := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&responseError)
		n.logger.Error(
			"Request failed",
			slog.String("requestUrl", req.URL.String()),
			slog.String("requestMethod", req.Method),
			slog.Int("statusCode", resp.StatusCode),
			slog.Any("responseBody", responseError),
		)
		return fmt.Errorf("request falied with status code %d", resp.StatusCode)
	}

	respBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&respBody)
	return nil
}
