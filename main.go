package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"os"
	actions "server01prober/pkgs/actions"
	probe "server01prober/pkgs/probe"
	"strings"
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

var baseCommands = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/status"),
		tgbotapi.NewKeyboardButton("/containerLog")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/containerRestart"),
		tgbotapi.NewKeyboardButton("/stats"),
	),
)

func generateContainersCommand(ac *actions.ActionsData) tgbotapi.InlineKeyboardMarkup {
	var kMarkup tgbotapi.InlineKeyboardMarkup

	names := ac.GetContainersName()

	if len(names) > 0 {
		kbuttons := make([][]tgbotapi.InlineKeyboardButton, len(names))

		for i, name := range names {
			kbuttons[i] = append(kbuttons[i], tgbotapi.NewInlineKeyboardButtonData(name[1:], name[1:]))
		}
		newKboard := tgbotapi.NewInlineKeyboardMarkup(kbuttons...)
		kMarkup = newKboard
	}

	return kMarkup
}

func callbackHandler(update tgbotapi.Update, acd *actions.ActionsData, bot *tgbotapi.BotAPI) bool {
	var ret bool
	switch strings.ToLower((update.CallbackQuery.Message.Text)) {
	case "restart":
		log.Println("Container id: ", acd.GetContainerID(update.CallbackQuery.Data))
		ret = acd.RestartContainer("/" + update.CallbackQuery.Data)
	case "log":
		if res := acd.GetContainerlog("/" + update.CallbackQuery.Data); res != nil {
			os.WriteFile("/tmp/"+update.CallbackQuery.Data+".log", res, 0644)
			file := tgbotapi.NewDocument(update.CallbackQuery.From.ID, tgbotapi.FilePath("/tmp/"+update.CallbackQuery.Data+".log"))
			bot.Send(file)
			ret = true
		}
	default:
		log.Println("Cannot perform action ", update.CallbackQuery.Message.Text)
	}
	return ret
}

func parseCommand(update tgbotapi.Update, acd *actions.ActionsData) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "start":
		msg.Text = "Welcome " + update.SentFrom().UserName
		msg.ReplyMarkup = baseCommands
	case "help":
		msg.Text = "Help requested"
	case "stats":
		msg.Text = probe.Probe()
	case "status":
		msg.Text = acd.GetStatus()
	case "containerRestart":
		msg.Text = "Restart"
		msg.ReplyMarkup = generateContainersCommand(acd)
	case "containerLog":
		msg.Text = "Log"
		msg.ReplyMarkup = generateContainersCommand(acd)
	default:
		msg.Text = "Can't Understand"

	}
	return &msg
}

func parseCallback(update tgbotapi.Update, acd *actions.ActionsData, bot *tgbotapi.BotAPI) *tgbotapi.MessageConfig {
	ret := callbackHandler(update, acd, bot)
	log.Println(update.CallbackQuery.Message.Text+" "+update.CallbackQuery.Data, " Ret: ", ret)
	result := "Success"
	if !ret {
		result = "Error"
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, result)
	return &msg
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
				if update.Message.IsCommand() {
					msg := parseCommand(update, acd)
					bot.Send(msg)
				} else {
					continue
				}
			} else if update.CallbackQuery != nil {
				cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "Performing "+update.CallbackQuery.Message.Text)
				bot.Request(cb)
				msg := parseCallback(update, acd, bot)
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User "+update.SentFrom().UserName+" cannot use this bot")
			bot.Send(msg)
			log.Printf("Unauthenticated")
		}

	}

}
