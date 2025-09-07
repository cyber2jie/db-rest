package config

import "testing"

func TestParseSql(t *testing.T) {
	stmt, err := parseSql("select name,age from user")
	if err != nil {
		t.Error(err)
	}
	if !IsSelectStatement(*stmt) {
		t.Error("not select statement")
	}
}

func TestParseSql2(t *testing.T) {
	stmt, err := parseSql("select * from (select name,age from user) t")
	if err != nil {
		t.Error(err)
	}
	if !IsSelectStatement(*stmt) {
		t.Error("not select statement")
	}
}

func TestCheckSql(t *testing.T) {
	err := CheckSql("update user set name = 'ab' where id=? ", true)
	if err != nil {
		t.Error(err)
	}

}
