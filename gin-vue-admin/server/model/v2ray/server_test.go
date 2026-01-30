package v2ray

import (
	"testing"
	"testing/quick"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TestGetStatPort_Property1 属性测试: GetStatPort 端口回退逻辑
// **Validates: Requirements 1.2, 7.2, 7.3**
func TestGetStatPort_Property1(t *testing.T) {
	// 设置全局配置的默认端口
	global.GVA_CONFIG.STAT_PORT = 56611

	// 属性: 当 StatPort > 0 时，返回 StatPort；否则返回全局配置
	property := func(statPort int) bool {
		server := &Server{StatPort: statPort}
		result := server.GetStatPort()

		if statPort > 0 {
			return result == statPort
		}
		return result == int(global.GVA_CONFIG.STAT_PORT)
	}

	if err := quick.Check(property, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 1 failed: %v", err)
	}
}

// TestGetStatPort_PositivePort 测试正数端口返回自身
func TestGetStatPort_PositivePort(t *testing.T) {
	global.GVA_CONFIG.STAT_PORT = 56611

	testCases := []struct {
		name     string
		statPort int
		expected int
	}{
		{"custom port 8080", 8080, 8080},
		{"custom port 56612", 56612, 56612},
		{"custom port 1", 1, 1},
		{"custom port 65535", 65535, 65535},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := &Server{StatPort: tc.statPort}
			if got := server.GetStatPort(); got != tc.expected {
				t.Errorf("GetStatPort() = %v, want %v", got, tc.expected)
			}
		})
	}
}

// TestGetStatPort_ZeroOrNegative 测试零或负数端口回退到全局配置
func TestGetStatPort_ZeroOrNegative(t *testing.T) {
	global.GVA_CONFIG.STAT_PORT = 56611

	testCases := []struct {
		name     string
		statPort int
	}{
		{"zero port", 0},
		{"negative port -1", -1},
		{"negative port -100", -100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := &Server{StatPort: tc.statPort}
			if got := server.GetStatPort(); got != int(global.GVA_CONFIG.STAT_PORT) {
				t.Errorf("GetStatPort() = %v, want %v (global config)", got, global.GVA_CONFIG.STAT_PORT)
			}
		})
	}
}
