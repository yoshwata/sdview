package screwdriver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	apiVersion        = "v4"
	authEndpoint      = "auth/token"
	buildsEndpoint    = "builds"
	jobsEndpoint      = "jobs"
	pipelinesEndpoint = "pipelines"
)

type SD struct {
	JWT        string
	baseAPIURL string
	client     *http.Client
}

type Job struct {
	Name       string `json:"name"`
	PipelineId int    `json:"pipelineId"`
}

type Build struct {
	JobID     int    `json:"jobId"`
	Container string `json:"container"`
}

type Pipeline struct {
	Name    string `json:"name"`
	ScmRepo struct {
		Name    string `json:"name"`
		Branch  string `json:"branch"`
		URL     string `json:"url"`
		RootDir string `json:"rootDir"`
	} `json:"scmRepo"`
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
	apiUrl := sd.makeURL(authEndpoint)
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

// func (sd *SD) Build(id string) *Build {
func (sd *SD) Build(id string) (interface{}, error) {
	apiUrl := sd.makeURL(buildsEndpoint)
	apiUrl.Path = path.Join(apiUrl.Path, id)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
		// os.Exit(1)
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response interface{}
	// strBody, _ := json.Marshal(res.Body)
	// byteBody := []byte(strBody)

	// var unMarshaledBody interface{}
	// json.Unmarshal(byteBody, &unMarshaledBody)
	json.Unmarshal(body, &response)

	fmt.Println("=================")
	// fmt.Printf("%#v\n", response)

	// pathedBi, _ := jsonpath.Read(response, "$.buildClusterName")
	// fmt.Printf("%s\n", pathedBi)
	// _ = pathedPod

	return response, nil
	// build := new(Build)
	// err = json.NewDecoder(res.Body).Decode(build)
	// if err != nil {
	// 	logrus.Error(err)
	// 	os.Exit(1)
	// }

	// return build
}

func (sd *SD) makeURL(endpoint string) *url.URL {
	u, err := url.Parse(sd.baseAPIURL)
	if err != nil {
		logrus.Error(err)
	}

	u.Path = path.Join(u.Path, apiVersion, endpoint)

	return u
}

func (sd *SD) Job(jobId string) (interface{}, error) {
	apiUrl := sd.makeURL(jobsEndpoint)
	apiUrl.Path = path.Join(apiUrl.Path, jobId)
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var response interface{}

	json.Unmarshal(body, &response)

	// job := new(Job)
	// err = json.NewDecoder(res.Body).Decode(job)
	// if err != nil {
	// 	logrus.Error(err)
	// }

	return response, nil
}

func (sd *SD) Events(eventId string) (interface{}, error) {
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

	return response, nil
}

func (sd *SD) Pipeline(pipelineId int) *Pipeline {
	apiUrl := sd.makeURL(pipelinesEndpoint)
	apiUrl.Path = path.Join(apiUrl.Path, strconv.Itoa(pipelineId))
	res, err := sd.request(http.MethodGet, apiUrl.String(), nil, true)
	if err != nil {
		logrus.Error(err)
	}

	defer res.Body.Close()
	pipeline := new(Pipeline)
	err = json.NewDecoder(res.Body).Decode(pipeline)
	if err != nil {
		logrus.Error(err)
	}

	return pipeline
}
