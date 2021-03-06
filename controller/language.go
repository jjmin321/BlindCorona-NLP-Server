package controller

import (
	"context"
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	language "cloud.google.com/go/language/apiv1"
	"github.com/labstack/echo"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Text is parameter form struct
type Text struct {
	Text string `json:"text" form:"text" query:"text"`
}

// Set Default Value
var (
	location = "전국"
	date     string
	wg       sync.WaitGroup
)

// ResetData reset Default Value
func ResetData() {
	location = "전국"
	date = ""
}

// RecoverServer recover server
func RecoverServer() {
	err := recover()
	if err != nil {
		log.Print("Recovered", err)
		debug.PrintStack()
	}
}

// Analyze method return analyzed result of text
func Analyze(c echo.Context) error {
	ResetData()
	u := new(Text)
	if err := c.Bind(u); err != nil {
		return err
	}
	u.Text = strings.Trim(u.Text, "안녕")
	u.Text = CheckDate(u.Text)
	wg.Add(1)
	go CheckEntities(u.Text)
	if location == "전국" {
		CheckAllWord(u.Text)
	}
	wg.Wait()
	return c.JSON(200, map[string]interface{}{
		"status":   200,
		"Location": location,
		"Date":     date,
	})
}

// CheckDate method set Date if text has 어제 or 오늘
func CheckDate(Text string) string {
	if strings.ContainsAny(Text, "어제") {
		now := time.Now()
		convHours, _ := time.ParseDuration("24h")
		custom := now.Add(-convHours).Format("20060102")
		date = custom
		Text = strings.Replace(Text, "어제", "", 1)
	} else if strings.ContainsAny(Text, "오늘") {
		now := time.Now()
		custom := now.Format("20060102")
		date = custom
		Text = strings.Replace(Text, "오늘", "", 1)
	}
	return Text
}

// CheckEntities method Check Location and Date
func CheckEntities(Text string) {
	defer wg.Done()
	defer RecoverServer()
	ctx := context.Background()
	client, _ := language.NewClient(ctx)
	entities, _ := client.AnalyzeEntities(ctx, &languagepb.AnalyzeEntitiesRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: Text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
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
}

// CheckAllWord method Range All Word for search Location
func CheckAllWord(Text string) {
	word := strings.Split(Text, " ")
	for i := 0; i < len(word); i++ {
		if word[i] != "" {
			wg.Add(1)
			go CheckEntities(word[i])
		}
	}
}
