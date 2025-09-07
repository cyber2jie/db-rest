package service

type DbQuery struct {
	Collection string
	Table      string
	Form       *DbQueryForm
}

type DbQueryForm struct {
	Page     int        `form:"page" json:"page"`
	PageSize int        `form:"pageSize" json:"pageSize"`
	Query    *QueryForm `form:"query" json:"query"`
}

type QueryForm struct {
	QueryType string   `json:"query_type"`
	Queries   []*Query `json:"queries"`
}

type Query struct {
	Field string `json:"field"`
	Op    string `json:"op"`
	Value string `json:"value"`
}

type ListResult struct {
	Total int64        `json:"total"`
	Data  ResultDbData `json:"data"`
}

type ResultDbData struct {
	DataStructs []*DataStruct `json:"data_structs"`
	DataRows    []*DataRow    `json:"data_rows"`
}

type DataStruct struct {
	Name string `json:"name"`
}
type DataRow struct {
	RowNum  int     `json:"row_num"`
	RowData RowData `json:"row_data"`
}

type RowData = map[string]string
