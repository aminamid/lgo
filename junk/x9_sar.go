package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func main() {
	file, err := os.Create("x.sb")
	if err != nil {
		fmt.Printf("failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	for i := 0; i < 6; i++ {
		cpuPerc, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Printf("failed to get CPU info: %v\n", err)
			return
		}

		vmStat, err := mem.VirtualMemory()
		if err != nil {
			fmt.Printf("failed to get memory info: %v\n", err)
			return
		}

		diskStat, err := disk.Usage("/")
		if err != nil {
			fmt.Printf("failed to get disk info: %v\n", err)
			return
		}

		hostStat, err := host.Info()
		if err != nil {
			fmt.Printf("failed to get host info: %v\n", err)
			return
		}

		netStat, err := net.IOCounters(false)
		if err != nil {
			fmt.Printf("failed to get network info: %v\n", err)
			return
		}

		if err := binary.Write(file, binary.LittleEndian, cpuPerc); err != nil {
			fmt.Printf("failed to write CPU data: %v\n", err)
			return
		}
		if err := binary.Write(file, binary.LittleEndian, vmStat); err != nil {
			fmt.Printf("failed to write memory data: %v\n", err)
			return
		}
		if err := binary.Write(file, binary.LittleEndian, diskStat); err != nil {
			fmt.Printf("failed to write disk data: %v\n", err)
			return
		}
		if err := binary.Write(file, binary.LittleEndian, hostStat); err != nil {
			fmt.Printf("failed to write host data: %v\n", err)
			return
		}
		if err := binary.Write(file, binary.LittleEndian, netStat); err != nil {
			fmt.Printf("failed to write network data: %v\n", err)
			return
		}

		time.Sleep(10 * time.Second)
	}

	fmt.Println("data has been written to x.sb")
}
