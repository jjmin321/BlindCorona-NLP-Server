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
	var date string
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

	if strings.ContainsAny(u.Text, "어제") {
		now := time.Now()
		convHours, _ := time.ParseDuration("24h")
		custom := now.Add(-convHours).Format("20060102")
		date = custom
		u.Text = strings.Replace(u.Text, "어제", "", 1)
	} else if strings.ContainsAny(u.Text, "오늘") {
		now := time.Now()
		custom := now.Format("20060102")
		date = custom
		u.Text = strings.Replace(u.Text, "오늘", "", 1)
	}

	log.Print(u.Text)
	for _, v := range entities.Entities {
		if v.Type.String() == "LOCATION" {
			location = v.Name
		} else if v.Type.String() == "DATE" {
			var year, month, day string
			s1 := strings.Split(v.Name, "년 ")
			year = s1[0]
			s2 := strings.Split(s1[1], "월 ")
			if len(s2[0]) == 1 {
				month = "0" + s2[0]
			} else {
				month = s2[0]
			}
			s3 := strings.Split(s2[1], "일")
			if len(s3[0]) == 1 {
				day = "0" + s3[0]
			} else {
				day = s3[0]
			}
			date = year + month + day
		}
	}

	return c.JSON(200, map[string]interface{}{
		"status":   200,
		"Location": location,
		"Date":     date,
	})
}
