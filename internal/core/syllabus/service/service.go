package service

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	"study_buddy/internal/model"

	"github.com/ledongthuc/pdf"
)

var _ Service = (*SyllabusService)(nil)

type SyllabusService struct {
	repo SyllabusProvider
}

func NewSyllabusService(repo SyllabusProvider) *SyllabusService {
	return &SyllabusService{
		repo: repo,
	}
}

type Service interface {
	GetSyllabus(ctx context.Context) (*model.Syllabus, error)
	SaveSyllabus(ctx context.Context, File *multipart.FileHeader) (*model.Syllabus, error)
	DeleteSyllabus(ctx context.Context) error
}

type SyllabusProvider interface {
	GetSyllabus(ctx context.Context) ([]string, error)
	SaveSyllabus(ctx context.Context, syllabus []model.Schedule) ([]string, error)
	DeleteSyllabus(ctx context.Context) error
}

func (s *SyllabusService) GetSyllabus(ctx context.Context) (*model.Syllabus, error) {
	response, err := s.repo.GetSyllabus(ctx)
	if err != nil {
		return nil, fmt.Errorf("[syllabysService][repo.GetSyllabus]: %w", err)
	}
	return &model.Syllabus{
		Topics: response,
	}, nil
}

func (s *SyllabusService) DeleteSyllabus(ctx context.Context) error {
	err := s.repo.DeleteSyllabus(ctx)
	if err != nil {
		return fmt.Errorf("[syllabysService][repo.DeleteSyllabus]: %w", err)
	}
	return nil
}

func (s *SyllabusService) SaveSyllabus(ctx context.Context, File *multipart.FileHeader) (*model.Syllabus, error) {
	schedule, err := extractCourseScheduleFromUpload(File)
	if err != nil {
		return nil, fmt.Errorf("[syllabysService][extractCourseScheduleFromUpload]: %w", err)
	}
	response, err := s.repo.SaveSyllabus(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("[syllabysService][repo.SaveSyllabus]: %w", err)
	}
	return &model.Syllabus{
		Topics: response,
	}, nil
}

func extractCourseScheduleFromUpload(fileHeader *multipart.FileHeader) ([]model.Schedule, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pdfReader, err := pdf.NewReader(file, fileHeader.Size)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	textReader, err := pdfReader.GetPlainText()
	if err != nil {
		return nil, err
	}
	buf.ReadFrom(textReader)
	raw := buf.String()

	re := regexp.MustCompile(`(?i)week\s*[-:]?\s*(\d{1,2})\s*[:-]?\s*([A-Za-z0-9 ,.\-–—()\[\]]+)`)
	matches := re.FindAllStringSubmatch(raw, -1)

	var result []model.Schedule
	for _, m := range matches {
		if len(m) >= 3 {
			result = append(result, model.Schedule{
				Week:  m[1],
				Topic: strings.TrimSpace(m[2]),
			})
		}
	}
	return result, nil
}
