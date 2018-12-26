package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	)

var db = make(map[string]string)

type labels struct {
	AlertName	string	`json:"alertname"`
	Instance 	string	`json:"instance"`
	Job 		string	`json:"job"`
	Servity 	string	`json:"servity"`
}

type annotations struct {
	Description 	string 	`json:"description"`
	Summary 		string 	`json:"summary"`
}

type alerts struct {
	Status 			string			`json:"status"`
	Labels 			labels 			`json:"labels"`
	Annotations 	annotations 	`json:"annotations"`
	StartsAt 		string 			`json:"startsAt"`
	EndsAt 			string			`json:"endsAt"`
}


type alertMessage struct {
	Alert 	[]alerts 	`json:"alerts"`
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
		var Status string = AlertMess.Alert[0].Status
		c.JSON(http.StatusOK, gin.H{"status": Status})
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:9999
	r.Run(":9999")
}
