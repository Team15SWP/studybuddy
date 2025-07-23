package delivery

import (
	"fmt"
	"net/http"

	taskUseCase "study_buddy/internal/core/task/service"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"

	"github.com/gin-gonic/gin"
)

var _ Handlers = (*Handler)(nil)

type Handler struct {
	service taskUseCase.Service
}

func NewTaskHandler(service taskUseCase.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type Handlers interface {
	GenerateTask(g *gin.Context)
	EvaluateCodeForTask(g *gin.Context)
	GetStatistics(g *gin.Context)

	GetNotificationSettings(g *gin.Context)
	SetNotificationSettings(g *gin.Context)
}

type GenerateTaskRequest struct {
	Topic      string `json:"topic" form:"topic"`
	Difficulty string `json:"difficulty" form:"difficulty"`
}

func (h *Handler) GenerateTask(g *gin.Context) {
	g.Set(constants.HandlerName, "GenerateTask")
	var request GenerateTaskRequest
	if err := g.BindJSON(&request); err != nil {
		_ = g.Error(fmt.Errorf("impossible to unmarshall: %w", errlist.ErrBadRequest))
		return
	}

	userId, ok := g.Get(constants.UserID)
	if !ok {
		_ = g.Error(fmt.Errorf("user not found"))
		return
	}

	response, err := h.service.GenerateTask(g.Request.Context(), userId.(int64), request.Topic, request.Difficulty)

	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.GenerateTask: %w", err))
		return
	}

	g.JSON(http.StatusOK, response)
}

type EvaluateCodeRequest struct {
	Task string `json:"task" form:"task"`
	Code string `json:"code" form:"code"`
}

func (h *Handler) EvaluateCodeForTask(g *gin.Context) {
	g.Set(constants.HandlerName, "EvaluateCodeForTask")
	var request EvaluateCodeRequest
	if err := g.BindJSON(&request); err != nil {
		_ = g.Error(fmt.Errorf("impossible to unmarshall: %w", errlist.ErrBadRequest))
		return
	}

	userId, ok := g.Get(constants.UserID)
	if !ok {
		_ = g.Error(fmt.Errorf("user not found"))
		return
	}

	response, err := h.service.EvaluateCodeForTask(g.Request.Context(), userId.(int64), request.Task, request.Code)

	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.EvaluateCodeForTask: %w", err))
		return
	}

	g.JSON(http.StatusOK, response)
}

func (h *Handler) GetStatistics(g *gin.Context) {
	g.Set(constants.HandlerName, "GetStatistics")
	userId, ok := g.Get(constants.UserID)
	if !ok {
		_ = g.Error(errlist.ErrUnauthorized)
		return
	}

	response, err := h.service.GetStatistics(g.Request.Context(), userId.(int64))
	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.GetStatistics: %w", err))
		return
	}
	g.JSON(http.StatusOK, response)
}
