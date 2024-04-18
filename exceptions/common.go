package exceptions

type ErrorWithStatus struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
