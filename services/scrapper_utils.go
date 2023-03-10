package services

import (
	"encoding/json"
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"os"
	"strings"
	time2 "time"
)

var err error

func getTimestamp(obj selenium.WebElement) (timestamp time2.Time) {
	timeTemp, err := obj.FindElement(selenium.ByXPATH, ".//time")
	CheckError("error while finding the element:", err)
	Timestamp, err := timeTemp.GetAttribute("datetime")
	CheckError("error while finding the attribute:", err)
	timestamp, err = time2.Parse("2006-01-02T15:04:05.000Z", Timestamp)
	CheckError("error while parsing time:", err)
	return
}

func getChannelName(obj selenium.WebElement) (channelName string) {
	channelTemp, err := obj.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div[2]/div[1]/div/a//span")
	CheckError("error while finding element:", err)
	channelName, err = channelTemp.Text()
	CheckError("error while finding text in element:", err)
	return
}

func getTweetContent(obj selenium.WebElement) (tweetContent string) {
	mainTweet, err := obj.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]")
	CheckError("error while finding element:", err)
	tweetContent, err = mainTweet.Text()
	CheckError("error while finding text in element:", err)
	return
}

func getUserHandlesAndHashtags(obj selenium.WebElement, tweetContent string, tweetId string) (hashtags []models.Hashtags, userHandles []models.UserHandles) {
	tokens := strings.FieldsFunc(tweetContent, splitCondition)
	for _, obj := range tokens {
		if strings.HasPrefix(obj, "#") {
			obj = strings.TrimLeft(obj, "#")
			hashtags = append(hashtags, models.Hashtags{TweetId: tweetId, Hashtags: "#" + obj})
		} else if strings.HasPrefix(obj, "@") {
			obj = strings.TrimLeft(obj, "@")
			userHandles = append(userHandles, models.UserHandles{TweetId: tweetId, UserHandle: "@" + obj})
		}
	}
	//fmt.Println(hashtags, userHandles)
	return
}

func getExternalURL(obj selenium.WebElement) (externalURL []models.Media) {
	ExteralUrl, err := obj.FindElements(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]//a[starts-with(@href,\"http\")]")
	CheckError("err", err)
	for _, obj := range ExteralUrl {
		URL, err := obj.Text()
		CheckError("error:", err)
		externalURL = append(externalURL, models.Media{Type: "link", URL: URL})
	}
	//fmt.Println()
	return
}

func getCardLinks(obj selenium.WebElement) (cards []models.Media) {
	CardURL, err := obj.FindElements(selenium.ByXPATH, ".//div[@data-testid=\"card.layoutLarge.media\"]/a")
	CheckError("error:", err)
	for _, obj := range CardURL {
		URL, err := obj.GetAttribute("href")
		CheckError("error:", err)
		cards = append(cards, models.Media{Type: "link", URL: URL})
	}
	return
}

func getImages(obj selenium.WebElement) (images []models.Media) {
	imgs, err := obj.FindElements(selenium.ByCSSSelector, "img.css-9pa8cd")
	CheckError("error while finding element:", err)
	for _, obj := range imgs {
		img, err := obj.GetAttribute("src")
		CheckError("error:", err)
		//utf8.EncodeRune(imgs)
		if !strings.Contains(img, "profile_images") {
			images = append(images, models.Media{Type: "image", URL: img})
		}
	}
	return
}

func getVideos(obj selenium.WebElement) (videos []models.Media) {
	video, err := obj.FindElements(selenium.ByCSSSelector, "video")
	CheckError("error while finding element:", err)
	for _, obj := range video {
		video, err := obj.GetAttribute("src")
		CheckError("error:", err)
		videos = append(videos, models.Media{Type: "video", URL: video})
	}
	return
}

func getTweetId(obj selenium.WebElement) (tweetId string) {
	id, err := obj.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div[2]/div/div[3]/a")
	CheckError("error while finding element:", err)
	tweetId, err = id.GetAttribute("href")
	CheckError("error while getting attribute:", err)
	tokens := strings.Split(tweetId, "/")
	tweetId = tokens[len(tokens)-1]
	return
}

func getTweetText(obj selenium.WebElement) models.NewsHandler {

	tweetData := models.NewsHandler{}

	//tweet id
	//.//div[@data-testid=\"User-Names\"]/div[2]/div[1]/div/a//span"
	tweetData.TweetId = getTweetId(obj)

	//retrieving timestamp
	tweetData.Timestamp = getTimestamp(obj)
	//channel name
	tweetData.ChannelName = getChannelName(obj)
	//main tweet
	tweetData.TweetContent = getTweetContent(obj)
	//HashTags & User Handles
	tweetData.HashTags, tweetData.UserHandle = getUserHandlesAndHashtags(obj, tweetData.TweetContent, tweetData.TweetId)
	//external urls
	tweetData.Media = getExternalURL(obj)
	//card Links
	tweetData.Media = append(tweetData.Media, getCardLinks(obj)...)
	//images
	tweetData.Media = append(tweetData.Media, getImages(obj)...)
	//videos
	tweetData.Media = append(tweetData.Media, getVideos(obj)...)

	return tweetData
}

func NewsScrapperNDTV(driver selenium.WebDriver, channelName string) (newsScrapped []models.NewsHandler) {

	driver.Get("https://twitter.com/" + channelName)

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, "//article[@data-testid=\"tweet\"]"))
	time2.Sleep(2 * time2.Second)
	tweetTexts, err := driver.FindElements(selenium.ByXPATH, "//article[@data-testid=\"tweet\"]")
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println(len(tweetTexts))

		for _, obj := range tweetTexts {
			newsScrapped = append(newsScrapped, getTweetText(obj))
		}
	}

	return newsScrapped
}

func ConvertToJSON(dataset []models.NewsHandler) {
	jsonData, err := json.MarshalIndent(dataset, "", "\t")
	CheckError("err:", err)
	err = ioutil.WriteFile("twitterData1.json", jsonData, os.ModePerm)
	CheckError("error occured while writing into file:", err)
	time2.Sleep(10 * time2.Second)
	fmt.Println("Successfully inserted into database")
}

func splitCondition(r rune) bool {
	return r == '\t' || r == '\n' || r == ' '
}

func loadData() []models.NewsHandler {
	var data []models.NewsHandler

	dataset, err := ioutil.ReadFile("twitterData1.json")
	CheckError("error while reading file", err)
	err = json.Unmarshal(dataset, &data)
	CheckError("error while parsing data", err)
	return data
}
