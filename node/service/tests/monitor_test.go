package service

import (
	"context"
	"go-node/models"
	"go-node/service"
	"log"
	"testing"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func TestGopsutilsUsage(t *testing.T) {

	a, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range a {
		log.Println("cpu", k, ":", v, "%")
	}

	cpuTotal, err := cpu.Times(true)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range cpuTotal {
		log.Println(v.CPU, ":", v)
	}

	g, err := cpu.Info()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range g {
		log.Println(v)
	}
}

func TestMem(t *testing.T) {
	for i := 0; i < 30; i++ {
		vm, err := mem.VirtualMemory()
		if err != nil {
			t.Fatal(err)
		}
		cpu.Info()
		ci, err := cpu.Percent(500*time.Millisecond, true)
		if err != nil {
			t.Fatal(err)
		}
		var count float64
		for _, v := range ci {
			count += v
		}

		log.Printf("cpu:%.4f count:%v", count/12.0, count)
		var mb uint64 = 1024 * 1024
		var gb uint64 = mb * 1024
		log.Printf("mem:%.4f", (float64(vm.Available/gb) / float64(vm.Total/gb)))

		is, err := net.IOCounters(false)
		if err != nil {
			t.Fatal(err)
		}
		in := is[0]

		log.Print("bytes:", is[0].BytesRecv/mb, is[0].BytesSent/mb)
		log.Print("packet:", in.PacketsRecv/mb, in.PacketsSent/mb)

		// time.Sleep(1 * time.Second)
	}
}

func TestMemUsage(t *testing.T) {
	var a []models.MonitorPacket
	for i := 0; i < 1600; i++ {
		a = append(a, models.MonitorPacket{})
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(15 * time.Second)
		t.Log("解放内存")
		a = a[:0]
		cancel()
	}()
	t.Log(len(a))

	select {
	case <-ctx.Done():
	}

	select {}

}

func TestNet(t *testing.T) {
	i, _ := net.IOCounters(false)
	t.Log(len(i))
	i, _ = net.IOCounters(true)
	t.Log(len(i))
}

func TestMonitorPacket(t *testing.T) {
	m := models.MonitorPacket{}
	o, _ := m.Marshal()
	t.Logf("%s", o)

}

func TestMonitor(t *testing.T) {
	nm := service.NewNodeMonitor(1600, 1*time.Second)
	nm.Send(func(mp *models.MonitorPacket) {
		log.Println(mp)
	})
}
