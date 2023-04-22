package response

type StatChartResponse struct {
	DataAxis []int   `json:"data_axis"`
	Data     []int64 `json:"data"`
}

type StatRankResponse struct {
	Rank     []int64  `json:"rank"`
	RankAxis []string `json:"rank_axis"`
}
