package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var operationRecordService = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService

var respPool sync.Pool

func init() {
	respPool.New = func() interface{} {
		return make([]byte, 1024)
	}
}

var operationTargetKeys = map[string]struct{}{
	"ID":           {},
	"IDS":          {},
	"Ids":          {},
	"apiId":        {},
	"api_id":       {},
	"authorityId":  {},
	"authority_id": {},
	"email":        {},
	"id":           {},
	"ids":          {},
	"menuId":       {},
	"menu_id":      {},
	"name":         {},
	"nickName":     {},
	"path":         {},
	"serverID":     {},
	"serverId":     {},
	"server_id":    {},
	"target":       {},
	"targets":      {},
	"userID":       {},
	"userId":       {},
	"userName":     {},
	"user_id":      {},
	"username":     {},
	"uuid":         {},
	"UUID":         {},
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userId int
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.GVA_LOG.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		}
		claims, _ := utils.GetClaims(c)
		if claims.ID != 0 {
			userId = int(claims.ID)
		} else {
			id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
			if err != nil {
				userId = 0
			}
			userId = id
		}
		record := system.SysOperationRecord{
			Ip:      c.ClientIP(),
			Method:  c.Request.Method,
			Path:    c.Request.URL.Path,
			Agent:   c.Request.UserAgent(),
			Body:    string(body),
			Targets: extractOperationTargets(c.Request.URL.Query(), body),
			UserID:  userId,
		}

		// 上传文件时候 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			if len(record.Body) > 1024 {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, record.Body)
				record.Body = string(newBody)
				defer respPool.Put(newBody[:0])
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()

		if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
			strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
			strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
			strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
			strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {
			if len(record.Resp) > 1024 {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, record.Resp)
				record.Resp = string(newBody)
				defer respPool.Put(newBody[:0])
			}
		}

		if err := operationRecordService.CreateSysOperationRecord(record); err != nil {
			global.GVA_LOG.Error("create operation record error:", zap.Error(err))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func extractOperationTargets(query url.Values, body []byte) string {
	targets := make(map[string]interface{})

	for key, values := range query {
		if !isOperationTargetKey(key) {
			continue
		}
		if value, ok := operationQueryTargetValue(values); ok {
			targets[key] = value
		}
	}

	if len(body) > 0 {
		var payload map[string]interface{}
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.UseNumber()
		if err := decoder.Decode(&payload); err == nil {
			for key, value := range payload {
				if !isOperationTargetKey(key) || isEmptyOperationTargetValue(value) {
					continue
				}
				targets[key] = value
			}
		}
	}

	if len(targets) == 0 {
		return ""
	}
	data, err := json.Marshal(targets)
	if err != nil {
		return ""
	}
	return string(data)
}

func isOperationTargetKey(key string) bool {
	_, ok := operationTargetKeys[key]
	return ok
}

func operationQueryTargetValue(values []string) (interface{}, bool) {
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			cleaned = append(cleaned, value)
		}
	}
	if len(cleaned) == 0 {
		return nil, false
	}
	if len(cleaned) == 1 {
		return cleaned[0], true
	}
	return cleaned, true
}

func isEmptyOperationTargetValue(value interface{}) bool {
	switch v := value.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(v) == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		return false
	}
}
