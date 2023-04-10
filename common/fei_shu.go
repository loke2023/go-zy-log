package common

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FSMsgContent struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type FSMessage struct {
	Msg_type string       		`json:"msg_type"`
	Content struct{
		Post struct{
			Zh_cn 	 struct{
				Title 	 string 	  		`json:"title"`
				Content  [][]FSMsgContent 	`json:"content"`
			} 							`json:"zh_cn"`
		}`json:"post"`
	}`json:"content"`
}

type FSCardelementsText struct {
	Content  string	`json:"content"`
	Tag 	 string	`json:"tag"`
}

type FSCardelements struct {
	Tag 	string `json:"tag"`
	Text 	FSCardelementsText `json:"text,omitempty"`
}

type FSCardMessage struct {
	Msg_type string       		`json:"msg_type"`
	Card struct{
		Config 	struct{
			Wide_screen_mode 	bool 	`json:"wide_screen_mode"`
		} `json:"config"`
		Elements []FSCardelements `json:"elements"`
		Header struct{
			Template 	string `json:"template"`
			Title 		struct{
				Content	string `json:"content"`
				Tag 	string `json:"tag"`
			}`json:"title"`
		}`json:"header"`
	}`json:"card"`
}

func (c *FSCardMessage)Init(title string)  {
	c.Msg_type = "interactive"
	c.Card.Config.Wide_screen_mode = true
	c.Card.Header.Template = "blue"
	c.Card.Header.Title.Tag = "plain_text"
	c.Card.Header.Title.Content = title
}


func (c *FSCardMessage)AddElements(txt string)  {
	if txt == "-" {
		c.Card.Elements = append(c.Card.Elements, FSCardelements{
			Tag:"hr",
		})
	}else{
		c.Card.Elements = append(c.Card.Elements, FSCardelements{
			Tag: "div",
			Text: FSCardelementsText{Content: txt, Tag: "lark_md"},
		})
	}
}

func NotifyStructAdMessageToFeiShu(msg interface{}, url string) {
	jsongMsg, _ := json.Marshal(msg)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsongMsg))
	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
		}
	}
}

func SendMessageToFeiShu(msg string, lg *logrus.Logger) {
	type FSMsgContent struct {
		Text string `json:"text"`
	}

	type FSMessage struct {
		Msg_type string       `json:"msg_type"`
		Content  FSMsgContent `json:"content"`
	}

	var fsMsg FSMessage
	fsMsg.Msg_type = "text"
	fsMsg.Content.Text = msg
	jsongMsg, _ := json.Marshal(fsMsg)
	req, err := http.NewRequest("POST",
		"https://open.feishu.cn/open-apis/bot/v2/hook/5ed2111f-5f19-4d65-a3ff-8de40a6fed55",
		bytes.NewBuffer(jsongMsg))
	if err == nil {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				if lg != nil {
					lg.Warn("send message to feishu faild, ret:", resp.StatusCode)
				}
			}
		} else {
			if lg != nil {
				lg.Warn("send message to feishu client.Do faild, err:", err.Error())
			}
		}
	} else {
		if lg != nil {
			lg.Warn("send message to feishu NewRequest faild, err:", err.Error())
		}
	}
}
