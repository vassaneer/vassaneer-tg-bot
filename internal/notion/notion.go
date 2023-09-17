package notion

import (
	"log/slog"
	"net/http"
)

type Notion struct {
	http       *http.Client
	DatabaseID string
	Token      string
	logger     *slog.Logger
}

func NewNotion(databaseId string, token string, logger *slog.Logger) *Notion {
	logger.Info("Notion client created")
	return &Notion{
		http:       &http.Client{},
		DatabaseID: databaseId,
		Token:      token,
		logger:     logger,
	}
}

type Parent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
}

type NotionPostPayload struct {
	Parent     Parent                 `json:"parent"`
	Properties map[string]interface{} `json:"properties"`
}

func NewTitleProperty(title string) map[string]interface{} {
	return map[string]interface{}{
		"type": "title",
		"title": []map[string]interface{}{
			{
				"type": "text",
				"text": map[string]interface{}{
					"content": title,
				},
			},
		},
	}
}

func NewNumberProperty(amount float64) map[string]interface{} {
	return map[string]interface{}{
		"type":   "number",
		"number": amount,
	}
}

func NewUrlProperty(s string) map[string]interface{} {
	return map[string]interface{}{
		"type": "url",
		"url":  s,
	}
}

func NewSelectProperty(category string) map[string]interface{} {
	return map[string]interface{}{
		"type": "select",
		"select": map[string]interface{}{
			"name": category,
		},
	}
}

func NewRelationProperty(database string, fieldName string, fieldId string) map[string]interface{} {
	return map[string]interface{}{
		"type": "relation",
		"select": map[string]interface{}{
			"database_id":          database,
			"synced_property_name": fieldName,
			"synced_property_id":   fieldId,
		},
	}
}
