package services

import (
	"errors"
	"fmt"
	"github.com/SunnyYadav16/News_Scrapper/models"
	conditions "github.com/serge1peshcoff/selenium-go-conditions"
	"github.com/tebeka/selenium"
	"strings"
	"time"
)

func NewsScrapper(driver selenium.WebDriver, channelName string) (newsScrapped []models.NewsHandler) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error: ", r)
			fmt.Println("Closing Driver and exiting program.")
			CloseService(driver)
		}
	}()

	var (
		articles []selenium.WebElement //ARTICLE ELEMENT FINDER
		err      error
		length   int    //LENGTH OF TOTAL ELEMENTS FOUND
		channel  string //NAME OF NEWS HANDLE CHANNEL
	)

	//GET CHANNEL NAME UNCOMMENT BELOW CODE FOR GENERALIZED
	switch channelName {
	case "timesofindia":
		channel = "Times Of India"
	//case "ndtv":
	//	channel = "NDTV"
	//case "indiatoday":
	//	channel = "India Today"
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
		tweetContent            string
		tweetContentSpans       []selenium.WebElement
		spanText                string
		length                  int
		tweetImages             []selenium.WebElement
		tweetImageLink          string
		tweetExternalSource     selenium.WebElement
		tweetExternalSourceLink string
		tweetVideo              selenium.WebElement
		tweetVideoLink          string
		tweetChannel            selenium.WebElement
	)

	tweetChannel, err = article.FindElement(selenium.ByXPATH, "//div[1]/div/a/div/div[1]/span/span")
	CheckError("Tweet Channel Name Not Found", err)

	scrappedNews.ChannelName, err = tweetChannel.Text()
	CheckError("Error Finding Channel Name", err)

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

	//FINDING TWEET CONTENT'S ELEMENT
	tweetContentDiv, err = article.FindElement(selenium.ByCSSSelector, "div[data-testid=tweetText]")
	CheckError("Tweet Contents Div Element Not Found", err)

	//FETCHING TWEET CONTENT
	tweetContent, err = tweetContentDiv.Text()
	CheckError("Tweet Content Not Found", err)

	//IF NO CONTENT INSIDE DIV TAG
	if len(strings.TrimSpace(tweetContent)) == 0 {

		//FINDING SPAN ELEMENT FOR TWEET CONTENT
		tweetContentSpans, err = tweetContentDiv.FindElements(selenium.ByCSSSelector, ".css-901oao.css-16my406.r-poiln3.r-bcqeeo.r-qvutc0")
		CheckError("Tweet Content's Span Element Not Found", err)
		length = len(tweetContentSpans)
		for i := 0; i < length; i++ {

			//FETCHING TWEET CONTENT
			spanText, err = tweetContentSpans[i].Text()
			CheckError("Tweet Content Not Found Inside Span Element", err)
			tweetContent += spanText
		}
	}

	if len(strings.TrimSpace(tweetContent)) == 0 {
		err = errors.New("tweet content empty")
		CheckError("Empty Tweet Content", err)
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
			if !strings.Contains(tweetImageLink, "profile_images") {
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

	//SPLITTING TWEETS TO FIND
	splitter := func(characters rune) bool {
		return characters == ' ' || characters == '\t' || characters == '\n'
	}

	//SPLITTING WORDS
	finalFunc := func(characters rune) bool {
		return characters == ')' || characters == ',' || characters == '.' || characters == '!' || characters == '?' || characters == ':' || characters == ';' || characters == '`' || characters == '"' || characters == '\'' || characters == '-' || characters == ']' || characters == '}' || characters == '*' || characters == '|'
	}

	words := strings.FieldsFunc(tweetContent, splitter)

	//TOTAL LENGTH OF WORDS OBTAINED
	length := len(words)
	for i := 0; i < length; i++ {
		var (
			media      models.Media
			hashTag    models.HashTag
			userHandle models.UserHandle
		)
		hashTag.NewsHandlers = append(hashTag.NewsHandlers, scrappedNews)
		userHandle.NewsHandlers = append(hashTag.NewsHandlers, scrappedNews)
		word := words[i]
		word = strings.TrimSpace(word)
		if word[:1] == "#" { //IF IT IS A MENTION
			final := strings.FieldsFunc(word, finalFunc)
			hashTag.TagName = final[0]
			scrappedNews.HashTags = append(scrappedNews.HashTags, &hashTag)
		} else if word[:1] == "@" { //IF IT IS A USER-HANDLE
			final := strings.FieldsFunc(word, finalFunc)
			userHandle.Name = final[0]
			scrappedNews.UserHandles = append(scrappedNews.UserHandles, &userHandle)
		} else if strings.Contains(word, "#") { //IF A WORD CONTAINS A MENTION E.G. THE#aajtak
			str := strings.Split(word, "#")
			final := strings.FieldsFunc(str[1], finalFunc)
			hashTag.TagName = "#" + final[0]
			scrappedNews.HashTags = append(scrappedNews.HashTags, &hashTag)
		} else if strings.Contains(word, "@") { //IF A WORD CONTAINS A USER-HANDEL of@user
			str := strings.Split(word, "@")
			final := strings.FieldsFunc(str[1], finalFunc)
			userHandle.Name = "@" + final[0]
			scrappedNews.UserHandles = append(scrappedNews.UserHandles, &userHandle)
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
