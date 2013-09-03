package main
import(
    "h2eroUtil/ar"
    "fmt"
)
func t(s string, num int) {
    fmt.Println("/--------------------------------")
    fmt.Println("")
    fmt.Println(s)
    fmt.Println("")
}
func main() {
    num := 0
    //test 1.
    test := ar.New()
    test.Select("id", "name").From("t1", "t2").Join("t3", "LEFT").On("t2.id", "=", ar.Expr("t3.uid % 2")).On("t2.id", "=", "t3.uid").Where("id", "=", "100").Build()
    t(test.Sql, num)

    //test 2.
    test = ar.New()
    test.Select("id","name").From("table1").Build()
    t(test.Sql, num)

    //test 3.
    test = ar.New()
    test.Select(ar.Expr("sum(id)"),"name").From("table1").Where("id", ">", ar.Expr("id")).Build()
    t(test.Sql, num)

    //test 4.
    test = ar.New()
    test.Select(ar.Expr("sum(id)"),"name").From("table1").
    Where("id", ">", ar.Expr("id")).
    WhereOpen().Where("id", "=", ar.Expr("id % 2")).Where("name", "=", "h2ero").Where("login_date", "=", "20120101").WhereClose().
    OrWhere("id", "<", "10").OrWhereOpen().Where("id", "=", "33"). Where("name", "=", "h2eros").WhereClose().
    Build()
    t(test.Sql, num)

    //test 5.
    test = ar.New()
    test.Select("id").From("t1", "t2").Where("id", "=", "1").OrderBy("id", "DESC").Limit(1,10).Offset(10).Build()
    t(test.Sql, num)

    //test 6.
    test = ar.New()
    test.Insert("t1", []string{"name", "date"}).Values([]interface{}{"h2ero", 1990}).Build();
    t(test.Sql, num)

    //test 7.
    test = ar.New()
    test.Update("t1").Set(map[string]interface{}{"name":"h2ero", "cash":25, "topup":ar.Expr("topup-20")}).Where("id", "=", "3").Build()
    t(test.Sql, num)

    //test 8.
    test = ar.New()
    test.Delete("t1").Where("id", "=", "3").Build()
    t(test.Sql, num)

    //test 9. join(expr)
    test = ar.New()
    subSql := " (SELECT `in_t1`.`group_id`, SUM(CASE `in_t1`.`delete_flag` WHEN 0 THEN 1 ELSE 0 END) as count'FROM `hoges` AS `in_t1`  GROUP BY `in_t1`.`group_id`) as `t1`"
    test.Select("id").From("t1").Join(ar.Expr(subSql)).On("`t0`.`id`", "=", "`t1`.group_id").Build()
    t(test.Sql, num)

}
