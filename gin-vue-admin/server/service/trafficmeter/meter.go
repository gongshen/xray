package trafficmeter

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	DefaultStatURL       = "http://127.0.0.1:56611"
	DefaultTag           = "1"
	DefaultFlushInterval = 10 * time.Second
)

type Meter struct {
	mu            sync.Mutex
	up            uint64
	down          uint64
	statURL       string
	tag           string
	flushInterval time.Duration
	client        *fasthttp.Client
	logger        *zap.Logger
	stop          chan struct{}
	startOnce     sync.Once
	stopOnce      sync.Once
}

type trafficEvent struct {
	Tag         string `json:"tag"`
	Down        uint64 `json:"down"`
	Up          uint64 `json:"up"`
	CollectedAt int64  `json:"collected_at"`
}

type trafficEventBatch struct {
	Events []trafficEvent `json:"events"`
}

func NewFromConfig(cfg config.TrafficMeter, client *fasthttp.Client, logger *zap.Logger) *Meter {
	if !cfg.Enable {
		return nil
	}
	statURL := strings.TrimRight(strings.TrimSpace(cfg.StatURL), "/")
	if statURL == "" {
		statURL = DefaultStatURL
	}
	tag := strings.TrimSpace(cfg.Tag)
	if tag == "" {
		tag = DefaultTag
	}
	flushInterval := DefaultFlushInterval
	if cfg.FlushInterval != "" {
		parsed, err := time.ParseDuration(strings.TrimSpace(cfg.FlushInterval))
		if err != nil || parsed <= 0 {
			if logger != nil {
				logger.Warn("invalid traffic-meter flush-interval, use default",
					zap.String("flush_interval", cfg.FlushInterval),
					zap.Duration("default", DefaultFlushInterval),
				)
			}
		} else {
			flushInterval = parsed
		}
	}
	if client == nil {
		client = &fasthttp.Client{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
	}
	return &Meter{
		statURL:       statURL,
		tag:           tag,
		flushInterval: flushInterval,
		client:        client,
		logger:        logger,
		stop:          make(chan struct{}),
	}
}

func (m *Meter) Add(up uint64, down uint64) {
	if m == nil || (up == 0 && down == 0) {
		return
	}
	m.mu.Lock()
	m.up += up
	m.down += down
	m.mu.Unlock()
}

func (m *Meter) Start() {
	if m == nil {
		return
	}
	m.startOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(m.flushInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := m.Flush(); err != nil && m.logger != nil {
						m.logger.Warn("flush xray-admin traffic meter failed", zap.Error(err))
					}
				case <-m.stop:
					if err := m.Flush(); err != nil && m.logger != nil {
						m.logger.Warn("flush xray-admin traffic meter on stop failed", zap.Error(err))
					}
					return
				}
			}
		}()
	})
}

func (m *Meter) Stop() {
	if m == nil {
		return
	}
	m.stopOnce.Do(func() {
		close(m.stop)
	})
}

func (m *Meter) Flush() error {
	if m == nil {
		return nil
	}
	up, down := m.take()
	if up == 0 && down == 0 {
		return nil
	}
	err := m.postEvent(up, down, time.Now().Unix())
	if err != nil {
		m.Add(up, down)
		return err
	}
	return nil
}

func (m *Meter) take() (uint64, uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	up, down := m.up, m.down
	m.up, m.down = 0, 0
	return up, down
}

func (m *Meter) postEvent(up uint64, down uint64, collectedAt int64) error {
	body, err := json.Marshal(trafficEventBatch{
		Events: []trafficEvent{
			{
				Tag:         m.tag,
				Down:        down,
				Up:          up,
				CollectedAt: collectedAt,
			},
		},
	})
	if err != nil {
		return err
	}

	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetRequestURI(m.statURL + "/stat/traffic/event")
	req.SetBody(body)
	if err := m.client.Do(req, resp); err != nil {
		return err
	}
	if status := resp.StatusCode(); status < fasthttp.StatusOK || status >= fasthttp.StatusMultipleChoices {
		return fmt.Errorf("stat traffic event returned status %d: %s", status, string(resp.Body()))
	}
	return nil
}

func Middleware(m *Meter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m == nil {
			c.Next()
			return
		}
		up := estimateRequestBytes(c.Request)
		c.Next()
		down := estimateResponseBytes(c.Writer)
		m.Add(up, down)
	}
}
