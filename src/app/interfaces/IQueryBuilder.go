package interfaces

//IQueryBuilder the interface is needed for an additional
// level of abstraction over the method of creating requests to the ORM
type IQueryBuilder interface {
	//Where
	//	Add new AND condition to the query
	//Params
	// fieldName - name of the field
	// comparator - logical comparison operator
	// value - value of the field
	//Return
	// updated query
	Where(fieldName, comparator string, value interface{}) IQueryBuilder
	//WhereQuery
	//	Add new complicated AND condition to the query based on another query
	//Params
	// builder - query builder
	//Return
	// updated query
	WhereQuery(builder IQueryBuilder) IQueryBuilder
	//Or
	//	Add new OR condition to the query
	//Params
	// fieldName - name of the field
	// comparator - logical comparison operator
	// value - value of the field
	//Return
	// updated query
	Or(fieldName, comparator string, value interface{}) IQueryBuilder
	//OrQuery
	//	Add new complicated OR condition to the query based on another query
	//Params
	// builder - query builder
	//Return
	// updated query
	OrQuery(builder IQueryBuilder) IQueryBuilder
	//Build
	//	Build a slice of query arguments
	//Return
	// slice of query arguments
	// error - if an error occurred, otherwise nil
	Build() (interface{}, error)
}
