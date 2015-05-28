// process.go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	//"strconv"
	"strings"
)

type Process struct {
	user string
	pid  string
	cpu  string
	mem  string
	cmd  string
	time string
}

func main() {
	cmd := exec.Command("ps", "-aux", "--sort=-%cpu")
	//cmd := exec.Command("top")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	processes := make([]*Process, 0)
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		//line = strings.Replace(line, "\t", "", -1)
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)

		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}

		//log.Println(len(ft), ft)
		user := ft[0]
		pid := ft[1]
		cpu := ft[2]
		mem := ft[3]
		time := ft[10]
		cmd := ft[9]

		processes = append(processes, &Process{user, pid, cpu, mem, time, cmd})
	}
	for i := 1; i < 11; i++ {
		//log.Println(processes[i])
		statPath := filepath.Join("/", "proc", processes[i].pid, "status")
		contents, err := ioutil.ReadFile(statPath)
		if err != nil {
			log.Println(err)
		}
		lines := strings.Split(string(contents), "\n")
		name := strings.Split(string(lines[0]), ":")

		//log.Println(name[1])
		//log.Println("User:", processes[i].user, "Pid:", processes[i].pid, " %CPU:", processes[i].cpu, "%Mem:", processes[i].mem, "Time:", processes[i].time, "Name:", name[1])
		process_info := fmt.Sprintln("User:", processes[i].user, "Pid:", processes[i].pid, " %CPU:", processes[i].cpu, "%Mem:", processes[i].mem, "Time:", processes[i].time, "Name:", name[1])
		process_info = strings.Replace(process_info, "\t", "", -1)
		log.Println(process_info)
	}
}
