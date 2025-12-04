package result

type Nil struct{}

// Result 统一返回结果
// T 泛型类型，Nil 时不返回 data 字段
type Result[T any] struct {
	Code    int    `json:"code"` // 业务逻辑: 0 表示成功，非 0 表示失败
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"` // data 可能为空
}

// Success 返回成功结果, data 为空
func Success[T any](msg string) Result[T] {
	return Result[T]{
		Code:    0,
		Message: msg,
		Data:    nil,
	}
}

// SuccessWithData 返回成功结果, 带 data
func SuccessWithData[T any](msg string, data T) Result[T] {
	return Result[T]{
		Code:    0,
		Message: msg,
		Data:    &data,
	}
}

// Error 返回失败结果, message 为默认值 "操作失败"
func Error[T any](msg string) Result[T] {
	return Result[T]{
		Code:    1,
		Message: msg,
		Data:    nil,
	}
}
