package parser_json

import (
	"encoding/json"
	"github.com/gohouse/e"
	"github.com/gohouse/i18n"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ParserJson struct {
	opts *i18n.Options
	// 示例: /zh-cn/error.json
	// {
	//  "params_format_error": "参数格式有误",
	//  "params_missing": "参数缺失",
	//  "err2": {
	//    "aa": "aaxx",
	//    "bb": "bbxx"
	//  }
	//}
	// map["zh-cn"]["error"]["params_format_error"]
	val map[string]map[string]interface{}
}

var _ i18n.IParser = &ParserJson{}

func NewParserJson() *ParserJson {
	return &ParserJson{val: make(map[string]map[string]interface{})}
}

func (pj *ParserJson) SetOptions(opts *i18n.Options) {
	pj.opts = opts
}

func (pj *ParserJson) Parse() e.E {
	// 获取lang目录的所有文件并解析
	var s []string
	fileAll, err := GetAllFile(pj.opts.LangDirectory, s)
	if err != nil {
		return e.New(err.Error())
	}

	// 去掉目录前缀, 获取语言和文件
	// 解析文件内容, 放入结果集中(pj.val)
	for _, item := range fileAll {
		fileSuf := strings.Replace(item, pj.opts.LangDirectory, "", 1)
		fileSuf = strings.TrimLeft(fileSuf, "/")
		//fileAllReal = append(fileAllReal, fileSuf)

		// 解析语言和文件名
		split := strings.Split(fileSuf, "/")
		if len(split) != 2 {
			return e.New("目录格式错误")
		}
		fileNameStr := strings.TrimRight(split[1], ".json")

		// 解析内容
		bytes, err := pj.ReadBytesFromFile(item)
		if err != nil {
			return e.New(err.Error())
		}

		var js map[string]interface{}
		err = json.Unmarshal([]byte(string(bytes)), &js)
		if err != nil {
			return e.New(err.Error())
		}

		// 保存到pj.val内存中
		langKey := StringToKey(split[0])
		if _, ok := pj.val[langKey]; !ok {
			pj.val[langKey] = make(map[string]interface{})
		}
		pj.val[langKey][fileNameStr] = js
	}
	return nil
}

func StringToKey(str string) string {
	return strings.Replace(str, "-", "_", -1)
}

func (pj *ParserJson) ReadBytesFromFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func (pj *ParserJson) Load(key string, defaultVal ...string) interface{} {
	if key == "" {
		return nil
	}
	var split []string
	// 如果key包含了点,则为多级调用
	if strings.Contains(key, ".") {
		split = strings.Split(key, ".")
	} else { // 如果key不包含点, 则就是直接一级调用
		split = []string{key}
	}

	// 取指定语言的配置
	var currentVal interface{} = pj.val[StringToKey(pj.opts.DefaultLang)]
	for _, item := range split {
		if v, ok := currentVal.(map[string]interface{}); ok {
			currentVal = v[item]
		} else {
			currentVal = nil
		}
	}

	// 如果没有取到且传入了默认值, 则返回默认值
	if currentVal == nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return currentVal
}

func GetAllFile(dirname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := dirname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				log.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := dirname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}