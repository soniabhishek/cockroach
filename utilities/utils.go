package utilities

import (
	"time"
	"net/http"
	"gitlab.com/playment-main/angel/app/models"
	"encoding/json"
	"fmt"
	"gitlab.com/playment-main/angel/app/models/status_codes"
	"io/ioutil"
)

const (
	Empty              = ""
	Star               = "*"
	Hyphen             = "-"
	UnderScore         = "_"
	WhiteSpace         = " "
	Colon              = ":"
	Dot                = "."
	Table_Referencer   = "." //Yeah, they both are same.
	Comma              = ","
	Spaced_Comma       = " , "
	Column_Quote       = "\""
	Place_Holder       = "%s"
	Place_Holder_Cover = "(%s)"

	Left_Parentheses  = "("
	Right_Parentheses = ")"
	Parentheses       = "()"
	Left_Braces       = "{"
	Right_Braces      = "{"
	Braces            = "{}"
	Left_Bracket      = "["
	Right_Bracket     = "]"
	Bracket           = "[]"
)

func TimeInMillis()  int64{
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	return umillisec
}

func TimeDiff(absolute bool, times ...int64) int64{

	var newTime, oldTime int64

	if len(times) > 1{
		oldTime = times[0]
		newTime = times[1]
	}else{
		oldTime = times[0]
		newTime = TimeInMillis()
	}

	if absolute{
		return Abs(newTime - oldTime)
	}else{
		return newTime - oldTime
	}
}

func Abs(n int64) int64{
	if n>0{
		return n
	}else{
		return 0-n
	}
}

func GetRandomID() int64 {
	return time.Now().UnixNano()
}

func IsEmptyOrNil(str string) bool {
	return str == Empty
}

func IsValidError(err error) bool {
	return err != nil
}

func ParseFluResponse(resp *http.Response) *models.Response{
	fluResp := &models.Response{}
	fluResp.HttpStatusCode = resp.StatusCode


	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	fmt.Println("response Headers:", resp)
	fmt.Println("response Body:", string(body))
	err := json.Unmarshal(body, fluResp)
	if err != nil {
		//TODO what to do with error
		fmt.Println(err)
		return fluResp
	}
	return fluResp
}


func HttpCodeForCallback(httpStatusCode int) bool {
	switch httpStatusCode {
	case
		http.StatusNotFound,
		http.StatusRequestTimeout,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

func IsValidInternalError(internalCode string) bool{
	switch internalCode {
	case
		status_codes.FF_FluIdNotPresent,
		status_codes.FF_RefIdNotPresent,
		status_codes.FF_TagIdNotPresent,
		status_codes.FF_ResultInvalid,
		status_codes.FF_Other:
		return true
	}
	return false
}
