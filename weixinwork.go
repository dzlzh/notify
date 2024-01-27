package notify

import (
	"encoding/json"
	"errors"

	"github.com/dzlzh/httpc"
	"github.com/tidwall/gjson"
)

type WeiXinWork struct {
	corpid     string // 企业ID
	agentid    string // 应用 ID
	corpsecret string // 应用 Secret
}

func NewWeiXinWork(corpid, agentid, corpsecret string) *WeiXinWork {
	return &WeiXinWork{
		corpid:     corpid,
		agentid:    agentid,
		corpsecret: corpsecret,
	}
}

func (w *WeiXinWork) accessToken() (string, error) {
	request := httpc.NewRequest(httpc.NewClient())
	request.SetMethod("GET").SetURL("https://qyapi.weixin.qq.com/cgi-bin/gettoken")
	request.SetQuery("corpid", w.corpid)
	request.SetQuery("corpsecret", w.corpsecret)
	request.Send()
	_, res, err := request.End()
	if err != nil || !gjson.ValidBytes(res) {
		return "", err
	}
	r := gjson.ParseBytes(res)
	if r.Get("errcode").Bool() {
		return "", errors.New(r.Get("errmsg").String())
	}
	return r.Get("access_token").String(), nil
}

func (w *WeiXinWork) Send(subject, message string) error {
	j, err := json.Marshal(struct {
		Touser  string `json:"touser"`
		Msgtype string `json:"msgtype"`
		Agentid string `json:"agentid"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		Touser:  "@all",
		Msgtype: "text",
		Agentid: w.agentid,
		Text: struct {
			Content string `json:"content"`
		}{
			Content: message,
		},
	})
	if err != nil {
		return err
	}

	accessToken, err := w.accessToken()
	if err != nil {
		return err
	}
	request := httpc.NewRequest(httpc.NewClient())
	request.SetMethod("POST").SetURL("https://qyapi.weixin.qq.com/cgi-bin/message/send")
	request.SetQuery("access_token", accessToken)
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
