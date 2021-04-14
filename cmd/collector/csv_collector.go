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
	"sync"
)

type cvsCollector struct {
	pushUrl string
	jobName string
	csvFiePath string
	metricsPrefix string
	labelCollumns []string
	collectColumns []string
}

func CsvCollector(pushUrl string, jobName string, csvFiePath string, metricsPrefix string, labelCollumns []string, collectColumns []string) *cvsCollector {
	return &cvsCollector{
		pushUrl: pushUrl,
		jobName: jobName,
		csvFiePath: csvFiePath,
		metricsPrefix: metricsPrefix,
		labelCollumns: labelCollumns,
		collectColumns: collectColumns,
	}
}

func (collector *cvsCollector) Collector() error {
	columnNames, csvContents, err := collector.readCvsFile()
	if err != nil {
		return err
	}

	err = collector.collectCsvContent(columnNames, csvContents)
	return err

	for start:=0; start<len(csvContents)-1; {
		end := start + 1000
		if end > len(csvContents)-1 {
			end = len(csvContents)-1
		}
		collectContent := csvContents[start:end]
		start = end

		err = collector.collectCsvContent(columnNames, collectContent)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func (collector *cvsCollector) collectCsvContent(columnNames []string, csvContents [][]string) error {

	pusher := push.New(collector.pushUrl, collector.jobName)

	for i:=0; i<len(csvContents); i++ {
		row := csvContents[i]

		lables := prometheus.Labels{}
		for j:=0; j<len(collector.labelCollumns); j++ {
			exists, index := in_array(collector.labelCollumns[j], columnNames)
			if exists {
				//pusher.Grouping(collector.labelCollumns[j], row[index])
				lables[collector.labelCollumns[j]] = row[index]
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
					Name: metrics_name,Help:metrics_name,ConstLabels:lables,
				})
				gauge.Set(metics_value)
				pusher.Collector(gauge)
				//log.Info("collect: ", column, " ", metics_value)
			} else {
				log.Info("ignore: ", column)
			}

		}

	}

	err := pusher.Push()

	return err
}


func (collector *cvsCollector) readCvsFile() ([]string, [][]string, error) {
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

func _push(pusher *push.Pusher, wg *sync.WaitGroup) {
	err := pusher.Push()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("done")
	}

	wg.Done()
	return
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