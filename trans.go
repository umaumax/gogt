package gogt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/tidwall/gjson"
)

type Client struct {
	vq      string
	vqMutex sync.RWMutex
}

func (c *Client) CacheVq() (err error) {
	_, err = c.getVq()
	return
}

func (c *Client) getVq() (vq string, err error) {
	c.vqMutex.RLock()
	if c.vq != "" {
		vq = c.vq
		c.vqMutex.RUnlock()
		return
	}
	c.vqMutex.RUnlock()
	c.vqMutex.Lock()
	defer c.vqMutex.Unlock()

	var resp *http.Response
	if c.vq == "" {
		resp, err = http.Get("https://translate.google.com")
		if err != nil {
			return
		}
		defer resp.Body.Close()
		var data []byte
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		c.vq, err = getVq(string(data))
		if err != nil {
			return
		}
		vq = c.vq
	}
	return
}

func (c *Client) TranslateToEn(text string) (ret string, err error) {
	return c.Translate("en", text)
}
func (c *Client) TranslateToJa(text string) (ret string, err error) {
	return c.Translate("ja", text)
}

func (c *Client) Translate(lang string, text string) (ret string, err error) {
	var resp *http.Response
	vq, err := c.getVq()
	if err != nil {
		return
	}
	var tk string
	tk, err = calTK(vq, text)
	if err != nil {
		return
	}

	u, _ := url.Parse("https://translate.google.com/translate_a/single")
	param := u.Query()
	param.Set("client", "t")
	param.Set("sl", "auto")
	param.Set("tl", lang)
	param.Set("dt", "t")
	param.Set("ie", "UTF-8")
	param.Set("oe", "UTF-8")
	param.Set("source", "btn")
	param.Set("ssel", "3")
	param.Set("tsel", "3")
	param.Set("kc", "0")
	param.Set("tk", tk)
	param.Set("q", text)
	u.RawQuery = param.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return
	}
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36"
	req.Header.Set("User-Agent", ua)
	client := new(http.Client)
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	jr := gjson.Parse(string(data))
	if !jr.Get("..0.0.0").Exists() {
		err = fmt.Errorf("Invalid response\n%s", string(data))
		return
	}
	jsonResult := jr.Get("..0.0").Array()
	for _, r := range jsonResult {
		ret += r.Get("..0.0").String()
	}
	return
}
