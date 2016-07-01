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
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/api/auther"
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
	"runtime"
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

func mainLog() {

	Trace("This is spartaa", "Another")
}

type fluOutputStruct struct {
	ID          uuid.UUID   `json:"flu_id"`
	ReferenceId string      `json:"reference_id"`
	Tag         string      `json:"tag"`
	Status      string      `json:"status"`
	Result      interface{} `json:"result"`
}

/**

{
    "feed_line_units": [
        {
            "flu_id": "dummy_flu_id",
            "reference_id": "dummy_review_id",
            "tag": "FLIPKART_REVIEW_MODERATION",
            "status": "COMPLETED",
            "result": {
                "action": "approved",
                "reason": ""
            }
        }
    ]
}

**/
func maincheck() {
	arr := make([]fluOutputStruct, 0)
	arr = append(arr, fluOutputStruct{
		ID:          uuid.NewV4(),
		ReferenceId: "someRef",
		Tag:         "flp",
		Status:      "completed",
		Result:      models.JsonFake{"1": "One"},
	})
	arr = append(arr, fluOutputStruct{
		ID:          uuid.NewV4(),
		ReferenceId: "someRef",
		Tag:         "flp",
		Status:      "completed",
		Result:      models.JsonFake{"1": "One"},
	})
	fmt.Println(arr)
	sendResp := make(map[string][]fluOutputStruct)
	sendResp["feed_line_units"] = arr
	bty, err := json.Marshal(sendResp)
	fmt.Println(bty, err, string(bty))
}

func Trace(tag string, args ...interface{}) {

	fmt.Println(tag)
	fmt.Println(args)
	x, fn, line, y := runtime.Caller(1)
	fmt.Println(x, fn, line, y)
}
func mainjson() {

	var jsMap models.JsonFake = make(map[string]interface{})
	jsMap["1"] = 2
	fmt.Println(jsMap)
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
func mainuuid() {
	idStr := "2073e7b2-c6c6-4523-9c68-77ccbd220332"
	uu, _ := uuid.FromString(idStr)
	fmt.Println([]byte(idStr))
	fmt.Println(auther.StdProdAuther.GetAPIKey(uu))
}

func main() {
	mp := map[string]string{"key": "The Knight & Day"}
	fmt.Println(strings.Index(mp["key"], "&"))
	fmt.Println(mp)
	bty3, err := json.Marshal(mp)
	fmt.Println(bty3, err)
	fmt.Println(string(bty3))

	hxe := hex.EncodeToString(bty3)
	fmt.Println(hxe)
	hxd, err := hex.DecodeString(hxe)
	str := string(hxd)
	fmt.Println(hxd, err, str)
	fmt.Println(strings.Index(str, "\u0026"))
	buff := bytes.NewBuffer(bty3)
	fmt.Println(buff.String())

	/*rc, ok := buff..(io.ReadCloser)
	if !ok && buff != nil {
		rc = ioutil.NopCloser(buff)
	}*/
	bty := make([]byte, len(bty3))
	fmt.Println(buff.Read(bty))
	fmt.Println(bty)
	fmt.Println(string(bty))
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mainpaytm() {

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
	body := `"{"feed_line_units":[{"flu_id":"3c119ff0-fcac-4165-aea0-1b8716af7618","reference_id":"54396881-1467267881677","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"accept","product_id":"54396881"}},{"flu_id":"8852805a-348f-4943-917b-5c9b836d4620","reference_id":"54396883-1467267881674","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_color":"Red","product_id":"54396883"}},{"flu_id":"f1ddc2ee-4a83-404f-8f66-eca79512eb70","reference_id":"54396885-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"true","action":"accept","attributes_upper_material":"Leather","product_id":"54396883"}},{"flu_id":"4bec0aa3-54be-406e-8c47-effaa042b7fd","reference_id":"54396882-1467267881678","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","message":"Error - Image is incomplete | Action - Provide complete image of the product","product_id":"54396882"}},{"flu_id":"1312a75a-b56f-4b0d-a896-f3389b6d37d9","reference_id":"54396884-1467267881675","tag":"PAYTM_5413","status":"COMPLETED","result":{"accept_with_edit":"false","action":"reject","category_id":"5030","message":"Error - Wrong mapping | Action - Should be mapped in to The Nightwear \u0026 Nighties","product_id":"54396884"}}]}"`
	hash := utilities.GetHMAC(body, "diyqc")
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("qc-uuid", hash)

	fmt.Println(body)
	fmt.Println(hash)
	//client := &http.Client{}
	/*resp, err := client.Do(req)
	fmt.Println("Err:", err)
	fmt.Println("Resp:", resp)
	response, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Err:", err)
	fmt.Println("RespBody:", string(response))*/
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

func maintest() {
	err, u := UnmarshalText([]byte("d6a1e0ad-0795-4e77-9f1e-8e7a96706c27"))
	fmt.Println(err, u)
}

var (
	urnPrefix  = []byte("urn:uuid:")
	byteGroups = []int{8, 4, 4, 4, 12}
)

type UUID [16]byte

func UnmarshalText(text []byte) (err error, u *UUID) {
	if len(text) < 32 {
		err = fmt.Errorf("uuid: UUID string too short: %s", text)
		return
	}

	t := text[:]

	if bytes.Equal(t[:9], urnPrefix) {
		t = t[9:]
	} else if t[0] == '{' {
		t = t[1:]
	}

	//b := u[:]
	b := make([]byte, 16)

	for _, byteGroup := range byteGroups {
		if t[0] == '-' {
			t = t[1:]
		}

		if len(t) < byteGroup {
			err = fmt.Errorf("uuid: UUID string too short: %s", text)
			return
		}

		_, err = hex.Decode(b[:byteGroup/2], t[:byteGroup])

		if err != nil {
			return
		}

		t = t[byteGroup:]
		b = b[byteGroup/2:]
	}

	fmt.Println(b)
	return
}
