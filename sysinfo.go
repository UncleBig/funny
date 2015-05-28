// sysinfo
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"time"
)

func main() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.CPUInfo()
	d, _ := disk.DiskUsage("/work")
	l, _ := load.LoadAvg()
	fmt.Printf("        Load                     : load1:%v, load5:%v,load15:%v\n", l.Load1, l.Load5, l.Load15)
	c_useage, _ := cpu.CPUPercent(1, false)
	netbefore, _ := net.NetIOCounters(false)
	time.Sleep(1 * time.Second)
	netLast, _ := net.NetIOCounters(false)
	fmt.Printf("        Net                       :  recv:%v M , sent:%v M\n", netbefore[0].BytesRecv/1024/1024, netbefore[0].BytesSent/1024/1024)
	fmt.Printf("        Net                       :  SentBytePersec: %v, RecvBytePersec:%v\n", netLast[0].BytesSent-netbefore[0].BytesSent, netLast[0].BytesRecv-netbefore[0].BytesRecv)
	fmt.Printf("        CPU_Useage        :   usepercent :%.2f%%\n", c_useage[0])
	n, _ := host.HostInfo()

	fmt.Printf("        Mem                     : %v GB  Free: %v MB Usage:%f%%\n", v.Total/1073741824, v.Free/1048576, v.UsedPercent)
	fmt.Printf("        Mem                     : %v  Free: %v  Usage:%.2f%%\n", v.Total, v.Free, v.UsedPercent)

	sub_cpu := c[0]
	modelname := sub_cpu.ModelName
	cores := sub_cpu.Cores
	fmt.Printf("        CPU_INFO           : %v   %v cores \n", modelname, cores)
	fmt.Printf("        HD                        : %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	fmt.Printf("        OS                        : %v   %v  \n", n.OS, n.PlatformVersion)
	fmt.Printf("        Hostname            : %v  \n", n.Hostname)

}
