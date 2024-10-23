package mt_util

type MtResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type MtGetTimeResponse struct {
	MtResponse
	Data int64 `json:"data"`
}

type MtPreSendCouponResponse struct {
	MtResponse
	Data []*MtPreSendCouponResult
}

type MtPreSendCouponResult struct {
	Code string
	Msg  string
}
