package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func InsertBulkRecord[T any](session *gocql.Session, dbName, table string, records []T) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in InsertBulkRecord:", r)
		}
	}()

	var queries []string
	for _, record := range records {
		str, err := json.Marshal(record)
		if err != nil {
			continue
		}
		query := fmt.Sprintf("INSERT INTO %v.%v JSON '%v'", dbName, table, string(str))
		queries = append(queries, query)
	}

	const (
		goroutines = 8
		batchSize  = 50 // Reduce batch size to avoid "Batch too large" error
	)

	var wg sync.WaitGroup
	in := make(chan *gocql.Batch, goroutines)

	// Launch worker goroutines
	for i := 0; i < goroutines; i++ {
		go ProcessBatches(session, in, &wg)
	}

	counter := 0
	b := session.NewBatch(gocql.LoggedBatch)

	for i, qry := range queries {
		b.Query(qry)
		counter++

		// Send batch when limit is reached or at the last query
		if counter == batchSize || i == len(queries)-1 {
			wg.Add(1)
			in <- b
			b = session.NewBatch(gocql.LoggedBatch)
			counter = 0
		}
	}

	close(in)
	wg.Wait()

	return nil
}

// ProcessBatches executes batches concurrently
func ProcessBatches(session *gocql.Session, in chan *gocql.Batch, wg *sync.WaitGroup) {
	for batch := range in {
		if err := session.ExecuteBatch(batch); err != nil {
			log.Printf("Couldn't execute batch: %v\n", err)
		}
		wg.Done()
	}
}
func InsertRecord(session *gocql.Session, dbName, table string, record any) error {
	str, err := json.Marshal(record)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("insert into %v.%v  JSON '%v' ", dbName, table, string(str))
	return session.Query(query).Exec()
}
func FindRecord[T any](session *gocql.Session, query string, outs T) (T, error) {
	records, err := FetchData2(session, query, outs)
	if err != nil {
		fmt.Println("::error FetchData> ", err)
		return outs, err
	}
	var rec T
	for _, record := range records {
		rec = record
	}
	return rec, nil
}
func FetchRecord[T any](session *gocql.Session, query string, outs T) (T, error) {
	records, err := FetchData2(session, query, outs)
	if err != nil {
		fmt.Println("::error FetchData> ", err)
		return outs, err
	}
	if len(records) == 0 {
		return outs, fmt.Errorf("no records found")
	}
	var rec T
	for _, record := range records {
		rec = record
	}
	return rec, nil
}
func FetchData[T any](session *gocql.Session, query string, outs T) ([]T, error) {
	iter := session.Query(query).Iter()
	rows, err := iter.SliceMap()
	if err != nil {
		return []T{}, err
	}
	b, _ := json.Marshal(rows)
	var data []T
	_ = json.Unmarshal(b, &data)
	return data, nil
}
func FetchData2[T any](session *gocql.Session, query string, outs T) ([]T, error) {
	iter := session.Query(query).Iter()
	rows, err := iter.SliceMap()
	if err != nil {
		return []T{}, err
	}
	b, _ := json.Marshal(rows)
	var data []T
	_ = json.Unmarshal(b, &data)
	return data, nil
}
func FetchRecordWithConditions[T any](session *gocql.Session, dbName, table string, conditions map[string]interface{}, outs T, allowFiltering ...string) ([]T, error) {
	qry := fmt.Sprintf("SELECT * FROM %v.%s ", dbName, table)

	if conditions != nil && len(conditions) > 0 {
		qry = qry + " where "
		for k, v := range conditions {
			if _, ok := v.(string); ok {
				qry = qry + fmt.Sprintf(" %s ='%v' and", k, v)
			} else {
				qry = qry + fmt.Sprintf(" %s =%v and", k, v)
			}
		}
		qry = strings.TrimSuffix(qry, "and")
	}

	if len(allowFiltering) > 0 {
		qry = qry + allowFiltering[0]
	}

	data, err := FetchData2(session, qry, outs)
	return data, err
}
func ExecuteQuery(session *gocql.Session, query string) error {
	return session.Query(query).Exec()
}
func GenerateSequenceNumber(session *gocql.Session, dbName, table, fieldName, prefixCode string, prefixStart int) (nexCode string, err error) {
	query := fmt.Sprintf(`select %v as code from %v.%v  `, fieldName, dbName, table)
	iter := session.Query(query).Iter()
	rows, err := iter.SliceMap()
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(rows)
	type PrefixConf struct {
		Code string
	}

	var codes []PrefixConf
	_ = json.Unmarshal(b, &codes)

	var ls []int
	for _, code := range codes {
		arr := strings.Split(code.Code, prefixCode)
		value, _ := strconv.Atoi(arr[1])
		ls = append(ls, value)
	}
	sort.Ints(ls)

	var value = int64(prefixStart)
	if len(ls) > 0 {
		lastIndex := len(ls) - 1
		value = int64(ls[lastIndex])
	}
	nextValue := fmt.Sprintf("%v%v", prefixCode, value+1)
	return nextValue, nil
}
func WhereClauseBuilder(conditions map[string]interface{}) (string, error) {
	qry := ""
	if conditions == nil {
		return qry, errors.New("conditions is nil")
	}
	if len(conditions) == 0 {
		return qry, errors.New("conditions is empty")
	}
	qry = qry + " where "
	for k, v := range conditions {
		if _, ok := v.(string); ok {
			qry = qry + fmt.Sprintf(" %s ='%v' and", k, v)
		} else {
			qry = qry + fmt.Sprintf(" %s =%v and", k, v)
		}
	}
	qry = strings.TrimSuffix(qry, "and")
	return qry, nil
}
func UpdateClauseBuilder(conditions map[string]interface{}) (string, error) {
	if conditions == nil {
		return "", errors.New("update conditions is nil")
	}
	if len(conditions) == 0 {
		return "", errors.New("update conditions is empty")
	}

	var sb strings.Builder
	sb.WriteString(" ")

	for k, v := range conditions {
		switch val := v.(type) {

		case string:
			sb.WriteString(fmt.Sprintf("%s = '%s', ", k, val))

		case map[string]interface{}:
			// build a UDT literal: {field1:val1,field2:'val2',...}
			innerKeys := make([]string, 0, len(val))
			for ik := range val {
				innerKeys = append(innerKeys, ik)
			}
			sort.Strings(innerKeys)

			parts := make([]string, 0, len(innerKeys))
			for _, ik := range innerKeys {
				iv := val[ik]
				switch s := iv.(type) {
				case string:
					// wrap string in single quotes
					parts = append(parts, fmt.Sprintf("%s:'%s'", ik, s))
				default:
					// numbers, booleans, etc.
					parts = append(parts, fmt.Sprintf("%s:%v", ik, s))
				}
			}

			udt := "{" + strings.Join(parts, ",") + "}"
			sb.WriteString(fmt.Sprintf("%s = %s, ", k, udt))

		default:
			sb.WriteString(fmt.Sprintf("%s = %v, ", k, val))
		}
	}

	qry := strings.TrimSuffix(sb.String(), ", ")
	fmt.Println("Update query -:) ", qry)
	return qry, nil
}
