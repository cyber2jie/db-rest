package service

import (
	"db-rest/db"
	"strings"
)

func buildResultData(datas []*db.DbData, apiConfig db.DbApiConfig) (*ResultDbData, error) {
	dataStruct := []*DataStruct{}
	dataRows := []*DataRow{}

	for _, column := range strings.Split(apiConfig.Columns, ",") {
		if column != "" {
			dataStruct = append(dataStruct, &DataStruct{
				Name: column,
			})
		}
	}

	for _, data := range datas {
		rowNum := data.RowNum
		var dataRow *DataRow
		for _, row := range dataRows {
			if row.RowNum == rowNum {
				dataRow = row
				break
			}
		}

		if dataRow == nil {
			dataRow = &DataRow{
				RowNum:  rowNum,
				RowData: RowData{},
			}
			dataRows = append(dataRows, dataRow)
		}

		dataRow.RowData[data.Column] = data.Value

	}

	return &ResultDbData{
		DataStructs: dataStruct,
		DataRows:    dataRows,
	}, nil
}
