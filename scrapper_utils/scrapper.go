package scrapper_utils

import (
	"encoding/json"
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	"github.com/SunnyYadav16/News_Scrapper/services"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"strings"
	"time"
)

var (
	tweets               []selenium.WebElement
	err                  error
	tweetIdPath          selenium.WebElement
	channelNamePath      selenium.WebElement
	channelName          string
	tweetContentPath     []selenium.WebElement
	tweetContent         string
	timeStampPath        selenium.WebElement
	timeStamp            string
	tweetImagePath       []selenium.WebElement
	imgURL               string
	tweetVideoPath       selenium.WebElement
	videoURL             string
	imageExternalURLPath selenium.WebElement
	imgExtURL            string
	tweetExternalURLPath selenium.WebElement
	tweetExtURL          string
	hashtagsPath         []selenium.WebElement
	hashtags             string
	mentionsPath         []selenium.WebElement
	mentions             string
	newsHandler          []models.NewsHandler
)

func TwitterLandingPage(driver selenium.WebDriver) {
	fmt.Println("Successfully logged in to twitter")
	time.Sleep(6 * time.Second)
	driver.Get("https://twitter.com/indiatoday")
	fmt.Println(driver.CurrentURL())
}
func NewsScrapper(driver selenium.WebDriver) (newshandler []models.NewsHandler) {
	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, "//article[@data-testid='tweet']"))
	tweets, err = driver.FindElements(selenium.ByXPATH, "//article[@data-testid='tweet']")
	services.CheckError("error finding tweet path", err)
	fmt.Println("-----------------------------")
	fmt.Println("No of tweets fetched: ", len(tweets))

	for i := 0; i < len(tweets); i++ {
		var handle models.NewsHandler
		fmt.Println("-----------------------------")
		FindingTweetID(tweets[i], driver, &handle)
		FindingChannelName(tweets[i], driver, &handle)
		FindingTweetContent(tweets[i], driver, &handle)
		FindingTweetTime(tweets[i], driver, &handle)
		FindingTweetHashtags(tweets[i], &handle)
		FindingTweetMentions(tweets[i], &handle)
		FindingTweetTextExternalLink(tweets[i], &handle)
		FindingTweetImages(tweets[i], &handle)
		FindingImageExternalLink(tweets[i], &handle)
		FindingTweetVideos(tweets[i], &handle)
		fmt.Println("-----------------------------")
		//APPENDING IT TO THE NewsHandler STRUCT
		newsHandler = append(newsHandler, handle)
	}

	return newsHandler
}
func FindingTweetID(tweets selenium.WebElement, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING TWEET ID

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div[2]/div/div[3]/a"))
	tweetIdPath, err = tweets.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div[2]/div/div[3]/a")
	services.CheckError("Error finding tweet id", err)
	tweetIdAttr, _ := tweetIdPath.GetAttribute("href")
	tweetId := strings.Split(tweetIdAttr, string('/'))
	handle.TweetId = tweetId[len(tweetId)-1]
	fmt.Println("Tweet ID: ", handle.TweetId)
}
func FindingChannelName(tweets selenium.WebElement, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING CHANNEL NAME

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div/div/a/div[1]/div/span/span"))
	channelNamePath, err = tweets.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div/div/a/div[1]/div/span/span")
	services.CheckError("Error finding channel name", err)
	channelName, _ = channelNamePath.Text()
	handle.ChannelName = channelName
	fmt.Println("Channel name: ", channelName)

}
func FindingTweetContent(tweets selenium.WebElement, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING TWEET CONTENT
	fmt.Println("Tweet content: ")
	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]"))
	tweetContentPath, err = tweets.FindElements(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]")
	services.CheckError("Error finding tweet content", err)

	for i := 0; i < len(tweetContentPath); i++ {
		tweetContent, _ = tweetContentPath[i].Text()
		if tweetContent != "" {
			handle.TweetContent = tweetContent
			fmt.Print(" ", tweetContent)
		}
	}
	fmt.Printf("\n")
}
func FindingTweetTime(tweets selenium.WebElement, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING TWEET TIME

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//time"))
	timeStampPath, err = tweets.FindElement(selenium.ByXPATH, ".//time")
	services.CheckError("Error finding time stamp", err)
	timeStamp, err = timeStampPath.GetAttribute("datetime")
	services.CheckError("Error getting datetime attr", err)
	handle.Timestamp, _ = time.Parse("2006-01-02T15:04:05.000Z", timeStamp)
	fmt.Println("Tweet posting time:  ", timeStamp)
}
func FindingTweetImages(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING TWEET IMAGES
	time.Sleep(2 * time.Second)
	for {
		tweetImagePath, err = tweets.FindElements(selenium.ByCSSSelector, "img.css-9pa8cd")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					services.CheckError("Error locating the image element", err)
				}
			} else {
				fmt.Println("No images present")
				break
			}

		} else {
			var m []models.Media
			for j := 1; j < len(tweetImagePath); j++ {
				imgURL, err = tweetImagePath[j].GetAttribute("src")
				services.CheckError("Error getting src", err)
				if !strings.Contains(imgURL, "/profile_images/") {
					fmt.Println("Image link: ", imgURL)
					m = append(m, models.Media{Type: "Image", URL: imgURL, TweetId: handle.TweetId})
				}
			}
			bytesArr, _ := json.MarshalIndent(m, "", " ")
			handle.MediaLinks += string(bytesArr)
		}
		break
	}

}
func FindingTweetVideos(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING TWEET VIDEOS
	time.Sleep(2 * time.Second)
	for {
		tweetVideoPath, err = tweets.FindElement(selenium.ByCSSSelector, "video")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					services.CheckError("Error locating the video element", err)
				}
			} else {
				fmt.Println("No videos present")
				break
			}
		} else {
			var m []models.Media
			videoURL, err = tweetVideoPath.GetAttribute("src")
			services.CheckError("Error getting src", err)
			fmt.Println("Video link: ", videoURL)
			m = append(m, models.Media{Type: "Video", URL: videoURL, TweetId: handle.TweetId})
			bytesArr, _ := json.MarshalIndent(m, "", " ")
			handle.MediaLinks += string(bytesArr)
		}
		break
	}

}
func FindingImageExternalLink(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING IMAGE ATTACHED ARTICLE EXTERNAL SOURCE LINK
	time.Sleep(2 * time.Second)
	for {
		imageExternalURLPath, err = tweets.FindElement(selenium.ByCSSSelector, "div[data-testid=\"card.layoutLarge.media\"]>a")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					services.CheckError("Error locating the image attached article external URL path", err)
				}
			} else {
				fmt.Println("No image attached external url present")
				break
			}
		} else {
			var m []models.Media
			imgExtURL, _ = imageExternalURLPath.GetAttribute("href")
			fmt.Println("Image Attached Article External URL: ", imgExtURL)
			m = append(m, models.Media{Type: "Image attached article Link", URL: imgExtURL, TweetId: handle.TweetId})
			bytesArr, _ := json.MarshalIndent(m, "", " ")
			handle.MediaLinks += string(bytesArr)
		}
		break
	}

}
func FindingTweetTextExternalLink(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING TWEET TEXT ATTACHED ARTICLE EXTERNAL SOURCE LINK
	time.Sleep(2 * time.Second)
	for {
		tweetExternalURLPath, err = tweets.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]/a[starts-with(@href,'https:')]")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					services.CheckError("Error locating the tweet external URL path", err)
				}
			} else {
				fmt.Println("No external url present in tweet text")
				break
			}

		} else {
			var m []models.Media
			tweetExtURL, _ = tweetExternalURLPath.GetAttribute("href")
			fmt.Println("Tweet External URL: ", tweetExtURL)
			m = append(m, models.Media{Type: "Tweet attached external link", URL: tweetExtURL, TweetId: handle.TweetId})
			bytesArr, _ := json.MarshalIndent(m, "", " ")
			handle.MediaLinks += string(bytesArr)
		}
		break
	}

}
func FindingTweetHashtags(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING ALL HASHTAGS IN TWEET
	time.Sleep(2 * time.Second)
	fmt.Print("Hashtags:")
	hashtagsPath, err = tweets.FindElements(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]/span/a[starts-with(@href,'/hashtag')]")
	if err != nil {
		services.CheckError("Error finding hashtags", err)
	} else {
		var h []models.Hashtags
		flag := 0
		for j := 0; j < len(hashtagsPath); j++ {
			hashtags, _ = hashtagsPath[j].Text()
			if hashtags != "" {
				fmt.Println(hashtags)
				h = append(h, models.Hashtags{HashTags: hashtags, TweetId: handle.TweetId})

				flag += 1
			}
		}
		bytesArr, _ := json.MarshalIndent(h, "", " ")
		handle.Hashtags = string(bytesArr)
		if flag == 0 {
			fmt.Printf("No hashtags present in tweet\n")
		}
	}
}
func FindingTweetMentions(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING ALL MENTIONS IN TWEET
	time.Sleep(2 * time.Second)
	fmt.Print("Mentions:")
	mentionsPath, err = tweets.FindElements(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]/div/span/a")
	if err != nil {
		services.CheckError("Error finding mentions", err)
	} else {
		var m []models.UserHandles

		flag := 0
		for j := 0; j < len(mentionsPath); j++ {
			mentions, _ = mentionsPath[j].Text()
			if mentions != "" {
				fmt.Println(mentions)
				m = append(m, models.UserHandles{UserHandle: mentions, TweetId: handle.TweetId})
				flag += 1
			}
		}
		bytesArr, _ := json.MarshalIndent(m, "", " ")
		handle.UserHandles = string(bytesArr)
		if flag == 0 {
			fmt.Printf("No mentions present in tweet\n")
		}
	}
}
func WriteIntoJSONFILE(h *[]models.NewsHandler) {
	//WRITING THE OUTPUT IN JSON FORMAT AND STORE IT IN FILE
	file, _ := json.MarshalIndent(h, "", "\t")
	_ = ioutil.WriteFile("scrapper_utils/tweet_data.json", file, 0644)
}
