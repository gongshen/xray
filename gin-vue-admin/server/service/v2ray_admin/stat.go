package v2ray_admin

import (
	"bytes"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	v2rayReq "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/request"
	v2rayResp "github.com/flipped-aurora/gin-vue-admin/server/model/v2ray/response"
	"strconv"
	"time"
)

var (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
	KB = 1024
)

type StatService struct{}

// CreateStat 创建Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) CreateStat(stat *v2ray.Stat) (err error) {
	err = global.GVA_DB.Create(stat).Error
	return err
}

// DeleteStat 删除Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) DeleteStat(stat v2ray.Stat) (err error) {
	err = global.GVA_DB.Delete(&stat).Error
	return err
}

// DeleteStatByIds 批量删除Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) DeleteStatByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]v2ray.Stat{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateStat 更新Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) UpdateStat(stat v2ray.Stat) (err error) {
	err = global.GVA_DB.Save(&stat).Error
	return err
}

// GetStat 根据id获取Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) GetStat(id uint) (stat v2ray.Stat, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&stat).Error
	return
}

// GetStatInfoList 分页获取Stat记录
// Author [piexlmax](https://github.com/piexlmax)
func (statService *StatService) GetStatInfoList(info v2rayReq.StatSearch) (list []*v2rayResp.Stat, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&v2ray.Stat{})
	var stats []v2ray.Stat
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		startCreatedAt := info.StartCreatedAt.Year()*10000 + int(info.StartCreatedAt.Month())*100 + info.StartCreatedAt.Day()
		endCreatedAt := info.EndCreatedAt.Year()*10000 + int(info.EndCreatedAt.Month())*100 + info.EndCreatedAt.Day()
		db = db.Where("created_at BETWEEN ? AND ?", startCreatedAt, endCreatedAt)
	}
	if info.Category != "" {
		db = db.Where("category = ?", info.Category)
	}
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	if info.ServerIp != "" {
		db = db.Where("server_ip = ?", info.ServerIp)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&stats).Error

	ans := make([]*v2rayResp.Stat, 0, len(stats))
	for _, stat := range stats {
		if stat.Down+stat.Up <= 0 {
			continue
		}
		createdAt := stat.CreatedAt
		day := createdAt % 100
		createdAt /= 100
		month := createdAt % 100
		createdAt /= 100
		ans = append(ans, &v2rayResp.Stat{
			ID:        stat.ID,
			Category:  stat.Category,
			Tag:       stat.Tag,
			Down:      format(stat.Down),
			Up:        format(stat.Up),
			Total:     format(stat.Down + stat.Up),
			CreatedAt: fmt.Sprintf("%d年%d月%d日", createdAt, month, day),
			ServerIp:  stat.ServerIp,
		})
	}
	// 补充用户名
	if info.Category == "user" {
		nickNameMap := make(map[string]string)
		tagsMap := make(map[string]struct{})
		for _, stat := range stats {
			tagsMap[stat.Tag] = struct{}{}
		}
		tags := make([]string, 0, len(tagsMap))
		for tag := range tagsMap {
			tags = append(tags, tag)
		}
		nickNameMap, err = findUserNameByIds(tags)
		for _, v := range ans {
			v.Username = nickNameMap[v.Tag]
		}
	}
	if err != nil {
		return
	}
	return ans, total, err
}

func findUserNameByIds(ids []string) (map[string]string, error) {
	u := make([]*system.SysUser, 0)
	if err := global.GVA_DB.Select("id,nick_name").Where("`id` in (?)", ids).Find(&u).Error; err != nil {
		return nil, err
	}
	ans := make(map[string]string)
	for _, v := range u {
		ans[strconv.Itoa(int(v.ID))] = v.NickName
	}
	return ans, nil
}

func format(used uint64) string {
	if used < uint64(GB) {
		if used < uint64(MB) {
			if used < uint64(KB) {
				return fmt.Sprintf("%.2fB", float64(used))
			} else {
				return fmt.Sprintf("%.2fKB", float64(used)/float64(KB))
			}
		} else {
			return fmt.Sprintf("%.2fMB", float64(used)/float64(MB))
		}
	} else {
		return fmt.Sprintf("%.2fGB", float64(used)/float64(GB))
	}
}

func (statService *StatService) StatsCollector(statsMap map[string]*v2ray.Stat) (err error) {
	buf := bytes.Buffer{}
	buf.WriteString("INSERT INTO v2ray_stat (category,tag,down,up,created_at,server_ip) VALUES ")
	args := make([]interface{}, 0, len(statsMap)*6)
	i := 0
	for _, stat := range statsMap {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("(?,?,?,?,?,?)")
		args = append(args, stat.Category, stat.Tag, stat.Down, stat.Up, stat.CreatedAt, stat.ServerIp)
		i++
	}
	buf.WriteString(" ON DUPLICATE KEY UPDATE down=down+VALUES(down),up=up+VALUES(up)")
	err = global.GVA_DB.Exec(buf.String(), args...).Error
	return
}

func (statService *StatService) GetStatCharts(info *v2rayReq.StatSearch) (*response.StatChartResponse, error) {
	resp := &response.StatChartResponse{
		DataAxis: make([]int, 13),
		Data:     make([]int64, 13),
	}
	// 转换为东八时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(location)
	for i := 12; i >= 0; i-- {
		date := now.AddDate(0, -i, 0)
		month := date.Month()
		year := date.Year()
		key := year*100 + int(month)
		resp.DataAxis[12-i] = key
	}
	// 创建db
	db := global.GVA_DB.Debug().Model(&v2ray.Stat{}).Select("down,up,created_at")
	stats := make([]*v2ray.Stat, 0)
	// 如果有条件搜索 下方会自动创建搜索语句
	db = db.Where("created_at BETWEEN ? AND ?", (now.Year()-1)*10000+int(now.Month())*100+now.Day(), now.Year()*10000+int(now.Month())*100+now.Day())
	if info.Category != "" {
		db = db.Where("category = ?", info.Category)
	}
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	if info.ServerIp != "" {
		db = db.Where("server_ip = ?", info.ServerIp)
	}
	if err := db.Find(&stats).Error; err != nil {
		return nil, err
	}
	if len(stats) <= 0 {
		return resp, nil
	}
	monthCount := make(map[int]int64)
	for _, stat := range stats {
		monthCount[stat.CreatedAt/100] += int64(stat.Down + stat.Up)
	}
	for i, key := range resp.DataAxis {
		resp.Data[i] = monthCount[key]
	}
	return resp, nil
}