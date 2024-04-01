package models

import (
	"encoding/json"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type MonitorPacket struct {
	Memory   *mem.VirtualMemoryStat `json:"memory"`
	Cpu      []cpu.InfoStat         `json:"cpuInfo"`
	CpuUsage float64                `json:"cpuUsage"`
	Net      net.IOCountersStat     `json:"netCounter"`
}

func (m *MonitorPacket) Marshal() ([]byte, error) {
	return json.Marshal(m)
}
