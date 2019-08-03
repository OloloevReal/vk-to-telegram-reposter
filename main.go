package main

import (
	"fmt"
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
	redisAddress := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,  // use default Addr
		Password: redisPassword, // no password set
	})

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

	val, err := redisdb.Get("vk_last_post_id").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	vk_last_post_id, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	for _, post := range wall.Posts {
		if post.IsPinned == 1 || post.MarkedAsAd == 1 || post.ID <= vk_last_post_id {
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

		bot.Send(msg)
		updateDb(redisdb, post.ID)
	}
}

func updateDb(client *redis.Client, i int) {
	err := client.Set("vk_last_post_id", i, 0).Err()
	if err != nil {
		panic(err)
	}
}

func getInvisibleLink(text string, url string) string {
	return "[" + text + "](" + url + ")"
}
