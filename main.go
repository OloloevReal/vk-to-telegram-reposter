package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/himidori/golang-vk-api"
	"log"
	"net/url"
	"os"
)

func main() {
	vkServiceToken := os.Getenv("VK_SERVICE_TOKEN")
	vkGroup := os.Getenv("VK_GROUP")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	tgChannelId := os.Getenv("TG_CHANNEL_ID")

	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Panic(err)
	}

	client, err := vkapi.NewVKClientWithToken(vkServiceToken, &vkapi.TokenOptions{
		ValidateOnStart: true,
		ServiceToken:    true,
	})

	if err != nil {
		log.Println(err)
	}

	params := url.Values{}
	wall, _ := client.WallGet(vkGroup, 10, params)

	if err != nil {
		panic(err)
	}

	for _, post := range wall.Posts {
		if post.IsPinned == 1 || post.MarkedAsAd == 1 {
			continue
		}

		text := post.Text

		log.Println(post.Text)
		log.Println(post.ID)

		if post.Attachments != nil {
			for _, attachment := range post.Attachments {
				if attachment.Type == "photo" {
					text += getInvisibleLink(" ", attachment.Photo.Photo604) // пустой символ чтобы прикрепить картинку в markdown без текста
				} else if attachment.Type == "video" {
					videoUrl := "https://vk.com/video" + string(attachment.Video.OwnerID) + "_" + string(attachment.Video.ID)
					text += getInvisibleLink("\nВидео", videoUrl)
				}
				break
			}
		}
		msg := tgbotapi.NewMessageToChannel(tgChannelId, post.Text)
		vkPostUrl := "https://vk.com/wall" + vkGroup + "-" + string(post.ID)

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Пост", vkPostUrl)),
		)

		_, err := bot.Send(msg)
		if err != nil {
			panic(err)
		}
		err = os.Setenv("LAST_POST_ID", string(post.ID))
		if err != nil {
			panic(err)
		}
	}
}

func getInvisibleLink(text string, url string) string {
	return "[" + text + "](" + url + ")"
}
