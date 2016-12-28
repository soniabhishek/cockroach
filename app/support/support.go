package support

import (
	"os"
	"runtime"
	"strings"
)

var angelRoot string
var externalDir string

func init() {

	var goPath string
	// windows has ; separator vs linux has :
	if runtime.GOOS == "windows" {
		goPath = strings.Split(os.Getenv("GOPATH"), ";")[0]
	} else {
		goPath = strings.Split(os.Getenv("GOPATH"), ":")[0]
	}

	// Derive the app root directory
	angelRoot = goPath + "/src/github.com/crowdflux/angel"
	externalDir = angelRoot + "/external"
}

func GetAngelRootDir() string {
	return angelRoot
}

func GetExposedDir() string {
	return externalDir
}
