package config

// basic config types
type WorkSpaceConfig struct {
	Name        string
	BasicConfig BasicConfig `yaml:"basic_config"`
}
type BasicConfig struct {
	Prefix string     // Api Prefix, if empty, default is /
	DbList []DbConfig `yaml:"db_list"`
	Auth   AuthConfig
	Bind   string // :8080
}
type DbConfig struct {
	Name    string // unique Db name
	DbType  string `yaml:"db_type"`
	LinkUrl string `yaml:"link_url"`
}
type AuthConfig struct {
	Enable bool `yaml:"enable"`
	User   string
	Pass   string
}

// db api config
type ApiConfigCollection struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ApiList     []ApiConfig `json:"api_list"`
}

type ApiColumn = string
type ApiConfig struct {
	Name            string      `json:"name"`
	Sql             string      `json:"sql"`
	Description     string      `json:"description"`
	DbName          string      `json:"db_name"`
	Columns         []ApiColumn `json:"columns"`
	TransformJsPath string      `json:"transform_js_path,omitempty"`
}

type WorkSpaceParsedConfig struct {
	WorkSpace     string
	Config        *WorkSpaceConfig
	ApiCollection []ApiConfigCollection
}
