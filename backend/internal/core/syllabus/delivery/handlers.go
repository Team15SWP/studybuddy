package delivery

import (
	"fmt"
	"mime/multipart"
	"net/http"

	syllabusUseCase "study_buddy/internal/core/syllabus/service"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/errlist"

	"github.com/gin-gonic/gin"
)

var _ Handlers = (*Handler)(nil)

type Handler struct {
	service syllabusUseCase.Service
}

func NewSyllabusHandler(service syllabusUseCase.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type Handlers interface {
	GetSyllabus(g *gin.Context)
	SaveSyllabus(g *gin.Context)
	DeleteSyllabus(g *gin.Context)
}

func (h *Handler) GetSyllabus(g *gin.Context) {
	g.Set(constants.HandlerName, "GetSyllabus")
	response, err := h.service.GetSyllabus(g.Request.Context())
	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.GetSyllabus: %w", err))
		return
	}

	g.JSON(http.StatusOK, response)
}

type LoadSyllabusRequest struct {
	File *multipart.FileHeader `form:"file" json:"-" binding:"required"`
}

func (h *Handler) SaveSyllabus(g *gin.Context) {
	g.Set(constants.HandlerName, "LoadSyllabus")
	var request LoadSyllabusRequest
	if err := g.ShouldBind(&request); err != nil {
		_ = g.Error(fmt.Errorf("impossible to unmarshall: %w", errlist.ErrBadRequest))
		return
	}

	response, err := h.service.SaveSyllabus(g.Request.Context(), request.File)
	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.SaveSyllabus: %w", err))
		return
	}

	g.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSyllabus(g *gin.Context) {
	g.Set(constants.HandlerName, "DeleteSyllabus")
	err := h.service.DeleteSyllabus(g.Request.Context())
	if err != nil {
		_ = g.Error(fmt.Errorf("h.service.DeleteSyllabus: %w", err))
		return
	}
	g.Status(http.StatusOK)
}
