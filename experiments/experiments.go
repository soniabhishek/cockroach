package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"bytes"
	"encoding/csv"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/manual_step"
	"gitlab.com/playment-main/angel/experiments/util"
	"gitlab.com/playment-main/angel/utilities"
	"io/ioutil"
	"net/http"
	"os"
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

func mainFlipkart() {

	//url := "http://localhost:8080/JServer/HelloServlet"
	url := "https://catalogadmin-staging.paytm.com/v1/tp/product/qc-status"

	//body := `{"feed_line_units":[{"flu_id":"9c019e3f-4e95-4837-90ee-0c25d82838fe","reference_id":"54396675-1466574896063","tag":"PAYTM_TSHIRT","status":"COMPLETED","result":{"action":"reject","product_id":"54396675","message":"Image check failed"}}]}`
	body := `{"feed_line_units":[{"flu_id":"dd0f4e24-2957-465d-8686-b28448c7f966","reference_id":"54396675-1467031774629","tag":"PAYTM_5030","status":"COMPLETED","result":{"action":"accept","product_id":"54396675"}}]}`
	hash := utilities.GetHMAC(body, "diyqc")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("qc-uuid", hash)

	fmt.Println(body)
	fmt.Println(hash)
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("Err:", err)
	fmt.Println("Resp:", resp)
	response, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Err:", err)
	fmt.Println("RespBody:", string(response))
}

func mainHit() {
	//file := `/Users/playment/Downloads/smalldata.txt`
	//file := `/Users/playment/Desktop/a10c5187-7c97-4485-90b7-a427769ceed8.csv`
	file := `/Users/playment/Desktop/a10c5187-7c97-4485-90b7-a427769ceed8.txt`
	url := `http://54.169.7.227/flats/`

	filename, err := manual_step.FlattenCSV(file, url, uuid.NewV4())
	fmt.Println("Err:", err)
	plog.Info("Sent file for upload: ", filename)
}

func main() {

	//url := "http://localhost:8080/JServer/HelloServlet"
	url := "https://catalogadmin-staging.paytm.com/v1/tp/product/qc-status"

	//body := js.String()
	//body := `{"feed_line_units":[{"flu_id":"PLAYMENT_1","reference_id":"PAYTM_QC_1","result":{"action":"accept","product_id":"ABC123","sleeve_type":"Half Sleeve"},"status":"COMPLETED","tag":"PAYTM_QC"},{"flu_id":"PLAYMENT_2","reference_id":"PAYTM_QC_2","result":{"action":"reject","message":"Image check failed","product_id":"XYZ321"},"status":"COMPLETED","tag":"PAYTM_QC"}]}`
	//body := `{"feed_line_units":[{"flu_id":"PLAYMENT_1","reference_id":"PAYTM_QC_1","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"accept","product_id":"ABC123","sleeve_type":"Half Sleeve"}},{"flu_id":"PLAYMENT_2","reference_id":"PAYTM_QC_2","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"reject","product_id":"XYZ321","message":"Image check failed"}}]}`
	//body := `{"feed_line_units":[{"flu_id":"2e2f12f6-5183-408e-834d-07faf1535e9e","reference_id":"54396695-1466594495186","tag":"PAYTM 5030","status":"COMPLETED","result":{"action":"accept","product_id":"54396695"}}]}`
	//body := `{"feed_line_units":[{"flu_id":"12255044-a5b1-44b6-a222-f96730bf82c2","reference_id":"54396690-1466594495183","tag":"PAYTM 5030","status":"COMPLETED","result":{"action":"reject","product_id":"54396690","message":"Image check failed"}}]}`
	//body := `{"feed_line_units":[{"flu_id":"ef6cf491-80dd-427d-b54a-e6ff43aa331b","reference_id":"54396730-1466684151464","tag":"PAYTM 5030","status":"COMPLETED","result":{"action":"reject","product_id":"54396730","message":"Image check failed"}}]}`
	//body := `{"feed_line_units":[{"flu_id":"9c019e3f-4e95-4837-90ee-0c25d82838fe","reference_id":"54396675-1466574896063","tag":"PAYTM_TSHIRT","status":"COMPLETED","result":{"action":"reject","product_id":"54396675","message":"Image check failed"}}]}`
	//body := `{"feed_line_units":[{"flu_id":"dd0f4e24-2957-465d-8686-b28448c7f966","reference_id":"54396675-1467031774629","tag":"PAYTM_5030","status":"COMPLETED","result":{"action":"accept","product_id":"54396675"}}]}`

	//body := `"{"feed_line_units":[{"flu_id":"3c119ff0-fcac-4165-aea0-1b8716af7618","reference_id":"54396881-1467267881677","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"accept","product_id":"54396881"}},{"flu_id":"8852805a-348f-4943-917b-5c9b836d4620","reference_id":"54396883-1467267881674","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_color":"Red","product_id":"54396883"}},{"flu_id":"f1ddc2ee-4a83-404f-8f66-eca79512eb70","reference_id":"54396885-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_upper_material":"Leather","product_id":"54396883"}},{"flu_id":"4bec0aa3-54be-406e-8c47-effaa042b7fd","reference_id":"54396882-1467267881678","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","message":"Error - Image is incomplete | Action - Provide complete image of the product","product_id":"54396882"}},{"flu_id":"1312a75a-b56f-4b0d-a896-f3389b6d37d9","reference_id":"54396884-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to The Nightwear & Nighties","product_id":"54396884"}}]}"`
	//body := `"{"feed_line_units":[{"flu_id":"3c119ff0-fcac-4165-aea0-1b8716af7618","reference_id":"54396881-1467267881677","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"accept","product_id":"54396881"}},{"flu_id":"8852805a-348f-4943-917b-5c9b836d4620","reference_id":"54396883-1467267881674","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_color":"Red","product_id":"54396883"}},{"flu_id":"f1ddc2ee-4a83-404f-8f66-eca79512eb70","reference_id":"54396885-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_upper_material":"Leather","product_id":"54396883"}},{"flu_id":"4bec0aa3-54be-406e-8c47-effaa042b7fd","reference_id":"54396882-1467267881678","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","message":"Error - Image is incomplete | Action - Provide complete image of the product","product_id":"54396882"}},{"flu_id":"1312a75a-b56f-4b0d-a896-f3389b6d37d9","reference_id":"54396884-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to The Nightwear \u0026 Nighties","product_id":"54396884"}}]}"`
	body := `{"feed_line_units":[{"flu_id":"f394ded5-1b8c-4167-a15f-a7fd26312214","reference_id":"59131192-1468849990418","tag":"PAYTM_5237","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to TShirt","product_id":"59131192"}},{"flu_id":"aa58c816-662c-480d-adea-e6f7ae534cdd","reference_id":"59415296-1469000729547","tag":"PAYTM_5237","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to TShirt","product_id":"59415296"}},{"flu_id":"9db023b5-4add-4af1-bc8d-14a77143a42d","reference_id":"59131191-1468849990419","tag":"PAYTM_5237","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to TShirt","product_id":"59131191"}}]}`
	hash := utilities.GetHMAC(body, "aa566750-49b6-11e6-beb8-9e71128cae77")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("qc-uuid", hash)

	fmt.Println(body)
	fmt.Println(hash)
	/*client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("Err:", err)
	fmt.Println("Resp:", resp)
	response, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Err:", err)
	fmt.Println("RespBody:", string(response))*/
}

func mainx() {
	//k1 := "diyqc"
	k2 := "a9a6bd8c-314c-11e6-ac61-9e71128cae77"
	body1 := `"{"feed_line_units":[{"flu_id":"3c119ff0-fcac-4165-aea0-1b8716af7618","reference_id":"54396881-1467267881677","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"accept","product_id":"54396881"}},{"flu_id":"4bec0aa3-54be-406e-8c47-effaa042b7fd","reference_id":"54396882-1467267881678","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","message":"Error - Image is incomplete | Action - Provide complete image of the product","product_id":"54396882"}},{"flu_id":"1312a75a-b56f-4b0d-a896-f3389b6d37d9","reference_id":"54396884-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to The Nightwear & Nighties","product_id":"54396884"}},{"flu_id":"8852805a-348f-4943-917b-5c9b836d4620","reference_id":"54396883-1467267881674","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_color":"Red","product_id":"54396883"}},{"flu_id":"f1ddc2ee-4a83-404f-8f66-eca79512eb70","reference_id":"54396885-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_upper_material":"Leather","product_id":"54396883"}}]}"`
	body2 := `{"feed_line_units":[{"flu_id":"3c119ff0-fcac-4165-aea0-1b8716af7618","reference_id":"54396881-1467267881677","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"accept","product_id":"54396881"}},{"flu_id":"4bec0aa3-54be-406e-8c47-effaa042b7fd","reference_id":"54396882-1467267881678","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","message":"Error - Image is incomplete | Action - Provide complete image of the product","product_id":"54396882"}},{"flu_id":"1312a75a-b56f-4b0d-a896-f3389b6d37d9","reference_id":"54396884-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to The Nightwear & Nighties","product_id":"54396884"}},{"flu_id":"8852805a-348f-4943-917b-5c9b836d4620","reference_id":"54396883-1467267881674","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_color":"Red","product_id":"54396883"}},{"flu_id":"f1ddc2ee-4a83-404f-8f66-eca79512eb70","reference_id":"54396885-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_upper_material":"Leather","product_id":"54396883"}}]}`
	hash1 := utilities.GetHMAC(body1, k2)
	hash2 := utilities.GetHMAC(body2, k2)

	fmt.Println(hash1)
	fmt.Println(hash2)
	fmt.Println("755572b69bf8506f031fb31b684a5706f9e45e9a3cf891e61171b8853c81905c")
}

func maineditor() {
	plaintext := []byte("a9a6bd8c-314c-11e6-ac61-9e71128cae77")
	fmt.Printf("%s\n", plaintext)
	ciphertext, err := utilities.Encrypt(string(plaintext))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%0x\n", ciphertext)
	result, err := utilities.Decrypt(string(ciphertext))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", result)
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
func maincsv() {
	csvfile, err := os.Create("/Users/gambler/Desktop/output.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close()

	records := [][]string{{"item1", "value1"}, {"item2", "value2"}, {"item3", "value3"}}

	writer := csv.NewWriter(csvfile)
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	writer.Flush()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

func main123() {
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
