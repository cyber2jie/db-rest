package config

import (
	"db-rest/util"
	"github.com/go-yaml/yaml"
	"io"
	"os"
	"path"
)

func InitWorkSpace(p string) error {
	//clean path
	workspace_db_path := path.Join(p, WORKSPACE_DB_NAME)
	workspace_config_path := path.Join(p, WORKSPACE_CONFIG)
	workspace_api_config_path := path.Join(p, WORKSPACE_API_CONFIG_DIR)

	if util.FileExists(p) {
		util.RemoveIfExists(workspace_db_path)
		util.RemoveIfExists(workspace_config_path)
		util.RemoveIfExists(workspace_api_config_path)
	} else {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//init config simple

	out, err := yaml.Marshal(simpleWorkSpaceConfig)
	if err != nil {
		return err
	}

	os.WriteFile(workspace_config_path, out, os.ModePerm)

	if !util.FileExists(workspace_api_config_path) {
		err := os.MkdirAll(workspace_api_config_path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	api1 := path.Join(workspace_api_config_path, "api1.json.simple")
	api2 := path.Join(workspace_api_config_path, "api2.json.simple")
	api1Out, err := util.PrettyJsonMarshal(&simpleApiConfigCollection)
	if err != nil {
		return err
	}
	os.WriteFile(api1, api1Out, os.ModePerm)
	api2Out, err := util.PrettyJsonMarshal(&simpleApiConfigCollection2)
	if err != nil {
		return err
	}
	os.WriteFile(api2, api2Out, os.ModePerm)

	//write embed file
	writeEmbedFile("internal", p)

	return nil
}

func writeEmbedFile(dirName, p string) {
	entry, err := SIMPLE.ReadDir(dirName)
	if err != nil {
		util.LogError("read embed file error: %v", err)
		return
	}
	for _, e := range entry {
		if e.IsDir() {
			nextP := path.Join(p, e.Name())
			if !util.FileExists(nextP) {
				os.MkdirAll(nextP, os.ModePerm)
			}
			writeEmbedFile(path.Join(dirName, e.Name()), nextP)
			continue
		}
		file, err := SIMPLE.Open(path.Join(dirName, e.Name()))
		if err != nil {
			util.LogError("open embed file error: %v", err)
			continue
		}
		defer file.Close()

		out, _ := os.OpenFile(path.Join(p, e.Name()), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			util.LogError("write file error: %v", err)
			continue
		}
	}

}
