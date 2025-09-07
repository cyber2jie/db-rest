package db

import "db-rest/util"

type DbDialect interface {
	getCountSql(sql string) string
	getPaginationSql(sql string, columns []string, page, pageSize int) string
}

var (
	sqliteDbDialect        = SqliteDbDialect{}
	mysqlDbDialect         = MysqlDbDialect{}
	postgresDbDialect      = PostgresDbDialect{}
	oracleDbDialect        = OracleDbDialect{}
	oracle11gDbDialect     = Oracle11gDbDialect{}
	sqlServerDbDialect     = SqlServerDbDialect{}
	sqlServer2012DbDialect = SqlServer2012DbDialect{}
)

// SqliteDbDialect
type SqliteDbDialect struct{}

func (dialect SqliteDbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect SqliteDbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	return "select * from (" + sql + ") t limit " + util.ToString(pageSize) + " offset " + util.ToString((page-1)*pageSize)
}

// MysqlDbDialect
type MysqlDbDialect struct{}

func (dialect MysqlDbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect MysqlDbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	return "select * from (" + sql + ") t limit " + util.ToString(pageSize) + " offset " + util.ToString((page-1)*pageSize)
}

// PostgresDbDialect
type PostgresDbDialect struct{}

func (dialect PostgresDbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect PostgresDbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	return "select * from (" + sql + ") t limit " + util.ToString(pageSize) + " offset " + util.ToString((page-1)*pageSize)
}

// Oracle11g以前版本
type Oracle11gDbDialect struct{}

func (dialect Oracle11gDbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect Oracle11gDbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	//第一个列必须可排序，非虚拟列
	if len(columns) > 0 {
		orderColumn := columns[0]
		offset := (page - 1) * pageSize
		return "select t2.* from (select t.*,ROW_NUMBER() OVER (ORDER BY (" + orderColumn + ")) AS RN from (" + sql + ") t ) t2 where t2.RN BETWEEN " + util.ToString(offset+1) + " AND " + util.ToString(offset+pageSize)
	}
	return ""
}

// OracleDbDialect
type OracleDbDialect struct{}

func (dialect OracleDbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect OracleDbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	//第一个列必须可排序，非虚拟列
	if len(columns) > 0 {
		orderColumn := columns[0]
		return "select * from (" + sql + ") t ORDER BY (" + orderColumn + ")  OFFSET " + util.ToString((page-1)*pageSize) + " ROWS  FETCH NEXT " + util.ToString(pageSize) + " ROWS ONLY "
	}
	return ""
}

// SqlServerDbDialect
type SqlServerDbDialect struct{}

func (dialect SqlServerDbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect SqlServerDbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	//第一个列必须可排序，非虚拟列
	orderColumn := "SELECT NULL" //可能出现获取数据
	if len(columns) > 0 {
		orderColumn = columns[0]
	}

	offset := (page - 1) * pageSize
	return "select * from (select *,ROW_NUMBER() OVER (ORDER BY (" + orderColumn + ")) AS RowNum from (" + sql + ") as  t )  as  t2 where t2.RowNum BETWEEN " + util.ToString(offset+1) + " AND " + util.ToString(offset+pageSize)
}

// SqlServer2012DbDialect
type SqlServer2012DbDialect struct{}

func (dialect SqlServer2012DbDialect) getCountSql(sql string) string {
	return "select count(*) as c from (" + sql + ") t"
}
func (dialect SqlServer2012DbDialect) getPaginationSql(sql string, columns []string, page, pageSize int) string {
	//第一个列必须可排序，非虚拟列
	orderColumn := "SELECT NULL" //可能出现获取数据
	if len(columns) > 0 {
		orderColumn = columns[0]
	}
	return "select * from (" + sql + ") as t ORDER BY (" + orderColumn + ")  OFFSET " + util.ToString((page-1)*pageSize) + " ROWS  FETCH NEXT " + util.ToString(pageSize) + " ROWS ONLY "
}
