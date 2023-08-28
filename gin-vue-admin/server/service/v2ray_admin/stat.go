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
	"sort"
	"strconv"
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
	db := global.GVA_DB.Debug().Model(&v2ray.Stat{})
	var stats []v2ray.Stat
	startCreatedAt := info.StartCreatedAt.Year()*10000 + int(info.StartCreatedAt.Month())*100 + info.StartCreatedAt.Day()
	endCreatedAt := info.EndCreatedAt.Year()*10000 + int(info.EndCreatedAt.Month())*100 + info.EndCreatedAt.Day()
	db = db.Where("created_at BETWEEN ? AND ?", startCreatedAt, endCreatedAt)
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
			Tag:       stat.Tag,
			Down:      format(stat.Down),
			Up:        format(stat.Up),
			Total:     format(stat.Down + stat.Up),
			CreatedAt: fmt.Sprintf("%d年%d月%d日", createdAt, month, day),
			ServerIp:  stat.ServerIp,
		})
	}
	// 补充用户名
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
	buf.WriteString("INSERT INTO v2ray_stat (tag,down,up,created_at,server_ip) VALUES ")
	args := make([]interface{}, 0, len(statsMap)*5)
	i := 0
	for _, stat := range statsMap {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("(?,?,?,?,?)")
		args = append(args, stat.Tag, stat.Down, stat.Up, stat.CreatedAt, stat.ServerIp)
		i++
	}
	buf.WriteString(" ON DUPLICATE KEY UPDATE down=down+VALUES(down),up=up+VALUES(up)")
	err = global.GVA_DB.Exec(buf.String(), args...).Error
	return
}

func (statService *StatService) GetStatCharts(info *v2rayReq.StatSearch) (*response.StatChartResponse, error) {
	resp := &response.StatChartResponse{}
	// 创建db
	db := global.GVA_DB.Debug().Model(&v2ray.Stat{}).Select("sum(down) as down,sum(up) as up,created_at")
	stats := make([]*v2ray.Stat, 0)
	startCreatedAt := info.StartCreatedAt.Year()*10000 + int(info.StartCreatedAt.Month())*100 + info.StartCreatedAt.Day()
	endCreatedAt := info.EndCreatedAt.Year()*10000 + int(info.EndCreatedAt.Month())*100 + info.EndCreatedAt.Day()
	db = db.Where("created_at BETWEEN ? AND ?", startCreatedAt, endCreatedAt)
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	if info.ServerIp != "" {
		db = db.Where("server_ip = ?", info.ServerIp)
	}
	if err := db.Group("created_at").Find(&stats).Error; err != nil {
		return nil, err
	}
	if len(stats) <= 0 {
		return resp, nil
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].CreatedAt < stats[j].CreatedAt
	})
	createdAt := make([]int, len(stats))
	createdAtCount := make([]int64, len(stats))
	for i, stat := range stats {
		createdAt[i] = stat.CreatedAt
		createdAtCount[i] = int64(stat.Down + stat.Up)
	}
	resp.Data = createdAtCount
	resp.DataAxis = createdAt
	return resp, nil
}

func (statService *StatService) GetStatRank(info *v2rayReq.StatSearch) (*response.StatRankResponse, error) {
	resp := &response.StatRankResponse{}
	// 创建db
	db := global.GVA_DB.Debug().Model(&v2ray.Stat{}).Select("sum(down) as down,sum(up) as up,tag")
	stats := make([]*v2ray.Stat, 0)
	startCreatedAt := info.StartCreatedAt.Year()*10000 + int(info.StartCreatedAt.Month())*100 + info.StartCreatedAt.Day()
	endCreatedAt := info.EndCreatedAt.Year()*10000 + int(info.EndCreatedAt.Month())*100 + info.EndCreatedAt.Day()
	db = db.Where("created_at BETWEEN ? AND ?", startCreatedAt, endCreatedAt)
	if info.Tag != "" {
		db = db.Where("tag = ?", info.Tag)
	}
	if info.ServerIp != "" {
		db = db.Where("server_ip = ?", info.ServerIp)
	}
	if err := db.Group("tag").Find(&stats).Error; err != nil {
		return nil, err
	}
	if len(stats) <= 0 {
		return resp, nil
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Down+stats[i].Up < stats[j].Down+stats[j].Up
	})
	tagCount := make([]int64, len(stats))
	tags := make([]string, len(stats))
	for i, stat := range stats {
		tags[i] = stat.Tag
		tagCount[i] = int64(stat.Down + stat.Up)
	}
	nickNameMap, err := findUserNameByIds(tags)
	if err != nil {
		return nil, err
	}
	for i, v := range tags {
		tags[i] = nickNameMap[v]
	}
	resp.RankAxis = tags
	resp.Rank = tagCount
	return resp, nil
}
