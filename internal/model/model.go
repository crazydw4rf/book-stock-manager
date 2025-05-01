package model

type WebResponse[T any] struct {
	Data T `json:"data"`
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
