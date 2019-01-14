package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"bytes"
	"fmt"
	"io/ioutil"
	"flag"
)

const postURLFormat = "https://%v/me/messages?access_token=%v"
var host string
var token string
var chatId int64
var filterSubstring string
var chatTopicFormat string

type HockeyRelease struct {
	AppVersion struct {
		Version string `json:"version"`
		ShortVersion string `json:"shortversion"`
		Title string `json:"title"`
		Notes string `json:"notes"`
	} `json:"app_version"`
	URL string `json:"url"`
}

type Message struct {
	Text string `json:"text"`
}

type Recipient struct {
	ChatID int64 `json:"chat_id"`
}

type TextMessage struct {
	Message `json:"message"`
	Recipient `json:"recipient"`
}

type ChatControl struct {
	Title string `json:"title"`
}

type ChangeTopicMessage struct {
	ChatControl `json:"chat_control"`
	Recipient `json:"recipient"`
}


func handler(w http.ResponseWriter, r *http.Request)  {
	decoder := json.NewDecoder(r.Body)
	var releaseInfo HockeyRelease
	err := decoder.Decode(&releaseInfo)
	if err != nil {
		fmt.Println("Failed to parse JSON")
		return
	}
	defer r.Body.Close()
	fmt.Println(releaseInfo)

	if strings.Contains(releaseInfo.AppVersion.Notes, filterSubstring) {
		messageText := fmt.Sprintf("%v %v(%v) %v", releaseInfo.AppVersion.Title,
			releaseInfo.AppVersion.ShortVersion, releaseInfo.AppVersion.Version, releaseInfo.URL)


		postURL := fmt.Sprintf(postURLFormat, host, token)
		message := TextMessage { Message {messageText}, Recipient{chatId} }
		postMessage(postURL, message)

		chatTopic := fmt.Sprintf(chatTopicFormat, releaseInfo.AppVersion.ShortVersion, releaseInfo.AppVersion.Version)
		chatTopicMessage := ChangeTopicMessage{ ChatControl{chatTopic}, Recipient{chatId} }
		postMessage(postURL, chatTopicMessage)
	}
}

func postMessage(url string, message interface{}) {
	json, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Failed to write JSON")
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to POST")
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main()  {
	flag.StringVar(&token, "token", "", "bot api token")
	flag.StringVar(&host, "host", "botapi.tamtam.chat", "bot api host")
	flag.Int64Var(&chatId, "chat", 0, "destination chat id");
	flag.StringVar(&chatTopicFormat, "topic", "Integration version %v (%v)", "chat topic format")
	flag.StringVar(&filterSubstring, "substring", "Branch: integration", "substring search string in app notes")
	flag.Parse()

	if len(host) > 0 && len(token) > 0 && chatId != 0 {
		http.HandleFunc("/ok-hockeyapp", handler)
		http.ListenAndServe(":8080", nil)
	}
}
