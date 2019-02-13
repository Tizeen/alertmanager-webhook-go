package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

type alertMessages struct {
	Alert []alerts `json:"alerts"`
}

type RespBodyStruct struct {
	Code    int    `json:"code, string"`
	Message string `json:"message"`
	// omitempty：在序列化的时候忽略零值或者空值
	Data    string `json:"data, omitempty"`
	Created string `json:"created"`
}

var (
	logFile = "/var/log/alertwebhook.log"
	sendKey = "xxxx-12ec757194444ca97236c686510a4126"
)

func handler(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("err opening file: %v\n", err)
	}
	defer f.Close()

	log.SetOutput(f)

	decoder := json.NewDecoder(r.Body)
	alertMess := alertMessages{}

	err = decoder.Decode(&alertMess)
	if err != nil {
		panic(err)
	}

	var alertString string

	for i := 0; i < len(alertMess.Alert); i++ {
		log.Printf("告警内容：" + alertMess.Alert[i].Status + " " +
			alertMess.Alert[i].Labels.AlertName + " " +
			alertMess.Alert[i].Labels.Job + " " +
			alertMess.Alert[i].Labels.Instance + " " +
			alertMess.Alert[i].Annotations.Summary,
		)
		alertString = alertString + "\n\n" + "**" + alertMess.Alert[i].Labels.AlertName + "**" + "\n\n" +
			alertMess.Alert[i].Status + "\n\n" +
			alertMess.Alert[i].Labels.Job + "\n\n" +
			alertMess.Alert[i].Labels.Instance + "\n\n" +
			alertMess.Alert[i].Annotations.Description + "\n\n" +
			alertMess.Alert[i].Annotations.Summary
	}

	title := "来告警信息了"

	status, err := sendMessage(title, alertString)
	if err != nil {
		log.Printf(err.Error() + "\n")
	}
	if status == true {
		log.Printf("告警信息推送成功\n")
	}

}

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

// 传入发送标题和内容，返回bool类型和err信息，为true表示发送成功，false表示发送失败
func sendMessage(title, message string) (status bool, err error) {

	postUrl := "https://pushbear.ftqq.com/sub"

	parameters := url.Values{}
	parameters.Add("sendkey", sendKey)
	parameters.Add("text", title)
	parameters.Add("desp", message)

	resp, err := http.PostForm(postUrl, parameters)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	respBD := &RespBodyStruct{}

	err = json.Unmarshal(body, respBD)
	if err != nil {
		return false, err
	}

	if respBD.Code != 0 {
		return false, errors.New(respBD.Message)
	}

	return true, nil

}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/alert", handler)
	log.Printf("Start to listening the incoming requests on http address: %s", ":9999")
	log.Fatal(http.ListenAndServe(":9999", nil))
}
