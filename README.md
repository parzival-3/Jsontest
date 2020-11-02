#Jsontest
利用dnslog批量测试fastjson漏洞。
# 用法
./Jsontest  -h
Usage: Jsontest [-r] x.txt [-d] x.dnslog.com [-ssl]
  -d    DNSlog地址
  -r    数据包文件
  -ssl
        是否HTTPS

  ./JsonDataTest -r 1.txt -d baidu.com -ssl
    在burp中发现使用json格式的数据包，保存为1.txt,使用$exp$替换掉对应json数据，替换方式如下：

## POST型请求

```
POST /exploit HTTP/1.1
Host: 192.168.0.2:8080
Pragma: no-cache
Cache-Control: no-cache
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Connection: close
Content-Type: application/json
Content-Length: 5

$exp$
```

## GET型请求：

```
GET /exploit?data=$exp$ HTTP/1.1
Host: 192.168.0.2:8080
Pragma: no-cache
Cache-Control: no-cache
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Accept-Encoding: gzip, deflate
Accept-Language: zh-CN,zh;q=0.9
Connection: close
Content-Type: application/json
Content-Length: 5
```

