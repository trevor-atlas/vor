package git

import (
	"github.com/trevor-atlas/vor/utils"
	"io/ioutil"
	"fmt"
	"bytes"
	"net/http"
)

func Post (url string, json []byte) {
	githubAPIKey := utils.GetStringEnv("github.apikey")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
    req.Header.Set("Authorization", "token " + githubAPIKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}