package db

import (
	"testing"
)

func TestGetVm(t *testing.T) {
	script := `function transform(row){
    for (var i = 0; i < row.columns.length; i++) {
        var column = row.columns[i];
        if (column.name == "name"){
            column.value=row.row_num+":"+column.value.toUpperCase();
        }
    }
}`
	_, fn, err := GetVm(script)
	if err != nil {
		t.Error(err)
	}

	rows := []*DataRow{
		&DataRow{
			Columns: []*DataColumn{
				{Name: "name", Value: "cyber"},
				{Name: "age", Value: "18"},
			},
			RowNum: 1,
		},
		&DataRow{
			Columns: []*DataColumn{
				{Name: "name", Value: "bob"},
				{Name: "age", Value: "20"},
			},
			RowNum: 2,
		},
	}

	for _, row := range rows {
		fn(row)
	}

	if rows[0].Columns[0].Value != "1:CYBER" {
		t.Error("transform failed")
	}
	if rows[1].Columns[0].Value != "2:BOB" {
		t.Error("transform failed")
	}
}

func TestGetVm2(t *testing.T) {
	script := `function transform(row){
    for (var i = 0; i < row.columns.length; i++) {
        var column = row.columns[i];
        if (column.name == "createAt"){
            column.value=new Date().toLocaleDateString()
        }
    }
}`
	_, fn, err := GetVm(script)
	if err != nil {
		t.Error(err)
	}

	rows := []*DataRow{
		&DataRow{
			Columns: []*DataColumn{
				{Name: "name", Value: "cyber"},
				{Name: "age", Value: "18"},
				{Name: "createAt", Value: ""},
			},
			RowNum: 1,
		},
	}

	for _, row := range rows {
		fn(row)
	}
	t.Logf("%s", rows[0].Columns[2].Value)
}

func TestGetVm3(t *testing.T) {
	script := `function transform(row) {
    for (var i = 0; i < row.columns.length; i++) {
        var column = row.columns[i];
        if (column.value) {
            column.value = column.value.toUpperCase()
        }
    }
}`
	_, fn, err := GetVm(script)
	if err != nil {
		t.Error(err)
	}

	rows := []*DataRow{
		&DataRow{
			Columns: []*DataColumn{
				{Name: "name", Value: "cyber"},
				{Name: "pwd", Value: "q18B"},
				{Name: "location", Value: "Unknow"},
			},
			RowNum: 1,
		},
	}

	for _, row := range rows {
		fn(row)
	}
	if rows[0].Columns[0].Value != "CYBER" {
		t.Error("transform failed")
	}
}
