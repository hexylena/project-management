package sqlish

import (
	"fmt"
	"testing"
)

func fieldTest(c int, t *testing.T, field string, expected string) {
	if field != expected {
		t.Errorf("[c%d] Mismatch\no:«%s»\n!=\ne:«%s»", c, field, expected)
	}
}

func fieldTesti(c int, t *testing.T, field int, expected int) {
	if field != expected {
		t.Errorf("[c%d] Mismatch\no:«%d»\n\t!=\n\te:«%d»", c, field, expected)
	}
}

func TestSqlParsingSimple(t *testing.T) {
	ans := ParseSqlQuery("SELECT id FROM asdf")
	fieldTest(1, t, ans.Select, "id")
	fieldTest(2, t, ans.From, "asdf")

	ans = ParseSqlQuery("SELECT id, title FROM asdf")
	fieldTest(3, t, ans.Select, "id, title")
	fieldTest(4, t, ans.From, "asdf")

	ans = ParseSqlQuery("SELECT a,b,c,d,e FROM qq")
	fieldTest(5, t, ans.Select, "a, b, c, d, e")
	fieldTest(6, t, ans.From, "qq")

	ans = ParseSqlQuery("SELECT id FROM asdf")
	fieldTest(7, t, ans.Select, "id")
	fieldTest(8, t, ans.From, "asdf")

	ans = ParseSqlQuery("SELECT id, title FROM asdf")
	fieldTest(9, t, ans.Select, "id, title")
	fieldTest(10, t, ans.From, "asdf")

	ans = ParseSqlQuery("SELECT a,b,c,d,e FROM qq")
	fieldTest(11, t, ans.Select, "a, b, c, d, e")
	fieldTest(12, t, ans.From, "qq")

	ans = ParseSqlQuery("SELECT id FROM asdf WHERE id=1")
	fieldTest(13, t, ans.Select, "id")
	fieldTest(14, t, ans.From, "asdf")
	fieldTest(15, t, ans.Where, "id = 1")

	ans = ParseSqlQuery("SELECT id, title FROM asdf WHERE id=1")
	fieldTest(16, t, ans.Select, "id, title")
	fieldTest(17, t, ans.From, "asdf")
	fieldTest(18, t, ans.Where, "id = 1")

	ans = ParseSqlQuery("SELECT a,b,c,d,e FROM qq WHERE id=1")
	fieldTest(19, t, ans.Select, "a, b, c, d, e")
	fieldTest(20, t, ans.From, "qq")
	fieldTest(21, t, ans.Where, "id = 1")

	ans = ParseSqlQuery("SELECT q FROM x WHERE id<=1")
	fieldTest(22, t, ans.Where, "id <= 1")

	ans = ParseSqlQuery("SELECT q FROM x WHERE id>4200")
	fieldTest(23, t, ans.Where, "id > 4200")

	ans = ParseSqlQuery("SELECT q FROM x WHERE id!=2")
	fieldTest(24, t, ans.Where, "id != 2")

	ans = ParseSqlQuery("SELECT q FROM x GROUP BY id")
	fieldTest(25, t, ans.GroupBy, "id")

	ans = ParseSqlQuery("SELECT q FROM x WHERE tag=gaming    GROUP BY id")
	fmt.Println(ans)
	fieldTest(26, t, ans.GroupBy, "id")
	// ans = ParseSqlQuery("SELECT id FROM asdf WHERE id=1 GROUP BY x ORDER BY y ASC LIMIT 10")

	ans = ParseSqlQuery("SELECT id FROM asdf WHERE id=1 GROUP BY x ORDER BY y ASC LIMIT 10")
	fieldTest(27, t, ans.Select, "id")
	fieldTest(28, t, ans.Where, "id = 1")
	fieldTest(29, t, ans.GroupBy, "x")
	fieldTest(30, t, ans.OrderBy[0].Field, "y")
	// fieldTest(30, t, ans.OrderBy[0].ASC, "ASC")
	// TODO: update tests.
	fieldTesti(31, t, ans.Limit, 10)

	ans = ParseSqlQuery("SELECT id FROM asdf WHERE id=1 GROUP BY x ORDER BY y DESC, z ASC LIMIT 10")
	fieldTest(32, t, ans.Select, "id")
	fieldTest(33, t, ans.Where, "id = 1")
	fieldTest(34, t, ans.GroupBy, "x")
	fieldTest(35, t, ans.OrderBy[0].Field, "y")
	fieldTest(35, t, ans.OrderBy[1].Field, "z")
	fieldTesti(36, t, ans.Limit, 10)

	ans = ParseSqlQuery("SELECT id, title FROM tasks WHERE project = \"4fca94d6-cdd9-4540-8b0e-6370eba448b7\" GROUP BY status")
	fieldTest(37, t, ans.Select, "id, title")
	fieldTest(38, t, ans.From, "tasks")
	fieldTest(39, t, ans.Where, "project = 4fca94d6-cdd9-4540-8b0e-6370eba448b7")
	fieldTest(40, t, ans.GroupBy, "status")
	// fieldTest(41, t, ans.OrderBy, nil)
	fieldTesti(42, t, ans.Limit, -1)

	ans = ParseSqlQuery("SELECT id, title FROM tasks WHERE project = '4fca94d6-cdd9-4540-8b0e-6370eba448b7' GROUP BY status")
	fieldTest(43, t, ans.Where, "project = 4fca94d6-cdd9-4540-8b0e-6370eba448b7")

	ans = ParseSqlQuery("SELECT id FROM asdf ORDER BY y DESC")
	fieldTest(44, t, ans.OrderBy[0].Field, "y")
	fieldTesti(45, t, ans.Limit, -1)

	ans = ParseSqlQuery("SELECT id FROM asdf ORDER BY y ASC")
	fieldTest(46, t, ans.OrderBy[0].Field, "y")
	fieldTesti(47, t, ans.Limit, -1)

	ans = ParseSqlQuery("SELECT short_id, title, Status, created FROM tasks WHERE project = '4fca94d6-cdd9-4540-8b0e-6370eba448b7' ORDER BY created ASC")
	fieldTest(48, t, ans.Select, "short_id, title, Status, created")

	ans = ParseSqlQuery("SELECT short_id FROM tasks WHERE project = '4fca94d6-cdd9-4540-8b0e-6370eba448b7' AND type = 'log'")
	fieldTest(49, t, ans.Select, "short_id")
	fieldTest(50, t, ans.Where, "project = 4fca94d6-cdd9-4540-8b0e-6370eba448b7 AND type = log")

	ans = ParseSqlQuery("SELECT short_id FROM tasks WHERE project = '4fca94d6-cdd9-4540-8b0e-6370eba448b7' AND type = 'log' ORDER BY short_id ASC")
	fieldTest(51, t, ans.Select, "short_id")
	fieldTest(52, t, ans.Where, "project = 4fca94d6-cdd9-4540-8b0e-6370eba448b7 AND type = log")
	fieldTest(53, t, ans.OrderBy[0].Field, "short_id")

	ans = ParseSqlQuery("select id, short_id, title, type, icon from asdf group by 'type' limit 3")
	fieldTest(54, t, ans.Select, "id, short_id, title, type, icon")
	fieldTest(55, t, ans.From, "asdf")
	fieldTest(56, t, ans.GroupBy, "type")
	fieldTesti(57, t, ans.Limit, 3)

	ans = ParseSqlQuery("select id FROM asdf WHERE project is null")
	fieldTest(58, t, ans.Where, "project is null")

	ans = ParseSqlQuery("select id FROM asdf WHERE project is NOT null")
	fieldTest(58, t, ans.Where, "project is not null")

}

// func TestSqlParsingSimple(t *testing.T) {
//     ans := ParseSqlQuery("SELECT id FROM asdf")
//     fieldTest(t, ans.Select, "id")
//     fieldTest(t, ans.From, "asdf")
//
//
// }
