package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/himidori/golang-vk-api"
	"log"
	"net/url"
	"os"
	"strconv"
)

func main() {
	vkServiceToken := os.Getenv("VK_SERVICE_TOKEN")
	vkGroup := os.Getenv("VK_GROUP")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	tgChannelId := os.Getenv("TG_CHANNEL_ID")
	redisUrl := os.Getenv("REDISCLOUD_URL")

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}
	redisDb := redis.NewClient(opt)

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

	vkGroupId, err := strconv.Atoi(vkGroup)
	if err != nil {
		panic(err)
	}

	params := url.Values{}
	wall, err := client.WallGet(vkGroupId, 3, params)
	if err != nil {
		panic(err)
	}

	vkLastPostId, err := redisDb.Get("VK_LAST_POST_ID").Result()
	if err != nil {
		panic(err)
	}
	vkLastPostIdstr, err := strconv.Atoi(vkLastPostId)
	if err != nil {
		panic(err)
	}

	for _, post := range wall.Posts {
		if post.IsPinned == 1 || post.MarkedAsAd == 1 || post.ID <= vkLastPostIdstr {
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
					videoUrl := "https://vk.com/video" + strconv.Itoa(attachment.Video.OwnerID) + "_" + strconv.Itoa(attachment.Video.ID)
					text += getInvisibleLink("\nВидео", videoUrl)
				}
				break
			}
		}
		msg := tgbotapi.NewMessageToChannel(tgChannelId, text)
		vkPostUrl := "https://vk.com/wall" + vkGroup + "_" + strconv.Itoa(post.ID)
		params := url.Values{}
		params.Set("url", vkPostUrl)
		params.Set("private", "1")

		var shortLink *ShortLink
		vkShortenedLink, _ := client.MakeRequest("utils.getShortLink", params)
		json.Unmarshal(vkShortenedLink.Response, &shortLink)

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("Пост", shortLink.ShortUrl)),
		)
		msg.ParseMode = "Markdown"

		_, err := bot.Send(msg)
		if err != nil {
			panic(err)
		}
		redisDb.Set("VK_LAST_POST_ID", strconv.Itoa(post.ID), 0)
	}
}

func getInvisibleLink(text string, url string) string {
	return "[" + text + "](" + url + ")"
}

type ShortLink struct {
	ShortUrl  string `json:"short_url"`
	AccessKey string `json:"access_key"`
	Key       string `json:"key"`
	Url       string `json:"url"`
}
