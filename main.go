package main

import (
	"github.com/himidori/golang-vk-api"

	"log"
	"net/url"
	"os"
)

func main() {
	serviceToken := os.Getenv("SERVICE_TOKEN")
	group := os.Getenv("VK_GROUP")

	client, err := vkapi.NewVKClientWithToken(serviceToken, &vkapi.TokenOptions{
		ValidateOnStart: true,
		ServiceToken:    true,
	})

	if err != nil {
		log.Println(err)
	}

	params := url.Values{}
	wall, _ := client.WallGet(group, 10, params)

	for _, post := range wall.Posts {
		if post.IsPinned == 1 || post.MarkedAsAd == 1 {
			continue
		}

		log.Println(post.Text)
		log.Println(post.ID)

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
