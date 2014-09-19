package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type Tweet struct {
	Text string `json:"text"`
}

type Creds struct {
	ConsumerToken  string `json:"consumer_token"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	AccessSecret   string `json:"secret_token"`
}

func getCreds() Creds {
	// Loads the creds out of creds.json
	file, _ := os.Open("creds.json")
	contents, _ := ioutil.ReadAll(file)
	var creds Creds
	json.Unmarshal(contents, &creds)
	return creds
}

func getTweets(username string, count string, consumer oauth.Credentials, access oauth.Credentials) (tweets []Tweet, err error) {
	// Pull a X number of tweets associated with a user.
	// Docs: https://dev.twitter.com/rest/reference/get/statuses/user_timeline
	client := oauth.Client{Credentials: consumer}
	params := url.Values{"screen_name": {username}, "count": {count}}
	resp, err := client.Get(
		http.DefaultClient,
		&access,
		"https://api.twitter.com/1.1/statuses/user_timeline.json",
		params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tweets)
	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func main() {
	// Parse the user out of the command line flag -user.
	username := flag.String(
		"user", "cameronmaske", "The user to retrieve the latest tweet from.")
	flag.Parse()

	// Set up an OAuth client with the various tokens.
	creds := getCreds()
	consumer := oauth.Credentials{
		Token:  creds.ConsumerToken,
		Secret: creds.ConsumerSecret}
	access := oauth.Credentials{
		Token:  creds.AccessToken,
		Secret: creds.AccessSecret}

	// Get the lastest tweets for the username.
	tweets, _ := getTweets(*username, "1", consumer, access)
	fmt.Println(tweets[0])
}
