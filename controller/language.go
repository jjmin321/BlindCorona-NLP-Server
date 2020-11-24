package controller

import (
	"context"
	"log"
	"strings"
	"time"

	language "cloud.google.com/go/language/apiv1"
	"github.com/labstack/echo"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Text is parameter form struct
type Text struct {
	Text string `json:"text" form:"text" query:"text"`
}

// Analyze method return analyzed result of text
func Analyze(c echo.Context) error {
	location := "전국"
	var date []string
	ctx := context.Background()
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	u := new(Text)
	if err := c.Bind(u); err != nil {
		return err
	}

	entities, err := client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: u.Text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	for _, v := range entities.Entities {
		if v.Type.String() == "LOCATION" {
			location = v.Name
		} else if v.Type.String() == "DATE" {
			date = append(date, v.Name)
		} else if strings.Contains(u.Text, "어제") {
			now := time.Now()
			convHours, _ := time.ParseDuration("24h")
			custom := now.Add(-convHours).Format("2006년 01월 02일")
			date = append(date, custom)
			u.Text = strings.Replace(u.Text, "어제", "", 1)
		} else if strings.Contains(u.Text, "오늘") {
			now := time.Now()
			custom := now.Format("2006년 01월 02일")
			date = append(date, custom)
			u.Text = strings.Replace(u.Text, "오늘", "", 1)
		}
	}

	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
		return c.JSON(500, map[string]interface{}{
			"status":  500,
			"message": "텍스트 분석 중에 오류가 생겼습니다",
		})
	}

	return c.JSON(200, map[string]interface{}{
		"status":   200,
		"Location": location,
		"Date":     date,
		// "Entities": entities.Entities,
	})
}
