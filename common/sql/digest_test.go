package sql_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/common/sql"
)

var queries = []struct{ ID, Input, Expected string }{
	{"comment_case_1",
		"-- select 1;",
		""},
	{"comment_case_2",
		"-- comment \nset @foo = 1;",
		"set @foo = ?;"},
	{"comment_case_3",
		"# comment \nset @foo = 1;",
		"set @foo = ?;"},
	{"comment_case_4",
		"/* comment */ set @foo = 1;",
		" set @foo = ?;"},
	{"comment_case_5",
		"/* comment */\nset @foo = 1;",
		" set @foo = ?;"},
	{"comment_case_6",
		"/*\n* comment\n*/\nset @foo = ?;",
		" set @foo = ?;"},
	{"comment_case_7",
		"set @foo = ?; -- test",
		"set @foo = ?;"},
	{"string_1",
		"SELECT id FROM foo WHERE email = 'aaa@aaa.aaa';",
		"select id from foo where email = '?';"},
	{"string_2",
		"SELECT id FROM foo WHERE email = \"aaa@aaa.aaa\";",
		"select id from foo where email = '?';"},
	{"string_3",
		`SELECT "<foo='test'/>don't bar</foo>";`,
		"select '?';"},
	{"string_4",
		`SELECT '<foo=\'test\'/>don\'t bar</foo>';`,
		"select '?';"},
	{"string_5",
		`SELECT '<foo="test"/>don\'t bar</foo>';`,
		"select '?';"},
	{"string_6",
		`SELECT "<foo=\"test\"/>don\'t bar</foo>";`,
		"select '?';"},
	{"string_7",
		`SELECT "2015-06-19";`,
		"select '?';"},
	{"string_8",
		`SELECT "2015-06-19 00:00:00";`,
		"select '?';"},
	{"string_9",
		`SELECT "foo"`,
		"select '?'"},
	{"string_number_1",
		`SELECT 1 FROM 1foo1.bar1 WHERE id = 12;`,
		"select ? from 1foo1.bar1 where id = ?;"},
	{"string_number_2",
		"select  1.1^1 from `1foo1`.`bar1`;",
		"select ?^? from 1foo1.bar1;"},
	{"string_number_alias_1",
		`select if(foo = "3", 1, 2) as "test";`,
		"select if(foo = '?', ?, ?) as '?';"},
	{"number_1",
		`select 1234;`,
		"select ?;"},
	{"number_2",
		`select .1;`,
		"select ?;"},
	{"number_3",
		`select 0.1;`,
		"select ?;"},
	{"number_4",
		`select -1;`,
		"select -?;"},
	{"number_5",
		`select -0.1;`,
		"select -?;"},
	{"number_6",
		`SELECT -.1;`,
		"select -?;"},
	{"number_7",
		`select - 1;`,
		"select - ?;"},
	{"number_8",
		`select (id + 1);`,
		"select (id + ?);"},
	{"function_1",
		"select abs(1);",
		"select abs(?);"},
	{"function_2",
		"select if(1, 1, 0);",
		"select if(?, ?, ?);"},
	{"list_1",
		`select if(4 in (1, 2, 3), true, false) from foo where id in (1, 2, 3);`,
		"select if(? in (?), true, false) from foo where id in (?);"},
	{"list_2",
		`SELECT IF(code = 'IN', 1, 0) FROM foo WHERE id in  (1, 2, 3);`,
		"select if(code = '?', ?, ?) from foo where id in (?);"},
	{"list_3",
		"SELECT IF(code = 'IN', 1, 0) FROM foo WHERE id in (\n1,\n2,\n3\n);",
		"select if(code = '?', ?, ?) from foo where id in (?);"},
	{"list_4",
		"select * from foo where id in ('A', 'B', 'C');",
		"select * from foo where id in (?);"},
	{"strange_case_1",
		"-",
		"-"},
	{"strange_case_2",
		"",
		""},
	{"subquery_case_1",
		"SELECT count(*) AS total FROM foo JOIN (SELECT * FROM bar WHERE fk = 1);",
		"select count(*) as total from foo join (select * from bar where fk = ?);"},
	{"subquery_case_2",
		"SELECT count(*) AS total FROM foo JOIN (SELECT * FROM bar WHERE fk = 1);",
		"select count(*) as total from foo join (select * from bar where fk = ?);"},
	{"subquery_case_3",
		"SELECT count(*) AS total FROM foo JOIN (SELECT DISTINCT * FROM bar WHERE fk = 1);",
		"select count(*) as total from foo join (select distinct * from bar where fk = ?);"},
	{"subquery_case_4",
		"SELECT * FROM foo INNER JOIN (SELECT * FROM bar WHERE fk = 1);",
		"select * from foo inner join (select * from bar where fk = ?);"},
}

func TestDigest(t *testing.T) {
	for _, test := range queries {
		actual := sql.Digest(test.Input)

		if test.Expected != actual {
			t.Errorf("Test %d - Expected: '%#v', got: '%#v'.", test.ID, test.Expected, actual)
		}
	}
}
