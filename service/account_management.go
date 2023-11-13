package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mohamedabdifitah/processor/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ResetPassword(email string) {
	generate := utils.GenerateNumbers(6)
	err := RedisClient.Set(Ctx, "otp"+":"+email, generate, 4*time.Minute).Err()
	if err != nil {
		fmt.Println(err)
	}
	templates := utils.CurrentTemplates()
	templates.LoadTemplates("assets/json/template.json", "")
	message, err := templates.TempelateInjector(
		"ResetPassword",
		map[string]string{
			"ExpireTime": "4",
			"Otp":        generate,
			"Unit":       "minutes",
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	err = SendEmailNotification(message, "reset your password", email)
	if err != nil {
		fmt.Println(err)
	}
}
func HandleResetPassword(data amqp.Delivery) {
	var info map[string]string = make(map[string]string)
	err := json.Unmarshal(data.Body, &info)
	if err != nil {
		fmt.Println(err)
	}
	email, ok := info["email"]
	if !ok {
		return
	}
	ResetPassword(email)
}
func HandleVerification(data amqp.Delivery) {
	var info map[string]string = make(map[string]string)
	err := json.Unmarshal(data.Body, &info)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
	Type, ok := info["type"]
	if !ok {
		return
	}
	rec, ok := info["rec"]
	if !ok {
		return
	}
	if Type == "email" {
		generate := utils.GenerateNumbers(6)
		err = RedisClient.Set(Ctx, "otp"+":"+rec, generate, 20*time.Minute).Err()
		if err != nil {
			fmt.Println(err)
		}
		templates := utils.CurrentTemplates()
		templates.LoadTemplates("assets/json/template.json", "")
		message, err := templates.TempelateInjector(
			"OtpTemplate",
			map[string]string{
				"ExpireTime": "20",
				"Otp":        generate,
				"Unit":       "minutes",
			},
		)
		if err != nil {
			fmt.Println(err)
		}
		err = SendEmailNotification(message, "please verify your email address", rec)
		if err != nil {
			fmt.Println(err)
		}
	} else if Type == "phone" {
		generate := utils.GenerateNumbers(6)
		err = RedisClient.Set(Ctx, "otp"+":"+rec, generate, 2*time.Minute).Err()
		if err != nil {
			fmt.Println(err)
		}
		message := "your somfood verification code is:" + generate
		SendNotifictionSms(message, rec)
	}
}
