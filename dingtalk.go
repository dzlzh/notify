package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/dzlzh/httpc"
	"github.com/tidwall/gjson"
)

type Dingtalk struct {
	accessToken string
	secret      string
}

func NewDingtalk(accessToken, secret string) *Dingtalk {
	return &Dingtalk{
		accessToken: accessToken,
		secret:      secret,
	}
}

func (d *Dingtalk) sign(timestamp string) string {
	data := timestamp + "\n" + d.secret
	h := hmac.New(sha256.New, []byte(d.secret))
	h.Write([]byte(data))
	return url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func (d *Dingtalk) Send(subject, message string) error {
	j, err := json.Marshal(struct {
		Msgtype  string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
	}{
		Msgtype: "markdown",
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{
			Title: subject,
			Text:  message,
		},
	})
	if err != nil {
		return err
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sign := d.sign(timestamp)
	request := httpc.NewRequest(httpc.NewClient())
	request.SetMethod("POST").SetURL("https://oapi.dingtalk.com/robot/send")
	request.SetQuery("access_token", d.accessToken)
	request.SetQuery("timestamp", timestamp)
	request.SetQuery("sign", sign)
	request.SetJson(j)
	request.Send()
	_, res, err := request.End()
	if err != nil || !gjson.ValidBytes(res) {
		return err
	}
	r := gjson.ParseBytes(res)
	if r.Get("errcode").Bool() {
		return errors.New(r.Get("errmsg").String())
	}
	return nil
}
