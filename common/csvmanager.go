package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)


type CsvTable struct {
	FileName string
	Records  []CsvRecord
}

type CsvRecord struct {
	Record map[string]string
}

func (c *CsvRecord) GetInt(field string) int64 {
	r, _ := strconv.ParseInt(c.Record[field], 10, 64)
	return r
}

func (c *CsvRecord) GetFloat(field string) float64 {
	r, _ := strconv.ParseFloat(c.Record[field], 64)
	return r
}

func (c *CsvRecord) GetString(field string) string {
	data, ok := c.Record[field]
	if ok {
		return data
	} else {
		return ""
	}
}

func LoadCsvCfg(filename string, withHead bool, ca rune) (*CsvTable, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()


	reader := csv.NewReader(file)
	if reader == nil {
		return nil, fmt.Errorf("NewReader return nil, file:", file)
	}
	reader.Comma = ca
	reader.LazyQuotes = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 0 {
		r := &CsvTable{
			FileName:filename,
		}
		return r, nil
	}

	colNum := len(records[0])
	recordNum := len(records)
	var allRecords []CsvRecord
	i := 0
	if withHead {
		i = 1
	}
	for ; i < recordNum; i++ {
		record := &CsvRecord{make(map[string]string)}
		for k := 0; k < colNum; k++ {
			record.Record[records[0][k]] = records[i][k]
		}
		allRecords = append(allRecords, *record)
	}
	var result = &CsvTable{
		filename,
		allRecords,
	}
	return result, nil
}

