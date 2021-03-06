package model

// Request structure for callback
type Request struct {
	SequenceID string `json:"sequenceId"`
	TimeStamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	MsgType    string `json:"msgType"`
	Status     string `json:"status"`
	FailedCode string `json:"failedCode,omitempty"`
}

// Response structure for callback
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

type RequestMessage struct {

	// An unique ID represent the upper stream system which is the API caller.<br/> E.g.: dkt_cube_cn / dkt_membership / dkt_newton
	SystemId string `json:"system_id"`

	// An unique ID generated by the upper stream system for tracking purpose.
	SequenceId string `json:"sequence_id"`

	// The country code (2 char) of upper stream system.<br/> E.g.: CubeCN trigger the message sending, the country code should be CN.<br/> Refer to [`ISO 3166-1`](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2)
	CountryCode string `json:"country_code"`

	// BizData *BizData `json:"biz_data"`
	BizData map[string]interface{} `json:"biz_data"`
}
