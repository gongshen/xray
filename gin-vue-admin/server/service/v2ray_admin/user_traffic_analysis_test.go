package v2ray_admin

import (
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/stretchr/testify/require"
)

func TestValidateUserTrafficAnalysisProxyRequestAllowsTwoHours(t *testing.T) {
	req, err := ValidateUserTrafficAnalysisProxyRequest(UserTrafficAnalysisProxyRequest{
		ServerID: 1,
		UserTag:  "8",
		Date:     "20260617",
		Start:    "8:10",
		End:      "10:10",
	})

	require.NoError(t, err)
	require.Equal(t, "20260617", req.Date)
	require.Equal(t, "8:10", req.Start)
	require.Equal(t, "10:10", req.End)
}

func TestValidateUserTrafficAnalysisProxyRequestRejectsLongRange(t *testing.T) {
	_, err := ValidateUserTrafficAnalysisProxyRequest(UserTrafficAnalysisProxyRequest{
		ServerID: 1,
		UserTag:  "8",
		Date:     "20260617",
		Start:    "8:10",
		End:      "10:11",
	})

	require.Error(t, err)
}

func TestBuildUserTrafficAnalysisStatURL(t *testing.T) {
	req, err := ValidateUserTrafficAnalysisProxyRequest(UserTrafficAnalysisProxyRequest{
		ServerID: 1,
		UserTag:  "user+8",
		Date:     "20260617",
		Start:    "8:10",
		End:      "9:00",
	})
	require.NoError(t, err)

	got := BuildUserTrafficAnalysisStatURL(&v2ray.Server{Ip: "127.0.0.1", StatPort: 56611}, req)

	require.True(t, strings.HasPrefix(got, "http://127.0.0.1:56611/stat/traffic/user-minute?"))
	require.Contains(t, got, "user_tag=user%2B8")
	require.Contains(t, got, "date=20260617")
	require.Contains(t, got, "start=8%3A10")
	require.Contains(t, got, "end=9%3A00")
}
