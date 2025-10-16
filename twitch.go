package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

const (
	GraphqlURL = "https://gql.twitch.tv/gql"
	UserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36"
)

type TwitchRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type TwitchIsOnlineResponse struct {
	Data struct {
		User struct {
			Stream struct {
				ID string `json:"id"`
			} `json:"stream"`
		} `json:"user"`
	} `json:"data"`
}

type TwitchLastVodResponse struct {
	Data struct {
		User struct {
			Videos struct {
				Edges []struct {
					Node struct {
						ID string `json:"id"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"videos"`
		} `json:"user"`
	} `json:"data"`
}

func IsOnline(channel string) (bool, error) {
	requestPayload := TwitchRequest{
		Query: "query($login: String!){ user(login: $login){ stream { id } } }",
		Variables: map[string]any{
			"login": channel,
		},
	}
	requestBody, _ := json.Marshal(requestPayload)

	request, err := http.NewRequest("POST", GraphqlURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return false, err
	}
	request.Header.Set("Client-ID", Config.TwitchClientID)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	responseData := &TwitchIsOnlineResponse{}
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		return false, err
	}

	return responseData.Data.User.Stream.ID != "", nil
}

func GetLatestVODUrl(channel string) (string, error) {
	requestPayload := TwitchRequest{
		Query: "query($login: String!, $first: Int!){ user(login: $login){ videos(types: [ARCHIVE], first: $first) { edges { node {id} } } } }",
		Variables: map[string]any{
			"login": channel,
			"first": 1,
		},
	}

	requestBody, _ := json.Marshal(requestPayload)

	request, err := http.NewRequest("POST", GraphqlURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	request.Header.Set("Client-ID", Config.TwitchClientID)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	responseData := &TwitchLastVodResponse{}
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		return "", err
	}

	if len(responseData.Data.User.Videos.Edges) == 0 {
		return "", fmt.Errorf("no vods found")
	}

	return "https://www.twitch.tv/videos/" + responseData.Data.User.Videos.Edges[0].Node.ID, nil
}

func DownloadVOD(url string, height uint) error {
	formatArg := fmt.Sprintf("--format=best[height<=%d]", height)

	cmd := exec.Command(
		"yt-dlp",
		"--paths", "home:"+DOWNLOAD_DIR,
		"--paths", "temp:./tmp",
		"--live-from-start",
		formatArg,
		url,
	)

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("yt-dlp exited with error: %s", exitErr.Stderr)
		}
		return fmt.Errorf("failed to execute yt-dlp: %w", err)
	}

	return nil
}
