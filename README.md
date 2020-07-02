# ctfdget

Ever started playing a ctf hosted using CTFd that felt like the requests were
processed manually by a goat? Not being able to access the challenges, but
being able to access the api? Just wanting to download all challenges at once
for not having to deal with CTFd until you want to submit your flags?

Here it is: ctfdget, a simple tool for fetching all challenges with their
included files.

## Building 

```
go build ./...
```

## Usage

```
./ctfdget --help
Usage of ./ctfdget:
  -out string
    	The name of the folder to dump the files to (default "challdump")
  -session string
    	The session (the value of the cookie named 'session') (default "9e8831af-ce30-48c3-8663-4b27262f43f1.pjKPVCYufDhuA9GPJAlc_xh45M8")
  -url string
		The root URL of the CTFd instance (default "https://ctf.example.com")
```

Example:

```
$ go run ./... -session 9e8831af-ce30-48c3-8663-4b27262f43f1.pjKPVCYufDhuA9GPJAlc_xh45M8 -url https://cscml.zenysec.com -out cscmlctf
Fetching all challenges using the ctf api...
Done fetching all challenges
Downloading the files included in the challenges
→ cscmlctf/misc/Censored
→ cscmlctf/web/WebServer I
→ cscmlctf/misc/Thanos
→ cscmlctf/web/Scrolls List
→ cscmlctf/pwn/WebServer II
→ cscmlctf/machine learning/Prosopagnosia
→ cscmlctf/misc/JSGame
…
```

## Features

- Dump all files from all challenges
- A simple directory structure get's created sorting the challenges into the corresponding categories (`<ctfname>/<category>/<challengename>/<challengefiles>`)
- README.md generation containing the challenge description and some other info:

Example `README.md`, `generated in cscmlctf/pwn/Strcmp/README.md`
```md
# Strcmp

Category: pwn

Files:
- chall

## Description

Oh jeez, this is the first time I wrote code, I really hope its good and the stack overflow guys won't be too harsh

nc ctf.cscml.zenysec.com 20005

## Writeup
```

(Writeup section included to encourage you to publish writeups :D)

## Contribution

Just open issues, pull requests or whatever and we'll work something out
