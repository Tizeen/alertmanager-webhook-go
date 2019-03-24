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
	sendKey = "xxxx-92cf757194444ca97236c686510a4126"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

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
		log.Printf(err.Error())
	}

	var alertString string

	for i := 0; i < len(alertMess.Alert); i++ {
		log.Printf("告警内容：" + alertMess.Alert[i].Status + " " +
			alertMess.Alert[i].Labels.AlertName + " " +
			alertMess.Alert[i].Labels.Job + " " +
			alertMess.Alert[i].Labels.Instance + " " +
			alertMess.Alert[i].Annotations.Summary,
		)
		alertString = fmt.Sprintf("%s \n\n **%s** \n\n %s \n\n %s \n\n %s \n\n %s \n\n %s", alertString,
			alertMess.Alert[i].Labels.AlertName,
			alertMess.Alert[i].Status,
			alertMess.Alert[i].Labels.Job,
			alertMess.Alert[i].Labels.Instance,
			alertMess.Alert[i].Annotations.Description,
			alertMess.Alert[i].Annotations.Summary)
	}

	title := "来告警信息了~~~"

	if err := sendMessage(title, alertString); err != nil {
		log.Printf(err.Error() + "\n")
	} else {
		log.Printf("告警信息推送成功~~~\n")
	}

}

// 传入发送标题和内容，返回bool类型和err信息，为true表示发送成功，false表示发送失败
func sendMessage(title, message string) (err error) {

	postUrl := "https://pushbear.ftqq.com/sub"

	parameters := url.Values{}
	parameters.Add("sendkey", sendKey)
	parameters.Add("text", title)
	parameters.Add("desp", message)

	resp, err := http.PostForm(postUrl, parameters)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respBD := &RespBodyStruct{}

	err = json.Unmarshal(body, respBD)
	if err != nil {
		return err
	}

	if respBD.Code != 0 {
		return errors.New(respBD.Message)
	}

	return nil

}

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/alert", handler)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
