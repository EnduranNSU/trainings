package postgres

import (
	"context"
	"database/sql"

	"github.com/EnduranNSU/trainings/internal/adapter/out/postgres/gen"
	"github.com/EnduranNSU/trainings/internal/domain"
	"github.com/EnduranNSU/trainings/internal/logging"
	"github.com/shopspring/decimal"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type TrainingRepositoryImpl struct {
	q  *gen.Queries
	db *sql.DB
}

func NewTrainingRepository(db *sql.DB) domain.TrainingRepository {
	return &TrainingRepositoryImpl{
		q:  gen.New(db),
		db: db,
	}
}

func (r *TrainingRepositoryImpl) GetTrainingsByUser(ctx context.Context, userID uuid.UUID) ([]*domain.Training, error) {
	trainings, err := r.q.GetTrainingsByUser(ctx, userID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"user_id": userID.String(),
		})
		logging.Error(err, "GetTrainingsByUser", jsonData, "failed to get trainings by user")
		return nil, err
	}

	result := make([]*domain.Training, len(trainings))
	for i, t := range trainings {
		result[i] = r.toDomainTraining(t)
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"user_id":         userID.String(),
		"trainings_count": len(result),
	})
	logging.Debug("GetTrainingsByUser", jsonData, "successfully retrieved user trainings")

	return result, nil
}

func (r *TrainingRepositoryImpl) GetTrainingWithExercises(ctx context.Context, trainingID int64) (*domain.Training, error) {
	training, err := r.q.GetTrainingWithExercises(ctx, trainingID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": trainingID,
		})
		logging.Error(err, "GetTrainingWithExercises", jsonData, "failed to get training with exercises")
		return nil, err
	}

	domainTraining := r.toDomainTrainingFromJoined(training)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id":     trainingID,
		"exercises_count": len(domainTraining.Exercises),
	})
	logging.Debug("GetTrainingWithExercises", jsonData, "successfully retrieved training with exercises")

	return domainTraining, nil
}

func (r *TrainingRepositoryImpl) CreateTraining(ctx context.Context, training *domain.Training) (*domain.Training, error) {
	params := gen.CreateTrainingParams{
		UserID:            training.UserID,
		IsDone:            training.IsDone,
		PlannedDate:       training.PlannedDate,
		ActualDate:        null.TimeFromPtr(training.ActualDate).NullTime,
		StartedAt:         null.TimeFromPtr(training.StartedAt).NullTime,
		FinishedAt:        null.TimeFromPtr(training.FinishedAt).NullTime,
		TotalDuration:     durationToNullInt64(training.TotalDuration),
		TotalRestTime:     durationToNullInt64(training.TotalRestTime),
		TotalExerciseTime: durationToNullInt64(training.TotalExerciseTime),
		Rating:            null.Int32FromPtr(training.Rating).NullInt32,
	}

	created, err := r.q.CreateTraining(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"user_id": training.UserID.String(),
			"planned": training.PlannedDate,
			"is_done": training.IsDone,
		})
		logging.Error(err, "CreateTraining", jsonData, "failed to create training")
		return nil, err
	}

	domainTraining := r.toDomainTrainingFromGen(created)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id": domainTraining.ID,
		"user_id":     domainTraining.UserID.String(),
	})
	logging.Debug("CreateTraining", jsonData, "successfully created training")

	return domainTraining, nil
}

func (r *TrainingRepositoryImpl) UpdateTraining(ctx context.Context, training *domain.Training) (*domain.Training, error) {
	params := gen.UpdateTrainingParams{
		IsDone:            training.IsDone,
		PlannedDate:       training.PlannedDate,
		ActualDate:        null.TimeFromPtr(training.ActualDate).NullTime,
		StartedAt:         null.TimeFromPtr(training.StartedAt).NullTime,
		FinishedAt:        null.TimeFromPtr(training.FinishedAt).NullTime,
		TotalDuration:     durationToNullInt64(training.TotalDuration),
		TotalRestTime:     durationToNullInt64(training.TotalRestTime),
		TotalExerciseTime: durationToNullInt64(training.TotalExerciseTime),
		Rating:            null.Int32FromPtr(training.Rating).NullInt32,
		ID:                training.ID,
	}

	updated, err := r.q.UpdateTraining(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": training.ID,
			"is_done":     training.IsDone,
			"rating":      training.Rating,
		})
		logging.Error(err, "UpdateTraining", jsonData, "failed to update training")
		return nil, err
	}

	domainTraining := r.toDomainTrainingFromGen(updated)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id": domainTraining.ID,
	})
	logging.Debug("UpdateTraining", jsonData, "successfully updated training")

	return domainTraining, nil
}

func (r *TrainingRepositoryImpl) DeleteTrainingAndExercises(ctx context.Context, trainingID int64) error {
	err := r.q.DeleteTrainingAndExercises(ctx, trainingID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": trainingID,
		})
		logging.Error(err, "DeleteTrainingAndExercises", jsonData, "failed to delete training and exercises")
		return err
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"training_id": trainingID,
	})
	logging.Debug("DeleteTrainingAndExercises", jsonData, "successfully deleted training and exercises")

	return nil
}

func (r *TrainingRepositoryImpl) AddExerciseToTraining(ctx context.Context, exercise *domain.TrainedExercise) (*domain.TrainedExercise, error) {
	weight := exercise.Weight.String()
	params := gen.AddExerciseToTrainingParams{
		TrainingID: exercise.TrainingID,
		ExerciseID: exercise.ExerciseID,
		Weight:     null.StringFromPtr(&weight).NullString,
		Approaches: null.Int32FromPtr(exercise.Approaches).NullInt32,
		Reps:       null.Int32FromPtr(exercise.Reps).NullInt32,
		Time:       durationToNullInt64(exercise.Time),
		Doing:      durationToNullInt64(exercise.Doing),
		Rest:       durationToNullInt64(exercise.Rest),
		Notes:      null.StringFromPtr(exercise.Notes).NullString,
	}
	created, err := r.q.AddExerciseToTraining(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"training_id": exercise.TrainingID,
			"exercise_id": exercise.ExerciseID,
		})
		logging.Error(err, "AddExerciseToTraining", jsonData, "failed to add exercise to training")
		return nil, err
	}

	domainExercise := r.toDomainTrainedExercise(created)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"trained_exercise_id": domainExercise.ID,
		"training_id":         domainExercise.TrainingID,
		"exercise_id":         domainExercise.ExerciseID,
	})
	logging.Debug("AddExerciseToTraining", jsonData, "successfully added exercise to training")

	return domainExercise, nil
}

func (r *TrainingRepositoryImpl) UpdateTrainedExercise(ctx context.Context, exercise *domain.TrainedExercise) (*domain.TrainedExercise, error) {
	weight := exercise.Weight.String()
	params := gen.UpdateTrainedExerciseParams{
		Weight:     null.StringFromPtr(&weight).NullString,
		Approaches: null.Int32FromPtr(exercise.Approaches).NullInt32,
		Reps:       null.Int32FromPtr(exercise.Reps).NullInt32,
		Time:       durationToNullInt64(exercise.Time),
		Doing:      durationToNullInt64(exercise.Doing),
		Rest:       durationToNullInt64(exercise.Rest),
		Notes:      null.StringFromPtr(exercise.Notes).NullString,
		ID:         exercise.ID,
	}

	updated, err := r.q.UpdateTrainedExercise(ctx, params)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"trained_exercise_id": exercise.ID,
			"training_id":         exercise.TrainingID,
		})
		logging.Error(err, "UpdateTrainedExercise", jsonData, "failed to update trained exercise")
		return nil, err
	}

	domainExercise := r.toDomainTrainedExercise(updated)

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"trained_exercise_id": domainExercise.ID,
	})
	logging.Debug("UpdateTrainedExercise", jsonData, "successfully updated trained exercise")

	return domainExercise, nil
}

func (r *TrainingRepositoryImpl) DeleteExerciseFromTraining(ctx context.Context, exerciseID, trainingID int64) error {
	err := r.q.DeleteExerciseFromTraining(ctx, gen.DeleteExerciseFromTrainingParams{
		ID:         exerciseID,
		TrainingID: trainingID,
	})
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"trained_exercise_id": exerciseID,
			"training_id":         trainingID,
		})
		logging.Error(err, "DeleteExerciseFromTraining", jsonData, "failed to delete exercise from training")
		return err
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"trained_exercise_id": exerciseID,
		"training_id":         trainingID,
	})
	logging.Debug("DeleteExerciseFromTraining", jsonData, "successfully deleted exercise from training")

	return nil
}

func (r *TrainingRepositoryImpl) GetUserTrainingStats(ctx context.Context, userID uuid.UUID) (*domain.TrainingStats, error) {
	trainings, err := r.GetTrainingsByUser(ctx, userID)
	if err != nil {
		jsonData := logging.MarshalLogData(map[string]interface{}{
			"user_id": userID.String(),
		})
		logging.Error(err, "GetUserTrainingStats", jsonData, "failed to get user training stats")
		return nil, err
	}

	stats := &domain.TrainingStats{}
	var totalRating int64
	var completedCount int

	for _, t := range trainings {
		stats.TotalTrainings++
		if t.IsDone {
			completedCount++
			if t.Rating != nil {
				totalRating += int64(*t.Rating)
			}
		}
	}

	stats.CompletedTrainings = int64(completedCount)
	if completedCount > 0 {
		stats.AverageRating = float64(totalRating) / float64(completedCount)
	}

	jsonData := logging.MarshalLogData(map[string]interface{}{
		"user_id":             userID.String(),
		"total_trainings":     stats.TotalTrainings,
		"completed_trainings": stats.CompletedTrainings,
		"average_rating":      stats.AverageRating,
	})
	logging.Debug("GetUserTrainingStats", jsonData, "successfully calculated user training stats")

	return stats, nil
}

// Вспомогательные методы для преобразования данных (остаются без изменений)
func (r *TrainingRepositoryImpl) toDomainTraining(t gen.GetTrainingsByUserRow) *domain.Training {
	training := &domain.Training{
		ID:                t.ID,
		UserID:            t.UserID,
		IsDone:            t.IsDone,
		PlannedDate:       t.PlannedDate,
		ActualDate:        nullTimeFromSQL(t.ActualDate),
		StartedAt:         nullTimeFromSQL(t.StartedAt),
		FinishedAt:        nullTimeFromSQL(t.FinishedAt),
		TotalDuration:     sqlNullInt64ToDuration(t.TotalDuration),
		TotalRestTime:     sqlNullInt64ToDuration(t.TotalRestTime),
		TotalExerciseTime: sqlNullInt64ToDuration(t.TotalExerciseTime),
		Rating:            nullIntFromSQL32(t.Rating),
	}

	// Преобразование упражнений
	if t.Exercises != nil {
		if exercisesSlice, ok := t.Exercises.([]gen.TrainedExercise); ok {
			exercises := make([]domain.TrainedExercise, len(exercisesSlice))
			for i, ex := range exercisesSlice {
				weight, _ := decimal.NewFromString(ex.Weight.String)
				exercises[i] = domain.TrainedExercise{
					ID:         ex.ID,
					TrainingID: training.ID,
					ExerciseID: ex.ExerciseID,
					Weight:     &weight,
					Approaches: nullIntFromSQL32(ex.Approaches),
					Reps:       nullIntFromSQL32(ex.Reps),
					Time:       sqlNullInt64ToDuration(ex.Time),
					Doing:      sqlNullInt64ToDuration(ex.Doing),
					Rest:       sqlNullInt64ToDuration(ex.Rest),
					Notes:      nullStringFromSQL(ex.Notes),
				}
			}
			training.Exercises = exercises
		}
	}

	return training
}

func (r *TrainingRepositoryImpl) toDomainTrainingFromJoined(t gen.GetTrainingWithExercisesRow) *domain.Training {
	training := &domain.Training{
		ID:                t.ID,
		UserID:            t.UserID,
		IsDone:            t.IsDone,
		PlannedDate:       t.PlannedDate,
		ActualDate:        nullTimeFromSQL(t.ActualDate),
		StartedAt:         nullTimeFromSQL(t.StartedAt),
		FinishedAt:        nullTimeFromSQL(t.FinishedAt),
		TotalDuration:     sqlNullInt64ToDuration(t.TotalDuration),
		TotalRestTime:     sqlNullInt64ToDuration(t.TotalRestTime),
		TotalExerciseTime: sqlNullInt64ToDuration(t.TotalExerciseTime),
		Rating:            nullIntFromSQL32(t.Rating),
	}

	// Преобразование упражнений
	if t.Exercises != nil {
		if exercisesSlice, ok := t.Exercises.([]gen.TrainedExercise); ok {
			exercises := make([]domain.TrainedExercise, len(exercisesSlice))
			for i, ex := range exercisesSlice {
				weight, _ := decimal.NewFromString(ex.Weight.String)
				exercises[i] = domain.TrainedExercise{
					ID:         ex.ID,
					TrainingID: training.ID,
					ExerciseID: ex.ExerciseID,
					Weight:     &weight,
					Approaches: nullIntFromSQL32(ex.Approaches),
					Reps:       nullIntFromSQL32(ex.Reps),
					Time:       sqlNullInt64ToDuration(ex.Time),
					Doing:      sqlNullInt64ToDuration(ex.Doing),
					Rest:       sqlNullInt64ToDuration(ex.Rest),
					Notes:      nullStringFromSQL(ex.Notes),
				}
			}
			training.Exercises = exercises
		}
	}

	return training
}

func (r *TrainingRepositoryImpl) toDomainTrainingFromGen(t gen.Training) *domain.Training {
	return &domain.Training{
		ID:                t.ID,
		UserID:            t.UserID,
		IsDone:            t.IsDone,
		PlannedDate:       t.PlannedDate,
		ActualDate:        nullTimeFromSQL(t.ActualDate),
		StartedAt:         nullTimeFromSQL(t.StartedAt),
		FinishedAt:        nullTimeFromSQL(t.FinishedAt),
		TotalDuration:     sqlNullInt64ToDuration(t.TotalDuration),
		TotalRestTime:     sqlNullInt64ToDuration(t.TotalRestTime),
		TotalExerciseTime: sqlNullInt64ToDuration(t.TotalExerciseTime),
		Rating:            nullIntFromSQL32(t.Rating),
	}
}

func (r *TrainingRepositoryImpl) toDomainTrainedExercise(ex gen.TrainedExercise) *domain.TrainedExercise {
	weight, _ := decimal.NewFromString(ex.Weight.String)
	return &domain.TrainedExercise{
		ID:         ex.ID,
		TrainingID: ex.TrainingID,
		ExerciseID: ex.ExerciseID,
		Weight:     &weight,
		Approaches: nullIntFromSQL32(ex.Approaches),
		Reps:       nullIntFromSQL32(ex.Reps),
		Time:       sqlNullInt64ToDuration(ex.Time),
		Doing:      sqlNullInt64ToDuration(ex.Doing),
		Rest:       sqlNullInt64ToDuration(ex.Rest),
		Notes:      nullStringFromSQL(ex.Notes),
	}
}
