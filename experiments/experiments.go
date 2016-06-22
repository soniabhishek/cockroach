package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step/manual_step"
	"gitlab.com/playment-main/angel/experiments/util"
	"gitlab.com/playment-main/angel/utilities"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
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

func mainexperiment() (resp *http.Response, err error) {
	hc := http.Client{}
	//req, err := http.NewRequest("POST", APIURL, nil)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)

	fmt.Println(fileContents, fi, body)

	form := url.Values{}
	form.Add("ln", "")
	req, err := http.NewRequest("POST", "http://54.169.7.227/flats", strings.NewReader(form.Encode()))
	fmt.Println("ERR", err)
	//req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("form was %v", form)
	resp, err = hc.Do(req)

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		plog.Error("Error", err)
	}

	fmt.Println(string(response))
	return
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", uri, body)
}

func mainx() {
	path, _ := os.Getwd()
	path += "/test.pdf"
	extraParams := map[string]string{
		"playment_request_id": "12345",
	}
	//request, err := newfileUploadRequest("http://54.169.7.227/flats", extraParams, "files", "/Users/playment/Desktop/a10c5187-7c97-4485-90b7-a427769ceed8.csv")
	request, err := newfileUploadRequest("http://54.169.7.227/flats",
		extraParams, "files",
		"/Users/playment/Downloads/smalldata.txt")
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	} else {
		//var bodyContent []byte
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		/*resp.Body.Read(bodyContent)
		resp.Body.Close()
		fmt.Println(bodyContent)*/

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			plog.Error("Error", err)
		}

		fmt.Println(string(response))
	}
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

func mainUrl() {
	url := "http://54.169.7.227/flats"
	//url := "http://localhost:8080/JServer/HelloServlet"

	body := `{
  "feed_line_units": [
    {
      "flu_id": "PLAYMENT_1",
      "reference_id": "PAYTM_QC_1",
      "tag": "PAYTM_QC",
      "status": "COMPLETED",
      "result": {
        "action": "accept",
        "product_id": "ABC123",
        "sleeve_type": "Half Sleeve"
      }
    },
    {
      "flu_id": "PLAYMENT_2",
      "reference_id": "PAYTM_QC_2",
      "tag": "PAYTM_QC",
      "status": "COMPLETED",
      "result": {
        "action": "reject",
        "product_id": "XYZ321",
        "message": "Image check failed"
      }
    }
  ]
}`
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("Err:", err)
	fmt.Println("Resp:", resp)
	response, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Err:", err)
	fmt.Println("RespBody:", string(response))
}

func mainHmac() {

	//url := "http://localhost:8080/JServer/HelloServlet"
	url := "https://catalogadmin-staging.paytm.com/v1/tp/product/qc-status"

	js := models.JsonFake{}
	js.Scan(`{
  "feed_line_units": [
    {
      "flu_id": "PLAYMENT_1",
      "reference_id": "PAYTM_QC_1",
      "tag": "PAYTM_QC",
      "status": "COMPLETED",
      "result": {
        "action": "accept",
        "product_id": "ABC123",
        "sleeve_type": "Half Sleeve"
      }
    },
    {
      "flu_id": "PLAYMENT_2",
      "reference_id": "PAYTM_QC_2",
      "tag": "PAYTM_QC",
      "status": "COMPLETED",
      "result": {
        "action": "reject",
        "product_id": "XYZ321",
        "message": "Image check failed"
      }
    }
  ]
}`)

	//body := js.String()
	//body := `{"feed_line_units":[{"flu_id":"PLAYMENT_1","reference_id":"PAYTM_QC_1","result":{"action":"accept","product_id":"ABC123","sleeve_type":"Half Sleeve"},"status":"COMPLETED","tag":"PAYTM_QC"},{"flu_id":"PLAYMENT_2","reference_id":"PAYTM_QC_2","result":{"action":"reject","message":"Image check failed","product_id":"XYZ321"},"status":"COMPLETED","tag":"PAYTM_QC"}]}`
	body := `{"feed_line_units":[{"flu_id":"PLAYMENT_1","reference_id":"PAYTM_QC_1","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"accept","product_id":"ABC123","sleeve_type":"Half Sleeve"}},{"flu_id":"PLAYMENT_2","reference_id":"PAYTM_QC_2","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"reject","product_id":"XYZ321","message":"Image check failed"}}]}`
	hash := utilities.GetHMAC(body, "diyqc")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("qc-uuid", hash)

	fmt.Println(body)
	fmt.Println(hash)

}

func main112() {
	jf := models.JsonFake{}
	jf.Scan(`{"feed_line_units":[{"flu_id":"PLAYMENT_1","reference_id":"PAYTM_QC_1","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"accept","product_id":"ABC123","sleeve_type":"Half Sleeve"}},{"flu_id":"PLAYMENT_2","reference_id":"PAYTM_QC_2","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"reject","product_id":"XYZ321","message":"Image check failed"}}]}`)
	fmt.Println(jf)

	hmacKeyStr := "diyqc"
	//hmacKeyStr := "a9a6bd8c-314c-11e6-ac61-9e71128cae77"
	key := []byte(hmacKeyStr)
	sig := hmac.New(sha256.New, key)

	fmt.Println(jf.String())
	fmt.Println(jf.String() == `{"feed_line_units":[{"flu_id":"PLAYMENT_1","reference_id":"PAYTM_QC_1","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"accept","product_id":"ABC123","sleeve_type":"Half Sleeve"}},{"flu_id":"PLAYMENT_2","reference_id":"PAYTM_QC_2","tag":"PAYTM_QC","status":"COMPLETED","result":{"action":"reject","product_id":"XYZ321","message":"Image check failed"}}]}`)
	sig.Write([]byte(jf.String()))
	myHash := hex.EncodeToString(sig.Sum(nil))
	s1 := "09732f0c328225c99caa166ea237fc09c77e68e6a865e72086676425a0886eca"
	s2 := "2674fe4b59546d02b2689958e1c72122b385c1a93a0a1663abf0545ab0d2568c"
	s3 := "2674fe4b59546d02b2689958e1c72122b385c1a93a0a1663abf0545ab0d2568c"
	fmt.Println(s1 == s2)
	fmt.Println(s2 == s3)
	fmt.Println(s3 == s1)

	fmt.Println(myHash)
	fmt.Println(myHash == s1)
	fmt.Println(myHash == s2)
	fmt.Println(myHash == s3)
}

func mainJson() {
	jf := models.JsonFake{}
	jf.Scan(`{

	"feed_line_units" : [

		{

			"flu_id" 	        : "PLAYMENT_1",

			"reference_id" 	: "PAYTM_QC_1",

			"tag" 		: "PAYTM_QC",

			"status" 		: "COMPLETED",

			"result" 		: {

				"action" 	      : "accept",

				"product_id"   : "ABC123",

				"sleeve_type" : "Half Sleeve"


			}

		},

		{

			"flu_id" 	: "PLAYMENT_2",

			"reference_id" 	: "PAYTM_QC_2",

			"tag" 		: "PAYTM_QC",

			"status" 		: "COMPLETED",

			"result" 		: {

				"action" 	      : "reject",

				"product_id"   : "XYZ321",

				"message"     :  "Image check failed"


			}

		}

	]

}`)

	fmt.Println(jf.String())
	hmacKeyStr := "diyqc"
	//hmacKeyStr := "a9a6bd8c-314c-11e6-ac61-9e71128cae77"
	key := []byte(hmacKeyStr)
	sig := hmac.New(sha256.New, key)

	jfBytes, _ := json.Marshal(jf.String())
	sig.Write(jfBytes)
	fmt.Println(hex.EncodeToString(sig.Sum(nil)))
}

func mainstaticfiles() {
	router := gin.Default()
	router.Static("/assets", "/Users/gambler/Desktop/tmp/")
	router.StaticFS("/downloadedfiles", http.Dir("/Users/gambler/Desktop/tmp/"))
	router.StaticFile("/test.txt", "/Users/gambler/Desktop/tmp/test.txt")

	// Listen and server on 0.0.0.0:8080
	router.Run(":8080")
}

func mainMapCheck() {
	var mp map[string]string = make(map[string]string)
	mp["k"] = "u"
	mp["p"] = "n"
	val1, ok1 := mp["k"]
	fmt.Println(val1, ok1)
	val2, ok2 := mp["x"]
	fmt.Println(val2, ok2)
}
func mainMongo() {
	fluInputQueue := feed_line_repo.NewInputQueue()

	err := fluInputQueue.MarkFinished()

	fmt.Println(err)
}

func main_gin() {
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.Run(":8080")
}

func mainhttp() {
	fs := http.FileServer(http.Dir("model"))
	http.Handle("/", fs)

	fmt.Println("Listening...")
	http.ListenAndServe(":3000", nil)
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

var path = "/Users/playment/Desktop/test.txt"

func mainFile() {
	createFile()
	writeFile()
	readFile()
	//deleteFile()
}

func createFile() {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		checkError(err)
		defer file.Close()
	}
}

func writeFile() {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// write some text to file
	_, err = file.WriteString("halo\n")
	checkError(err)
	_, err = file.WriteString("mari belajar golang\n")
	checkError(err)

	// save changes
	err = file.Sync()
	checkError(err)
}

func readFile() {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			checkError(err)
		}
		if n == 0 {
			break
		}
	}
	fmt.Println(string(text))
	checkError(err)
}

func deleteFile() {
	// delete file
	var err = os.Remove(path)
	checkError(err)
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
