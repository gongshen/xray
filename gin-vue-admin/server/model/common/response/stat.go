package response

type StatChartResponse struct {
	DataAxis []int   `json:"data_axis"`
	Data     []int64 `json:"data"`
}
