package translator

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/anhhuy1010/DATN-cms-customer/config"
	"github.com/anhhuy1010/DATN-cms-customer/helpers/util"
)

var dataTranslation map[string]map[string]string
var localeDefault = []string{"en", "my", "id", "th", "ms", "jp", "vn", "ko"}

func GetLocale(ctx context.Context) string {
	locale, ok := util.GetKeyFromContext(ctx, "locale")
	if !ok {
		cfg := config.GetConfig()
		locale = cfg.GetString("server.locale")
		//SetLocale(ctx, locale.(string))
	}

	return locale.(string)
}

func SetLocale(ctx context.Context, locale string) context.Context {
	if !IsLocaleSupported(locale) {
		cfg := config.GetConfig()
		locale = cfg.GetString("server.locale")
	}

	ctx = context.WithValue(ctx, "locale", locale)

	return ctx
}

func LoadFileTranslation() {
	dataTranslation = make(map[string]map[string]string)

	var files []string

	root := "languages"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Println("[ERROR] filepath.Walk", err)
	}
	for _, file := range files {
		locate := strings.ReplaceAll(file, root+"/", "")
		locate = strings.ReplaceAll(locate, ".json", "")
		if locate == "" || locate == "languages" {
			continue
		}
		jsonFile, err := os.Open(file)
		if err != nil {
			fmt.Println("[ERROR] os.Open", err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var dataLocateTranslation map[string]string
		json.Unmarshal(byteValue, &dataLocateTranslation)
		dataTranslation[locate] = dataLocateTranslation
	}
}

func Trans(ctx context.Context, key string) string {
	locate, ok := util.GetKeyFromContext(ctx, "locale")
	if !ok {
		fmt.Println("can not get locate")
		return key
	}
	dataLocateTranslation, ok := dataTranslation[locate.(string)]
	if !ok {
		fmt.Println("can not get locate array")
		return key
	}
	val, ok := dataLocateTranslation[key]
	if !ok {
		return key
	}
	return val
}

func IsLocaleSupported(locale interface{}) bool {
	switch reflect.TypeOf(localeDefault).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(localeDefault)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(locale, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}
