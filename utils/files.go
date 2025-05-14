package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"os"
	"strconv"
	"strings"
	"time"
)

func CalculatePercent(amount, percent float64) float64 {
	return amount * percent / 100
}
func getLastDayOfMonth(year, month int) time.Time {
	// Get the first day of the next month
	firstDayNextMonth := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)
	// Subtract one day to get the last day of the current month
	return firstDayNextMonth.AddDate(0, 0, -1)
}
func ExtractYearAndMonth(dateIn string) (year, month int) {
	sYear, _ := ConvertDateStringIntoAnotherLayout(dateIn, "200601", "2006")
	sMonth, _ := ConvertDateStringIntoAnotherLayout(dateIn, "200601", "01")
	year, _ = strconv.Atoi(sYear)
	month, _ = strconv.Atoi(sMonth)
	return year, month
}
func CalculateTenure(start, end, resultFormat string) (float64, error) {
	// Start date: "YYYY-MM-DD"
	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return 0, fmt.Errorf("error parsing start date: %s > %v", start, err)
	}

	// End date: "YYYYMM"
	endYear, endMonth := ExtractYearAndMonth(end)

	// Get the last day of the month
	endDate := getLastDayOfMonth(endYear, endMonth)

	// Calculate tenure
	years := endDate.Year() - startDate.Year()
	months := int(endDate.Month()) - int(startDate.Month())
	days := endDate.Day() - startDate.Day()

	// Adjust if negative days
	if days < 0 {
		prevMonthLastDay := getLastDayOfMonth(endDate.Year(), int(endDate.Month()-1))
		days += prevMonthLastDay.Day()
		months -= 1
	}

	// Adjust if negative months
	if months < 0 {
		months += 12
		years -= 1
	}

	// Convert tenure based on requested format
	switch resultFormat {
	case "years":
		return float64(years), nil
	case "months":
		totalMonths := float64(years*12 + months)
		return totalMonths, nil
	case "days":
		duration := endDate.Sub(startDate).Hours() / 24
		return duration, nil
	default:
		return 0, fmt.Errorf("invalid result format: %s (valid options: 'years', 'months', 'days')", resultFormat)
	}
}
func ConvertExcelDateToNormal(dateIn string) string {
	excelDate, err := strconv.Atoi(dateIn) // 42736 // Example Excel date
	if err != nil {
		return dateIn
	}
	// Excel's base date (1900-01-01)
	baseDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

	// Adjust by subtracting 2 days
	goDate := baseDate.AddDate(0, 0, excelDate-2)

	// Format as YYYY-MM-DD
	formattedDate := goDate.Format("2006-01-02")
	//fmt.Println(formattedDate) // Output: 2016-12-21
	return formattedDate
}
func ConvertDateStringIntoAnotherLayout(dateStr, inputLayout, outputLayout string) (string, error) {
	// Define the layout of the input date
	//inputLayout := "02Jan2006"

	// Parse the input string
	t, err := time.Parse(inputLayout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return dateStr, err
	}

	// Format the date into "2006-01-02"
	//outputLayout := "2006-01-02"
	formattedDate := t.Format(outputLayout)

	return formattedDate, nil
}
func BuildCsvContent(arr []string, headers map[string]int) map[string]string {
	record := make(map[string]string)
	for item, position := range headers {
		value := strings.TrimSpace(arr[position])
		record[item] = value
	}
	return record
}
func BuildCsvHeader(arr []string) map[string]int {
	headers := make(map[string]int)
	for index, item := range arr {
		if item == "" {
			continue
		}
		item = strings.TrimSpace(item)
		item = strings.ToLower(item)
		headers[item] = index
	}
	return headers
}
func GetCurrentPeriod() string {
	date, _ := global.GetDateAndTimeString()
	return global.ChangeDateFormat(date, "200601")
}
func IsValidCsv(fileName string) error {
	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		//fmt.Println("Error opening CSV file:", err)
		return err
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from CSV
	_, err = reader.ReadAll()
	if err != nil {
		//fmt.Println("Error reading CSV:", err)
		return err
	}
	return nil
}
func WhichFileIsIt(fileName string) int {
	var cols []string
	lines, err := global.GetCsvFileContentUrl(fileName)
	if err != nil {
		return 0
	}
	for _, line := range lines {
		if line == "" {
			continue
		}
		var arr = strings.Split(line, ",")
		for _, col := range arr {
			col = strings.TrimSpace(col)
			if col == "" {
				continue
			}
			cols = append(cols, col)
		}
		break
	}
	return len(cols)
}
