package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/z1003031335/go-fileutil"
)

// QR indicates query result
type QR int

const (
	// TranslSucc represents success of translation
	TranslSucc = iota
	// TranslFail represents failure of translation
	TranslFail
)

// QueryResult is result which is obtained by netowork or local.
type QueryResult struct {
	Translation string
	Explains    []string
	Result      QR
}

// Query is query interface,in order to avoid using another netowork query interface in the future
type Query interface {
	Query(word string) QueryResult
}

// YoudaoImpl for querying word
type YoudaoImpl struct {
	apiKey  string
	keyfrom string
}

// NewYoudaoImpl new instance of yaodaoImpl
func NewYoudaoImpl(apiKey, keyfrom string) *YoudaoImpl {
	return &YoudaoImpl{
		apiKey:  apiKey,
		keyfrom: keyfrom,
	}
}

//http://fanyi.youdao.com/openapi.do?keyfrom=youdao-cli&key=1879868570&type=data&doctype=json&version=1.1&q=
//http://fanyi.youdao.com/openapi.do?keyfrom=yaodao-cli&key=1879868570&type=data&doctype=json&version=1.1&q=hello
func (youdao *YoudaoImpl) getAPIURL(word string) string {
	url := "http://fanyi.youdao.com/openapi.do"
	// url := "http://fanyi.youdao.com/openapi.do?keyfrom=" + youdao.keyfrom + "&key=" + youdao.apiKey
	// url += "&type=data&doctype=json&version=1.1&q=" + word
	return url
}

func getBytes(url string, data url.Values) ([]byte, error) {
	client := &http.Client{}
	url = url + "?" + data.Encode()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.101 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}

// Query implementing youdao
func (youdao *YoudaoImpl) Query(word string) *QueryResult {
	apirURL := youdao.getAPIURL(word)
	params := url.Values{}
	params.Add("keyfrom", youdao.keyfrom)
	params.Add("key", youdao.apiKey)
	params.Add("type", "data")
	params.Add("doctype", "json")
	params.Add("version", "1.1")
	params.Add("q", word)

	// &type=data&doctype=json&version=1.1&q=
	data, err := getBytes(apirURL, params)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(data))
	var dat map[string]interface{}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		log.Fatal(err)
	}
	trans := dat["translation"]
	transCon := trans.([]interface{})[0]
	var result QR = TranslSucc
	if transCon == word {
		result = TranslFail
	}
	// expl := dat["basic"].([]interface{})
	basic := dat["basic"]
	var explainStrs []string
	if basic != nil {
		var explains interface{}
		explains = basic.(map[string]interface{})["explains"]
		exp := reflect.ValueOf(explains)
		expLen := exp.Len()
		explainStrs = make([]string, expLen)
		for i := 0; i < expLen; i++ {
			explainStrs[i] = exp.Index(i).Interface().(string)
		}
	} else {
		explainStrs = make([]string, 0)
	}
	return &QueryResult{
		Translation: transCon.(string),
		Explains:    explainStrs,
		Result:      result,
	}
}

func query(word string) string {
	return "ok"
}

// Storage for storing and querying local data
type Storage interface {
	QueryWord(content string) (*Word, bool)
	AddWord(word *Word)
}

// FileStorage using file storage such as json,xml and so on.
type FileStorage struct {
	storage map[string]*Word
}

// QueryWord implement
func (storage *FileStorage) QueryWord(word string) (*Word, bool) {
	if localWord, ok := storage.storage[word]; ok {
		localWord.QueryCount++
		storage.serilize()
		return localWord, true
	}
	return nil, false
}

// AddWord implement
func (storage *FileStorage) AddWord(word *Word) {
	storage.storage[word.Content] = word
	storage.serilize()
}

// NewFileStorage creating a FileStorage and init it
func NewFileStorage() *FileStorage {
	storage := &FileStorage{}
	storage.Init()
	return storage
}
func createDirectoryIfNotExist(dir string) bool {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			os.MkdirAll(dir, os.ModePerm)
		} else {
			// other error
		}
		return false
	}
	return true
}
func createFileIfNotExist(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			log.Println("Creating work dir", file)
			// file does not exist
			os.Create(file)
		} else {
			// other error
		}
		return false
	}
	return true
}

func unmarshal(filename string, inter interface{}) {
	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	if len(configBytes) <= 0 {
		return
	}
	err = json.Unmarshal(configBytes, inter)
	if err != nil {
		log.Fatal(err)
	}
}

// InstallDir app install dir
var InstallDir string

func checkInstallDir() bool {
	InstallDir = os.Getenv("HOME") + "/app/go-word"
	return createDirectoryIfNotExist(InstallDir)
}

// Init map
func (storage *FileStorage) Init() {
	created := createFileIfNotExist(storage.getSerilizeLocation())
	if created {
		storage.storage = make(map[string]*Word)
		unmarshal(storage.getSerilizeLocation(), &storage.storage)
	}
}

func (storage *FileStorage) getSerilizeLocation() string {
	return InstallDir + "/data.json"
}

func (storage *FileStorage) serilize() {
	fileutil.Marshal(storage.getSerilizeLocation(), storage.storage)
}

// Export file
type Export interface {
	Export(wordList []*Word)
}

// ExportTarget exports target and sort it.
type ExportTarget struct {
	ExportObject int
}

const (
	// ExportWord exports original word
	ExportWord = 1 << iota
	// ExportTranslation export translation of word
	ExportTranslation
	// ExportExplains export explains of word
	ExportExplains
	// ExportQueryCount export query count of word
	ExportQueryCount
)

// Filter export data
type Filter func(*Word) bool

// ExportStrategy export strategy
type ExportStrategy struct {
	ExportTargets []*ExportTarget
	Separator     string
	Filter        Filter
}

// AddExportTarget add export
func (strategy *ExportStrategy) AddExportTarget(exportTarget *ExportTarget) {
	strategy.ExportTargets = append(strategy.ExportTargets, exportTarget)
}

// AbsExport common export
type AbsExport struct {
	ExportStrategy *ExportStrategy
}

// GetExportString :
func (export *AbsExport) getExportString(wordList []*Word, wordSeparator string) string {
	if len(wordList) == 0 {
		return ""
	}
	buffer := bytes.NewBufferString("")
	for i, word := range wordList {
		if word == nil || !export.ExportStrategy.Filter(word) {
			continue
		}
		fmt.Println("Exporting ", word)
		for _, exportTarget := range export.ExportStrategy.ExportTargets {
			insSep := false
			if (exportTarget.ExportObject & ExportWord) != 0 {
				buffer.WriteString(word.Content)
				insSep = true
			}
			if (exportTarget.ExportObject & ExportTranslation) != 0 {
				if insSep {
					buffer.WriteString(export.ExportStrategy.Separator)
				}
				buffer.WriteString(word.Content)
				insSep = true
			}
		}
		if i != len(export.ExportStrategy.ExportTargets)-1 {
			buffer.WriteString(wordSeparator)
		}
	}
	return buffer.String()
}

// TxtExport txt file
type TxtExport struct {
	AbsExport
	ExportLocation string
}

// Export txt file export impl
func (txtExport *TxtExport) Export(wordList []*Word) {
	exportStr := txtExport.getExportString(wordList, "\n")
	err := ioutil.WriteFile(InstallDir+"/export_data.txt", []byte(exportStr), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

// Word is mainly for recording query count.
type Word struct {
	Content           string
	TranslatedContent string
	QueryCount        int
	Explains          []string
}

// ExportToTxtFile export to file
func ExportToTxtFile(storage interface{}, strategy *ExportStrategy) {
	expLoc := InstallDir + "/export_data.txt"
	log.Println("Exporting to " + expLoc)
	export := &TxtExport{
		ExportLocation: expLoc,
	}
	export.ExportStrategy = strategy
	switch storage.(type) {
	case *FileStorage:
		fileStorage := storage.(*FileStorage)
		wordList := make([]*Word, len(fileStorage.storage))
		for _, wordEntry := range fileStorage.storage {
			wordList = append(wordList, wordEntry)
		}
		export.Export(wordList)
	default:
		log.Fatal("Not found storage", storage)
	}
	log.Println("Export done!")
}

func isOption(opt string) bool {
	opts := []string{"--export", "--filter", "--exclude-zh"}
	isOpt := false
	for _, val := range opts {
		isOpt = strings.HasPrefix(opt, val)
		if isOpt {
			return true
		}
	}
	return false
}

func getArgValue(i *int, args []string) (string, bool) {
	hasValue := (*i) != len(args)-1 && !isOption(args[(*i)+1])
	if hasValue {
		(*i)++
		return args[(*i)], hasValue
	}
	return "", false
}

func applyArgValue(i *int, args []string, arg *string) {
	if value, ok := getArgValue(i, args); ok {
		*arg = value
	}
}

// Condi export condition
type Condi int

const (
	// GreaterThan 1
	GreaterThan = iota
	// LessThan 2
	LessThan
	// GreaterOrEqual 3
	GreaterOrEqual
	// LessOrEqual 4
	LessOrEqual
	// Equal 5
	Equal
)

// FilterCondi filter condition
type FilterCondi struct {
	IsQueryCount bool
	Condi        Condi
	QCValue      int
}

// Arguments argument object
type Arguments struct {
	IsExport    bool
	ExportLoc   string
	FilterCondi *FilterCondi
	ExcludeZh   bool
}

// HandleExport handle export
func HandleExport(storage interface{}, argument *Arguments) {
	strategy := &ExportStrategy{}
	strategy.Filter = func(word *Word) bool {
		_, sz := utf8.DecodeRuneInString(word.Content)
		if argument.ExcludeZh && sz == 3 {
			return false
		}
		if !argument.FilterCondi.IsQueryCount {
			return true
		}
		condi := argument.FilterCondi.Condi
		qcValue := argument.FilterCondi.QCValue
		switch condi {
		case GreaterThan:
			return word.QueryCount > qcValue
		case GreaterOrEqual:
			return word.QueryCount >= qcValue
		case LessThan:
			return word.QueryCount < qcValue
		case LessOrEqual:
			return word.QueryCount <= qcValue
		case Equal:
			return word.QueryCount == qcValue
		}
		return true
	}
	strategy.AddExportTarget(&ExportTarget{
		ExportObject: ExportWord,
	})
	ExportToTxtFile(storage, strategy)
}

// ParseArgs parse arguments
func ParseArgs(args []string) *Arguments {
	arguments := &Arguments{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--export":
			arguments.IsExport = true
			applyArgValue(&i, args, &arguments.ExportLoc)
		case "--filter":
			if condi, ok := getArgValue(&i, args); ok {
				arguments.FilterCondi = parseCondition(condi)
			}
		case "--exclude-zh":
			arguments.ExcludeZh = true
		}
	}
	return arguments
}

func parseCondition(condiStr string) *FilterCondi {
	filterCondi := &FilterCondi{}
	var condi Condi
	if strings.HasPrefix(condiStr, "q") {
		filterCondi.IsQueryCount = true
		c := condiStr[1]
		qc := 0
		var err error
		if c == '>' {
			if condiStr[2] == '=' {
				condi = GreaterOrEqual
				qc, err = strconv.Atoi(condiStr[3:])
			} else {
				condi = GreaterThan
				qc, err = strconv.Atoi(condiStr[2:])
			}
		} else if c == '<' {
			if condiStr[2] == '=' {
				condi = LessOrEqual
				qc, err = strconv.Atoi(condiStr[3:])
			} else {
				condi = LessThan
				qc, err = strconv.Atoi(condiStr[2:])
			}
		} else if c == '=' {
			qc, err = strconv.Atoi(condiStr[2:])
		}
		if err != nil {
			log.Fatal(err)
		}
		filterCondi.Condi = condi
		filterCondi.QCValue = qc
	}
	return filterCondi
}

// Init check install dir,log and so on.
func Init() {
	checkInstallDir()
	logFile, err := os.OpenFile("/home/zgq/app/go-word/run.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	Init()
	storage := NewFileStorage()
	if len(os.Args) == 1 {
		fmt.Println("Please input word!")
		return
	}
	arguments := ParseArgs(os.Args[1:])
	if arguments.IsExport {
		HandleExport(storage, arguments)
	}
	queryWord := strings.Join(os.Args[1:], " ")
	fmt.Println(queryWord + " ~ fanyi.youdao.com")
	var word *Word
	var ok bool
	if word, ok = storage.QueryWord(queryWord); ok == false {
		youdao := NewYoudaoImpl("1879868570", "youdao-cli")
		result := youdao.Query(queryWord)
		word = &Word{
			Content:           queryWord,
			TranslatedContent: result.Translation,
			Explains:          result.Explains,
			QueryCount:        1,
		}
		if result.Result == TranslSucc {
			storage.AddWord(word)
		}
	}

	fmt.Println("translation:" + word.TranslatedContent)
	fmt.Println("explains:", word.Explains)
	fmt.Println("has been queried " + strconv.Itoa(word.QueryCount) + " times")

}
