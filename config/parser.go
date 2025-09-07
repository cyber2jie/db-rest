package config

import (
	"db-rest/util"
	"encoding/json"
	"github.com/go-yaml/yaml"
	"os"
	"path"
	"sort"
	"strings"
)

func ParseWorkSpaceConfig(p string) (*WorkSpaceParsedConfig, error) {

	workSpaceConfigPath := path.Join(p, WORKSPACE_CONFIG)
	apiCollectionDir := path.Join(p, WORKSPACE_API_CONFIG_DIR)

	out, err := os.ReadFile(workSpaceConfigPath)

	if err != nil {
		return nil, err
	}

	workSpaceConfig := WorkSpaceConfig{}

	err = yaml.Unmarshal(out, &workSpaceConfig)
	if err != nil {
		return nil, err
	}

	collection, err := parseCollections(apiCollectionDir)
	if err != nil {
		return nil, err
	}

	if collection == nil {
		collection = []ApiConfigCollection{}
	}

	return &WorkSpaceParsedConfig{
		WorkSpace:     p,
		ApiCollection: collection,
		Config:        &workSpaceConfig,
	}, nil
}

func parseCollections(p string) (apiCollection []ApiConfigCollection, err error) {
	if util.FileExists(p) {

		entry, err := os.ReadDir(p)
		if err != nil {
			return nil, err
		}

		sort.Slice(entry, func(l, r int) bool {
			lInfo, err1 := entry[l].Info()
			rInfo, err2 := entry[r].Info()
			if err1 != nil || err2 != nil {
				return false
			}
			return lInfo.ModTime().After(rInfo.ModTime())
		})

		checkSql := GetEnvValue[bool](VIPER_KEY_CHECK_SQL)
		useTidbParser := GetEnvValue[bool](VIPER_KEY_USE_TIDB_PARSER)

	entryFor:
		for _, e := range entry {

			//dir not support
			if e.IsDir() {
				continue
			}

			if !strings.HasSuffix(e.Name(), ".json") {
				continue
			}

			apiFile := path.Join(p, e.Name())
			out, err := os.ReadFile(apiFile)

			if err != nil {
				util.LogWarn("read api_config file %s occur error: %v", apiFile, err)
				continue
			}
			collection := ApiConfigCollection{}
			err = json.Unmarshal(out, &collection)
			if err != nil {
				util.LogWarn("parse api_config file %s occur error: %v", apiFile, err)
				continue
			}

			if strings.TrimSpace(collection.Name) == "" {
				util.LogWarn("api_config file %s name is empty", apiFile)
				continue
			}

			for _, c := range apiCollection {

				if util.StrEquals(c.Name, collection.Name) {
					continue entryFor
				}
			}

			//check api sql config
			canAppend := true

			if checkSql {
				util.Log("start check sql config is correct")
				for _, api := range collection.ApiList {
					err := CheckSql(api.Sql, useTidbParser)
					if err != nil {
						util.LogWarn("api_config file %s, api %s sql check error: %v", apiFile, api.Name, err)
						canAppend = false
						break
					}
				}
			}

			if canAppend {
				apiCollection = append(apiCollection, collection)
			}
		}
	}
	return apiCollection, nil
}
