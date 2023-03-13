package services

import (
	"encoding/json"
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	tweets               []selenium.WebElement
	err                  error
	tweetIdPath          selenium.WebElement
	channelNamePath      selenium.WebElement
	channelName          string
	tweetContentPath     selenium.WebElement
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
	fmt.Println(driver.CurrentURL())

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, "//article[@data-testid='tweet']"))
	tweets, err = driver.FindElements(selenium.ByXPATH, "//article[@data-testid='tweet']")
	CheckError("error finding tweet path", err)
	fmt.Println("-----------------------------")
	fmt.Println("No of tweets fetched: ", len(tweets))

	for i := 0; i < len(tweets); i++ {
		var handle models.NewsHandler
		fmt.Println("-----------------------------")
		FindingTweetID(tweets[i], driver, &handle)
		FindingChannelName(tweets[i], driver, &handle)
		FindingTweetContent(tweets, i, driver, &handle)
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
	CheckError("Error finding tweet id", err)
	tweetIdAttr, _ := tweetIdPath.GetAttribute("href")
	tweetId := strings.Split(tweetIdAttr, string('/'))
	handle.TweetId = tweetId[len(tweetId)-1]
	fmt.Println("Tweet ID: ", handle.TweetId)
}
func FindingChannelName(tweets selenium.WebElement, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING CHANNEL NAME

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div/div/a/div[1]/div/span/span"))
	channelNamePath, err = tweets.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"User-Names\"]/div/div/a/div[1]/div/span/span")
	CheckError("Error finding channel name", err)
	channelName, _ = channelNamePath.Text()
	handle.ChannelName = channelName
	fmt.Println("Channel name: ", channelName)

}
func FindingTweetContent(tweets []selenium.WebElement, i int, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING TWEET CONTENT

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]"))
	tweetContentPath, err = tweets[i].FindElement(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]")
	CheckError("Error finding tweet content", err)
	tweetContent, _ = tweetContentPath.Text()
	handle.TweetContent = tweetContent
	fmt.Println("Tweet", i+1, "content: ", tweetContent)
}
func FindingTweetTime(tweets selenium.WebElement, driver selenium.WebDriver, handle *models.NewsHandler) {
	//FINDING TWEET TIME

	driver.Wait(conditions.ElementIsLocated(selenium.ByXPATH, ".//time"))
	timeStampPath, err = tweets.FindElement(selenium.ByXPATH, ".//time")
	CheckError("Error finding time stamp", err)
	timeStamp, err = timeStampPath.GetAttribute("datetime")
	CheckError("Error getting datetime attr", err)
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
					log.Fatal("Error locating the image element", err)
				}
			} else {
				fmt.Println("No images present")
				break
			}

		} else {
			fmt.Println(len(tweetImagePath))
			for j := 1; j < len(tweetImagePath); j++ {
				imgURL, err = tweetImagePath[j].GetAttribute("src")
				CheckError("Error getting src", err)
				if !strings.Contains(imgURL, "/profile_images/") {
					fmt.Println("Image link: ", imgURL)
					handle.Media = append(handle.Media, models.Media{Type: "Image", URL: imgURL, TweetId: handle.TweetId})
				}
			}
			break
		}
	}

}
func FindingTweetVideos(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING TWEET VIDEOS
	for {
		tweetVideoPath, err = tweets.FindElement(selenium.ByCSSSelector, "video")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					log.Fatal("Error locating the image element", err)
				}
			} else {
				fmt.Println("No videos present")
				break
			}
		} else {
			videoURL, err = tweetVideoPath.GetAttribute("src")
			CheckError("Error getting src", err)
			fmt.Println("Video link: ", videoURL)
			handle.Media = append(handle.Media, models.Media{Type: "Video", URL: videoURL, TweetId: handle.TweetId})
		}
		break
	}

}
func FindingImageExternalLink(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING IMAGE ATTACHED ARTICLE EXTERNAL SOURCE LINK
	for {
		imageExternalURLPath, err = tweets.FindElement(selenium.ByCSSSelector, "div[data-testid=\"card.layoutLarge.media\"]>a")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					log.Fatal("Error locating the image attached article external URL path", err)
				}
			} else {
				fmt.Println("No image attached external url present")
				break
			}
		} else {
			imgExtURL, _ = imageExternalURLPath.GetAttribute("href")
			fmt.Println("Image Attached Article External URL: ", imgExtURL)
			handle.Media = append(handle.Media, models.Media{Type: "Image attached article Link", URL: imgExtURL, TweetId: handle.TweetId})
		}
		break
	}

}
func FindingTweetTextExternalLink(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING TWEET TEXT ATTACHED ARTICLE EXTERNAL SOURCE LINK
	for {
		tweetExternalURLPath, err = tweets.FindElement(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]/a[starts-with(@href,'https:')]")
		if err != nil {
			if !strings.Contains(err.Error(), "no such element") {
				if strings.Contains(err.Error(), "stale element reference") {
					continue
				} else {
					log.Fatal("Error locating the tweet external URL path", err)
				}
			} else {
				fmt.Println("No external url present in tweet text")
				break
			}

		} else {
			tweetExtURL, _ = tweetExternalURLPath.GetAttribute("href")
			fmt.Println("Tweet External URL: ", tweetExtURL)
			handle.Media = append(handle.Media, models.Media{Type: "Tweet attached external link", URL: tweetExtURL, TweetId: handle.TweetId})
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
		fmt.Println("Error finding hashtags", err)
	} else {
		flag := 0
		for j := 0; j < len(hashtagsPath); j++ {
			hashtags, _ = hashtagsPath[j].Text()
			if hashtags != "" {
				fmt.Println(hashtags)
				handle.HashTags = append(handle.HashTags, models.Hashtags{HashTags: hashtags, TweetId: handle.TweetId})
				flag += 1
			}
		}
		if flag == 0 {
			fmt.Printf("No hashtags present in tweet\n")
		}
	}
}
func FindingTweetMentions(tweets selenium.WebElement, handle *models.NewsHandler) {
	//FINDING ALL MENTIONS IN TWEET
	time.Sleep(2 * time.Second)
	mentionsPath, err = tweets.FindElements(selenium.ByXPATH, ".//div[@data-testid=\"tweetText\"]/div/span/a")
	if err != nil {
		fmt.Print("Error finding mentions", err)
	} else {
		fmt.Println("Mentions:")
		flag := 0
		for j := 0; j < len(mentionsPath); j++ {
			mentions, _ = mentionsPath[j].Text()
			if mentions != "" {
				fmt.Println(mentions)
				handle.UserHandle = append(handle.UserHandle, models.UserHandles{UserHandle: mentions, TweetId: handle.TweetId})
				flag += 1
			}
		}
		if flag == 0 {
			fmt.Printf("No mentions present in tweet\n")
		}
	}
}
func WriteIntoJSONFILE(h *[]models.NewsHandler) {
	//WRITING THE OUTPUT IN JSON FORMAT AND STORE IT IN FILE
	file, _ := json.MarshalIndent(h, "", "\t")
	_ = ioutil.WriteFile("services/tweet_data.json", file, 0644)
}
