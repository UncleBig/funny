//check_disk
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
	"github.com/shirou/gopsutil/disk"
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
	int_w, _ := strconv.Atoi(*warning)
	int_c, _ := strconv.Atoi(*critical)
	//fmt.Printf("post_url:%s\n", post_url)
	//fmt.Printf("warning:%d, critical:%d\n", int_w, int_c)

	//get cpu_info
	d, _ := disk.DiskUsage("/")
	check_status := "OK"

	if d.UsedPercent > float64(int_c) {

		check_status = "CRITICAL"
	} else if d.UsedPercent >= float64(int_w) {

		check_status = "WARNING"
	}
	use_percent := fmt.Sprintf("%.2f%%", d.UsedPercent)
	disk_total := fmt.Sprintf("%v GB", d.Total/1024/1024/1024)
	disk_free := fmt.Sprintf("%v GB", d.Free/1024/1024/1024)

	//create Json
	post_Json := make(map[string]interface{})
	r := make(map[string]interface{})
	post_Json["cmd"] = "CheckDiskReceiver"
	r["check_status"] = check_status
	r["use_percent"] = use_percent
	r["disk_total"] = disk_total
	r["disk_free"] = disk_free
	post_Json["body"] = r

	b, _ := json.Marshal(post_Json)
	fmt.Println(string(b))

	_, err := Post(*post_url, string(b), *private_key)
	if err != nil {
		fmt.Println(err)
	}
}
