package service

import "db-rest/util"

type ResultWorkSpace struct {
	Name           string                 `json:"name"`
	Prefix         string                 `json:"prefix"`
	ApiCollections []*ResultApiCollection `json:"api_collections"`
}

type ResultApiCollection struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	ApiConfigs  []*ResultApiConfig `json:"api_configs"`
}
type ResultApiConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Sql         string `json:"sql"`
	DbName      string `json:"db_name"`
	Columns     string `json:"columns"`
}

func ListWorkspace(sc *ServiceContext) (*ResultWorkSpace, error) {
	basicConfig := sc.WorkSpace.Config.BasicConfig

	apiCollection := []*ResultApiCollection{}

	for _, ac := range sc.WorkSpace.ApiCollection {

		apiConfig := []*ResultApiConfig{}
		for _, ac := range ac.ApiList {
			apiConfig = append(apiConfig, &ResultApiConfig{
				Name:        ac.Name,
				Description: ac.Description,
				Sql:         ac.Sql,
				DbName:      ac.DbName,
				Columns:     util.JoinStr(",", ac.Columns...),
			})
		}
		apiCollection = append(apiCollection, &ResultApiCollection{
			Name:        ac.Name,
			Description: ac.Description,
			ApiConfigs:  apiConfig,
		})
	}

	return &ResultWorkSpace{
		Name:           sc.WorkSpace.Config.Name,
		Prefix:         basicConfig.Prefix,
		ApiCollections: apiCollection,
	}, nil
}
