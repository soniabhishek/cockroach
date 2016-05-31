package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"gitlab.com/playment-main/angel/experiments/util"
)

type typeA struct {
	name1 string
}

type typeB struct {
	name2 string
}

type typeC struct {
	name3 string
}

func main1() {
	getType(&typeA{"typeAVal"})
	getType(&typeB{"typeBVal"})
	getType(&typeC{"typeCVal"})
}

func getType(t interface{}) {
	val := reflect.TypeOf(t)
	v, ok := t.(*typeA)
	v.name1 = "this"
	fmt.Println(val)
	fmt.Println(v, ok)
}
func main() {
	str := "ABCasd xyzABCasd xyzABCasd xyzABCasd xyz"
	fmt.Println(binary.Size([]byte(str)))
	fmt.Println(unsafe.Sizeof(str))

	/*fmt.Println(util.IsEmptyOrNil(""))
	fmt.Println(util.GetNestedQuery(`Select * from people`))
	fmt.Println(util.GetNestedQuery(`select t.* from (select * from people) t`))
	fmt.Println(util.GetNestedQuery(`select t.* from (table) t`))*/
	queryCall(`select people.* from (select * from people) people`)
	queryCall(`select t.* from (select * from people) t`)
	queryCall(`Select people.* from people`)
	queryCall(`SELECT p.*, c.* from people p inner join cities c on c.id = p.city_id`)
	queryCall(`SELECT * from people p inner join cities c on c.id = p.city_id`)
	queryCall(`SELECT p.ID, c.* from people as p inner join cities c on c.id = p.city_id`)
	queryCall(`select t.* from (select * from people) t`)
	queryCall(`Select id from (Select id from (Select id from people) people) people`)
	queryCall(`Select id from people ORDER BY first_name ASC`)

	/*SELECT p.*, city.id "city.id" , city.Name "city.name" FROM people p inner join cities city on city.id = p.city_id ORDER BY first_name ASC*/
	queryCall(`SELECT p.*, city.id , city.Name FROM people p inner join cities city on city.id = p.city_id ORDER BY first_name ASC`)
	queryCall(`select people.* , city.*  from people inner join cities on (cities.id = people.city_id) where cities.id = 1`)
	queryCall(`SELECT employees.*, boss.id , boss.name FROM employees JOIN employees AS boss ON employees.boss_id = boss.id`)
	queryCall(`astat sdqwoet`)
	queryCall(`this is spartaa!!!`)
	queryCall(`update table employees set a = 0`)
	queryCall(``)
	queryCall(`select asb.*.asdf from employees`)
}

func queryCall(query string) {
	fmt.Println("\nBaseQuery: ", query)
	finalQ, _ := util.ResolveSelectQuery(query)
	fmt.Println("FinalQuery: ", finalQ)
}
