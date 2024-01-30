package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/henrylee2cn/goutil/calendar/cron"
)

var client *resty.Client
var rToken string

func init() {
	client = resty.New().SetRetryCount(3).SetRetryWaitTime(3 * time.Second)

	rToken = os.Getenv("refresh-token")
	if rToken == "" {
		panic("refresh-token is empty")
	}
}

func main() {
	c := cron.New()

	c.AddFunc("0 0 2 * * *", func() {
		accessData, accessErr := getAccessToken()
		if accessErr != nil {
			printWithTime("getAccessToken error: %s", accessErr.Error())
			return
		} else {
			printWithTime("getAccessToken success: %s", accessData.UserName)
		}

		bToken := "Bearer " + accessData.AccessToken
		count, signinErr := signin(bToken)
		if signinErr != nil {
			printWithTime("signin error: %s", signinErr.Error())
			return
		} else {
			printWithTime("signin success: %d days signed in", count)
		}

		receiveData, receiveErr := receive(bToken, count)
		if receiveErr != nil {
			printWithTime("receive error: %s", receiveErr.Error())
			return
		} else {
			printWithTime("receive success: %s", receiveData.Notice)
		}
	})

	c.Start()
	printWithTime("timed task started, enjoy it!")

	select {}
}

type AccessTokenReq struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	UserName    string `json:"user_name"`
}

func getAccessToken() (*AccessTokenResp, error) {
	req := &AccessTokenReq{
		GrantType:    "refresh_token",
		RefreshToken: rToken,
	}
	resp := &AccessTokenResp{}
	res, err := client.R().SetBody(req).SetResult(resp).Post("https://auth.aliyundrive.com/v2/account/token")

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("%s", res.Status())
	}

	return resp, err
}

func printWithTime(msg string, args ...interface{}) {
	fmt.Printf("%s %s\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, args...))
}

type SigninResp struct {
	Result *SigninRespResult `json:"result"`
}
type SigninRespResult struct {
	SignInCount int `json:"signInCount"`
}

func signin(bToken string) (int, error) {
	resp := &SigninResp{}
	res, err := client.R().
		SetHeader("Authorization", bToken).
		SetBody(map[string]string{"_rx-s": "mobile"}).
		SetResult(resp).
		Post("https://member.aliyundrive.com/v1/activity/sign_in_list")
	if err != nil {
		return 0, err
	}

	if res.IsError() {
		return 0, fmt.Errorf("%s", res.Status())
	}

	return resp.Result.SignInCount, nil
}

type ReceiveReq struct {
	SignInDay int `json:"signInDay"`
}
type ReceiveResp struct {
	Result *ReceiveRespResult `json:"result"`
}
type ReceiveRespResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Notice      string `json:"notice"`
}

func receive(bToken string, count int) (*ReceiveRespResult, error) {
	req := &ReceiveReq{
		SignInDay: count,
	}
	resp := &ReceiveResp{}
	res, err := client.R().
		SetHeader("Authorization", bToken).
		SetBody(req).
		SetResult(resp).
		Post("https://member.aliyundrive.com/v1/activity/sign_in_reward?_rx-s=mobile")
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("%s", res.Status())
	}

	return resp.Result, nil
}
