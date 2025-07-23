package server

import (
	authRepo "study_buddy/internal/core/auth/repository"
	syllabusRepo "study_buddy/internal/core/syllabus/repository"
	taskRepo "study_buddy/internal/core/task/repository"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repoLayer struct {
	authRepo         authRepo.Repository
	syllabusRepo     syllabusRepo.Repository
	taskRepo         taskRepo.TaskRepository
	statsRepo        taskRepo.StatsRepository
	notificationRepo taskRepo.NotificationRepository
	trManager        *manager.Manager
}

func initRepoLayer(db *pgxpool.Pool) *repoLayer {
	return &repoLayer{
		authRepo:         authRepo.NewAuthRepo(db),
		syllabusRepo:     syllabusRepo.NewSyllabusRepo(db),
		taskRepo:         taskRepo.NewTaskRepo(db),
		statsRepo:        taskRepo.NewStatsRepo(db),
		notificationRepo: taskRepo.NewNotificationRepo(db),
		trManager:        manager.Must(trmpgx.NewDefaultFactory(db)),
	}
}
