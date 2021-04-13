build: cleanmod goget
	go build -tags=middleware --ldflags "-extldflags -static" -o csv-pushgateway cmd/main/main.go

# 自动提示
# 要在vim中自动提示，请先运行make autocompletor
autocompletor:
	go install ./...

# 清理mod缓存文件
cleanmod:
	go clean --modcache

goget:
	go env -w  GOPROXY=https://goproxy.cn,direct
	go env -w GO111MODULE=on
	go env | grep GOPROXY
	go get -u github.com/prometheus/client_golang@v1.9.0
	go mod tidy


# 开发环境运行
run: build
