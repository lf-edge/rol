package infrastructure

import (
	"fmt"
	"rol/app/interfaces"
	"strings"
)

//YamlQueryBuilder query builder struct for yaml
type YamlQueryBuilder struct {
	QueryString string
	Values      []interface{}
}

//NewYamlQueryBuilder is a constructor for YamlQueryBuilder
func NewYamlQueryBuilder() *YamlQueryBuilder {
	return &YamlQueryBuilder{}
}

func (y *YamlQueryBuilder) addQuery(condition, fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	if len(y.QueryString) > 0 {
		y.QueryString += fmt.Sprintf(" %s ", condition)
	}
	y.QueryString += fmt.Sprintf("%s %s ?", fieldName, comparator)
	y.Values = append(y.Values, value)
	return y
}

func (y *YamlQueryBuilder) addQueryBuilder(condition string, builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	if len(y.QueryString) > 0 {
		y.QueryString += fmt.Sprintf(" %s ", condition)
	}
	argsInterface, err := builder.Build()
	if err != nil {
		return y
	}
	argsArrInterface := argsInterface.([]interface{})
	switch v := argsArrInterface[0].(type) {
	case string:
		if len(argsArrInterface[0].(string)) < 1 {
			return y
		}
		y.QueryString += fmt.Sprintf("(%s)", strings.ReplaceAll(v, "WHERE ", ""))
	default:
		panic("[YamlQueryBuilder] can't add passed query builder to current builder, check what you pass YamlQueryBuilder")
	}
	for i := 1; i < len(argsArrInterface); i++ {
		y.Values = append(y.Values, argsArrInterface[i])
	}
	return y
}

//Where add new AND condition to the query
//Params
// fieldName - name of the field
// comparator - logical comparison operator
// value - value of the field
//Return
// updated query
func (y *YamlQueryBuilder) Where(fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	return y.addQuery("AND", fieldName, comparator, value)
}

//WhereQuery add new complicated AND condition to the query based on another query
//Params
// builder - query builder
//Return
// updated query
func (y *YamlQueryBuilder) WhereQuery(builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	return y.addQueryBuilder("AND", builder)
}

//Or add new OR condition to the query
//Params
// fieldName - name of the field
// comparator - logical comparison operator
// value - value of the field
//Return
// updated query
func (y *YamlQueryBuilder) Or(fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	return y.addQuery("OR", fieldName, comparator, value)
}

//OrQuery add new complicated OR condition to the query based on another query
//Params
// builder - query builder
//Return
// updated query
func (y *YamlQueryBuilder) OrQuery(builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	return y.addQueryBuilder("OR", builder)
}

//Build a slice of query arguments
//Return
// slice of query arguments
// error - if an error occurred, otherwise nil
func (y *YamlQueryBuilder) Build() (interface{}, error) {
	arr := make([]interface{}, 0)
	if len(y.QueryString) < 1 {
		return arr, nil
	}
	arr = append(arr, y.QueryString)
	arr = append(arr, y.Values...)
	return arr, nil
}
