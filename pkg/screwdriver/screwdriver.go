package screwdriver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	apiVersion = "v4"
)

type SD struct {
	JWT        string
	baseAPIURL string
	client     *http.Client
}

type tokenResponse struct {
	JWT string `json:"token"`
}

func New(token, baseAPIURL string) *SD {
	sd := new(SD)
	sd.baseAPIURL = baseAPIURL
	sd.client = new(http.Client)
	sd.JWT = sd.jwt(token)

	fmt.Println(sd.JWT)

	return sd
}

func (sd *SD) request(method, url string, body io.Reader, isRequestAuth bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	if isRequestAuth {
		prefix := "Bearer "
		req.Header.Add("Authorization", prefix+sd.JWT)
		req.Header.Add("Content-Type", "application/json")
	}

	return sd.client.Do(req)
}

func (sd *SD) jwt(token string) string {
	apiUrl := sd.makeURL("auth/token")
	q := apiUrl.Query()
	q.Set("api_token", token)
	apiUrl.RawQuery = q.Encode()

	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, false)
	if err != nil {
		logrus.Error(err)
	}
	defer res.Body.Close()

	tr := new(tokenResponse)
	err = json.NewDecoder(res.Body).Decode(tr)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	return tr.JWT
}

func (sd *SD) Build(id string) interface{} {
	apiUrl := sd.makeURL("builds")
	apiUrl.Path = path.Join(apiUrl.Path, id)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response interface{}
	json.Unmarshal(body, &response)

	return response
}

func (sd *SD) makeURL(endpoint string) *url.URL {
	u, err := url.Parse(sd.baseAPIURL)
	if err != nil {
		logrus.Error(err)
	}

	u.Path = path.Join(u.Path, apiVersion, endpoint)

	return u
}

func (sd *SD) Job(jobId string) interface{} {
	apiUrl := sd.makeURL("jobs")
	apiUrl.Path = path.Join(apiUrl.Path, jobId)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response interface{}

	json.Unmarshal(body, &response)

	return response
}

func (sd *SD) Events(eventId string) interface{} {
	apiUrl := sd.makeURL("events")
	apiUrl.Path = path.Join(apiUrl.Path, eventId)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response interface{}

	json.Unmarshal(body, &response)

	return response
}

func (sd *SD) Pipeline(pipelineId string) interface{} {
	apiUrl := sd.makeURL("pipelines")
	apiUrl.Path = path.Join(apiUrl.Path, pipelineId)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response interface{}

	json.Unmarshal(body, &response)

	return response
}
