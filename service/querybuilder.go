package service

import (
	"db-rest/config"
	"db-rest/db"
	"db-rest/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

var maxQueryCount = config.GetEnvValue[int](config.VIPER_KEY_MAX_QUERY_CONDITION)

type QueryOperator = func(query *Query, queryType string) (string, []any)

var queryOperators = map[string]QueryOperator{
	"eq":      eq,
	"ne":      ne,
	"null":    null,
	"notnull": notNull,
	"like":    like,
	"in":      in,
	"nin":     nin,
	"gt":      gt,
	"lt":      lt,
}

const (
	or  = "or"
	and = "and"
)

func isOr(queryType string) bool {
	return or == queryType
}

func buildQueryDb(Db *gorm.DB, config *db.DbApiConfig, query *QueryForm) (_db *gorm.DB, err error) {
	_db = Db

	if query != nil {

		columns := strings.Split(config.Columns, ",")
		queryType := query.QueryType

		if strings.ToLower(queryType) == and {
			queryType = and
		} else {
			queryType = or
		}

		queryCount := 0

		sqlConditions := []string{}
		sqlArgs := []any{}

		for _, query := range query.Queries {

			if queryCount > maxQueryCount {
				return nil, errors.New("too many query conditions")
			}

			if query.Field == "" {
				continue
			}
			field := query.Field

			contain := false

			for _, column := range columns {
				if strings.ToLower(field) == strings.ToLower(column) {
					contain = true
					break
				}
			}
			if !contain {
				util.LogWarn("query column not exist: %s", field)
				continue
			}

			if query.Op != "" {

				op := strings.ToLower(query.Op)

				operator := queryOperators[op]

				if operator != nil {
					sql, args := operator(query, queryType)
					sqlConditions = append(sqlConditions, sql)
					sqlArgs = append(sqlArgs, args...)
				}
				queryCount++
			}
		}

		if len(sqlConditions) > 0 {
			sql := strings.Join(sqlConditions, " "+queryType+" ")

			if queryType == and {
				_db = _db.Where(sql, sqlArgs...)
			} else {
				// or 需要使用row_num实现
				selectRowNumSql := fmt.Sprintf(" ( row_num  in (select distinct row_num from workspace_db_data where (%s)) )", sql)
				_db = _db.Where(selectRowNumSql, sqlArgs...)
			}

		}

	}

	return _db, err
}

// operators
// 使用子查询处理and类型，大数据量性能较差
func fieldToExist(field string, condition string) string {
	return fmt.Sprintf(`( EXISTS (
	select 1 from workspace_db_data w where w.row_num=workspace_db_data.row_num  and  w.column='%s' and ( w.value %s )
		) )`, field, condition)
}

func fieldToDbColumn(field string, condition string) string {
	return fmt.Sprintf(" ( column='%s' and value %s )", field, condition)
}

func eq(query *Query, queryType string) (string, []any) {
	v := query.Value

	if isOr(queryType) {

		return fieldToDbColumn(query.Field, " = ?"), []any{v}

	}

	return fieldToExist(query.Field, " = ?"), []any{v}

}
func ne(query *Query, queryType string) (string, []any) {
	v := query.Value
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " != ?"), []any{v}
	}
	return fieldToExist(query.Field, " != ?"), []any{v}
}
func null(query *Query, queryType string) (string, []any) {
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " IS NULL"), nil
	}
	return fieldToExist(query.Field, " IS NULL"), nil
}
func notNull(query *Query, queryType string) (string, []any) {
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " IS NOT NULL"), nil
	}
	return fieldToExist(query.Field, " IS NOT NULL"), nil
}
func like(query *Query, queryType string) (string, []any) {
	v := query.Value
	if v == "" {
		return "", nil
	}
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " LIKE ?"), []any{"%" + v + "%"}
	}
	return fieldToExist(query.Field, " LIKE ?"), []any{"%" + v + "%"}
}
func in(query *Query, queryType string) (string, []any) {
	v := strings.Split(query.Value, ",")
	if len(v) == 0 {
		return "", nil
	}
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " IN ?"), []any{v}
	}
	return fieldToExist(query.Field, " IN ?"), []any{v}
}
func nin(query *Query, queryType string) (string, []any) {
	v := strings.Split(query.Value, ",")
	if len(v) == 0 {
		return "", nil
	}
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " NOT IN ?"), []any{v}
	}
	return fieldToExist(query.Field, " NOT IN ?"), []any{v}
}

func gt(query *Query, queryType string) (string, []any) {
	v := query.Value

	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " > ?"), []any{v}
	}

	return fieldToExist(query.Field, " > ?"), []any{v}
}
func lt(query *Query, queryType string) (string, []any) {
	v := query.Value
	if isOr(queryType) {
		return fieldToDbColumn(query.Field, " < ?"), []any{v}
	}
	return fieldToExist(query.Field, " < ?"), []any{v}
}
