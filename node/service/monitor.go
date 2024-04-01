package service

import (
	"go-node/models"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/sirupsen/logrus"
)

// 系统监控
type NodeMonitor interface {
	Send(func(*models.MonitorPacket))
	GetInfoPacket() (*models.MonitorPacket, error)
	HttpHandle(http.ResponseWriter, *http.Request)
}

func NewNodeMonitor(maxSave int, gap time.Duration) NodeMonitor {
	return &NodeMonitorImplV1{
		gap:     gap,
		maxSave: maxSave,
	}
}

type NodeMonitorImplV1 struct {
	gap     time.Duration
	maxSave int
}

func (n *NodeMonitorImplV1) Send(recv func(*models.MonitorPacket)) {
	count := 0
	for {
		select {
		case <-time.After(n.gap):
			p, err := n.GetInfoPacket()
			if err != nil {
				logrus.Error("get info packet", err)
				continue
			}

			recv(p)
			count++
		}
		if count == n.maxSave {
			// 保存到数据库中
			logrus.Print("save to db")
		}
	}
}

func (n *NodeMonitorImplV1) HttpHandle(w http.ResponseWriter, req *http.Request) {
	// 相当于调一次GetInfoPacket，用http返回
	m, err := n.GetInfoPacket()
	if err != nil {
		logrus.Error("http handle", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := m.Marshal()
	if err != nil {
		logrus.Error("json marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func (n *NodeMonitorImplV1) GetInfoPacket() (*models.MonitorPacket, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	ci, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		return nil, err
	}
	var cpuUsage float64
	for _, v := range ci {
		cpuUsage += v
	}
	cpuUsage = cpuUsage / float64(len(ci))
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	is, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	return &models.MonitorPacket{
		Memory:   vm,
		CpuUsage: cpuUsage,
		Cpu:      cpuInfo,
		Net:      is[0],
	}, nil
}
