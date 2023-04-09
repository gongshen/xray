package response

type StatChartResponse struct {
	DataAxis []int64 `json:"data_axis"`
	Data     []int64 `json:"data"`
}
