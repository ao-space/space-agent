package identify

import (
	"agent/biz/model/dto/bind/identify"
	identifySvc "agent/biz/service/bind/identify"
	"agent/config"
	"agent/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetTicket godoc
// @Summary 获取临时ticket [客户端蓝牙/局域网调用/网关call调用]
// @Description 通道认证，client获取临时ticket
// @ID GetTicket
// @Tags Pair
// @Accept  plain
// @Produce  json
// @Param   ticketReq body identify.TicketReq true  "boxuuid"
// @Success 200 {object} dto.BaseRspStr{results=identify.TicketRsp} "code=AG-200 成功;"
// @Router /agent/v1/api/bind/identify/ticket [POST]
func GetTicket(c *gin.Context) {
	logger.AppLogger().Debugf("%+v", c.Request)

	var req identify.TicketReq
	svc := new(identifySvc.TicketService)
	if c.Request.Host == config.Config.Web.DockerLocalListenAddr {
		c.JSON(http.StatusOK, svc.InitGatewayService("", c.Request.Header, c).Enter(svc, &req))
	} else {
		c.JSON(http.StatusOK, svc.InitLanService("", c.Request.Header, c).Enter(svc, &req))
	}

}
