package zoomeye

type Options struct {
	Search string
	Page   int
}

type (
	Result struct {
		Total   int                      `json:"total"`
		Results []map[string]interface{} `json:"matches"`
	}
)
