package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/h2non/gentleman.v2"
    "gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

var session = flag.String("session", "9e8831af-ce30-48c3-8663-4b27262f43f1.pjKPVCYufDhuA9GPJAlc_xh45M8", "The session (the value of the cookie named 'session')")
var rootURL = flag.String("url", "https://ctf.example.com", "The root URL of the CTFd instance")
var outputFolder = flag.String("out", "challdump", "The name of the folder to dump the files to")

func main() {
	flag.Parse()

	// fetch the list of all challenges
	challs, err := fetchAllChallenges()
	if err != nil {
		panic(err)
	}

	// iterate over all challenges downloading all files
	fmt.Println("Downloading the files included in the challenges")
	for _, chall := range challs.Data {

		// define where to store the challenge
		filepath := fmt.Sprintf("%s/%s/%s", *outputFolder, chall.Category, chall.Name)
		fmt.Printf("→ %s\n", filepath)
		err := os.MkdirAll(filepath, os.ModePerm) // create the directory
		if err != nil {
			fmt.Println(err)
		}

		// fetch the challenge information
		challenge, err := fetchChallenge(chall.ID)
		if err != nil {
			fmt.Println(err)
		}

		// download all files
		for _, file := range challenge.Data.Files {
			err := Download(file, filepath)
			if err != nil {
				fmt.Println(err)
			}
		}

		// store the description of the challenge in a README.md file
		err = saveDescription(challenge, filepath)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// fetchAllChallenges fetches the list of all challs using the ctfs api.
func fetchAllChallenges() (Challenges, error) {
	fmt.Println("Fetching all challenges using the ctf api...")
	cli := gentleman.New()
	cli.URL(*rootURL)

	// define the timeouts outrageously long, as some CTFs hosted using CTFd are incredibly inresponsive.
	cli.Use(timeout.Request(1000 * time.Second))
	cli.Use(timeout.Dial(1000 * time.Second, 2000 * time.Second))

	req := cli.Request()
	req.Path("/api/v1/challenges")

	req.SetHeader("Cookie", fmt.Sprintf("session=%s", *session))

	// Perform the request
	res, err := req.Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return Challenges{}, err
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
		return Challenges{}, err
	}

	// unmarshal the resulting json into a Challenges struct
	var challenges Challenges
	if err := json.Unmarshal(res.Bytes(), &challenges); err != nil {
		return Challenges{}, err
	}
	fmt.Println("Done fetching all challenges")
	return challenges, nil
}

// fetchChallenge fetches a single challenge
func fetchChallenge(id int) (Challenge, error) {
	cli := gentleman.New()
	cli.URL(*rootURL)

	req := cli.Request()
	req.Path(fmt.Sprintf("/api/v1/challenges/%d", id))

	req.SetHeader("Cookie", fmt.Sprintf("session=%s", *session))

	// Perform the request
	res, err := req.Send()
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return Challenge{}, err
	}
	if !res.Ok {
		fmt.Printf("Invalid server response: %d\n", res.StatusCode)
		return Challenge{}, err
	}

	var challenge Challenge
	if err := json.Unmarshal(res.Bytes(), &challenge); err != nil {
		return Challenge{}, err
	}
	return challenge, nil
}

// Download downloads a file from the given URL and stores it at the given
// filepath
func Download(url string, filepath string) error {

	// So what the code below does, is it extracts the filename from the url by
	// first splitting the url at slashed and then at questionmarks (could have
	// used regex, but this is CTF code ¯\_(ツ)_/¯)
	a := strings.Split(url, "/")
	b := strings.Split(a[len(a)-1], "?")
	filename := b[0]

	prefix := *rootURL
	fullurl := fmt.Sprintf("%s%s", prefix, url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", fullurl, nil)
	req.Header.Add("Cookie", fmt.Sprintf("session=%s", *session))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/%s", filepath, filename))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func saveDescription(challenge Challenge, filepath string) error {
	path := fmt.Sprintf("%s/README.md", filepath)
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	// fill the readme with some content
	f.WriteString(fmt.Sprintf("# %s\n\n", challenge.Data.Name))
	f.WriteString(fmt.Sprintf("Category: %s\n\n", challenge.Data.Category))

	f.WriteString("Files:\n")
	for _, file := range challenge.Data.Files {

		// So what the code below does, is it extracts the filename from the url by
		// first splitting the url at slashed and then at questionmarks (could have
		// used regex, but this is CTF code ¯\_(ツ)_/¯)
		a := strings.Split(file, "/")
		b := strings.Split(a[len(a)-1], "?")
		filename := b[0]
		f.WriteString(fmt.Sprintf("- %s\n", filename))
	}
	f.WriteString("\n")
	f.WriteString("## Description\n\n")
	f.WriteString(fmt.Sprintf("%s\n\n", challenge.Data.Description))
	f.WriteString("## Writeup")

	return nil
}
