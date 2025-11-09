package httpin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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

	planned, err := time.Parse(time.RFC3339, req.Planned)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid planned time format"})
		return
	}

	var done *time.Time
	if req.Done != nil {
		doneTime, err := time.Parse(time.RFC3339, *req.Done)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid done time format"})
			return
		}
		done = &doneTime
	}

	var totalTime *time.Duration
	if req.TotalTime != nil {
		duration, err := time.ParseDuration(*req.TotalTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_time format"})
			return
		}
		totalTime = &duration
	}

	cmd := svctraining.CreateTrainingCmd{
		UserID:    uid,
		IsDone:    req.IsDone,
		Planned:   planned,
		Done:      done,
		TotalTime: totalTime,
		Rating:    req.Rating,
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

	planned, err := time.Parse(time.RFC3339, req.Planned)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid planned time format"})
		return
	}

	var done *time.Time
	if req.Done != nil {
		doneTime, err := time.Parse(time.RFC3339, *req.Done)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid done time format"})
			return
		}
		done = &doneTime
	}

	var totalTime *time.Duration
	if req.TotalTime != nil {
		duration, err := time.ParseDuration(*req.TotalTime)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid total_time format"})
			return
		}
		totalTime = &duration
	}

	cmd := svctraining.UpdateTrainingCmd{
		ID:        trainingID,
		IsDone:    req.IsDone,
		Planned:   planned,
		Done:      done,
		TotalTime: totalTime,
		Rating:    req.Rating,
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

	var exerciseTime *time.Time
	if req.Time != nil {
		timeVal, err := time.Parse(time.RFC3339, *req.Time)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid time format"})
			return
		}
		exerciseTime = &timeVal
	}

	cmd := svctraining.AddExerciseToTrainingCmd{
		TrainingID: req.TrainingID,
		ExerciseID: req.ExerciseID,
		Weight:     req.Weight,
		Approaches: req.Approaches,
		Reps:       req.Reps,
		Time:       exerciseTime,
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

	var exerciseTime *time.Time
	if req.Time != nil {
		timeVal, err := time.Parse(time.RFC3339, *req.Time)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid time format"})
			return
		}
		exerciseTime = &timeVal
	}

	cmd := svctraining.UpdateTrainedExerciseCmd{
		ID:         exerciseID,
		Weight:     req.Weight,
		Approaches: req.Approaches,
		Reps:       req.Reps,
		Time:       exerciseTime,
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

	var lastTrainingDate *string

	resp := dto.TrainingStatsResponse{
		TotalTrainings:     stats.TotalTrainings,
		CompletedTrainings: stats.CompletedTrainings,
		AverageRating:      stats.AverageRating,
		TotalDuration:      stats.TotalTime,
		LastTrainingDate:   lastTrainingDate,
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
	var done *string
	if training.Done != nil {
		doneStr := training.Done.Format(time.RFC3339)
		done = &doneStr
	}

	var totalTime *string
	if training.TotalTime != nil {
		durationStr := training.TotalTime.String()
		totalTime = &durationStr
	}

	var exercises []dto.TrainedExerciseResponse
	if training.Exercises != nil {
		exercises = make([]dto.TrainedExerciseResponse, 0, len(training.Exercises))
		for _, exercise := range training.Exercises {
			exercises = append(exercises, h.trainedExerciseToResponse(&exercise))
		}
	}


	return dto.TrainingResponse{
		ID:        training.ID,
		UserID:    training.UserID.String(),
		IsDone:    training.IsDone,
		Planned:   training.Planned.Format(time.RFC3339),
		Done:      done,
		TotalTime: totalTime,
		Rating:    training.Rating,
		Exercises: exercises,
	}
}

func (h *TrainingHandler) trainedExerciseToResponse(exercise *svctraining.TrainedExercise) dto.TrainedExerciseResponse {
	var timeStr *string
	if exercise.Time != nil {
		timeVal := exercise.Time.Format(time.RFC3339)
		timeStr = &timeVal
	}

	return dto.TrainedExerciseResponse{
		ID:         exercise.ID,
		TrainingID: exercise.TrainingID,
		ExerciseID: exercise.ExerciseID,
		Weight:     exercise.Weight,
		Approaches: exercise.Approaches,
		Reps:       exercise.Reps,
		Time:       timeStr,
		Notes:      exercise.Notes,
	}
}
