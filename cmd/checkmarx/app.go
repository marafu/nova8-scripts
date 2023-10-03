package checkmarx

import (
	"fmt"
	"log"

	"github.com/marafu/nova8-scripts/cmd/utils"
)

func App(filename string, projectName string) error {

	config, err := utils.ReadConfig(filename)

	token, err := GetRefreshToken(*config)

	if err != nil {
		return err
	}

	projects, err := GetProject(InputProject{
		ProjectName: projectName,
		AccessToken: token.AccessToken,
	}, *config)

	if err != nil {
		return err
	}

	scans, err := GetScan(InputGetScan{
		ProjectID:   projects[0].ID,
		AccessToken: token.AccessToken,
	}, *config)

	results, err := GetSastResult(InputSast{
		ScanID:      scans[0].ID,
		AccessToken: token.AccessToken,
	}, *config)

	for _, result := range results {

		log.Println(result.QueryName)
		for _, node := range result.Nodes {
			log.Println(node.FullName)

		}

		fmt.Println()
	}

	return nil
}
