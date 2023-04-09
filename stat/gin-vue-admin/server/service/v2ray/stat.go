package v2ray

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
	"sort"
	"strconv"
	"time"
)

var (
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
)

type StatService struct {
}

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
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt.Unix(), info.EndCreatedAt.Unix())
	}
	if info.Category != "" {
		db = db.Where("category = ?", info.Category)
	}
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&stats).Error

	ans := make([]*v2rayResp.Stat, 0, len(stats))
	for _, stat := range stats {
		ans = append(ans, &v2rayResp.Stat{
			ID:        stat.ID,
			Category:  stat.Category,
			Tag:       stat.Tag,
			Down:      format(stat.Down),
			Up:        format(stat.Up),
			Total:     format(stat.Down + stat.Up),
			CreatedAt: time.Unix(stat.CreatedAt, 0).Format("2006-01-02"),
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
		return fmt.Sprintf("%.2fMB", float64(used)/float64(MB))
	} else {
		return fmt.Sprintf("%.2fGB", float64(used)/float64(GB))
	}
}

func (statService *StatService) StatsCollector(statsMap map[string]*v2ray.Stat) (err error) {
	buf := bytes.Buffer{}
	buf.WriteString("INSERT INTO v2ray_stat (category,tag,down,up,created_at) VALUES ")
	args := make([]interface{}, 0, len(statsMap)*5)
	i := 0
	for _, stat := range statsMap {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("(?,?,?,?,?)")
		args = append(args, stat.Category, stat.Tag, stat.Down, stat.Up, stat.CreatedAt)
		i++
	}
	buf.WriteString(" ON DUPLICATE KEY UPDATE down=down+VALUES(down),up=up+VALUES(up)")
	err = global.GVA_DB.Exec(buf.String(), args...).Error
	return
}

func (statService *StatService) GetStatCharts(info *v2rayReq.StatSearch) (*response.StatChartResponse, error) {
	// 创建db
	db := global.GVA_DB.Model(&v2ray.Stat{}).Debug().Select("down,up,created_at")
	stats := make([]*v2ray.Stat, 0)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt.Unix(), info.EndCreatedAt.Unix())
	}
	if info.Category != "" {
		db = db.Where("category = ?", info.Category)
	}
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	if err := db.Find(&stats).Error; err != nil {
		return nil, err
	}
	if len(stats) <= 0 {
		return nil, nil
	}
	monthCount := make(map[int64]int64)
	for _, stat := range stats {
		date := time.Unix(stat.CreatedAt, 0)
		month := date.Month()
		year := date.Year()
		key := year*100 + int(month)
		monthCount[int64(key)] += int64(stat.Down + stat.Up)
	}
	var month []int64
	for key := range monthCount {
		month = append(month, key)
	}
	sort.Slice(month, func(i, j int) bool {
		if month[i]/100 > month[j]/100 {
			return true
		} else if month[i]%100 < month[j]%100 {
			return false
		}
		return false
	})
	data := make([]int64, len(month))
	for i, m := range month {
		data[i] = monthCount[m] / int64(MB)
	}
	return &response.StatChartResponse{Data: data, DataAxis: month}, nil
}
