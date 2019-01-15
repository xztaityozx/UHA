package cmd

import (
	"fmt"
	"time"

	"github.com/nlopes/slack"
)

type SlackConfig struct {
	Channel string
	User    string
	Token   string
}

type SlackMessage struct {
	Succsess     int
	Failed       int
	StartTime    time.Time
	FinishedTime time.Time
	SubMessage   string
}

var SlackNoNotify bool

func PostFailed(config SlackConfig, err error) error {
	if SlackNoNotify {
		return nil
	}
	api := slack.New(config.Token)
	text := fmt.Sprintf("<@%s> シミュレーションに失敗しました・・・ :face_with_rolling_eyes:\n```\n%v\n```", config.User, err)
	param := slack.PostMessageParameters{
		AsUser: true,
	}

	_, _, xx := api.PostMessage(config.Channel, slack.MsgOptionText(text, true), slack.MsgOptionPostMessageParameters(param))
	return xx
}

func Post(config SlackConfig, message SlackMessage) error {
	if SlackNoNotify {
		return nil
	}
	api := slack.New(config.Token)
	text := fmt.Sprintf("<@%s> こちらはUHA botです。シミュレーションが終わりましたよ。\n%s\n:seikou: %d 個 :sippai: %d 個\n開始時間 : %s\n終了時間 : %s", config.User, message.SubMessage, message.Succsess, message.Failed, message.StartTime.Format("2006/01/02/15:04.05"), message.FinishedTime.Format("2006/01/02/15:04.05"))
	param := slack.PostMessageParameters{
		AsUser: true,
	}

	_, _, err := api.PostMessage(config.Channel, slack.MsgOptionText(text, true), slack.MsgOptionPostMessageParameters(param))
	return err
}
