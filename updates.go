package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

var versionSlug = fmt.Sprintf("%s-%s", Version, GitCommit)

type updateManifest struct {
	LatestSlug string `json:"latest"`
	Commit     string `json:"commit"`
	FullCommit string `json:"full_commit"`
	Version    string `json:"version"`
}

func checkUpdate() bool {
	resp, err := http.Get("https://wonky.i-hate.science/index.json")
	if err != nil {
		fmt.Printf("error! failed to check for update :(\n%v\n", err)
		return false
	}
	defer resp.Body.Close()

	var m updateManifest
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		fmt.Printf("error! failed to decode update manifest :(\n%v\n", err)
		return false
	}

	if versionSlug != m.LatestSlug {
		fmt.Printf("It seems there is an update!\n Current version: %s | Latest version: %s\n", versionSlug, m.LatestSlug)
		fmt.Println("You can download the latest version here:")
		switch runtime.GOOS {
		case "darwin":
			fmt.Printf("https://wonky.i-hate.science/%s/Wonky-Shell-mac64\n", m.LatestSlug)
		case "windows":
			if runtime.GOARCH == "386" {
				fmt.Printf("https://wonky.i-hate.science/%s/Wonky-Shell-win32.exe\n", m.LatestSlug)
			} else if runtime.GOARCH == "amd64" {
				fmt.Printf("https://wonky.i-hate.science/%s/Wonky-Shell-win64.exe\n", m.LatestSlug)
			} else {
				fmt.Printf("Unkown windows arch: %s\n", runtime.GOOS)
			}
		case "linux":
			if runtime.GOARCH == "386" {
				fmt.Printf("https://wonky.i-hate.science/%s/Wonky-Shell-lin32\n", m.LatestSlug)
			} else if runtime.GOARCH == "amd64" {
				fmt.Printf("https://wonky.i-hate.science/%s/Wonky-Shell-lin64\n", m.LatestSlug)
			} else {
				fmt.Printf("Unkown linux arch: %s\n", runtime.GOOS)
			}
		}
		return true
	} else {
		return false
	}
}
