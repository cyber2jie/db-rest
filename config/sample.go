package config

var simpleDbList = []DbConfig{
	{Name: "sqlite", DbType: "sqlite", LinkUrl: Sqlite_LinkUrl},
	{Name: "mysql", DbType: "mysql", LinkUrl: Mysql_LinkUrl},
	{Name: "postgres", DbType: "postgres", LinkUrl: Pgsql_LinkUrl},
	{Name: "oracle", DbType: "oracle", LinkUrl: Oracle_LinkUrl},
	{Name: "sqlserver", DbType: "sqlserver", LinkUrl: Sqlserver_LinkUrl},
	{Name: "mysql1", DbType: "mysql", LinkUrl: Mysql_LinkUrl},
	{Name: "mysql2", DbType: "mysql", LinkUrl: Mysql_LinkUrl},
}

var simpleWorkSpaceConfig = &WorkSpaceConfig{
	Name: "simple",
	BasicConfig: BasicConfig{
		Prefix: "/api",
		Bind:   BIND,
		DbList: simpleDbList,
		Auth: AuthConfig{
			Enable: true,
			User:   "admin",
			Pass:   "admin",
		},
	},
}

var simpleApiConfigCollection = &ApiConfigCollection{
	Name:        "user",
	Description: "simple api collection",
	ApiList: []ApiConfig{
		{
			Name:        "users",
			Description: "users api",
			Columns:     []ApiColumn{"id", "name"},
			Sql:         "select id, name from users",
			DbName:      "sqlite",
		},
		{
			Name:        "contacts",
			Description: "contacts api",
			Columns:     []ApiColumn{"id", "user_id", "contact_name", "phone", "email"},
			Sql:         "select id,user_id,contact_name,phone,email from contacts",
			DbName:      "mysql",
		},
	},
}

var simpleApiConfigCollection2 = &ApiConfigCollection{
	Name:        "fruits",
	Description: "simple api collection",
	ApiList: []ApiConfig{
		{
			Name:            "fruit",
			Description:     "fruit api",
			Columns:         []ApiColumn{"id", "name"},
			Sql:             "select id, name from fruits",
			DbName:          "mysql2",
			TransformJsPath: "transform-js/simple.js",
		},
	},
}
