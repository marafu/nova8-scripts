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

type InputProject struct {
	AccessToken string
	ProjectName string
}

type Projects struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	TenantID  string `json:"tenantId"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Groups    []any  `json:"groups"`
	Tags      struct {
		WebGoat string `json:"#WebGoat"`
	} `json:"tags"`
	RepoURL          string `json:"repoUrl"`
	MainBranch       string `json:"mainBranch"`
	Criticality      int    `json:"criticality"`
	PrivatePackage   bool   `json:"privatePackage"`
	ImportedProjName string `json:"imported_proj_name"`
}

type ResponseProject struct {
	TotalCount         int        `json:"totalCount"`
	FilteredTotalCount int        `json:"filteredTotalCount"`
	Projects           []Projects `json:"projects"`
}

func GetProject(input InputProject, config utils.Config) ([]Projects, error) {

	uri := fmt.Sprintf("%s/api/projects?offset=0&limit=1&name=%s", config.Checkmarx.BaseUrl, input.ProjectName)

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

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		var response ResponseProject

		err = json.Unmarshal(body, &response)

		if err != nil {
			return nil, err
		}

		return response.Projects, nil
	}

	return nil, nil

}

func CreateProject(name string) {

}
