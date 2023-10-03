package checkmarx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/marafu/nova8-scripts/cmd/utils"
)

type InputScan struct {
	ProjectID                     int    `json:"projectId"`
	OverrideProjectSetting        string `json:"overrideProjectSetting"`
	IsIncremental                 string `json:"isIncremental"`
	IsPublic                      string `json:"isPublic"`
	ForceScan                     string `json:"forceScan"`
	Comment                       string `json:"comment"`
	PresetID                      int    `json:"presetId"`
	EngineConfigurationID         int    `json:"engineConfigurationId"`
	CustomFields                  string `json:"customFields"`
	PostScanActionID              int    `json:"postScanActionId"`
	RunPostScanOnlyWhenNewResults string `json:"runPostScanOnlyWhenNewResults"`
	RunPostScanMinSeverity        string `json:"runPostScanMinSeverity"`
	PostScanActionArguments       string `json:"postScanActionArguments"`
	ZippedSource                  string `json:"zippedSource"`
}

func CreateScan(input InputScan, config utils.Config) error {

	uri := fmt.Sprintf("%s/cxrestapi/help/sast/scanWithSettings", config.Checkmarx.BaseUrl)

	data, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json; version=1.0")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", config.Checkmarx.AccessToken))

	proxyUrl, err := url.Parse(config.General.Proxy)

	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		var response ResponseProject

		err = json.Unmarshal(body, &response)

		if err != nil {
			return err
		}

	}

	log.Println(string(body))

	return nil
}

type InputGetScan struct {
	ProjectID   string
	AccessToken string
}

type ScanType struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	StatusDetails []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Details string `json:"details"`
		Loc     int    `json:"loc,omitempty"`
	} `json:"statusDetails"`
	Branch      string `json:"branch"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	ProjectID   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	UserAgent   string `json:"userAgent"`
	Initiator   string `json:"initiator"`
	Tags        struct {
	} `json:"tags"`
	Metadata struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Handler struct {
			UploadHandler struct {
				Branch    string `json:"branch"`
				UploadURL string `json:"upload_url"`
			} `json:"UploadHandler"`
		} `json:"Handler"`
		Configs []struct {
			Type  string `json:"type"`
			Value struct {
				Incremental string `json:"incremental"`
			} `json:"value,omitempty"`
		} `json:"configs"`
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
		CreatedAt struct {
			Nanos   int `json:"nanos"`
			Seconds int `json:"seconds"`
		} `json:"created_at"`
	} `json:"metadata"`
	Engines      []string `json:"engines"`
	SourceType   string   `json:"sourceType"`
	SourceOrigin string   `json:"sourceOrigin"`
}

type ResponseScan struct {
	TotalCount         int        `json:"totalCount"`
	FilteredTotalCount int        `json:"filteredTotalCount"`
	Scans              []ScanType `json:"scans"`
}

func GetScan(input InputGetScan, config utils.Config) ([]ScanType, error) {

	uri := fmt.Sprintf("%s/api/scans/?project-id=%s", config.Checkmarx.BaseUrl, input.ProjectID)

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json; version=1.0")
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

	body, _ := io.ReadAll(res.Body)

	defer res.Body.Close()

	if res.StatusCode == 200 {
		var response ResponseScan

		err = json.Unmarshal(body, &response)

		if err != nil {
			return nil, err
		}

		return response.Scans, nil
	}

	return nil, nil

}

type InputGetScanDetails struct {
	ScanID      string
	AccessToken string
}

func GetScanDetails(input InputGetScanDetails, config utils.Config) error {

	uri := fmt.Sprintf("%s/api/scans/%s", config.Checkmarx.BaseUrl, input.ScanID)

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json; version=1.0")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", input.AccessToken))

	proxyUrl, err := url.Parse(config.General.Proxy)

	if err != nil {
		return err
	}

	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl), TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		var response ResponseProject

		err = json.Unmarshal(body, &response)

		if err != nil {
			return err
		}

	}

	log.Println(string(body))

	return nil

}
