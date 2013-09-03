package ar
import(
    // "h2eroUtil/ar"
    "testing"
    "fmt"
)
func T(s string, num int) {
    fmt.Println("/--------------------------------")
    fmt.Println("")
    fmt.Println(s)
    fmt.Println("")
}
func TestMain(t *testing.T) {
    num := 0
    //test 1.
    test := New()
    test.Select("id", "name").From("t1", "t2").Join("t3", "LEFT").On("t2.id", "=", Expr("t3.uid % 2")).On("t2.id", "=", "t3.uid").Where("id", "=", "100").Build()
    T(test.Sql, num)

    //test 2.
    test = New()
    test.Select("id","name").From("table1").Build()
    T(test.Sql, num)

    //test 3.
    test = New()
    test.Select(Expr("sum(id)"),"name").From("table1").Where("id", ">", Expr("id")).Build()
    T(test.Sql, num)

    //test 4.
    test = New()
    test.Select(Expr("sum(id)"),"name").From("table1").
    Where("id", ">", Expr("id")).
    WhereOpen().Where("id", "=", Expr("id % 2")).Where("name", "=", "h2ero").Where("login_date", "=", "20120101").WhereClose().
    OrWhere("id", "<", "10").OrWhereOpen().Where("id", "=", "33"). Where("name", "=", "h2eros").WhereClose().
    Build()
    T(test.Sql, num)

    //test 5.
    test = New()
    test.Select("id").From("t1", "t2").Where("id", "=", "1").OrderBy("id", "DESC").Limit(1,10).Offset(10).Build()
    T(test.Sql, num)

    //test 6.
    test = New()
    test.Insert("t1", []string{"name", "date"}).Values([]interface{}{"h2ero", 1990}).Build();
    T(test.Sql, num) 
    //test 7.
    test = New()
    test.Update("t1").Set(map[string]interface{}{"name":"h2ero", "cash":25, "topup":Expr("topup-20")}).Where("id", "=", "3").Build()
    T(test.Sql, num)

    //test 8.
    test = New()
    test.Delete("t1").Where("id", "=", "3").Build()
    T(test.Sql, num)

    //test 9. join(expr)
    test = New()
    subSql := " (SELECT `in_t1`.`group_id`, SUM(CASE `in_t1`.`delete_flag` WHEN 0 THEN 1 ELSE 0 END) as count'FROM `hoges` AS `in_t1`  GROUP BY `in_t1`.`group_id`) as `t1`"
    test.Select("id").From("t1").Join(Expr(subSql)).On("`t0`.`id`", "=", "`t1`.group_id").Build()
    T(test.Sql, num)

}
