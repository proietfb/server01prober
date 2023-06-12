package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	actions "server01prober/pkgs"
)

type Conf struct {
	BotAPIKey  string
	RestrictTo *[]Restrictions
}
type Restrictions struct {
	Username string
	ChatID   int64
}

func parseConf() *Conf {
	var conf Conf
	viper.SetConfigName("conf")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if err = viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
	return &conf
}

func (conf *Conf) auth(username string, chatID int64) bool {
	if conf != nil && conf.RestrictTo != nil {
		for _, r := range *conf.RestrictTo {
			if username == r.Username && chatID == r.ChatID {
				return true
			}
		}
	}
	return false
}

func main() {
	conf := parseConf()
	bot, err := tgbotapi.NewBotAPI(conf.BotAPIKey)

	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Println("Authorized on account ", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	acd := actions.Init()

	for update := range updates {
		if conf.auth(update.SentFrom().UserName, update.SentFrom().ID) {
			if update.Message != nil {
				log.Printf("Authenticated")
				str := acd.GetContainers()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User "+update.SentFrom().UserName+" cannot use this bot")
			bot.Send(msg)
			log.Printf("Unauthenticated")
		}

	}

}
