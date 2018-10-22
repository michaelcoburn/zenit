package tables

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type Table struct {
	schema    string
	table     string
	size      float64
	rows      float64
	increment float64
}

const QuerySQLTable = `
SELECT table_schema AS 'schema',
       table_name AS 'table',
       data_length + index_length AS 'size',
       table_rows AS 'rows',
       auto_increment AS 'increment'
FROM information_schema.tables
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
ORDER BY table_schema, table_name;
`

func Collect() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Tables - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(QuerySQLTable)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Tables - Impossible to execute query: " + err.Error())
	}

	var a = metrics.Load()

	for rows.Next() {
		var t Table

		rows.Scan(
			&t.schema,
			&t.table,
			&t.size,
			&t.rows,
			&t.increment)

		a.Add(metrics.Metric{
			Key: "mysql_stats_tables",
			Tags: []metrics.Tag{
				{"schema", t.schema},
				{"table", t.table}},
			Values: []metrics.Value{
				{"size", t.size},
				{"rows", t.rows},
				{"increment", t.increment}},
		})
	}
}