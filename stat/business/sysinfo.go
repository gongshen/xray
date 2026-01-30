package business

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/valyala/fasthttp"
)

// SysInfo 系统信息（精简版，节省流量）
type SysInfo struct {
	// 磁盘信息 (单位: MB)
	DiskTotal uint64 `json:"dt"`
	DiskUsed  uint64 `json:"du"`
	// 内存信息 (单位: MB)
	MemTotal uint64 `json:"mt"`
	MemUsed  uint64 `json:"mu"`
	// CPU使用率 (0-100)
	CPUPercent float64 `json:"cp"`
	// 时间戳
	Timestamp int64 `json:"ts"`
}

// CollectSysInfo 收集系统信息
func CollectSysInfo(reqCtx *fasthttp.RequestCtx) {
	info := SysInfo{
		Timestamp: time.Now().Unix(),
	}

	// 获取磁盘信息
	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:"
	}
	if diskStat, err := disk.Usage(diskPath); err == nil {
		info.DiskTotal = diskStat.Total / 1024 / 1024 // 转换为 MB
		info.DiskUsed = diskStat.Used / 1024 / 1024
	}

	// 获取内存信息
	if memStat, err := mem.VirtualMemory(); err == nil {
		info.MemTotal = memStat.Total / 1024 / 1024 // 转换为 MB
		info.MemUsed = memStat.Used / 1024 / 1024
	}

	// 获取 CPU 使用率 (采样 200ms)
	if cpuPercent, err := cpu.Percent(200*time.Millisecond, false); err == nil && len(cpuPercent) > 0 {
		info.CPUPercent = cpuPercent[0]
	}

	data, err := json.Marshal(info)
	if err != nil {
		reqCtx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	reqCtx.Success("application/json", data)
}
