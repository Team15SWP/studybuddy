package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"study_buddy/internal/config"
	"study_buddy/internal/model"
	"study_buddy/pkg/constants"
	"study_buddy/pkg/llm"
)

var _ Service = (*TaskService)(nil)

type TaskService struct {
	repo             TaskProvider
	statsRepo        StatsProvider
	notificationRepo NotificationProvider
	openAi           *config.OpenAI
	prompts          *config.Prompts
	LLM              llm.Client
}

func NewTaskService(repo TaskProvider, statsRepo StatsProvider, notificationRepo NotificationProvider, llmClient llm.Client, openAi *config.OpenAI, prompts *config.Prompts) *TaskService {
	return &TaskService{
		repo:             repo,
		statsRepo:        statsRepo,
		notificationRepo: notificationRepo,
		openAi:           openAi,
		prompts:          prompts,
		LLM:              llmClient,
	}
}

type Service interface {
	GenerateTask(ctx context.Context, userId int64, topic, difficulty string) (*model.Task, error)
	EvaluateCodeForTask(ctx context.Context, userId int64, task, code string) (*model.Feedback, error)
	GetStatistics(ctx context.Context, userId int64) (*model.Statistics, error)

	GetNotificationSettings(ctx context.Context, userId int64) (*model.Notification, error)
	SetNotificationSettings(ctx context.Context, notif *model.Notification) error
}

type TaskProvider interface {
	CreateTask(ctx context.Context, task *model.GeneratedTask) error
	GetTask(ctx context.Context, userId int64, taskName string) (*model.GeneratedTask, error)
	UpdateTaskSolved(ctx context.Context, task *model.GeneratedTask, stats *model.Statistics) error
}

type StatsProvider interface {
	GetStatisticsData(ctx context.Context, userId int64) (*model.Statistics, error)
}

type NotificationProvider interface {
	GetNotification(ctx context.Context, userId int64) (*model.Notification, error)
	CreateNotification(ctx context.Context, notif *model.Notification) error
	UpdateNotification(ctx context.Context, notif *model.Notification) error
}

func (t *TaskService) GenerateTask(ctx context.Context, userId int64, topic, difficulty string) (*model.Task, error) {
	prompt := fmt.Sprintf(t.prompts.GenerateTask, topic, difficulty)
	response, err := t.LLM.Complete(ctx, prompt)
	if err != nil {
		return nil, err
	}
	var task *model.GeneratedTask
	err = json.Unmarshal([]byte(response), &task)
	if err == nil && task != nil {
		task.Difficulty = constants.DifficultyToInt(difficulty)
		task.UserID = userId
		task.Solved = 0
		err = t.repo.CreateTask(ctx, task)
		if err != nil {
			return nil, fmt.Errorf("t.repo.CreateTask: %w", err)
		}
	}
	return task.ToServer(), nil
}

func (t *TaskService) EvaluateCodeForTask(ctx context.Context, userId int64, task, code string) (*model.Feedback, error) {
	fullTask, err := t.repo.GetTask(ctx, userId, task)
	if err != nil {
		return nil, fmt.Errorf("t.repo.GetTask: %w", err)
	}
	stats, err := t.statsRepo.GetStatisticsData(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("t.statsRepo.GetStatisticsData: %w", err)
	}

	if strings.Contains(code, "show solution") {
		return &model.Feedback{
			Feedback: fullTask.Solution,
		}, nil
	}

	taskJSON, err := json.MarshalIndent(fullTask, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	prompt := fmt.Sprintf(t.prompts.CheckCodeForTask, string(taskJSON), code)

	response, err := t.LLM.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("t.LLM.Complete: %w", err)
	}

	feedback := &model.Question{}
	feedback.Task = task
	feedback.Code = code
	err = json.Unmarshal([]byte(response), &feedback)
	if err != nil {
		return nil, fmt.Errorf("t.LLM.Complete: %w", err)
	}

	if !feedback.Request {
		if feedback.Verdict {
			fullTask.Solved++
			switch fullTask.Difficulty {
			case constants.Easy:
				stats.Easy++
			case constants.Medium:
				stats.Medium++
			case constants.Hard:
				stats.Hard++
			default:
			}
			stats.Total++
			err = t.repo.UpdateTaskSolved(ctx, fullTask, stats)
			if err != nil {
				return nil, fmt.Errorf("t.repo.UpdateTaskSolved: %w", err)
			}
		}
	}
	return feedback.ToFeedback(), nil
}

func (t *TaskService) GetStatistics(ctx context.Context, userId int64) (*model.Statistics, error) {
	response, err := t.statsRepo.GetStatisticsData(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("t.statsRepo: %w", err)
	}
	return response, nil
}
