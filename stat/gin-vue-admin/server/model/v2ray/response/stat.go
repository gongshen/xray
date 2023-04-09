package response

type Stat struct {
	ID        uint
	Category  string `json:"category" form:"category"`
	Tag       string `json:"tag" form:"tag"`
	Down      string `json:"down" form:"down""`
	Up        string `json:"up" form:"up"`
	Total     string `json:"total" form:"total"`
	CreatedAt string `json:"created_at" form:"created_at"`
	Username  string `json:"username" form:"username"`
}
