//check_load
package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/shirou/gopsutil/load"
)

func StrMd5(_strmd5 string) string {
	h := md5.New()
	h.Write([]byte(_strmd5)) // 需要加密的字符串
	//return fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Post(url, json_str, private_key string) (recv string, err error) {

	var js map[string]interface{}
	json.Unmarshal([]byte(json_str), &js)
	dana_time := time.Now().Unix()

	str := fmt.Sprintf("%d%s", dana_time, private_key)
	signature := StrMd5(str)
	js["dana_time"] = dana_time
	js["signature"] = signature
	js_str, err := json.Marshal(js)

	req := httplib.Post(url).SetTimeout(5*time.Second, 5*time.Second)
	// DANA是通过HTTPS进行交互的所以需要开启TLS忽略Key的有效性核对.
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	req.Header("Content-Type", "application/json;charset=UTF-8")
	req.Body(string(js_str))

	recv, err = req.String()

	//
	return recv, err

}

func main() {
	//get options
	warning := flag.String("w", "", "warning value")
	critical := flag.String("c", "", "critical value")
	post_url := flag.String("p", "", "post url")
	private_key := flag.String("k", "", "private_key")
	flag.Parse()
	int_w, _ := strconv.ParseFloat(*warning, 64)
	int_c, _ := strconv.ParseFloat(*critical, 64)
	//fmt.Printf("post_url:%s\n", post_url)
	//fmt.Printf("warning:%v, critical:%v\n", int_w, int_c)

	//get cpu_info
	l, _ := load.LoadAvg()
	check_status := "OK"

	if l.Load1 > float64(int_c) || l.Load5 > float64(int_c) || l.Load15 > float64(int_c) {

		check_status = "CRITICAL"
	} else if l.Load1 > float64(int_w) || l.Load5 > float64(int_w) || l.Load15 > float64(int_w) {

		check_status = "WARNING"
	}
	//fmt.Printf("        Load                     : load1:%v, load5:%v,load15:%v\n", l.Load1, l.Load5, l.Load15)

	//create Json
	post_Json := make(map[string]interface{})
	r := make(map[string]interface{})
	post_Json["cmd"] = "CheckLoadReceiver"
	r["check_status"] = check_status
	r["load1"] = fmt.Sprintf("%v", l.Load1)
	r["load5"] = fmt.Sprintf("%v", l.Load5)
	r["load15"] = fmt.Sprintf("%v", l.Load15)
	post_Json["body"] = r

	b, _ := json.Marshal(post_Json)
	//fmt.Println(string(b))

	_, err := Post(*post_url, string(b), *private_key)
	if err != nil {
		fmt.Println(err)
	}
}
