package service

import (
	"context"
	"errors"
	"time"

	"github.com/EnduranNSU/trainings/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidTrainingID = errors.New("invalid training id")
	ErrTrainingNotFound  = errors.New("training not found")
)

func NewTrainingService(repo domain.TrainingRepository) domain.TrainingService {
	return &trainingService{repo: repo}
}

type trainingService struct {
	repo domain.TrainingRepository
}

func (s *trainingService) GetTrainingsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Training, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetTrainingsByUser(ctx, userID)
}

func (s *trainingService) GetTrainingWithExercises(ctx context.Context, trainingID int64) (*domain.Training, error) {
	if trainingID <= 0 {
		return nil, ErrInvalidTrainingID
	}

	return s.repo.GetTrainingWithExercises(ctx, trainingID)
}

func (s *trainingService) CreateTraining(ctx context.Context, cmd domain.CreateTrainingCmd) (*domain.Training, error) {
	if cmd.UserID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	if cmd.Planned.IsZero() {
		return nil, errors.New("planned date is required")
	}

	// Если тренировка создается как выполненная, но не указана дата выполнения,
	// устанавливаем текущее время
	var done *time.Time
	if cmd.IsDone && cmd.Done == nil {
		now := time.Now().UTC()
		done = &now
	} else {
		done = cmd.Done
	}

	training := &domain.Training{
		UserID:    cmd.UserID,
		IsDone:    cmd.IsDone,
		Planned:   cmd.Planned,
		Done:      done,
		TotalTime: cmd.TotalTime,
		Rating:    cmd.Rating,
	}

	return s.repo.CreateTraining(ctx, training)
}

func (s *trainingService) UpdateTraining(ctx context.Context, cmd domain.UpdateTrainingCmd) (*domain.Training, error) {
	if cmd.ID <= 0 {
		return nil, ErrInvalidTrainingID
	}

	// Проверяем существование тренировки
	existing, err := s.repo.GetTrainingWithExercises(ctx, cmd.ID)
	if err != nil {
		return nil, ErrTrainingNotFound
	}

	// Обновляем только переданные поля
	if cmd.IsDone != nil {
		existing.IsDone = *cmd.IsDone
	}
	if !cmd.Planned.IsZero() {
		existing.Planned = cmd.Planned
	}
	if cmd.Done != nil {
		existing.Done = cmd.Done
	}
	if cmd.TotalTime != nil {
		existing.TotalTime = cmd.TotalTime
	}
	if cmd.Rating != nil {
		existing.Rating = cmd.Rating
	}

	return s.repo.UpdateTraining(ctx, existing)
}

func (s *trainingService) DeleteTraining(ctx context.Context, trainingID int64) error {
	if trainingID <= 0 {
		return ErrInvalidTrainingID
	}

	return s.repo.DeleteTrainingAndExercises(ctx, trainingID)
}

func (s *trainingService) AddExerciseToTraining(ctx context.Context, cmd domain.AddExerciseToTrainingCmd) (*domain.TrainedExercise, error) {
	if cmd.TrainingID <= 0 {
		return nil, ErrInvalidTrainingID
	}
	if cmd.ExerciseID <= 0 {
		return nil, ErrInvalidExerciseID
	}

	// Проверяем существование тренировки
	_, err := s.repo.GetTrainingWithExercises(ctx, cmd.TrainingID)
	if err != nil {
		return nil, ErrTrainingNotFound
	}

	exercise := &domain.TrainedExercise{
		TrainingID: cmd.TrainingID,
		ExerciseID: cmd.ExerciseID,
		Weight:     cmd.Weight,
		Approaches: cmd.Approaches,
		Reps:       cmd.Reps,
		Time:       cmd.Time,
		Notes:      cmd.Notes,
	}

	return s.repo.AddExerciseToTraining(ctx, exercise)
}

func (s *trainingService) UpdateTrainedExercise(ctx context.Context, cmd domain.UpdateTrainedExerciseCmd) (*domain.TrainedExercise, error) {
	if cmd.ID <= 0 {
		return nil, ErrInvalidExerciseID
	}

	// Для обновления нам нужен существующий объект упражнения
	// В реальном приложении здесь нужно получить упражнение по ID
	// Для простоты создаем новый объект с переданными данными
	exercise := &domain.TrainedExercise{
		ID:         cmd.ID,
		Weight:     cmd.Weight,
		Approaches: cmd.Approaches,
		Reps:       cmd.Reps,
		Time:       cmd.Time,
		Notes:      cmd.Notes,
	}

	return s.repo.UpdateTrainedExercise(ctx, exercise)
}

func (s *trainingService) RemoveExerciseFromTraining(ctx context.Context, trainingID, exerciseID int64) error {
	if trainingID <= 0 {
		return ErrInvalidTrainingID
	}
	if exerciseID <= 0 {
		return ErrInvalidExerciseID
	}

	return s.repo.DeleteExerciseFromTraining(ctx, exerciseID, trainingID)
}

func (s *trainingService) GetUserTrainingStats(ctx context.Context, userID uuid.UUID) (*domain.TrainingStats, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetUserTrainingStats(ctx, userID)
}

func (s *trainingService) CompleteTraining(ctx context.Context, trainingID int64, rating *int32) (*domain.Training, error) {
	if trainingID <= 0 {
		return nil, ErrInvalidTrainingID
	}

	// Получаем текущую тренировку
	training, err := s.repo.GetTrainingWithExercises(ctx, trainingID)
	if err != nil {
		return nil, ErrTrainingNotFound
	}

	// Помечаем как выполненную
	training.IsDone = true
	now := time.Now().UTC()
	training.Done = &now
	training.Rating = rating

	return s.repo.UpdateTraining(ctx, training)
}