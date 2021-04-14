package main

import (
	"csv-pushgateway/cmd/collector"
	"flag"
	"github.com/prometheus/common/log"
)


var (
	pushUrl  = flag.String("gatewayUrl", "http://pushgateway.alig.plt.babytree-inc.com/", "Push Gateway URL")
	csvPath    = flag.String("cvsFilePaht", "/Users/zhangling/go/src/csv-pushgateway/log", "Path of the csv file")
	jobName = flag.String("jobName", "nginx_metrics", "metrics job name")
	metricsPrefix = flag.String("metricsPrefix", "nginx_metrics_", "metrics name prefix")
)

func main() {
	flag.Parse()

	labelColumns := []string{"ng_request_domain","ng_request_url_short","ng_ipnet","internal"}
	collectColumns := []string{"count_all","ipv6_all","error_all","error_4xx","error_5xx","le10ms","le25ms","le50ms","le70ms","le100ms","le200ms","le300ms","le400ms","le500ms","le700ms","le1000ms","le1500ms","le2000ms","le3000ms","avg_request_time","median_request_time","p95_request_time","p99_request_time","sum_body_bytes"}


	collector := collector.CsvCollector(*pushUrl, *jobName, *csvPath, *metricsPrefix, labelColumns, collectColumns)

	err := collector.Collector()
	if err != nil {
		log.Fatal(err)
	}
}
