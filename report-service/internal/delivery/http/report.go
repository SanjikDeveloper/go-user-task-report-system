package http

import (
	"time"

	"github.com/gin-gonic/gin"
)

// localhost:8000/api/v1/report/?startDate='2025-03-20'&endDate='2025-04-20'

type GetReportRequest struct {
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

func (h *Handler) GetReport(c *gin.Context) {
	// TODO: нужен ли юзер?
	//usernameTypeless, _ := c.Get("username")
	//
	//username, ok := usernameTypeless.(string)
	//if !ok {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}

	//userIDStr := c.Param("userID")
	//userID, err := strconv.ParseInt(userIDStr, 10, 64)
	//if err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}

	var req GetReportRequest
	if err := c.BindQuery(&req); err != nil {
		h.logger.Error("GetReport err", err)
		return
	}

	// TODO: проверять на ошибки
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		h.logger.Error("start time err:", err)
	}

	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		h.logger.Error("end time err:", err)
	}

	reportDat, err := h.services.ReportByRange(c, 0, "", start, end)
	if err != nil {
		h.logger.Error("h.services.ReportByRange(c, 0, ..., start, end) err:", err)
		return
	}
	c.HTML(200, "report.html", reportDat)
}
