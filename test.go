package main
import(
    "h2eroUtil/ar"
    "fmt"
)
func main() {
    // activeRecord.New()
    // test.init()
    test := ar.New()
    // test.Select("id","name").From("table1").Build()
    // fmt.Println(test.Sql)
    // test.Select(ar.Expr("sum(id)"),"name").From("table1").Where("id", ">", ar.Expr("id")).
    // WhereOpen().Where("id", "=", ar.Expr("id % 2")).Where("name", "=", "h2ero").Where("login_date", "=", "20120101").WhereClose().Build()
    // OrWhere("id", "<", "10").
    // OrWhereOpen().Where("id", "=", "33"). Where("name", "=", "h2eros").WhereClose().
    // Build()
    test.Select("id", "name").From("t1", "t2").Join("t3", "LEFT").On("t2.id", "=", ar.Expr("t3.uid % 2")).On("t2.id", "=", "t3.uid").Where("id", "=", "100").Build()
    // // test.Select("id").From("t1", "t2").Where("id", "=", "1").OrderBy("id", "DESC").Limit(1,10).Offset(10).Build()
    // // test.Insert("t1", []string{"name", "date"}).Values([]string{"h2ero", "1990"}).Build();
    // test.Update("t1").Set(map[string]interface{}{"name":"h2ero", "cash":25, "topup":ar.Expr("topup-20")}).Where("id", "=", "3").Build()
    // // test.Delete("t1").Where("id", "=", "3").Build()
    fmt.Println(test.Sql)
}
