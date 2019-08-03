package main

import (
	"github.com/himidori/golang-vk-api"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serviceToken := os.Getenv("SERVICE_TOKEN")
	group := os.Getenv("VK_GROUP")

	client, err := vkapi.NewVKClientWithToken(serviceToken, &vkapi.TokenOptions{
		ValidateOnStart: true,
		ServiceToken:    true,
	})

	if err != nil {
		log.Panic(err)
	}

	params := url.Values{}
	wall, _ := client.WallGet(group, 10, params)

	for _, post := range wall.Posts {
		if post.IsPinned == 1 || post.MarkedAsAd == 1 {
			continue
		}

		log.Println(post.ID)
		log.Println(post.Text)

		if post.Attachments != nil {
			for _, attachment := range post.Attachments {
				switch attachment.Type {
				case "photo":
					log.Println(attachment.Photo.Photo1280)
				case "video":
					log.Println(attachment.Video.FirstFrame800)
				}
			}
		}
	}
}
