package main

import (
	"csv-pushgateway/cmd/collector"
	"flag"
)


var (
	pushUrl  = flag.String("gatewayUrl", "http://pushgateway.alig.plt.babytree-inc.com/", "Push Gateway URL")
	csvPath    = flag.String("cvsFilePaht", "/Users/zhangling/go/src/csv-pushgateway/log", "Path of the csv file")
	jobName = flag.String("jobName", "nginx_metrics", "metrics job name")
	metricsPrefix = flag.String("metricsPrefix", "nginx_metrics_", "metrics name prefix")
)

func main() {
	flag.Parse()

	labelColumns := []string{"ng_request_domain", "ng_request_url_short", "dc", "internal", "app", "module", "function", "master"}
	collectColumns := []string{"le10ms","le25ms","le50ms","le70ms","le100ms","le200ms","le300ms","le400ms","le500ms","le700ms","ipv6_all","le1000ms","le1500ms","le2000ms","le3000ms","count_all","error_4xx","error_5xx","error_all","median_request_time","sum_body_bytes","p99_request_tim","p95_request_time","avg_request_time"}


	collector := collector.CsvCollector(*pushUrl, *jobName, *csvPath, *metricsPrefix, labelColumns, collectColumns)

	collector.Collector()
}
