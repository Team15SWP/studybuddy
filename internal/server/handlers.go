package server

import (
	authDelivery "study_buddy/internal/core/auth/delivery"
	authUseCase "study_buddy/internal/core/auth/service"
	syllabusDelivery "study_buddy/internal/core/syllabus/delivery"
	syllabusUseCase "study_buddy/internal/core/syllabus/service"
	taskDelivery "study_buddy/internal/core/task/delivery"
	taskUseCase "study_buddy/internal/core/task/service"
	"study_buddy/internal/middlewares"
	"study_buddy/pkg/llm"
)

func (s *Server) setupMiddlewares() {
	s.router.Use(middlewares.CORSMiddleware)
	s.router.Use(middlewares.ErrorHandler(s.log))
}

func (s *Server) mapHandlers() {
	authService := authUseCase.NewAuthService(s.repoLayer.authRepo, s.repoLayer.notificationRepo, s.cfg)
	authHandlers := authDelivery.NewAuthHandler(authService)
	authDelivery.MapAuthRoutes(s.router, authHandlers)

	authMiddleware := middlewares.UserIdentity(&s.cfg.HashConfig)

	syllabusService := syllabusUseCase.NewSyllabusService(s.repoLayer.syllabusRepo)
	syllabusHandlers := syllabusDelivery.NewSyllabusHandler(syllabusService)
	syllabusDelivery.MapSyllabusRoutes(s.router, syllabusHandlers, authMiddleware)

	llmClient := llm.NewOpenRouterClient(&s.cfg.OpenAI)

	taskService := taskUseCase.NewTaskService(s.repoLayer.taskRepo, s.repoLayer.statsRepo, s.repoLayer.notificationRepo, llmClient, &s.cfg.OpenAI, &s.cfg.Prompts)
	taskHandlers := taskDelivery.NewTaskHandler(taskService)
	taskDelivery.MapTaskRoutes(s.router, taskHandlers, authMiddleware)
}
