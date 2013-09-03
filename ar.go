package ar
import(
    "database/sql"
    "strconv"
    "strings"
    // "fmt"
)

type where struct{
    con     string
    bracket string
    column  interface{}
    op      string
    value   interface{}
}
type join struct{
    table    string
    joinType string
    on       [][]interface{}
}
type dbExpr struct{
    value string
}
type ar struct{
    DB                sql.DB
    queryType         string
    selectColumn      []interface{}
    expr              []string
    from              []string
    join              []join
    groupBy           [][]string
    orderBy           []map[string]string
    where             []where
    limit             []int
    offset            int
    //update,insert,delete
    talbe             string
    column            []string
    value             []string
    set               map[string]interface{}

    Sql               string
    quoteReservedChar string
    quoteChar         string
    quoteQuoteChar    string
}

func (a *ar) init() {
    a.quoteReservedChar = "`"
    a.quoteChar = "'"
    //if sqlite need set quoteQuotChar with '
    a.quoteQuoteChar = "\\"
}
func New() *ar {
    ar := &ar{}
    ar.init()
    return ar
}
func (a *ar) setQueryType(queryType string){
    a.queryType = queryType
}
func (a *ar) Select(selectColumns ...interface{}) *ar {
    a.selectColumn = append(a.selectColumn, selectColumns...)
    a.setQueryType("SELECT")
    return a
}
func (a *ar) From(froms ...string) *ar{
    a.from = append(a.from, froms...)
    return a
}
func (a *ar) Where(column interface{}, op string, value interface{}) *ar{
    var tmpWhere where
    if len(a.where) != 0 && a.where[len(a.where)-1].bracket != "(" { 
        tmpWhere.con    = "AND"
    }
    tmpWhere.column = column
    tmpWhere.op     = op
    tmpWhere.value  = value
    a.where         = append(a.where, tmpWhere)
    return a
}

func (a *ar) OrWhere(column interface{}, op string, value interface{}) *ar{
    var tmpWhere where
    tmpWhere.con    = "OR"
    tmpWhere.column = column
    tmpWhere.op     = op
    tmpWhere.value  = value
    a.where         = append(a.where, tmpWhere)
    return a
}

func (a *ar) Limit(limit ...int) *ar{
    a.limit = limit
    return a
}
func (a *ar) Offset(offset int) *ar{
    a.offset = offset
    return a
}
func (a *ar) OrderBy(name string, sort string) *ar{
    tmp := map[string]string{"name":name,"sort":sort}
    a.orderBy = append(a.orderBy, tmp)
    return a
}

func (a *ar) Insert(table string, column ...[]string) *ar{
    a.talbe = table
    a.column = append(a.column, column[0]...)
    a.setQueryType("INSERT")
    return a
}
func (a *ar) Values(value []string) *ar{
    a.value = append(a.value, value...)
    return a
}

func (a *ar) Update(talbe string) *ar{
    a.talbe = talbe
    a.setQueryType("UPDATE")
    return a
}
func (a *ar) Set(set map[string]interface{}) *ar{
    a.set = set
    return a
}

func (a *ar) Delete(talbe string) *ar{
    a.talbe = talbe
    a.setQueryType("DELETE")
    return a
}

func (a *ar) Join(table string, joinType ...string) *ar{
    var tmp  join 
    tmp.table = table
    if joinType != nil {
        tmp.joinType = joinType[0]
    }
    a.join = append(a.join, tmp)
    return a
}
func (a *ar) On(on ...interface{}) *ar{
    var tmp join
    tmp = a.join[len(a.join)-1]
    if len(tmp.on) >= 1 {
        on = append(on , "AND")
    }
    on = append(on , " ")
    tmp.on = append(a.join[len(a.join)-1].on, on)
    a.join[len(a.join)-1] = tmp
    return a
}

func (a *ar) WhereOpen() *ar {
    var tmpWhere where
    tmpWhere.op= "AND"
    tmpWhere.bracket = "("
    a.where= append(a.where, tmpWhere)
    return a
}
func (a *ar) OrWhereOpen() *ar {
    var tmpWhere where
    tmpWhere.op= "OR"
    tmpWhere.bracket = "("
    a.where= append(a.where, tmpWhere)
    return a
}

func (a *ar) AndWhereOpen() *ar {
    var tmpWhere where
    tmpWhere.op= "AND"
    tmpWhere.bracket = "("
    a.where= append(a.where, tmpWhere)
    return a
}


func (a *ar) WhereClose() *ar {
    var tmpWhere where
    tmpWhere.bracket = ")"
    a.where= append(a.where, tmpWhere)
    return a
}

func (a *ar) Quote() {
}
func (a *ar) Build() *ar{
    switch a.queryType {
        case "SELECT": a.buildSelect().buildFrom().buildJoin().buildWhere().buildLimit().buildOrderBy()
        case "INSERT": a.buildInsert().buildValues()
        case "UPDATE": a.buildUpdate().buildSet().buildWhere()
        case "DELETE": a.buildDelete().buildWhere()
    }
    return a
}

func (a *ar) buildSelect() *ar{
    //selectColumn 
    tmp := []string{}
    if len(a.selectColumn) != 0 {
        a.Sql += Concat("SELECT")
        for _,v := range a.selectColumn {
            tmp = append(tmp, Concat("", a.buildExpr(v)))
        }
        a.Sql += strings.Join(tmp, ", ")
    }
    return a
}
func (a *ar) buildFrom() *ar{
    //from
    if len(a.from) != 0 {
        a.Sql += Concat("FROM")
        for k,v := range a.from{
            a.from[k] = Concat("", a.quote(v, a.quoteReservedChar))
        }
        a.Sql += strings.Join(a.from, ", ")
    }
    return a
}

func (a *ar) buildJoin() *ar{
    // join
    if len(a.join) != 0  {
        for jk,jv := range a.join {
            a.Sql += Concat(jv.joinType, "JOIN", a.quote(jv.table,a.quoteReservedChar))
            if len(a.join[jk].on) != 0 {
                a.Sql += Concat("ON")
                for _, ov := range a.join[jk].on {
                    a.Sql += Concat(ov[3].(string), a.buildExpr(ov[0]), ov[1].(string), a.buildExpr(ov[2]))
                }
            }
        }
    }
    return a
}
func (a *ar) buildWhere() *ar{
    //where
    if len(a.where) != 0 {
        a.Sql += Concat("WHERE")
        for _,v := range a.where {
            a.Sql += Concat("", v.con, a.buildExpr(v.column) , v.op , a.buildExpr(v.value)  , v.bracket)
        }
    }
    return a
}
func (a *ar) buildOrderBy() *ar{
    //order by
    if len(a.orderBy) != 0 {
        a.Sql += Concat("ORDER BY")
        for _,ov := range a.orderBy {
            a.Sql += Concat(ov["name"], ov["sort"])
        }
    }
    return a
}
func (a *ar) buildLimit() *ar{
    //limit
    if len(a.limit) == 1 {
        a.Sql += Concat("LIMIT", strconv.Itoa(a.limit[0]))
        //offset
        a.Sql += Concat("OFFSET", strconv.Itoa(a.offset))
    }else if len(a.limit) != 0{
        a.Sql += Concat("LIMIT", strconv.Itoa(a.limit[0]), ",", strconv.Itoa(a.limit[1]))
    }
    return a
}

func (a *ar) buildInsert() *ar{
    a.Sql += Concat("INSERT INTO", a.quote(a.talbe, a.quoteReservedChar))
    for k,cv := range a.column {
        a.column[k] = a.quote(cv,a.quoteReservedChar)
    }
    a.Sql += Concat("(", strings.Join(a.column, ", "), ")")
    return a
}

func (a *ar) buildValues() *ar{
    a.Sql += " VALUES "
    for k,vv := range a.value{
        a.value[k] = a.quote(vv)
    }
    a.Sql += Concat("(", strings.Join(a.value, ", "), ")")
    return a
}

func (a *ar) buildUpdate() *ar{
    a.Sql += Concat("UPDATE", a.quote(a.talbe,a.quoteReservedChar))
    return a
}
func (a *ar) buildSet() *ar{
    tmp := []string{}
    a.Sql += "SET "
    for k,v := range a.set {
        tmp = append(tmp, Concat(a.quote(k), "=", a.buildExpr(v)))
    }
    a.Sql += strings.Join(tmp, ", ")
    return a
}
func (a *ar) buildDelete() *ar{
    a.Sql += Concat("DELETE", a.quote(a.talbe , a.quoteReservedChar))
    return a
}

func (a *ar)quote(s string,char ...string) string {
    quoteChar := a.quoteChar
    if s == "" {
        return ""
    }
    if char != nil {
        quoteChar = char[0]
    }
    return quoteChar+strings.Replace(strings.Replace(s, a.quoteQuoteChar+quoteChar, quoteChar, -1),quoteChar, a.quoteQuoteChar+quoteChar, -1)+quoteChar
}
func Concat(firstWord string, words ...string) string {
    tmpString := ""
    if firstWord != "" {
        tmpString = firstWord + " "
    }
    for _,w := range words {
        if w != "" {
            tmpString += w + " "
        }
    }
    return tmpString
}
func (a *ar) Exec() *ar{ 
    return a
}
func (a *ar) buildExpr(i interface{}) string {
    var s string
    switch i := i.(type) {
        case dbExpr: s = i.value
        case string: s = a.quote(string(i))
        case int: s = strconv.Itoa(i)
    }
    return s
}
func Expr(s string) dbExpr{
    var expr dbExpr
    expr.value = s
    return expr
}


