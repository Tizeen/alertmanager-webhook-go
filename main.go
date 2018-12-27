package main

import (
	"net/http"

	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type labels struct {
	AlertName string `json:"alertname"`
	Instance  string `json:"instance"`
	Job       string `json:"job"`
	Servity   string `json:"servity"`
}

type annotations struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

type alerts struct {
	Status      string      `json:"status"`
	Labels      labels      `json:"labels"`
	Annotations annotations `json:"annotations"`
	StartsAt    string      `json:"startsAt"`
	EndsAt      string      `json:"endsAt"`
}

type alertMessage struct {
	Alert []alerts `json:"alerts"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/alert", func(c *gin.Context) {
		AlertMess := alertMessage{}
		c.BindJSON(&AlertMess)

		var alertMessage string

		for i := 0; i < len(AlertMess.Alert); i++ {
			alertMessage = AlertMess.Alert[i].Labels.AlertName + "\n" +
				AlertMess.Alert[i].Labels.Instance + "\n" +
				AlertMess.Alert[i].Labels.Job + "\n" +
				AlertMess.Alert[i].Annotations.Description + "\n" +
				AlertMess.Alert[i].Annotations.Summary
		}

		log.Println(alertMessage)

		pushResp, err := sendMessage("警报来了~~~", alertMessage)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(pushResp)
		c.JSON(http.StatusOK, gin.H{"code": 0})

	})

	return r
}

func sendMessage(text string, desp string) (pushResp string, err error) {

	uri := "https://pushbear.ftqq.com/sub"

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	writer.WriteField("sendkey", "7637-92cf757194444ca97236c686510a4126")
	writer.WriteField("text", text)
	writer.WriteField("desp", desp)

	req, _ := http.NewRequest(http.MethodPost, uri, buf)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return pushResp, err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	pushResp = string(data)

	return pushResp, nil
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:9999
	r.Run(":9999")
}
