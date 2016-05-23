package util

import (
	"errors"
	"fmt"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/experiments/model"
	"reflect"
	"strings"
)

// Table Names
const (
	DB_TAG         = "db"
	NO_table       = ""
	People_table   = "people"
	City_table     = "cities"
	Employee_table = "employees"
	MacroTaskTable = "macro_tasks"
	ProjectsTable  = "projects"
)

// SQL Keywords Constants
const (
	SELECT     = "select"
	FROM       = "from"
	WHERE      = "where"
	ON         = "on"
	AS         = "as"
	JOIN       = "join"
	INNER_JOIN = "inner"
	LEFT_JOIN  = "left"
	RIGHT_JOIN = "right"
	OUTER_JOIN = "outer"
)

//SQL Keywords Slice
var keywords = []string{
	SELECT,
	FROM,
	WHERE,
	ON,
	AS,
	JOIN,
	INNER_JOIN,
	LEFT_JOIN,
	RIGHT_JOIN,
	OUTER_JOIN,
}

//Table to Alias map
var aliasMap map[string]string = make(map[string]string)

//ToDo
func QueryMapper(generic interface{}, query string) {

}

/*
SELECT people.*, cities.* from people inner join cities on cities.id = people.city_id
*/
func ResolveSelectQuery(query string) (finalQuery string, err error) {
	return getResolvedSelectQuery(query, false)
}

func getResolvedSelectQuery(query string, isNested bool) (finalQuery string, err error) {
	//fmt.Println(query)
	if IsEmptyOrNil(query) {
		return Empty, errors.New("Empty String")
	}
	leftQ, nestedQ, rightQ, doesItHaveNested := GetNestedQuery(query)
	var resolvedQ, outerQ string
	if doesItHaveNested {
		resolvedQ, err = getResolvedSelectQuery(nestedQ, true)
		outerQ, err = getResolvedQuery(leftQ+WhiteSpace+Place_Holder_Cover+WhiteSpace+rightQ, isNested)
		finalQuery = fmt.Sprintf(outerQ, strings.TrimSpace(resolvedQ))
	} else {
		finalQuery, err = getResolvedQuery(query, isNested)
	}
	if IsValidError(err) {
		return query, err
	}
	return
}

func getResolvedQuery(query string, isNested bool) (finalQuery string, err error) {
	//fmt.Println("Query: ", query)
	if IsEmptyOrNil(query) {
		return Empty, errors.New("Empty String")
	}
	spell := strings.ToLower(query)
	selct := strings.Index(spell, SELECT) + len(SELECT)
	from := strings.Index(spell, FROM)
	on := strings.Index(spell, ON)
	where := strings.Index(spell, WHERE)
	var end int = -1
	if where != -1 {
		end = where
	}
	if on != -1 {
		end = on
	}

	var tables string
	if end != -1 {
		tables = spell[from+4 : end]
	} else {
		tables = spell[from+4:]
	}
	//fmt.Println(tables, spell)
	tableslist, err := getTables(tables)
	if IsValidError(err) {
		fmt.Println(err)
		return finalQuery, err
	}
	//fmt.Println(tableslist)

	sub := spell[selct:from]
	//fmt.Println(sub)
	searchQ := Empty

	if strings.TrimSpace(sub) != Star {
		subs := strings.Split(sub, Comma)
		//fmt.Println(subs)
		//fmt.Println(len(subs), len(tableslist))
		for v := range subs {
			if subs[v] == Place_Holder {
				continue
			}
			tmp := strings.Split(subs[v], Dot)
			if len(tmp) < 2 {
				return query, errors.New("Query needs NO expansion")
			}
			alias := tmp[0]
			wildCard := tmp[1]
			if strings.TrimSpace(wildCard) == Star {
				aliasStruct := tableslist[v]
				fin := reflect.TypeOf(aliasStruct)
				for i := 0; i < fin.NumField(); i++ {
					col := fin.Field(i)
					//searchQ += alias + Dot + col.Name + Spaced_Comma
					searchQ += getQ(alias, col.Tag.Get(DB_TAG), isNested)
					//searchQ += getQ(alias, col.Name)
				}
			} else {
				//searchQ += alias + Dot + wildCard + Spaced_Comma
				searchQ += getQ(alias, wildCard, isNested)
			}
		}
	} else {
		for i := range tableslist {
			aliasStruct := tableslist[i]
			fin := reflect.TypeOf(aliasStruct)
			aliasName := getAlias(getRelatedTable(aliasStruct))
			for i := 0; i < fin.NumField(); i++ {
				col := fin.Field(i)
				searchQ += getQ(aliasName, col.Tag.Get(DB_TAG), isNested)
				//searchQ += getQ(aliasName, col.Name)
			}
		}
	}

	//fmt.Println("Search query : ", searchQ)
	//fmt.Println("Search query : ", searchQ[:strings.LastIndex(searchQ, Comma)])

	finalQuery += query[:selct+1]
	finalQuery += searchQ[:strings.LastIndex(searchQ, Comma)] + WhiteSpace
	finalQuery += query[from:]
	//fmt.Println("FinalQuery: ", finalQuery)
	return
}

func getTables(str string) (tables []interface{}, err error) {
	tables = make([]interface{}, 0)
	var tags = strings.Split(str, WhiteSpace)
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		table := getRelatedStruct(tag)

		if table != nil {
			tables = append(tables, table)
			if i < len(tags)-1 {
				if !isKeyword(tags[i+1]) {
					aliasMap[tag] = tags[i+1]
					i++
				} else if areEqual_Without_Whitespaces(tags[i+1], AS) {
					aliasMap[tag] = tags[i+2]
					i += 2
				}
			}
		} else if IsValidTag(tag) {
			err = errors.New("Unknown Table [" + tag + "]")
			return
		}
	}
	/*fmt.Println("Given string [", str, "] :: Tables are: ", len(tables))
	for i := range tables {
		fmt.Println(reflect.TypeOf(tables[i]))
	}*/
	return tables, nil
}

func getRelatedStruct(tableName string) interface{} /*, err error)*/ {
	switch tableName {

	case People_table:
		return model.Person{}

	case City_table:
		return model.City{}

	case Employee_table:
		return model.Employee{}
	case MacroTaskTable:
		return models.MacroTask{}

	case ProjectsTable:
		return models.Project{}

	default:
		return nil
	}
}

func getRelatedTable(structName interface{}) string /*, err error)*/ {

	switch structName {

	case model.Person{}:
		return People_table

	case model.City{}:
		return City_table

	case model.Employee{}:
		return Employee_table

	case models.MacroTask{}:
		return MacroTaskTable

	case models.Project{}:
		return ProjectsTable
	default:
		return NO_table
	}
}

func getAlias(alias string) string {
	if val, ok := aliasMap[strings.ToLower(alias)]; ok {
		return val
	} else {
		return strings.ToLower(alias)
	}
}

func isKeyword(tag string) bool {
	for v := range keywords {
		if areEqual_Without_Whitespaces(keywords[v], strings.ToLower(tag)) {
			return true
		}
	}
	return false
}
func GetNestedQuery(query string) (leftQ string, nestedQ string, rightQ string, isNested bool) {
	l := strings.Index(query, Left_Parentheses)
	r := strings.LastIndex(query, Right_Parentheses)
	if l != -1 && r != -1 {
		insideQ := query[l+1 : r]
		nxt := strings.Index(insideQ, WhiteSpace)
		if nxt != -1 && isKeyword(strings.TrimSpace(insideQ[:nxt])) {
			nestedQ = insideQ
			leftQ = query[:l]
			rightQ = query[r+1:]
			isNested = true
			return
		}

	}
	return
}

func getQ(aliasName string, col string, isNested bool) string {
	if IsEmptyOrNil(aliasName) || IsEmptyOrNil(col) {
		return Empty
	}
	aliasName = strings.TrimSpace(aliasName)
	col = strings.TrimSpace(col)
	fmt.Println("IsNested: ", isNested)
	if isNested {
		return aliasName + Dot + col + Spaced_Comma
	} else {
		return aliasName + Dot + col + WhiteSpace + Column_Quote + aliasName + Dot + col + Column_Quote + Spaced_Comma
	}
}
