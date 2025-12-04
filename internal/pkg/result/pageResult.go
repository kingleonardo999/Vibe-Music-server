package result

type PageResult[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}
