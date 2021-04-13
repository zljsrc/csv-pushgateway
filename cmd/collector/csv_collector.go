package collector

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"strconv"
	"strings"
)

type csvCollector struct {
	pushUrl string
	jobName string
	csvFiePath string
	metricsPrefix string
	labelCollumns []string
	collectColumns []string
}

func CsvCollector(pushUrl string, jobName string, csvFiePath string, metricsPrefix string, labelCollumns []string, collectColumns []string) *csvCollector {
	return &csvCollector{
		pushUrl: pushUrl,
		jobName: jobName,
		csvFiePath: csvFiePath,
		metricsPrefix: metricsPrefix,
		labelCollumns: labelCollumns,
		collectColumns: collectColumns,
	}
}

func (collector *csvCollector) Collector() error {
	columnNames, csvContents, err := collector.readCsvFile()
	if err != nil {
		return err
	}

	err = collector.collectCsvContent(columnNames, csvContents)

	return err
}

func (collector *csvCollector) collectCsvContent(columnNames []string, csvContents [][]string) error {

	for i:=0; i<len(csvContents); i++ {
		row := csvContents[i]
		pusher := push.New(collector.pushUrl, collector.jobName)

		for j:=0; j<len(collector.labelCollumns); j++ {
			exists, index := in_array(collector.labelCollumns[j], columnNames)
			if exists {
				pusher.Grouping(collector.labelCollumns[j], row[index])
			}
		}

		for j:=0; j<len(columnNames); j++ {
			column := columnNames[j]

			//label的列
			exists, _ := in_array(column, collector.labelCollumns)
			if exists {
				continue
			}

			//该列是否收集
			exists, _ = in_array(column, collector.collectColumns)
			if exists {
				metrics_name := fmt.Sprintf("%v%v", collector.metricsPrefix, column)
				metics_value, err := strconv.ParseFloat(row[j], 64)
				if err != nil {
					metics_value = 0
				}

				gauge := prometheus.NewGauge(prometheus.GaugeOpts{
					Name: metrics_name,Help:metrics_name,
				})
				gauge.Set(metics_value)
				pusher.Collector(gauge)
			}

		}

		pusher.Push()
		log.Info("push: %v", row)
	}


	return nil
}

func (collector *csvCollector) readCsvFile() ([]string, [][]string, error) {
	cntb, err := ioutil.ReadFile(collector.csvFiePath)

	if err != nil {
		return nil, nil, err
	}

	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	rs, _ := r2.ReadAll()
	size := len(rs)

	if size < 2 {
		return nil, nil, errors.New("csv is emepty")
	}

	columnNames := rs[0]
	contents := rs[1:]
	csvContents := contents

	log.Info("csvFileLines:", len(rs))

	return columnNames, csvContents, nil
}

func in_array(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1;

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}

	return
}