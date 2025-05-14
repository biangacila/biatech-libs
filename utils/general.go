package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/biangacila/luvungula-go/global"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type XlsxHeader struct {
	Position      int
	Name          string
	DataType      string
	DateFormatIn  string
	DateFormatOut string
}

func CompareStringSlices(a, b []string) (missingInA, missingInB []string) {
	mapA := make(map[string]bool)
	mapB := make(map[string]bool)

	for _, item := range a {
		mapA[item] = true
	}
	for _, item := range b {
		mapB[item] = true
	}

	// Find items in B not in A
	for _, item := range b {
		if !mapA[item] {
			missingInA = append(missingInA, item)
		}
	}

	// Find items in A not in B
	for _, item := range a {
		if !mapB[item] {
			missingInB = append(missingInB, item)
		}
	}

	return
}
func IsValidDateWithFormat(input, format string) bool {
	_, err := time.Parse(format, input)
	return err == nil
}

func GetPeriodStartEndDate(period string) (startDate, endDate time.Time, err error) {
	// Parse the period like "202503"
	t, err := time.Parse("200601", period)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Start date is the first of the month
	startDate = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

	// End date is the last day of the month (start of next month minus 1 day)
	endDate = startDate.AddDate(0, 1, -1)

	return startDate, endDate, nil
}

func ExcelReaderWithContent(headerRowStart, contentRowStart int, sheetName string) ([]byte, error) {
	// todo open file
	f, err := excelize.OpenFile(sheetName)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	rows := f.GetRows(sheetName)
	var headers []string
	var headerMapKey = make(map[string]int)
	if len(rows) >= headerRowStart {
		headers = rows[headerRowStart]
		headerMapKey = ArrayToFileHeader(headers)
	}
	//global.DisplayObject("Headers", headers)
	//global.DisplayObject("headerMapKey", headerMapKey)
	contents := GetExcelContentBasedOnHeaders(rows, headerMapKey, contentRowStart)
	var contentRows []interface{}
	_ = json.Unmarshal(contents, &contentRows)
	//global.DisplayObject("ContentRows", contentRows[1])
	return contents, nil
}
func GetExcelContentBasedOnHeaders(data [][]string, headers map[string]int, startRow int) []byte {
	var maps []map[string]interface{}
	for index, rows := range data {
		if index < startRow {
			continue
		}
		var rec = make(map[string]interface{})
		for hName, hIndex := range headers {
			value := rows[hIndex]
			rec[hName] = value
		}
		maps = append(maps, rec)
	}
	str, _ := json.Marshal(maps)
	return str
}
func ArrayToFileHeader(arr []string) map[string]int {
	maps := make(map[string]int)
	for index, x := range arr {
		maps[x] = index
	}
	return maps
}
func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}
func CleanStringToCompare(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.ToLower(s)
	return s
}
func ConvertStringToDateTime2(dateIn string, timeIn string) time.Time {
	arr := strings.Split(dateIn, "-")
	year, _ := strconv.Atoi(arr[0])
	month, _ := strconv.Atoi(arr[1])
	day, _ := strconv.Atoi(arr[2])

	arr2 := strings.Split(timeIn, ":")
	hour, _ := strconv.Atoi(arr2[0])
	err, _ := strconv.Atoi(arr2[1])
	sec, _ := strconv.Atoi(arr2[2])

	date := time.Date(year, time.Month(month), day, hour, err, sec, 0, time.UTC)
	return date
}
func CondString(key string, value string) map[string]interface{} {
	cond := make(map[string]interface{})
	cond[key] = value
	return cond
}
func ByteToAny[T any](b []byte, t T) T {
	var newObj T
	_ = json.Unmarshal(b, &newObj)
	return newObj
}
func IsInArray(target string, arr []string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
func ConvertExcelToBase64(fName string) string {
	// Open the Excel file
	filePath := fName // Replace it with your file path
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content into memory
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Encode file content to Base64
	base64Encoded := base64.StdEncoding.EncodeToString(fileContent)

	// Print or use the base64 encoded string
	fmt.Println(base64Encoded)
	return base64Encoded
}
func GetExpiredAt(hour int64) time.Time {
	expiresAt := time.Now().Add(time.Hour * time.Duration(hour)).Unix()
	return time.Unix(expiresAt, 0)
}
func ExtractQueryParams(r *http.Request) (map[string]interface{}, error) {
	// Create a map to store the query parameters
	queryParams := make(map[string]interface{})

	// Get all query parameters from the request URL
	values := r.URL.Query()

	// Loop through all keys and their values
	for key, valueSlice := range values {
		if len(valueSlice) > 1 {
			// If there are multiple values for a key, store them as a slice in the map
			queryParams[key] = valueSlice
		} else {
			// Otherwise, store the single value as a string in the map
			queryParams[key] = valueSlice[0]
		}
	}

	// For debugging: Print the query parameters
	for key, value := range queryParams {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}

	return queryParams, nil
}
func ObjectToBufferReader(input interface{}) (*bytes.Buffer, error) {
	// Encode the parameters as JSON
	jsonData, err := json.Marshal(input)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return nil, err
	}
	return bytes.NewBuffer(jsonData), nil
}
func ObjectToMap(o interface{}) map[string]interface{} {
	var mapRecord map[string]interface{}
	str, _ := json.Marshal(o)
	_ = json.Unmarshal(str, &mapRecord)
	return mapRecord
}
func GenerateCodeBasedOnCurrentDateTime(prefix string) string {
	return strings.ReplaceAll(prefix+time.Now().Format("20060102150405.000"), ".", "")
}
func GenerateCodeBasedOnTimestamp(prefix string) string {
	return fmt.Sprintf("%v%v", prefix, time.Now().Unix())
}
func HttpResponseError(err error) string {
	var maps = map[string]interface{}{
		"error": err.Error(),
	}
	return MapToString(maps)
}
func HttpResponseErrors(errs []error) string {
	var maps []map[string]interface{}
	for _, err := range errs {
		maps = append(maps, map[string]interface{}{
			"error": err.Error(),
		})
	}
	b, _ := json.Marshal(maps)

	return string(b)
}
func MapToString(input map[string]interface{}) string {
	b, _ := json.Marshal(input)
	return string(b)
}
func FormatDateCsvBank1(input string) string {
	// Extract the date part (ignoring the first 8 characters)
	datePart := input[8:15] // "240819"
	// Rearrange the extracted date part into YYYY-MM-DD format
	formattedDate := "20" + datePart[0:2] + "-" + datePart[2:4] + "-" + datePart[4:6]
	return formattedDate
}
func GetCsvColData(filePath string, row, col int) (string, error) {
	lines, _ := global.GetCsvFileContentUrl(filePath)
	if len(lines) == 0 {
		return "", errors.New("file empty")
	}
	if len(lines[row])-1 < row {
		return "", errors.New("not enough lines")
	}
	line := lines[row]
	data := strings.Split(line, ",")
	if len(data) < col {
		return "", errors.New("column not match")
	}
	val := data[col]
	return val, nil
}
func CsvReader(filePath string, startRow int, headerInfo []XlsxHeader) ([]byte, error) {

	lines, err := global.GetCsvFileContentUrl(filePath)
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", filePath, err)
		return nil, err
	}
	highestPosition := GetHighestCol(headerInfo)
	var records []map[string]interface{}
	for index, line := range lines {
		if index < startRow {
			continue
		}
		row := strings.Split(line, ",")
		if len(row)-1 < highestPosition {
			continue
		}
		var rec = make(map[string]interface{})
		for _, col := range headerInfo {
			var value any
			value = row[col.Position]
			if col.DataType == "float64" {
				value, err = StringToFloat64(ToString(value))
			}
			if col.DataType == "date" {
				value = DateConvertor(ToString(value), col.DateFormatIn, col.DateFormatOut)
			}
			rec[col.Name] = value
		}

		records = append(records, rec)
	}
	b, err := json.Marshal(records)
	return b, err
}
func IsValidDate(dateStr string) bool {
	// Define the layout for the expected date format (YYYY-MM-DD)
	layout := "2006-01-02"

	// Try to parse the string according to the layout
	_, err := time.Parse(layout, dateStr)

	// If parsing returns an error, the date is not in the correct format
	return err == nil
}
func XlsxReader(filePath string, sheetName string, startRow int, headerInfo []XlsxHeader) ([]byte, error) {
	// todo open file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	rows := f.GetRows(sheetName)
	highestPosition := GetHighestCol(headerInfo)
	var records []map[string]interface{}
	for index, row := range rows {
		if index < startRow {
			continue
		}
		if len(row)-1 < highestPosition {
			continue
		}
		//
		var rec = make(map[string]interface{})
		for _, col := range headerInfo {
			var value any
			value = row[col.Position]
			if col.DataType == "float64" {
				value, err = StringToFloat64(ToString(value))
			}
			if col.DataType == "date" {
				value = DateConvertor(ToString(value), col.DateFormatIn, col.DateFormatOut)
			}
			rec[col.Name] = value
		}

		records = append(records, rec)
	}
	b, err := json.Marshal(records)
	return b, err
}
func GetHighestCol(headers []XlsxHeader) int {
	var arr []int
	for _, header := range headers {
		arr = append(arr, header.Position)
	}
	sort.Ints(arr)
	h := arr[len(arr)-1]
	return h
}
func GetCurrentDateTimeString() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	return formattedTime
}
func FilterData[T any](data []T, conditions map[string]interface{}, obj T) []T {
	var infos = make([]map[string]interface{}, len(data))
	var outs []T
	b, _ := json.Marshal(data)
	_ = json.Unmarshal(b, &infos)
	for _, info := range data {
		var infoMap = make(map[string]interface{})
		b, _ := json.Marshal(info)
		_ = json.Unmarshal(b, &infoMap)

		isFind := true
		for k, v := range conditions {
			if infoMap[k] != v {
				isFind = false
			}
		}
		if isFind {
			outs = append(outs, info)
		}
	}
	return outs
}
func StringToInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
}
func DateConvertor(dateStr string, formatIn, formatOut string) string {
	// Remove the quotes from the string
	dateStr = strings.Trim(dateStr, "\"")
	// Parse the date string into a time.Time object
	parsedDate, err := time.Parse(formatIn, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err, " | ", dateStr)
		return dateStr
	}
	// Format the parsed date into the desired format
	formattedDate := parsedDate.Format(formatOut)
	return formattedDate
}
func ToString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
func StringToFloat64(str string) (float64, error) {
	str, err := cleanString(str)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
func StringToFloat64_2(str string) float64 {
	str, err := cleanString(str)
	if err != nil {
		return 0
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return val
}

func cleanString(input string) (string, error) {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove commas
	input = strings.ReplaceAll(input, ",", "")

	// Remove any other non-numeric characters except '.' and '-'
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' || r == '-' {
			return r
		}
		return -1
	}, input)

	return cleaned, nil
}
