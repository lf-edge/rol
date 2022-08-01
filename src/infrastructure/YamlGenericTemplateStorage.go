package infrastructure

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"rol/app/errors"
	"rol/app/interfaces"
	"sort"
	"strconv"
	"strings"
	"time"
)

//YamlGenericTemplateStorage is a storage for yaml templates
type YamlGenericTemplateStorage[TemplateType interface{}] struct {
	//TemplatesDirectory is a directory where the templates are located
	TemplatesDirectory string
	logger             *logrus.Logger
	logSourceName      string
	Templates          *[]TemplateType
}

//QueryUnit represents bracketed expression that is part of query string
type QueryUnit struct {
	FieldName  string
	Comparator string
	ValueIndex int
}

//NewYamlGenericTemplateStorage is a constructor for YamlGenericTemplateStorage
//
//Params:
//	dirName - directory name where the templates are located
//	log - logrus.Logger
func NewYamlGenericTemplateStorage[TemplateType interface{}](dirName string, log *logrus.Logger) (interfaces.IGenericTemplateStorage[TemplateType], error) {
	model := new(TemplateType)
	executedFilePath, _ := os.Executable()
	templatesDirectory := path.Join(path.Dir(executedFilePath), "templates", dirName)
	if _, err := os.Stat(templatesDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(templatesDirectory, os.ModePerm)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to create template directory")
		}
	}
	storage := &YamlGenericTemplateStorage[TemplateType]{
		TemplatesDirectory: templatesDirectory,
		logger:             log,
		logSourceName:      fmt.Sprintf("YamlGenericTemplateStorage<%s>", reflect.TypeOf(*model).Name()),
		Templates:          &[]TemplateType{},
	}
	err := storage.reloadFromFiles()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "load from files error")
	}
	return storage, nil
}

func (y *YamlGenericTemplateStorage[TemplateType]) reloadFromFiles() error {
	y.Templates = &[]TemplateType{}

	files, err := ioutil.ReadDir(y.TemplatesDirectory)
	if err != nil {
		return errors.Internal.Wrap(err, "reading dir error")
	}
	for _, f := range files {
		template, err := y.getTemplateObjFromYaml(f.Name())
		if err != nil {
			return errors.Internal.Wrap(err, "yaml parsing error")
		}
		*y.Templates = append(*y.Templates, *template)
	}
	return nil
}

func (y *YamlGenericTemplateStorage[TemplateType]) getTemplateObjFromYaml(templateName string) (*TemplateType, error) {
	template := new(TemplateType)
	templateFilePath := path.Join(y.TemplatesDirectory, fmt.Sprintf(templateName))
	f, err := os.Open(templateFilePath)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "directory opening error")
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(template)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "yaml decoding error")
	}
	return template, nil
}

func (y *YamlGenericTemplateStorage[TemplateType]) sortTemplatesSlice(templates *[]TemplateType, orderBy, orderDirection string) error {
	if len(*templates) < 1 {
		return nil
	}
	if !isFieldExist((*templates)[0], orderBy) && orderBy != "" {
		return errors.Internal.Newf("there is no field with name '%s' at template", orderBy)
	}
	sort.Slice(*templates, func(i, j int) bool {
		firstElem := (*templates)[i]
		secondElem := (*templates)[j]
		firstReflect := reflect.ValueOf(firstElem).FieldByName(orderBy)
		secondReflect := reflect.ValueOf(secondElem).FieldByName(orderBy)
		switch firstReflect.Kind() {
		case reflect.String:
			if strings.ToLower(orderDirection) == "desc" || strings.ToLower(orderDirection) == "descending" {
				return reflect.Indirect(firstReflect).String() > reflect.Indirect(secondReflect).String()
			}
			return reflect.Indirect(firstReflect).String() < reflect.Indirect(secondReflect).String()

		case reflect.Int:
			if strings.ToLower(orderDirection) == "desc" || strings.ToLower(orderDirection) == "descending" {
				return reflect.Indirect(firstReflect).Int() > reflect.Indirect(secondReflect).Int()
			}
			return reflect.Indirect(firstReflect).Int() < reflect.Indirect(secondReflect).Int()

		case reflect.Struct:
			if firstReflect.Type().String() == "time.Time" {
				fTime := firstReflect.Interface().(time.Time)
				sTime := secondReflect.Interface().(time.Time)
				if strings.ToLower(orderDirection) == "desc" || strings.ToLower(orderDirection) == "descending" {
					return fTime.After(sTime)
				}
				return fTime.Before(sTime)
			}
			return false

		default:
			return false
		}
	})
	return nil
}

//GetByName gets template by name
//Params:
//	ctx - context is used only for logging
//	templateName - name of template
//Return:
//	*TemplateType - pointer to template
//	error - if an error occurs, otherwise nil
func (y *YamlGenericTemplateStorage[TemplateType]) GetByName(ctx context.Context, templateName string) (*TemplateType, error) {
	y.log(ctx, logrus.DebugLevel, fmt.Sprintf("GetByName: name = %s", templateName))
	queryBuilder := NewYamlQueryBuilder()
	queryBuilder.Where("Name", "==", templateName)
	query, err := queryBuilder.Build()
	if err != nil {
		return nil, errors.Internal.Wrap(err, "query building error")
	}
	queryArr := query.([]interface{})
	templates, err := y.handleQuery(*y.Templates, queryArr...)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "query handling error")
	}
	if len(*templates) == 1 {
		return &(*templates)[0], nil
	}
	return nil, nil
}

//GetList gets list of templates with filtering and pagination
//
//Params:
//	ctx - context is used only for logging
//	orderBy - order by string parameter
//	orderDirection - ascending or descending order
//	page - page number
//	size - page size
//	queryBuilder - query builder for filtering
//Return:
//	*[]TemplateType - pointer to array of templates
//	error - if an error occurs, otherwise nil
func (y *YamlGenericTemplateStorage[TemplateType]) GetList(ctx context.Context, orderBy, orderDirection string, page, pageSize int, queryBuilder interfaces.IQueryBuilder) (*[]TemplateType, error) {
	y.log(ctx, logrus.DebugLevel, fmt.Sprintf("GetList: IN: orderBy=%s, orderDirection=%s, page=%d, size=%d, queryBuilder=%s", orderBy, orderDirection, page, pageSize, queryBuilder))
	var (
		templates *[]TemplateType
		queryArr  []interface{}
		err       error
	)
	offset := (page - 1) * pageSize
	if queryBuilder != nil {
		query, err := queryBuilder.Build()
		if err != nil {
			return nil, errors.Internal.Wrap(err, "query building error")
		}
		queryArr = query.([]interface{})
	}
	if len(queryArr) > 1 {
		templates, err = y.handleQuery(*y.Templates, queryArr...)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "query handling error")
		}
	} else {
		templates = y.Templates
	}

	err = y.sortTemplatesSlice(templates, orderBy, orderDirection)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "templates sorting error")
	}
	paginatedSlice, err := y.getPaginatedSlice(*templates, offset, pageSize)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "templates pagination error")
	}

	return &paginatedSlice, nil
}

func (y *YamlGenericTemplateStorage[TemplateType]) getPaginatedSlice(templates []TemplateType, offset, limit int) ([]TemplateType, error) {
	limit += offset
	if offset > len(templates) {
		return nil, errors.Internal.Newf("paginated slice offset bounds out of range [%d:] with length %d", offset, len(templates))
	}
	if limit > len(templates) {
		return templates[offset:], nil
	}
	return templates[offset:limit], nil
}

//Count gets total count of templates with current query
//Params
//	ctx - context is used only for logging
//	queryBuilder - query for entities to count
//Return
//	int64 - number of entities
//	error - if an error occurs, otherwise nil
func (y *YamlGenericTemplateStorage[TemplateType]) Count(ctx context.Context, queryBuilder interfaces.IQueryBuilder) (int64, error) {
	y.log(ctx, logrus.DebugLevel, fmt.Sprintf("Count: IN: queryBuilder=%+v", queryBuilder))
	var templatesSlice []TemplateType
	files, err := ioutil.ReadDir(y.TemplatesDirectory)
	if err != nil {
		return 0, errors.Internal.Wrap(err, "get templates files error")
	}
	for _, f := range files {
		template, err := y.getTemplateObjFromYaml(f.Name())
		if err != nil {
			return 0, errors.Internal.Wrap(err, "error converting yaml to struct")
		}
		templatesSlice = append(templatesSlice, *template)
	}
	queryStr, err := queryBuilder.Build()
	if err != nil {
		return 0, errors.Internal.Wrap(err, "query building error")
	}
	queryArr := queryStr.([]interface{})
	foundTemplates, err := y.handleQuery(templatesSlice, queryArr...)
	if err != nil {
		return 0, errors.Internal.Wrap(err, "query handling error")
	}
	count := int64(len(*foundTemplates))
	y.log(ctx, logrus.DebugLevel, fmt.Sprintf("Count: OUT: count=%d", count))
	return count, nil
}

func (y *YamlGenericTemplateStorage[TemplateType]) log(ctx context.Context, level logrus.Level, message string) {
	if ctx != nil {
		actionID := uuid.UUID{}
		if ctx.Value("requestId") != nil {
			actionID = ctx.Value("requestId").(uuid.UUID)
		}

		entry := y.logger.WithFields(logrus.Fields{
			"actionID": actionID,
			"source":   y.logSourceName,
		})
		switch level {
		case logrus.ErrorLevel:
			entry.Error(message)
		case logrus.InfoLevel:
			entry.Info(message)
		case logrus.WarnLevel:
			entry.Warn(message)
		case logrus.DebugLevel:
			entry.Debug(message)
		}
	}
}

//NewQueryBuilder gets new query builder
//Params
//	ctx - context is used only for logging
//Return
//	interfaces.IQueryBuilder - new query builder
func (y *YamlGenericTemplateStorage[TemplateType]) NewQueryBuilder(ctx context.Context) interfaces.IQueryBuilder {
	y.log(ctx, logrus.DebugLevel, "Call method NewQueryBuilder")
	return NewYamlQueryBuilder()
}

func (y *YamlGenericTemplateStorage[TemplateType]) handleQuery(templatesSlice []TemplateType, args ...interface{}) (*[]TemplateType, error) {
	if len(args) < 1 {
		return &templatesSlice, nil
	}
	query := replaceQuestionsToIndexes(args[0].(string))
	queryValues := args[1:]
	finalSlice := &[]TemplateType{}
	for _, template := range templatesSlice {
		queryForTemplate := query
		startIndex, endIndex := findLowerQueryIndexes(queryForTemplate)
		for {
			if startIndex == -1 && endIndex == -1 {
				break
			}
			result, err := handleSimpleQuery(template, queryForTemplate[startIndex+1:endIndex-1], queryValues)
			if err != nil {
				return nil, errors.Internal.Wrap(err, "simple query handling error")
			}
			if result {
				queryForTemplate = replaceWithFakeTrueQuery(queryForTemplate, startIndex, endIndex)
			} else {
				queryForTemplate = replaceWithFakeFalseQuery(queryForTemplate, startIndex, endIndex)
			}
			startIndex, endIndex = findLowerQueryIndexes(queryForTemplate)
		}
		result, err := handleSimpleQuery(template, queryForTemplate, queryValues)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "simple query handling error")
		}
		if result {
			*finalSlice = append(*finalSlice, template)
		}
	}
	return finalSlice, nil
}

func handleSimpleQuery(template interface{}, query string, queryValues []interface{}) (bool, error) {
	condition := ""
	queryUnitString, lastParsedIndex := getQueryUnitString(query, 0)
	result := false
	for {
		if len(queryUnitString) < 3 {
			break
		}
		queryUnit, err := parseQueryUnitString(strings.Trim(queryUnitString, " "))
		if err != nil {
			return false, errors.Internal.Wrap(err, "error parsing query unit")
		}
		if !isFieldExist(template, queryUnit.FieldName) && queryUnit.FieldName != "FakeFalse" && queryUnit.FieldName != "FakeTrue" {
			return false, errors.Internal.Newf("there is no field with name '%s' at template", queryUnit.FieldName)
		}

		value := queryValues[queryUnit.ValueIndex]

		interimResult, err := getResultOfQueryUnit(template, queryUnit, value)
		if err != nil {
			return false, errors.Internal.Wrap(err, "error getting result of query unit")
		}
		if condition == "" {
			result = interimResult
		} else if condition == "AND" {
			result = result && interimResult
		} else if condition == "OR" {
			result = result || interimResult
		}
		// Get condition if exist for the next iteration
		condition, lastParsedIndex = getConditionString(query, lastParsedIndex)
		queryUnitString, lastParsedIndex = getQueryUnitString(query, lastParsedIndex)
	}
	return result, nil
}

func replaceQuestionsToIndexes(query string) string {
	count := strings.Count(query, "?")
	for i := 0; i < count; i++ {
		query = strings.Replace(query, "?", strconv.Itoa(i), 1)
	}
	return query
}

func findLowerQueryIndexes(query string) (int, int) {
	endIndexOfQueryGroup := strings.Index(query, ")")
	if endIndexOfQueryGroup < 1 {
		return -1, -1
	}
	endIndexOfQueryGroup = endIndexOfQueryGroup + 1
	startIndexOfQueryGroup := strings.LastIndex(query[0:endIndexOfQueryGroup], "(")
	return startIndexOfQueryGroup, endIndexOfQueryGroup
}

func findConditionIndexAndLen(query string, searchStartIndex int) (int, int) {
	searchAbleQuery := query[searchStartIndex:]
	andIndex := strings.Index(searchAbleQuery, " AND ")
	orIndex := strings.Index(searchAbleQuery, " OR ")
	if orIndex != -1 && andIndex != -1 {
		if andIndex < orIndex {
			return searchStartIndex + andIndex + 1, 3
		}
		return searchStartIndex + orIndex + 1, 2
	}
	if orIndex != -1 {
		return searchStartIndex + orIndex + 1, 2
	}
	if andIndex != -1 {
		return searchStartIndex + andIndex + 1, 3
	}
	return -1, -1
}

func getFieldValue(template interface{}, fieldName string) (interface{}, error) {
	valueOfTemplate := reflect.ValueOf(template)
	if valueOfTemplate.Kind() == reflect.Ptr {
		valueOfTemplate = valueOfTemplate.Elem()
	}
	fieldReflect := valueOfTemplate.FieldByName(fieldName)
	var fieldValue interface{}
	switch fieldReflect.Kind() {
	case reflect.String:
		fieldValue = valueOfTemplate.FieldByName(fieldName).String()
	case reflect.Int:
		fieldValue = int(valueOfTemplate.FieldByName(fieldName).Int())
	case reflect.Struct:
		if fieldReflect.Type().String() == "time.Time" {
			fieldValue = valueOfTemplate.FieldByName(fieldName).Interface().(time.Time)
			break
		}
		return nil, errors.Internal.New("wrong field type")
	default:
		return nil, errors.Internal.New("wrong field type")
	}
	return fieldValue, nil
}

func isFieldExist(template interface{}, fieldName string) bool {
	return reflect.ValueOf(template).FieldByName(fieldName).IsValid()
}

func parseQueryUnitString(queryUnit string) (QueryUnit, error) {
	queryUnitSlice := strings.Split(queryUnit, " ")
	fieldName, comparator := queryUnitSlice[0], queryUnitSlice[1]
	valueIndex, err := strconv.Atoi(queryUnitSlice[2])
	if err != nil {
		return QueryUnit{}, errors.Internal.Wrap(err, "error converted to type int")
	}
	return QueryUnit{
		FieldName:  fieldName,
		Comparator: comparator,
		ValueIndex: valueIndex,
	}, nil
}

func isBigger(first, second any) (bool, error) {
	switch first.(type) {
	case string:
		return first.(string) > second.(string), nil
	case int:
		return first.(int) > second.(int), nil
	case time.Time:
		fTime := first.(time.Time)
		sTime := second.(time.Time)
		return fTime.After(sTime), nil
	default:
		return false, errors.Internal.New("wrong type")
	}
}

func isBiggerOrEqual(first, second any) (bool, error) {
	switch first.(type) {
	case string:
		return first.(string) >= second.(string), nil
	case int:
		return first.(int) >= second.(int), nil
	case time.Time:
		fTime := first.(time.Time)
		sTime := second.(time.Time)
		return fTime.After(sTime) || fTime.Equal(sTime), nil
	default:
		return false, errors.Internal.New("wrong type")
	}
}

func isLesser(first, second any) (bool, error) {
	switch first.(type) {
	case string:
		return first.(string) < second.(string), nil
	case int:
		return first.(int) < second.(int), nil
	case time.Time:
		fTime := first.(time.Time)
		sTime := second.(time.Time)
		return fTime.Before(sTime), nil
	default:
		return false, errors.Internal.New("wrong type")
	}
}

func isLesserOrEqual(first, second any) (bool, error) {
	switch first.(type) {
	case string:
		return first.(string) <= second.(string), nil
	case int:
		return first.(int) <= second.(int), nil
	case time.Time:
		fTime := first.(time.Time)
		sTime := second.(time.Time)
		return fTime.Before(sTime) || fTime.Equal(sTime), nil
	default:
		return false, errors.Internal.New("wrong type")
	}
}

func getResultOfQueryUnit(template interface{}, queryUnit QueryUnit, value interface{}) (bool, error) {
	// This is a hack
	if queryUnit.FieldName == "FakeTrue" {
		return true, nil
	}
	if queryUnit.FieldName == "FakeFalse" {
		return false, nil
	}

	fieldValue, err := getFieldValue(template, queryUnit.FieldName)
	if err != nil {
		return false, errors.Internal.Wrap(err, "error getting a field value")
	}
	switch queryUnit.Comparator {
	case "==":
		return fieldValue == value, nil
	case "!=":
		return fieldValue != value, nil
	case ">":
		return isBigger(fieldValue, value)
	case "<":
		return isLesser(fieldValue, value)
	case ">=":
		return isBiggerOrEqual(fieldValue, value)
	case "<=":
		return isLesserOrEqual(fieldValue, value)
	case "LIKE":
		return strings.Contains(fieldValue.(string), value.(string)), nil
	default:
		return false, errors.Internal.New("invalid comparator")
	}
}

func replaceWithFakeTrueQuery(query string, start, end int) string {
	return query[:start] + "FakeTrue == 0" + query[end:]
}

func replaceWithFakeFalseQuery(query string, start, end int) string {
	return query[:start] + "FakeFalse == 0" + query[end:]
}

func getQueryUnitString(query string, lastParsedIndex int) (string, int) {
	condIndex, _ := findConditionIndexAndLen(query, lastParsedIndex)
	if condIndex != -1 {
		return query[lastParsedIndex : condIndex-1], condIndex - 1
	}
	return query[lastParsedIndex:], len(query)
}

func getConditionString(query string, lastParsedIndex int) (string, int) {
	condIndex, condLength := findConditionIndexAndLen(query, lastParsedIndex)
	if condIndex != -1 {
		return query[condIndex : condIndex+condLength], condIndex + condLength
	}
	return "", lastParsedIndex
}
