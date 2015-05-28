//check_net

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
	"github.com/shirou/gopsutil/net"
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

	//get net_info
	netbefore, _ := net.NetIOCounters(false)
	time.Sleep(1 * time.Second)
	netLast, _ := net.NetIOCounters(false)
	fmt.Printf("        Net                       :  recv:%v M , sent:%v M\n", netbefore[0].BytesRecv/1024/1024, netbefore[0].BytesSent/1024/1024)
	fmt.Printf("        Net                       :  SentBytePersec: %v, RecvBytePersec:%v\n", netLast[0].BytesSent-netbefore[0].BytesSent,
		netLast[0].BytesRecv-netbefore[0].BytesRecv)
	recv_total := fmt.Sprintf("%v MB", netbefore[0].BytesRecv/1024/1024)
	sent_total := fmt.Sprintf("%v MB", netbefore[0].BytesSent/1024/1024)
	sent_byte_persec := fmt.Sprintf("%v", netLast[0].BytesSent-netbefore[0].BytesSent)
	recv_byte_persec := fmt.Sprintf("%v", netLast[0].BytesRecv-netbefore[0].BytesRecv)

	check_status := "OK"

	if netLast[0].BytesSent-netbefore[0].BytesSent > uint64(int_c) || netLast[0].BytesRecv-netbefore[0].BytesRecv > uint64(int_c) {

		check_status = "CRITICAL"
	} else if netLast[0].BytesSent-netbefore[0].BytesSent > uint64(int_w) || netLast[0].BytesRecv-netbefore[0].BytesRecv > uint64(int_w) {

		check_status = "WARNING"
	}

	//create Json
	post_Json := make(map[string]interface{})
	r := make(map[string]interface{})
	post_Json["cmd"] = "CheckNetReceiver"
	r["check_status"] = check_status
	r["recv_total"] = recv_total
	r["sent_total"] = sent_total
	r["sent_byte_persec"] = sent_byte_persec
	r["recv_byte_persec"] = recv_byte_persec
	post_Json["body"] = r

	b, _ := json.Marshal(post_Json)
	//fmt.Println(string(b))

	_, err := Post(*post_url, string(b), *private_key)
	if err != nil {
		fmt.Println(err)
	}
}
