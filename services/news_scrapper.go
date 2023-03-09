package services

import (
	"errors"
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"os"
	"strings"
	"time"
)

func NewsScrapper(driver selenium.WebDriver, channelName string) (newsScrapped []models.NewsHandler) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Closing Driver and exiting program.")
			CloseService(driver)
			os.Exit(0)
		}
	}()

	var (
		articles []selenium.WebElement //ARTICLE ELEMENT FINDER
		err      error
		length   int    //LENGTH OF TOTAL ELEMENTS FOUND
		channel  string //NAME OF NEWS HANDLE CHANNEL
	)

	//GET CHANNEL NAME
	switch channelName {
	case "timesofindia":
		channel = "Times Of India"
	case "ndtv":
		channel = "NDTV"
	case "indiatoday":
		channel = "India Today"
	default:
		err = errors.New("invalid channel")
		CheckError("Scrapping Invalid Channel", err)
	}

	//GET CHANNEL HANDLE
	err = driver.Get("https://twitter.com/" + channelName)
	CheckError("Error Getting News Handle", err)

	//LOAD WAIT
	err = driver.Wait(conditions.ElementIsLocated(selenium.ByCSSSelector, "article[role=article]"))
	CheckError("Error Waiting For Article Elements To Load", err)

	//FIND NEWS TWEETS ARTICLES
	articles, err = driver.FindElements(selenium.ByCSSSelector, "article[role=article]")
	CheckError("Twitter Articles Elements Not Found", err)

	//TOTAL LENGTH OF ARTICLES FETCHED
	length = len(articles)

	//SCRAPING TWEET NEWS
	for i := 0; i < length; i++ {

		//SCRAPPING SINGLE NEWS TWEET
		var scrappedNews models.NewsHandler
		scrappedNews.ChannelName = channel
		scrapNews(&scrappedNews, articles[i])
		newsScrapped = append(newsScrapped, scrappedNews)
	}
	return
}

func scrapNews(scrappedNews *models.NewsHandler, article selenium.WebElement) {

	var (
		err                     error
		tweetIdElement          selenium.WebElement
		tweetIdLink             string
		tweetDateTimeElement    selenium.WebElement
		tweetDateTime           string
		dateTimeLayout          = "2006-01-02T15:04:05.000Z"
		tweetContentDiv         selenium.WebElement
		tweetContentSpans       []selenium.WebElement
		spanText                string
		tweetContent            string
		length                  int
		tweetImages             []selenium.WebElement
		tweetImageLink          string
		tweetExternalSource     selenium.WebElement
		tweetExternalSourceLink string
		tweetVideo              selenium.WebElement
		tweetVideoLink          string
	)

	//FINDING  A ELEMENT FOR TWEET ID
	tweetIdElement, err = article.FindElement(selenium.ByCSSSelector, "a.css-4rbku5.css-18t94o4.css-901oao.r-14j79pv.r-1loqt21.r-xoduu5.r-1q142lx.r-1w6e6rj.r-37j5jr.r-a023e6.r-16dba41.r-9aw3ui.r-rjixqe.r-bcqeeo.r-3s2u2q.r-qvutc0")
	CheckError("Tweet Id Element Not Found", err)

	//FETCHING LINK FOR TWEET ID
	tweetIdLink, err = tweetIdElement.GetAttribute("href")
	CheckError("Tweet Id Link Not Found", err)

	splitter := strings.Split(tweetIdLink, "/")

	//FETCHING TWEET ID
	scrappedNews.TweetId = splitter[len(splitter)-1]

	//FINDING TIME ELEMENT
	tweetDateTimeElement, err = article.FindElement(selenium.ByCSSSelector, "time")
	CheckError("Tweet Time Element Not Found", err)

	//FETCHING DATE AND TIME OF TWEET
	tweetDateTime, err = tweetDateTimeElement.GetAttribute("datetime")
	CheckError("Tweet Date Time Not Found", err)

	//PARSING IN DATETIME FORMAT
	scrappedNews.Timestamp, err = time.Parse(dateTimeLayout, tweetDateTime)
	CheckError("Cannot Convert Into Date Time Format", err)

	//FINDING TWEET CONTENT'S PARENT ELEMENT
	tweetContentDiv, err = article.FindElement(selenium.ByCSSSelector, "div[data-testid=tweetText]")
	CheckError("Tweet Contents Div Element Not Found", err)

	//FINDING TWEET CONTENT'S ELEMENT
	tweetContentSpans, err = tweetContentDiv.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
	CheckError("Tweet Contents Span Element Not Found", err)

	//TOTAL LENGTH OF SPAN ELEMENT FETCHED FOR TWEET CONTENT
	length = len(tweetContentSpans)
	for i := 0; i < length; i++ {

		//FETCHING TWEET CONTENT
		spanText, err = tweetContentSpans[i].Text()
		CheckError("Error Finding Text in Span Element", err)
		tweetContent += spanText
	}
	scrappedNews.TweetContent = tweetContent

	//EXTRACTING MENTIONS USER-HANDLES AND LINKS FROM TWEET CONTENT
	tagsHandlesAndLinksFinder(tweetContent, scrappedNews)

	//FINDING LINKS PRESENT IN A TWEET
	tweetExternalSource, err = article.FindElement(selenium.ByCSSSelector, "a[rel='noopener noreferrer nofollow']")
	if err == nil {

		//FETCHING LINKS IN A TWEET
		tweetExternalSourceLink, err = tweetExternalSource.GetAttribute("href")
		CheckError("Error Getting External Source Link", err)
		scrappedNews.Media = append(scrappedNews.Media, models.Media{Type: "Link", URL: tweetExternalSourceLink})
	} else {
		if !strings.Contains(err.Error(), "no such element") {
			CheckError("Something Went Wrong Fetching Link", err)
		}
	}

	//FINDING IMAGES IN A TWEET
	tweetImages, err = article.FindElements(selenium.ByCSSSelector, "img.css-9pa8cd")
	if err == nil {

		//TOTAL LENGTH OF IMAGE ELEMENT FETCHED
		length = len(tweetImages)
		for i := 0; i < length; i++ {

			//FETCHING IMAGE LINK
			tweetImageLink, err = tweetImages[i].GetAttribute("src")
			CheckError("Error Getting Image Link", err)

			//INSERTING LINK IF IT IS NOT A PROFILE PICTURE
			if !strings.Contains(tweetImageLink, "normal.jpg") && !strings.Contains(tweetImageLink, "normal.png") && !strings.Contains(tweetImageLink, "normal.jpeg") {
				scrappedNews.Media = append(scrappedNews.Media, models.Media{Type: "Image", URL: tweetImageLink})
			}
		}
	} else {
		if !strings.Contains(err.Error(), "no such element") {
			CheckError("Something Went Wrong Fetching Image", err)
		}
	}

	//FINDING VIDEO IN A TWEET
	tweetVideo, err = article.FindElement(selenium.ByCSSSelector, "video")
	if err == nil {

		//FETCHING VIDEO LINK
		tweetVideoLink, err = tweetVideo.GetAttribute("src")
		CheckError("Error Getting Video Link", err)
		scrappedNews.Media = append(scrappedNews.Media, models.Media{Type: "Video", URL: tweetVideoLink})
	} else {
		if !strings.Contains(err.Error(), "no such element") {
			CheckError("Something Went Wrong Fetching Video", err)
		}
	}
}

func tagsHandlesAndLinksFinder(tweetContent string, scrappedNews *models.NewsHandler) {
	var media models.Media

	//SPLITTING TWEETS TO FIND
	splitter := func(characters rune) bool {
		return characters == ' ' || characters == '\t' || characters == '\n'
	}
	words := strings.FieldsFunc(tweetContent, splitter)

	//TOTAL LENGTH OF WORDS OBTAINED
	length := len(words)
	for i := 0; i < length; i++ {
		word := words[i]
		word = strings.TrimSpace(word)
		if word[:1] == "#" { //IF IT IS A MENTION
			scrappedNews.HashTags = append(scrappedNews.HashTags, word)
		} else if word[:1] == "@" { //IF IT IS A USER-HANDLE
			scrappedNews.UserHandle = append(scrappedNews.UserHandle, word)
		} else if strings.Contains(word, "#") { //IF A WORD CONTAINS A MENTION E.G. THE#aajtak
			str := strings.Split(word, "#")
			scrappedNews.HashTags = append(scrappedNews.HashTags, "#"+str[1])
		} else if strings.Contains(word, "@") { //IF A WORD CONTAINS A USER-HANDEL of@user
			str := strings.Split(word, "@")
			scrappedNews.UserHandle = append(scrappedNews.UserHandle, "@"+str[1])
		} else if strings.Contains(word, "http") || strings.Contains(word, "https") { //IF IT IS A LINK
			media.Type = "Link"
			media.URL = word
			scrappedNews.Media = append(scrappedNews.Media, media)
			if i+1 != length && words[i+1] == word {
				i++
			}
		}
	}
}
