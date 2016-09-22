package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/basilboli/hackernewsbot/fetcher"
	"github.com/basilboli/hackernewsbot/models"
	"github.com/tucnak/telebot"
)

var bot *telebot.Bot

func listenMessages() {
	for message := range bot.Messages {

		log.Printf("Message %s from %d", message.Text, message.Chat.ID)
		switch message.Text {
		case "start", "/start", "/start start", "Start", "/Start":
			welcome, err := telebot.NewFile("welcome.png")
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Reading file %s %s", welcome.FileID, welcome.Local())
			photo := telebot.Photo{File: welcome}
			bot.SendPhoto(message.Chat, &photo, nil)
			bot.SendMessage(message.Chat, "Hello, this is hacker news bot. You can check hacker news top stories or subscribe for hourly digest!", buildLayout(message.Chat))
		case "Top Stories", "top stories", "top", "Check top stories":
			sendTopStoriesAllInOne(message.Chat)
		case "Subscribe", "subscribe":
			models.SubscribeForTopStories(message.Chat.ID)
			bot.SendMessage(message.Chat, "You're subscribed. You'll get 5 tasty links every hour ;)", hideKeyboard())
			bot.SendMessage(message.Chat, "Remember, you can always type menu to get more options!", nil)
		case "Unsubscribe", "unsubscribe":
			models.UnSubscribeForTopStories(message.Chat.ID)
			bot.SendMessage(message.Chat, "You're unsubscribed. No more hassle ;)", hideKeyboard())
			bot.SendMessage(message.Chat, "Remember, you can always type menu to get more options and re-subscribe!", nil)
		case "Menu", "menu", "help":
			bot.SendMessage(message.Chat, "You can do the following : ", buildLayout(message.Chat))
		default:
			bot.SendMessage(message.Chat, "I can't answer this question. Give me another try.", buildLayout(message.Chat))
		}
	}
}

func sendTopStories(chat telebot.Chat) {
	// topFive, err := models.GetTopFive()
	bot.SendMessage(chat, "Hacker News latest trending stories!", hideKeyboard())
	posts, err := models.FindTopFive()
	if err != nil {
		log.Fatal(err)
	}
	for _, post := range posts {
		msg := fmt.Sprintf("%s %s", post.Title, post.URL)
		bot.SendMessage(chat, msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
	}
	if err != nil {
		log.Fatal(err)
	}
	bot.SendMessage(chat, "Remember, you can always type menu to get more options!", nil)
}

func sendTopStoriesAllInOne(chat telebot.Chat) {
	posts, err := models.FindTopFive()
	if err != nil {
		log.Fatal(err)
	}
	var result, msg string
	for i, post := range posts {
		msg = fmt.Sprintf("%d - %s %s \n", i, post.Title, post.URL)
		result = result + msg
	}
	bot.SendMessage(chat, "Hacker News latest trending stories!", hideKeyboard())
	bot.SendMessage(chat, result, hideKeyboard())
	bot.SendMessage(chat, "You can always type menu to get more options.", nil)
}

func buildLayout(chat telebot.Chat) *telebot.SendOptions {
	isSubscribed, err := models.IsSubscribed(chat.ID)
	if err != nil {
		return nil
	}
	if isSubscribed {
		return &telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:         true,
				Selective:          true,
				HideCustomKeyboard: true,

				CustomKeyboard: [][]string{
					[]string{"Top Stories"},
					[]string{"Unsubscribe"},
				},
			},
		}
	}

	return &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			ForceReply:         true,
			Selective:          true,
			HideCustomKeyboard: true,

			CustomKeyboard: [][]string{
				[]string{"Top Stories"},
				[]string{"Subscribe"},
			},
		},
	}
}

func hideKeyboard() *telebot.SendOptions {
	return &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			HideCustomKeyboard: true,
		},
	}
}
func showKeyboard() *telebot.SendOptions {
	return &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			HideCustomKeyboard: false,
		},
	}
}

// SendNotifications sends notifications by chat to subscribed users one in a hour
func sendNotifications() {
	for {
		fmt.Println("Notifying")
		users, err := models.GetAllSubscribedUsers()
		if err != nil {
			log.Fatal(err)
		}
		for _, user := range users {
			log.Printf("Notifying user %s", user)
			id, err := strconv.ParseInt(user, 10, 64)
			if err != nil {
				return
			}
			chat := telebot.Chat{ID: id}
			sendTopStoriesAllInOne(chat)
		}
		time.Sleep(time.Hour * 1)
	}
}

func main() {
	flag.Parse()
	token := flag.Arg(0)

	if token != "" {
		fmt.Println("Found flag token")
	} else {
		fmt.Println("Token flag is empty, trying environment variable")
		token = os.Getenv("CF_TELEGRAM_TOKEN")
		if token == "" {
			fmt.Println("No token provided")
			os.Exit(0)
		} else {
			fmt.Printf("Found env token %s\n", token)
		}
	}

	fmt.Println("Starting bot ...")

	if newBot, err := telebot.NewBot(token); err != nil {
		fmt.Println(err)
		return
	} else {
		// shadowing, remember?
		bot = newBot
	}

	bot.Messages = make(chan telebot.Message)

	go listenMessages()
	go sendNotifications()
	go fetcher.FetchTopStories()

	bot.Start(1 * time.Millisecond)

}
