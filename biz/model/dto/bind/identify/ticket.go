package identify

type TicketRsp struct {
	Ticket    string `json:"ticket"`    // 临时ticket用于通道实名认证，web携带后可以临时访问平台
	ExpiresAt string `json:"expiresAt"` // 过期时间，秒时间戳
}

type TicketReq struct {
	BoxUuid     string `json:"boxUuid"`
	PlatformUrl string `json:"platformUrl"` //平台的url
}
