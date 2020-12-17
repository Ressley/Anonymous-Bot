package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Config - gets token from json file
type Config struct {
	TelegramBotToken string
}

var (
	//FirstID stores address of First Chat
	FirstID int64 = -00000000 //<-Here goes id of First chat->
	//SecondID stores address of Second Chat
	SecondID int64 = -00000000 // <-Here goes id of Second chat->
	id         int64 = 0
	now        time.Time
	//FirstChatTimer stores limit time for first chat to avoid spam
	FirstChatTimer map[int64]time.Time
	//SecondChatTimer stores limit time for second chat to avoid spam
	SecondChatTimer map[int64]time.Time
)

var chatsTab = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("First Chat", "Message is send to First Chat"),
		tgbotapi.NewInlineKeyboardButtonData("Second Chat", "Message is send to Second Chat"),
	),
)
var nullTab = tgbotapi.NewInlineKeyboardMarkup()

//Command - works with command operations
func Command(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Command() {
	case "start":
		msg.Text = "Hello, I am an anonymous bot, send me messages, then I will send it anonymously to a given chat."
	default:
		msg.Text = "Send me messages, then I will send it anonymously to a given chat"
	}
}

//SendMessage - sends message to server
func SendMessage(update *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(id, update.Message.ReplyToMessage.Text)
	if update.Message.ForwardFrom != nil && update.Message.ForwardFromChat != nil {
		bot.Send(tgbotapi.NewForward(id, update.Message.ForwardFromChat.ID, update.Message.ForwardFromMessageID))
	}

	if update.Message.ReplyToMessage.Photo != nil {
		photo := *update.Message.ReplyToMessage.Photo
		fileid := photo[len(photo)-1].FileID
		if update.Message.ReplyToMessage.Caption != "" {
			msg.Text = update.Message.ReplyToMessage.Caption
		}
		bot.Send(tgbotapi.NewPhotoShare(id, fileid))
	}
	if update.Message.ReplyToMessage.Sticker != nil {
		if update.Message.ReplyToMessage.Caption != "" {
			msg.Text = update.Message.ReplyToMessage.Caption
		}
		bot.Send(tgbotapi.NewStickerShare(id, update.Message.ReplyToMessage.Sticker.FileID))
	}
	if update.Message.ReplyToMessage.Document != nil {
		if update.Message.ReplyToMessage.Caption != "" {
			msg.Text = update.Message.ReplyToMessage.Caption
		}
		bot.Send(tgbotapi.NewDocumentShare(id, update.Message.ReplyToMessage.Document.FileID))
	}
	if update.Message.ReplyToMessage.Video != nil {
		if update.Message.ReplyToMessage.Caption != "" {
			msg.Text = update.Message.ReplyToMessage.Caption
		}
		bot.Send(tgbotapi.NewVideoShare(id, update.Message.ReplyToMessage.Video.FileID))
	}
	if update.Message.ReplyToMessage.VideoNote != nil {
		if update.Message.ReplyToMessage.Caption != "" {
			msg.Text = update.Message.ReplyToMessage.Caption
		}
		bot.Send(tgbotapi.NewVideoNoteShare(id, update.Message.ReplyToMessage.VideoNote.Length, update.Message.ReplyToMessage.VideoNote.FileID))
	}
	if update.Message.ReplyToMessage.Voice != nil {
		if update.Message.ReplyToMessage.Caption != "" {
			msg.Text = update.Message.ReplyToMessage.Caption
		}
		bot.Send(tgbotapi.NewVoiceShare(id, update.Message.ReplyToMessage.Voice.FileID))
	}
	bot.Send(msg)
	log.Printf("Message: %v", update.Message.ReplyToMessage.Text)
}

func reply(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "choose a server")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = chatsTab
	bot.Send(msg)
}

func create(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery != nil {
		if strings.Contains(update.CallbackQuery.Data, "First Chat") {
			if _, ok := FirstChatTimer[update.CallbackQuery.Message.Chat.ID]; ok {
				if time.Until(FirstChatTimer[update.CallbackQuery.Message.Chat.ID]) < -5*time.Second {
					FirstChatTimer[update.CallbackQuery.Message.Chat.ID] = time.Now()
					id = FirstID
					SendMessage(update.CallbackQuery, bot)
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
					edit := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					bot.Send(edit)
				} else {
					n := (float64)(time.Until(FirstChatTimer[update.CallbackQuery.Message.Chat.ID]) / time.Second)
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "wait "+strconv.Itoa(30-(int)(math.Abs(n)))+" seconds please"))
				}
			} else {
				FirstChatTimer[update.CallbackQuery.Message.Chat.ID] = time.Now()
				id = FirstID
				SendMessage(update.CallbackQuery, bot)
				edit := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				bot.Send(edit)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			}
		} else if strings.Contains(update.CallbackQuery.Data, "Second Chat") {
			if _, ok := SecondChatTimer[update.CallbackQuery.Message.Chat.ID]; ok {
				if time.Until(SecondChatTimer[update.CallbackQuery.Message.Chat.ID]) < -30*time.Second {
					SecondChatTimer[update.CallbackQuery.Message.Chat.ID] = time.Now()
					id = SecondID
					SendMessage(update.CallbackQuery, bot)
					edit := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					bot.Send(edit)
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				} else {
					n := (float64)(time.Until(SecondChatTimer[update.CallbackQuery.Message.Chat.ID]) / time.Second)
					bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "wait "+strconv.Itoa(30-(int)(math.Abs(n)))+" seconds please"))
				}
			} else {
				SecondChatTimer[update.CallbackQuery.Message.Chat.ID] = time.Now()
				id = SecondID
				SendMessage(update.CallbackQuery, bot)
				edit := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				bot.Send(edit)
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
			}
		}
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	if update.Message.IsCommand() {
		Command(&update, &msg)
		bot.Send(msg)
	} else {
		reply(&update, bot)
	}
}

func main() {
	now = time.Now()
	FirstChatTimer = make(map[int64]time.Time)
	SecondChatTimer = make(map[int64]time.Time)
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)

	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil && (update.Message.Chat.ID == FirstID || update.Message.Chat.ID == SecondID) {
			continue
		}
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.Message != nil && update.Message.Chat.ID == FirstID {
			continue
		}
		go create(update, bot)
	}
}
