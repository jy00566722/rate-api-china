package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const baiduTransAPIBaseURL = "https://fanyi-api.baidu.com/api/trans/vip/translate"

// TransResult represents the translation result in the Baidu API response.
type TransResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

// BaiduTransResponse represents the structure of the response from the Baidu translate API.
type BaiduTransResponse struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []TransResult `json:"trans_result"`
	ErrorCode   string        `json:"error_code"`
	ErrorMsg    string        `json:"error_msg"`
}

// translate uses the Baidu Translate API to translate text from one language to another.
func translate(text string, fromLang string, toLang string) (string, error) {
	appID := os.Getenv("BAIDU_APP_ID")         // Set BAIDU_APP_ID environment variable before calling
	secretKey := os.Getenv("BAIDU_SECRET_KEY") // Set BAIDU_SECRET_KEY environment variable before calling

	if appID == "" || secretKey == "" {
		return "", errors.New("appid or secret key not set in the environment variables")
	}

	salt := strconv.FormatInt(time.Now().UnixNano(), 10)
	signStr := appID + text + salt + secretKey
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signStr)))

	// Prepare the query parameters.
	data := url.Values{
		"q":     {text},
		"from":  {fromLang},
		"to":    {toLang},
		"appid": {appID},
		"salt":  {salt},
		"sign":  {sign},
	}

	// Make the request.
	response, err := http.PostForm(baiduTransAPIBaseURL, data)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", response.Status)
	}

	// Read the body.
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Parse the response.
	var resp BaiduTransResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}

	if resp.ErrorCode != "" {
		return "", fmt.Errorf("error from Baidu Translate API: %s,%s", resp.ErrorMsg, resp.ErrorCode)
	}

	// Return the translation result.
	if len(resp.TransResult) > 0 {
		return resp.TransResult[0].Dst, nil
	}

	return "", errors.New("no translation found in API response")
}
