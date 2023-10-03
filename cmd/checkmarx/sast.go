package checkmarx

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/marafu/nova8-scripts/cmd/utils"
)

type InputSast struct {
	ScanID      string
	AccessToken string
}

type Result struct {
	QueryID       json.Number `json:"queryID"`
	QueryName     string      `json:"queryName"`
	Severity      string      `json:"severity"`
	CweID         int         `json:"cweID"`
	SourceLine    int         `json:"sourceLine"`
	SourceNode    string      `json:"sourceNode"`
	SinkLine      int         `json:"sinkLine"`
	SinkNode      string      `json:"sinkNode"`
	NumberOfNodes int         `json:"numberOfNodes"`
	SimilarityID  int         `json:"similarityID"`
	UniqueID      int         `json:"uniqueID"`
	Nodes         []struct {
		Column     int    `json:"column"`
		FileName   string `json:"fileName"`
		FullName   string `json:"fullName"`
		Length     int    `json:"length"`
		Line       int    `json:"line"`
		MethodLine int    `json:"methodLine"`
		Method     string `json:"method"`
		Name       string `json:"name"`
		DomType    string `json:"domType"`
	} `json:"nodes"`
	Group        string   `json:"group"`
	Compliances  []string `json:"compliances"`
	PathSystemID string   `json:"pathSystemID"`
	ResultHash   string   `json:"resultHash"`
	LanguageName string   `json:"languageName"`
	FirstScanID  string   `json:"firstScanID"`
	FirstFoundAt string   `json:"firstFoundAt"`
	FoundAt      string   `json:"foundAt"`
	Status       string   `json:"status"`
	State        string   `json:"state"`
}

type ResponseResult struct {
	Results    []Result `json:"results"`
	TotalCount int      `json:"totalCount"`
}

func GetSastResult(input InputSast, config utils.Config) ([]Result, error) {
	uri := fmt.Sprintf("%s/api/sast-results?scan-id=%s", config.Checkmarx.BaseUrl, input.ScanID)

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", input.AccessToken))

	proxyUrl, err := url.Parse(config.General.Proxy)

	if err != nil {
		return nil, err
	}

	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode == 200 {
		var response ResponseResult

		err = json.Unmarshal(body, &response)

		if err != nil {
			return nil, err
		}

		return response.Results, nil
	}

	return nil, nil
}
