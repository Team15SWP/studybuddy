package delivery

import (
	"fmt"
	"net/http"
	"time"

	"study_buddy/internal/model"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetNotificationSettings(g *gin.Context) {
	g.Set(constants.HandlerName, "GetNotificationSettings")
	userId, ok := g.Get(constants.UserID)
	if !ok {
		_ = g.Error(fmt.Errorf("user not found"))
		return
	}

	response, err := h.service.GetNotificationSettings(g.Request.Context(), userId.(int64))
	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.EvaluateCodeForTask: %w", err))
		return
	}

	var resp struct {
		Enabled          bool   `json:"enabled"`
		NotificationTime string `json:"notification_time"`
		NotificationDays []int  `json:"notification_days"`
	}

	resp.Enabled = response.Enabled
	resp.NotificationTime = response.Time24.Format("15:04")
	resp.NotificationDays = response.Days

	g.JSON(http.StatusOK, resp)
}

type NotificationRequest struct {
	Enabled          bool   `json:"enabled"`
	NotificationTime string `json:"notification_time"`
	NotificationDays []int  `json:"notification_days"`
}

func (h *Handler) SetNotificationSettings(g *gin.Context) {
	g.Set(constants.HandlerName, "SetNotificationSettings")
	var request NotificationRequest
	if err := g.BindJSON(&request); err != nil {
		_ = g.Error(fmt.Errorf("impossible to unmarshall: %w", errlist.ErrBadRequest))
		return
	}

	userId, ok := g.Get(constants.UserID)
	if !ok {
		_ = g.Error(fmt.Errorf("user not found"))
		return
	}

	tt, err := time.Parse("15:04", request.NotificationTime)
	if err != nil {
		_ = g.Error(fmt.Errorf("invalid notification_time: %w", err))
		return
	}

	err = h.service.SetNotificationSettings(g.Request.Context(), &model.Notification{
		UserID:  userId.(int64),
		Enabled: request.Enabled,
		Time24:  tt,
		Days:    request.NotificationDays,
	})
	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.EvaluateCodeForTask: %w", err))
		return
	}

	g.Status(http.StatusOK)
}
