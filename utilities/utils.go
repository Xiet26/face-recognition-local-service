package utilities

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

func LoadEnvFromFile(config interface{}, configPrefix, envPath string) (err error) {
	godotenv.Load(envPath)
	err = envconfig.Process(configPrefix, config)
	return
}

func LoadEnv(config interface{}, prefix string, source string) error {
	if err := LoadEnvFromDir(config, prefix, source); err != nil {
		return LoadEnvFromFile(config, prefix, source)
	}

	return nil
}

func LoadEnvFromDir(config interface{}, configPrefix, dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	filePaths := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(dir, f.Name()))
	}

	if err := godotenv.Load(filePaths...); err != nil {
		return err
	}
	return envconfig.Process(configPrefix, config)
}

func GetListFilePathFromDir(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filePaths := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(dir, f.Name()))
	}

	return filePaths, nil
}

func StringInArray(str string, arr []string) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if strings.TrimSpace(str) == strings.TrimSpace(val) {
			return true
		}
	}
	return false
}

func IntInArray(num int, arr []int) bool {
	if len(arr) == 0 {
		return false
	}

	for _, val := range arr {
		if num == val {
			return true
		}
	}
	return false
}

func Running() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	running := true
	for running {
		select {
		case <-sig:
			log.Println("\nSignal received, stopping")
			running = false
			break
		}
	}
}

func BindJSON(r *http.Request, obj interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, obj)
}

func RemoveStringInArrayString(array []string, value string) []string {
	if !StringInArray(value, array) {
		return array
	}

	var result []string
	for _, v := range array {
		if v != value {
			result = append(result, v)
		}
	}

	return result
}
func RandomPassWord() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}