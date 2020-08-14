package trello

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

const maxRetries = 3

type HTTPHandler struct {
	key   string
	token string
}

func newHTTPHandler(key, token string) *HTTPHandler {
	return &HTTPHandler{
		key:   key,
		token: token,
	}
}

func (h *HTTPHandler) call(method, url string, obj interface{}) error {
	log.Debug("Calling ", url)
	resp, err := h.makeAPICall(url, method)
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&obj)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("%s: reading response. %s", url, err)
		return err
	}
	return nil
}

func (h *HTTPHandler) getJSON(url string, obj interface{}) error {
	return h.call("GET", url, obj)
}

func (h *HTTPHandler) postJSON(URL string, args map[string]string, obj interface{}) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	for k, v := range args {
		u = addQ(u, k, v)
	}

	return h.call("POST", u.String(), obj)
}

func (h *HTTPHandler) addAuth(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	u = addQ(u, "key", h.key)
	u = addQ(u, "token", h.token)
	return u.String(), nil
}

func addQ(u *url.URL, key, value string) *url.URL {
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u
}

// makeAPICall performs an HTTP call to the given url, returning the response
func (h *HTTPHandler) makeAPICall(URL, method string) (*http.Response, error) {
	client := &http.Client{}

	var resp *http.Response
	var errorsList []error
	authUrl, err := h.addAuth(URL)
	if err != nil {
		return nil, err
	}
	for i := 1; i <= maxRetries; i++ {
		req, err := http.NewRequest(method, authUrl, nil)

		r, err := client.Do(req)
		if err != nil {
			log.Debugf("#%d %s: %s\n", i, authUrl, err)
			errorsList = append(errorsList, err)
			if i >= maxRetries {
				for _, e := range errorsList {
					log.Error(e)
				}
				return nil, errors.New("Too many errors")
			}
			// Go on and try again after a little pause
			time.Sleep(2 * time.Second)
			continue
		}

		if r.StatusCode >= 400 {
			errorsList = append(errorsList, errors.New(r.Status))
			if i >= maxRetries {
				for _, e := range errorsList {
					log.Error(e)
				}
				return nil, errors.New("Too many errors")
			}

			if r.StatusCode == 429 {
				// Header Retry-After tells the number of seconds until the end of the current window
				log.Error("Got 429 too many requests, let's try to wait 10 seconds...")
				log.Errorf("Retry-After header: %s\n", r.Header.Get("Retry-After"))
				time.Sleep(10 * time.Second)
			}
			continue
		}
		resp = r
		break

	}
	return resp, nil
}
