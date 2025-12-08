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

func (s *trainingService) GetUserTrainingStats(ctx context.Context, userID uuid.UUID) (*domain.TrainingStats, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetUserTrainingStats(ctx, userID)
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

	if cmd.PlannedDate.IsZero() {
		return nil, errors.New("planned date is required")
	}

	training := &domain.Training{
		UserID:            cmd.UserID,
		IsDone:            cmd.IsDone,
		PlannedDate:       cmd.PlannedDate,
		ActualDate:        cmd.ActualDate,
		StartedAt:         cmd.StartedAt,
		FinishedAt:        cmd.FinishedAt,
		TotalDuration:     cmd.TotalDuration,
		TotalRestTime:     cmd.TotalRestTime,
		TotalExerciseTime: cmd.TotalExerciseTime,
		Rating:            cmd.Rating,
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
	if !cmd.PlannedDate.IsZero() {
		existing.PlannedDate = cmd.PlannedDate
	}
	if cmd.ActualDate != nil {
		existing.ActualDate = cmd.ActualDate
	}
	if cmd.StartedAt != nil {
		existing.StartedAt = cmd.StartedAt
	}
	if cmd.FinishedAt != nil {
		existing.FinishedAt = cmd.FinishedAt
	}
	if cmd.TotalDuration != nil {
		existing.TotalDuration = cmd.TotalDuration
	}
	if cmd.TotalRestTime != nil {
		existing.TotalRestTime = cmd.TotalRestTime
	}
	if cmd.TotalExerciseTime != nil {
		existing.TotalExerciseTime = cmd.TotalExerciseTime
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
		Doing:      cmd.Doing,
		Rest:       cmd.Rest,
		Notes:      cmd.Notes,
	}

	return s.repo.AddExerciseToTraining(ctx, exercise)
}

func (s *trainingService) UpdateTrainedExercise(ctx context.Context, cmd domain.UpdateTrainedExerciseCmd) (*domain.TrainedExercise, error) {
	if cmd.ID <= 0 {
		return nil, ErrInvalidExerciseID
	}

	exercise := &domain.TrainedExercise{
		ID:         cmd.ID,
		Weight:     cmd.Weight,
		Approaches: cmd.Approaches,
		Reps:       cmd.Reps,
		Time:       cmd.Time,
		Doing:      cmd.Doing,
		Rest:       cmd.Rest,
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

func (s *trainingService) StartTraining(ctx context.Context, trainingID int64) (*domain.Training, error) {
	if trainingID <= 0 {
		return nil, ErrInvalidTrainingID
	}

	training, err := s.repo.GetTrainingWithExercises(ctx, trainingID)
	if err != nil {
		return nil, ErrTrainingNotFound
	}

	if training.StartedAt != nil {
		return training, nil // Уже начата
	}

	now := time.Now().UTC()
	training.StartedAt = &now

	return s.repo.UpdateTraining(ctx, training)
}

func (s *trainingService) CompleteTraining(ctx context.Context, trainingID int64, rating *int32) (*domain.Training, error) {
	if trainingID <= 0 {
		return nil, ErrInvalidTrainingID
	}

	training, err := s.repo.GetTrainingWithExercises(ctx, trainingID)
	if err != nil {
		return nil, ErrTrainingNotFound
	}

	training.IsDone = true
	now := time.Now().UTC()
	training.ActualDate = &now
	training.FinishedAt = &now
	training.Rating = rating

	// Если не было начато, устанавливаем время начала как текущее минус предполагаемая длительность
	if training.StartedAt == nil {
		startedAt := now.Add(-*training.TotalDuration)
		training.StartedAt = &startedAt
	}

	return s.repo.UpdateTraining(ctx, training)
}