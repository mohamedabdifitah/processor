package notification

import (
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMS struct {
	To      string
	Message string
}

func (sms SMS) Send() error {
	accountSid := os.Getenv("accountSid")
	authToken := os.Getenv("authToken")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	// client := twilio.NewRestClient()

	params := &api.CreateMessageParams{}
	params.SetBody(sms.Message)
	params.SetFrom(os.Getenv("MY_PHONE"))
	params.SetTo(sms.To)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	if resp.Sid != nil {
		fmt.Println(*resp.Sid)
	} else {
		fmt.Println(resp.Sid)
	}
	return nil
}
