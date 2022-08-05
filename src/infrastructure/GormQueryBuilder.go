package infrastructure

import (
	"errors"
	"fmt"
	"regexp"
	"rol/app/interfaces"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

//ToSnakeCase converts camelCase field name to snake_case
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

//GormQueryBuilder query builder struct for gorm
type GormQueryBuilder struct {
	QueryString string
	Values      []interface{}
}

//NewGormQueryBuilder Gets new query builder instance
func NewGormQueryBuilder() *GormQueryBuilder {
	return &GormQueryBuilder{}
}

func (g *GormQueryBuilder) addQuery(condition, fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	if len(g.QueryString) > 0 {
		g.QueryString += fmt.Sprintf(" %s ", condition)
	}
	finComparator := comparator
	if comparator == "==" {
		finComparator = "="
	}
	if value == nil {
		if comparator == "=" {
			g.QueryString += fmt.Sprintf("%s IS NULL", ToSnakeCase(fieldName))
		} else if comparator == "!=" {
			g.QueryString += fmt.Sprintf("%s IS NOT NULL", ToSnakeCase(fieldName))
		}
		return g
	}
	g.QueryString += fmt.Sprintf("%s %s ?", ToSnakeCase(fieldName), finComparator)
	g.Values = append(g.Values, value)
	return g
}

func (g *GormQueryBuilder) addQueryBuilder(condition string, builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	if len(g.QueryString) > 0 {
		g.QueryString += fmt.Sprintf(" %s ", condition)
	}
	argsInterface, err := builder.Build()
	if err != nil {
		return g
	}
	argsArrInterface := argsInterface.([]interface{})
	switch v := argsArrInterface[0].(type) {
	case string:
		if len(argsArrInterface[0].(string)) < 1 {
			return g
		}
		g.QueryString += fmt.Sprintf("(%s)", strings.ReplaceAll(v, "WHERE ", ""))
	default:
		panic("[GormQueryBuilder] can't add passed query builder to current builder, check what you pass GormQueryBuilder")
	}
	for i := 1; i < len(argsArrInterface); i++ {
		g.Values = append(g.Values, argsArrInterface[i])
	}
	return g
}

//Where add new AND condition to the query
//Params
//	fieldName - name of the field
//	comparator - comparison logical operator
//	value - value of the field
//Return
//	interfaces.IQueryBuilder - updated query builder
func (g *GormQueryBuilder) Where(fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	return g.addQuery("AND", fieldName, comparator, value)
}

//WhereQuery add new complicated AND condition to the query based on another query
//Params
//	interfaces.IQueryBuilder - a ready-builder to the query
//Return
//	interfaces.IQueryBuilder - updated query builder
func (g *GormQueryBuilder) WhereQuery(builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	return g.addQueryBuilder("AND", builder)
}

//Or add new OR condition to the query
//Params
//	fieldName - name of the field
//	comparator - comparison logical operator
//	value - value of the field
//Return
//	interfaces.IQueryBuilder - updated query builder
func (g *GormQueryBuilder) Or(fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	return g.addQuery("OR", fieldName, comparator, value)
}

//OrQuery add new complicated OR condition to the query based on another query
//Params
//	interfaces.IQueryBuilder - a ready-builder to the query
//Return
//	interfaces.IQueryBuilder - updated query builder
func (g *GormQueryBuilder) OrQuery(builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	return g.addQueryBuilder("OR", builder)
}

//Build build a slice of query arguments
//Return
//	interface{} - slice of interface{}
//	error - if error occurs return error, otherwise nil
func (g *GormQueryBuilder) Build() (interface{}, error) {
	if len(g.QueryString) < 1 {
		return nil, errors.New("queryBuilder is empty")
	}
	arr := make([]interface{}, 0)
	arr = append(arr, g.QueryString)
	arr = append(arr, g.Values...)
	return arr, nil
}
