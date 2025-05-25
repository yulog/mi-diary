package shared

type QueryParams struct {
	Page  int    `url:"page"`
	S     string `url:"s,omitempty"`
	Color string `url:"color,omitempty"`
}
