package httpin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/EnduranNSU/trainings/internal/adapter/in/http/dto"
	svctraining "github.com/EnduranNSU/trainings/internal/domain"
)

type TrainingHandler struct {
	svc svctraining.TrainingService
}

func NewTrainingHandler(svc svctraining.TrainingService) *TrainingHandler {
	return &TrainingHandler{svc: svc}
}

// GetTrainingsByUser получает все тренировки пользователя
// @Summary      Получить тренировки пользователя
// @Description  Возвращает все тренировки указанного пользователя
// @Tags         trainings
// @Produce      json
// @Param        user_id query string true "User ID"
// @Success      200  {array}   dto.TrainingResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings [get]
func (h *TrainingHandler) GetTrainingsByUser(c *gin.Context) {
	uidStr := c.Query("user_id")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	trainings, err := h.svc.GetTrainingsByUser(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get trainings"})
		return
	}

	if len(trainings) == 0 {
		c.JSON(http.StatusOK, []dto.TrainingResponse{})
		return
	}

	resp := make([]dto.TrainingResponse, 0, len(trainings))
	for _, training := range trainings {
		resp = append(resp, h.trainingToResponse(training))
	}

	c.JSON(http.StatusOK, resp)
}

// GetTrainingWithExercises получает тренировку с упражнениями
// @Summary      Получить тренировку с упражнениями
// @Description  Возвращает информацию о тренировке со списком упражнений
// @Tags         trainings
// @Produce      json
// @Param        id path int64 true "Training ID"
// @Success      200  {object}  dto.TrainingResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings/{id} [get]
func (h *TrainingHandler) GetTrainingWithExercises(c *gin.Context) {
	trainingID, err := parseInt64Param(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid training id"})
		return
	}

	training, err := h.svc.GetTrainingWithExercises(c.Request.Context(), trainingID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: "training not found"})
		return
	}

	c.JSON(http.StatusOK, h.trainingToResponse(training))
}

// CreateTraining создает новую тренировку
// @Summary      Создать тренировку
// @Description  Создает новую тренировку для пользователя
// @Tags         trainings
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTrainingRequest true "Данные тренировки"
// @Success      201  {object}  dto.TrainingResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings [post]
func (h *TrainingHandler) CreateTraining(c *gin.Context) {
	var req dto.CreateTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad json"})
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	plannedDate, err := time.Parse(time.RFC3339, req.PlannedDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid planned_date time format"})
		return
	}

	var actualDate, startedAt, finishedAt *time.Time
	var totalDuration, totalRestTime, totalExerciseTime *time.Duration

	// Parse optional time fields
	if req.ActualDate != nil {
		t, err := time.Parse(time.RFC3339, *req.ActualDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid actual_date time format"})
			return
		}
		actualDate = &t
	}
	if req.StartedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.StartedAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid started_at time format"})
			return
		}
		startedAt = &t
	}
	if req.FinishedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.FinishedAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid finished_at time format"})
			return
		}
		finishedAt = &t
	}

	// Parse optional duration fields
	if req.TotalDuration != nil {
		duration, err := time.ParseDuration(*req.TotalDuration)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_duration format"})
			return
		}
		totalDuration = &duration
	}
	if req.TotalRestTime != nil {
		duration, err := time.ParseDuration(*req.TotalRestTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_rest_time format"})
			return
		}
		totalRestTime = &duration
	}
	if req.TotalExerciseTime != nil {
		duration, err := time.ParseDuration(*req.TotalExerciseTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_exercise_time format"})
			return
		}
		totalExerciseTime = &duration
	}

	cmd := svctraining.CreateTrainingCmd{
		UserID:            uid,
		IsDone:            req.IsDone,
		PlannedDate:       plannedDate,
		ActualDate:        actualDate,
		StartedAt:         startedAt,
		FinishedAt:        finishedAt,
		TotalDuration:     totalDuration,
		TotalRestTime:     totalRestTime,
		TotalExerciseTime: totalExerciseTime,
		Rating:            req.Rating,
	}

	training, err := h.svc.CreateTraining(c.Request.Context(), cmd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to create training"})
		return
	}

	c.JSON(http.StatusCreated, h.trainingToResponse(training))
}

// UpdateTraining обновляет тренировку
// @Summary      Обновить тренировку
// @Description  Обновляет информацию о тренировке
// @Tags         trainings
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Training ID"
// @Param        request body dto.UpdateTrainingRequest true "Данные для обновления"
// @Success      200  {object}  dto.TrainingResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings/{id} [put]
func (h *TrainingHandler) UpdateTraining(c *gin.Context) {
	trainingID, err := parseInt64Param(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid training id"})
		return
	}

	var req dto.UpdateTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad json"})
		return
	}

	plannedDate, err := time.Parse(time.RFC3339, req.PlannedDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid planned_date time format"})
		return
	}

	var actualDate, startedAt, finishedAt *time.Time
	var totalDuration, totalRestTime, totalExerciseTime *time.Duration

	if req.ActualDate != nil {
		t, err := time.Parse(time.RFC3339, *req.ActualDate)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid actual_date time format"})
			return
		}
		actualDate = &t
	}
	if req.StartedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.StartedAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid started_at time format"})
			return
		}
		startedAt = &t
	}
	if req.FinishedAt != nil {
		t, err := time.Parse(time.RFC3339, *req.FinishedAt)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid finished_at time format"})
			return
		}
		finishedAt = &t
	}

	if req.TotalDuration != nil {
		duration, err := time.ParseDuration(*req.TotalDuration)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_duration format"})
			return
		}
		totalDuration = &duration
	}
	if req.TotalRestTime != nil {
		duration, err := time.ParseDuration(*req.TotalRestTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_rest_time format"})
			return
		}
		totalRestTime = &duration
	}
	if req.TotalExerciseTime != nil {
		duration, err := time.ParseDuration(*req.TotalExerciseTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_exercise_time format"})
			return
		}
		totalExerciseTime = &duration
	}

	cmd := svctraining.UpdateTrainingCmd{
		ID:                trainingID,
		IsDone:            req.IsDone,
		PlannedDate:       plannedDate,
		ActualDate:        actualDate,
		StartedAt:         startedAt,
		FinishedAt:        finishedAt,
		TotalDuration:     totalDuration,
		TotalRestTime:     totalRestTime,
		TotalExerciseTime: totalExerciseTime,
		Rating:            req.Rating,
	}

	training, err := h.svc.UpdateTraining(c.Request.Context(), cmd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update training"})
		return
	}

	c.JSON(http.StatusOK, h.trainingToResponse(training))
}

// DeleteTraining удаляет тренировку
// @Summary      Удалить тренировку
// @Description  Удаляет тренировку по ID
// @Tags         trainings
// @Produce      json
// @Param        id path int64 true "Training ID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings/{id} [delete]
func (h *TrainingHandler) DeleteTraining(c *gin.Context) {
	trainingID, err := parseInt64Param(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid training id"})
		return
	}

	err = h.svc.DeleteTraining(c.Request.Context(), trainingID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to delete training"})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddExerciseToTraining добавляет упражнение к тренировке
// @Summary      Добавить упражнение к тренировке
// @Description  Добавляет упражнение к существующей тренировке
// @Tags         training-exercises
// @Accept       json
// @Produce      json
// @Param        request body dto.AddExerciseToTrainingRequest true "Данные упражнения"
// @Success      201  {object}  dto.TrainedExerciseResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/training-exercises [post]
func (h *TrainingHandler) AddExerciseToTraining(c *gin.Context) {
	var req dto.AddExerciseToTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad json"})
		return
	}

	var weight *decimal.Decimal
	if req.Weight != nil {
		w := decimal.NewFromFloat(*req.Weight)
		weight = &w
	}

	var timeVal, doing, rest *time.Duration
	if req.Time != nil {
		duration, err := time.ParseDuration(*req.Time)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid time format, use duration format like '1h30m'"})
			return
		}
		timeVal = &duration
	}
	if req.Doing != nil {
		duration, err := time.ParseDuration(*req.Doing)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid doing format, use duration format like '1h30m'"})
			return
		}
		doing = &duration
	}
	if req.Rest != nil {
		duration, err := time.ParseDuration(*req.Rest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid rest format, use duration format like '1h30m'"})
			return
		}
		rest = &duration
	}

	// Convert approaches and reps from *int64 to *int32
	var approaches, reps *int32
	if req.Approaches != nil {
		a := int32(*req.Approaches)
		approaches = &a
	}
	if req.Reps != nil {
		r := int32(*req.Reps)
		reps = &r
	}

	cmd := svctraining.AddExerciseToTrainingCmd{
		TrainingID: req.TrainingID,
		ExerciseID: req.ExerciseID,
		Weight:     weight,
		Approaches: approaches,
		Reps:       reps,
		Time:       timeVal,
		Doing:      doing,
		Rest:       rest,
		Notes:      req.Notes,
	}

	exercise, err := h.svc.AddExerciseToTraining(c.Request.Context(), cmd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to add exercise to training"})
		return
	}

	c.JSON(http.StatusCreated, h.trainedExerciseToResponse(exercise))
}

// UpdateTrainedExercise обновляет выполненное упражнение
// @Summary      Обновить выполненное упражнение
// @Description  Обновляет информацию о выполненном упражнении в тренировке
// @Tags         training-exercises
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Trained Exercise ID"
// @Param        request body dto.UpdateTrainedExerciseRequest true "Данные для обновления"
// @Success      200  {object}  dto.TrainedExerciseResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/training-exercises/{id} [put]
func (h *TrainingHandler) UpdateTrainedExercise(c *gin.Context) {
	exerciseID, err := parseInt64Param(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid exercise id"})
		return
	}

	var req dto.UpdateTrainedExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad json"})
		return
	}

	var weight *decimal.Decimal
	if req.Weight != nil {
		w := decimal.NewFromFloat(*req.Weight)
		weight = &w
	}

	var timeVal, doing, rest *time.Duration
	if req.Time != nil {
		duration, err := time.ParseDuration(*req.Time)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid time format, use duration format like '1h30m'"})
			return
		}
		timeVal = &duration
	}
	if req.Doing != nil {
		duration, err := time.ParseDuration(*req.Doing)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid doing format, use duration format like '1h30m'"})
			return
		}
		doing = &duration
	}
	if req.Rest != nil {
		duration, err := time.ParseDuration(*req.Rest)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid rest format, use duration format like '1h30m'"})
			return
		}
		rest = &duration
	}

	// Convert approaches and reps from *int64 to *int32
	var approaches, reps *int32
	if req.Approaches != nil {
		a := int32(*req.Approaches)
		approaches = &a
	}
	if req.Reps != nil {
		r := int32(*req.Reps)
		reps = &r
	}

	cmd := svctraining.UpdateTrainedExerciseCmd{
		ID:         exerciseID,
		Weight:     weight,
		Approaches: approaches,
		Reps:       reps,
		Time:       timeVal,
		Doing:      doing,
		Rest:       rest,
		Notes:      req.Notes,
	}

	exercise, err := h.svc.UpdateTrainedExercise(c.Request.Context(), cmd)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, h.trainedExerciseToResponse(exercise))
}

// RemoveExerciseFromTraining удаляет упражнение из тренировки
// @Summary      Удалить упражнение из тренировки
// @Description  Удаляет упражнение из тренировки
// @Tags         training-exercises
// @Produce      json
// @Param        training_id query int64 true "Training ID"
// @Param        exercise_id query int64 true "Exercise ID"
// @Success      204
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/training-exercises [delete]
func (h *TrainingHandler) RemoveExerciseFromTraining(c *gin.Context) {
	trainingID, err := parseInt64Query(c, "training_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid training_id"})
		return
	}

	exerciseID, err := parseInt64Query(c, "exercise_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid exercise_id"})
		return
	}

	err = h.svc.RemoveExerciseFromTraining(c.Request.Context(), trainingID, exerciseID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to remove exercise from training"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserTrainingStats получает статистику тренировок пользователя
// @Summary      Получить статистику тренировок
// @Description  Возвращает статистику тренировок пользователя
// @Tags         trainings
// @Produce      json
// @Param        user_id query string true "User ID"
// @Success      200  {object}  dto.TrainingStatsResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings/stats [get]
func (h *TrainingHandler) GetUserTrainingStats(c *gin.Context) {
	uidStr := c.Query("user_id")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
		return
	}

	stats, err := h.svc.GetUserTrainingStats(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get training stats"})
		return
	}

	resp := dto.TrainingStatsResponse{
		TotalTrainings:     stats.TotalTrainings,
		CompletedTrainings: stats.CompletedTrainings,
		AverageRating:      stats.AverageRating,
		TotalDuration:      stats.TotalTime,
	}

	c.JSON(http.StatusOK, resp)
}

// CompleteTraining завершает тренировку
// @Summary      Завершить тренировку
// @Description  Отмечает тренировку как завершенную и устанавливает рейтинг
// @Tags         trainings
// @Accept       json
// @Produce      json
// @Param        id path int64 true "Training ID"
// @Param        request body dto.CompleteTrainingRequest true "Данные для завершения"
// @Success      200  {object}  dto.TrainingResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/v1/trainings/{id}/complete [patch]
func (h *TrainingHandler) CompleteTraining(c *gin.Context) {
	trainingID, err := parseInt64Param(c, "id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid training id"})
		return
	}

	var req dto.CompleteTrainingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad json"})
		return
	}

	training, err := h.svc.CompleteTraining(c.Request.Context(), trainingID, req.Rating)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to complete training"})
		return
	}

	c.JSON(http.StatusOK, h.trainingToResponse(training))
}

// Вспомогательные методы

func (h *TrainingHandler) trainingToResponse(training *svctraining.Training) dto.TrainingResponse {
	var actualDate, startedAt, finishedAt *string
	if training.ActualDate != nil {
		s := training.ActualDate.Format(time.RFC3339)
		actualDate = &s
	}
	if training.StartedAt != nil {
		s := training.StartedAt.Format(time.RFC3339)
		startedAt = &s
	}
	if training.FinishedAt != nil {
		s := training.FinishedAt.Format(time.RFC3339)
		finishedAt = &s
	}

	var totalDuration, totalRestTime, totalExerciseTime *string
	if training.TotalDuration != nil {
		s := formatDuration(*training.TotalDuration)
		totalDuration = &s
	}
	if training.TotalRestTime != nil {
		s := formatDuration(*training.TotalRestTime)
		totalRestTime = &s
	}
	if training.TotalExerciseTime != nil {
		s := formatDuration(*training.TotalExerciseTime)
		totalExerciseTime = &s
	}

	var exercises []dto.TrainedExerciseResponse
	if training.Exercises != nil {
		exercises = make([]dto.TrainedExerciseResponse, 0, len(training.Exercises))
		for _, exercise := range training.Exercises {
			exercises = append(exercises, h.trainedExerciseToResponse(&exercise))
		}
	}

	return dto.TrainingResponse{
		ID:                training.ID,
		UserID:            training.UserID.String(),
		IsDone:            training.IsDone,
		PlannedDate:       training.PlannedDate.Format(time.RFC3339),
		ActualDate:        actualDate,
		StartedAt:         startedAt,
		FinishedAt:        finishedAt,
		TotalDuration:     totalDuration,
		TotalRestTime:     totalRestTime,
		TotalExerciseTime: totalExerciseTime,
		Rating:            training.Rating,
		Exercises:         exercises,
	}
}

func (h *TrainingHandler) trainedExerciseToResponse(exercise *svctraining.TrainedExercise) dto.TrainedExerciseResponse {
	var weight *float64
	if exercise.Weight != nil {
		f, _ := exercise.Weight.Float64()
		weight = &f
	}

	var timeStr, doingStr, restStr *string
	if exercise.Time != nil {
		s := formatDuration(*exercise.Time)
		timeStr = &s
	}
	if exercise.Doing != nil {
		s := formatDuration(*exercise.Doing)
		doingStr = &s
	}
	if exercise.Rest != nil {
		s := formatDuration(*exercise.Rest)
		restStr = &s
	}

	return dto.TrainedExerciseResponse{
		ID:         exercise.ID,
		TrainingID: exercise.TrainingID,
		ExerciseID: exercise.ExerciseID,
		Weight:     weight,
		Approaches: exercise.Approaches,
		Reps:       exercise.Reps,
		Time:       timeStr,
		Doing:      doingStr,
		Rest:       restStr,
		Notes:      exercise.Notes,
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}