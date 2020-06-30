package providers

type (
	// TODO: Add page path
	Pagination struct {
		PreviousPage int `json:"previous_page"`
		CurrentPage  int `json:"current_page"`
		NextPage     int `json:"next_page"`
	}
)
