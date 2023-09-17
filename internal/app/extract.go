package app

import (
	"log/slog"
	"strconv"
	"strings"
)

func ExpenseCommandExtract(subMatchs []string, srv *Service) map[string]Field {
	if len(subMatchs) < 3 {
		srv.logger.Info("Cannot extract amount, category and title from command message",
			slog.String("subMathches", strings.Join(subMatchs, ",")))
		return map[string]Field{}
	}
	amountStr := subMatchs[1]
	category := getCategory(subMatchs[2])
	title := ""
	if len(subMatchs) > 3 && subMatchs[3] != "" {
		title = subMatchs[3]
	}

	amount, _ := strconv.ParseFloat(amountStr, 64)

	return map[string]Field{"amount": Field{DataType: Number, Value: amount}, "title": Field{DataType: Title, Value: title}, "category": Field{DataType: Select, Value: category}}
}

func RepCommandExtract(subMatchs []string, srv *Service) map[string]Field {
	if len(subMatchs) < 3 {
		srv.logger.Info("Cannot extract name, weight and rep from command message",
			slog.String("subMathches", strings.Join(subMatchs, ",")))
		return map[string]Field{}
	}
	Name := ""
	if len(subMatchs) > 3 && subMatchs[3] != "" {
		Name = subMatchs[1]
	}

	Weight, _ := strconv.ParseFloat(subMatchs[2], 64)
	Rep, _ := strconv.ParseFloat(subMatchs[3], 64)

	return map[string]Field{"Name": Field{DataType: Title, Value: Name}, "Weight": Field{DataType: Number, Value: Weight}, "Rep": Field{DataType: Number, Value: Rep}}
}
