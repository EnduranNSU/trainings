package dto

// CreateTrainingRequest представляет запрос на создание тренировки
type CreateTrainingRequest struct {
	UserID            string  `json:"user_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid" description:"UUID пользователя"`
	IsDone            bool    `json:"is_done" example:"false" description:"Завершена ли тренировка"`
	PlannedDate       string  `json:"planned_date" binding:"required" example:"2023-10-05T15:00:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Запланированная дата и время тренировки"`
	ActualDate        *string `json:"actual_date,omitempty" example:"2023-10-05T16:30:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Фактическая дата и время выполнения тренировки (опционально)"`
	StartedAt         *string `json:"started_at,omitempty" example:"2023-10-05T15:00:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Время начала тренировки (опционально)"`
	FinishedAt        *string `json:"finished_at,omitempty" example:"2023-10-05T16:30:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Время окончания тренировки (опционально)"`
	TotalDuration     *string `json:"total_duration,omitempty" example:"1h30m" description:"Общее время тренировки (опционально)"`
	TotalRestTime     *string `json:"total_rest_time,omitempty" example:"30m" description:"Общее время отдыха (опционально)"`
	TotalExerciseTime *string `json:"total_exercise_time,omitempty" example:"1h" description:"Общее время выполнения упражнений (опционально)"`
	Rating            *int32  `json:"rating,omitempty" example:"5" minimum:"1" maximum:"5" description:"Оценка тренировки от 1 до 5 (опционально)"`
}

// UpdateTrainingRequest представляет запрос на обновление тренировки
type UpdateTrainingRequest struct {
	IsDone            *bool   `json:"is_done,omitempty" example:"true" description:"Завершена ли тренировка (опционально)"`
	PlannedDate       string  `json:"planned_date" binding:"required" example:"2023-10-05T15:00:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Запланированная дата и время тренировки"`
	ActualDate        *string `json:"actual_date,omitempty" example:"2023-10-05T16:30:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Фактическая дата и время выполнения тренировки (опционально)"`
	StartedAt         *string `json:"started_at,omitempty" example:"2023-10-05T15:00:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Время начала тренировки (опционально)"`
	FinishedAt        *string `json:"finished_at,omitempty" example:"2023-10-05T16:30:00Z" pattern:"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$" description:"Время окончания тренировки (опционально)"`
	TotalDuration     *string `json:"total_duration,omitempty" example:"1h30m" description:"Общее время тренировки (опционально)"`
	TotalRestTime     *string `json:"total_rest_time,omitempty" example:"30m" description:"Общее время отдыха (опционально)"`
	TotalExerciseTime *string `json:"total_exercise_time,omitempty" example:"1h" description:"Общее время выполнения упражнений (опционально)"`
	Rating            *int32  `json:"rating,omitempty" example:"5" minimum:"1" maximum:"5" description:"Оценка тренировки от 1 до 5 (опционально)"`
}

// AddExerciseToTrainingRequest представляет запрос на добавление упражнения к тренировке
type AddExerciseToTrainingRequest struct {
	TrainingID int64    `json:"training_id" binding:"required" example:"1" minimum:"1" description:"ID тренировки"`
	ExerciseID int64    `json:"exercise_id" binding:"required" example:"1" minimum:"1" description:"ID упражнения"`
	Weight     *float64 `json:"weight,omitempty" example:"50.5" minimum:"0" maximum:"1000" description:"Вес в килограммах (опционально)"`
	Approaches *int64   `json:"approaches,omitempty" example:"3" minimum:"1" maximum:"20" description:"Количество подходов (опционально)"`
	Reps       *int64   `json:"reps,omitempty" example:"10" minimum:"1" maximum:"100" description:"Количество повторений (опционально)"`
	Time       *string  `json:"time,omitempty" example:"1h30m" description:"Общее время упражнения в формате duration (опционально)"`
	Doing      *string  `json:"doing,omitempty" example:"1h" description:"Время выполнения упражнения в формате duration (опционально)"`
	Rest       *string  `json:"rest,omitempty" example:"30m" description:"Время отдыха в формате duration (опционально)"`
	Notes      *string  `json:"notes,omitempty" example:"Тяжело далось" description:"Заметки к упражнению (опционально)"`
}

// UpdateTrainedExerciseRequest представляет запрос на обновление выполненного упражнения
type UpdateTrainedExerciseRequest struct {
	Weight     *float64 `json:"weight,omitempty" example:"55.0" minimum:"0" maximum:"1000" description:"Вес в килограммах (опционально)"`
	Approaches *int64   `json:"approaches,omitempty" example:"4" minimum:"1" maximum:"20" description:"Количество подходов (опционально)"`
	Reps       *int64   `json:"reps,omitempty" example:"12" minimum:"1" maximum:"100" description:"Количество повторений (опционально)"`
	Time       *string  `json:"time,omitempty" example:"1h45m" description:"Общее время упражнения в формате duration (опционально)"`
	Doing      *string  `json:"doing,omitempty" example:"1h15m" description:"Время выполнения упражнения в формате duration (опционально)"`
	Rest       *string  `json:"rest,omitempty" example:"30m" description:"Время отдыха в формате duration (опционально)"`
	Notes      *string  `json:"notes,omitempty" example:"Стало легче" description:"Заметки к упражнению (опционально)"`
}

// TrainingResponse представляет ответ с информацией о тренировке
type TrainingResponse struct {
	ID                int64                     `json:"id" example:"1" description:"ID тренировки"`
	UserID            string                    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" description:"UUID пользователя"`
	IsDone            bool                      `json:"is_done" example:"true" description:"Завершена ли тренировка"`
	PlannedDate       string                    `json:"planned_date" example:"2023-10-05T15:00:00Z" description:"Запланированная дата и время"`
	ActualDate        *string                   `json:"actual_date,omitempty" example:"2023-10-05T16:30:00Z" description:"Фактическая дата и время выполнения"`
	StartedAt         *string                   `json:"started_at,omitempty" example:"2023-10-05T15:00:00Z" description:"Время начала тренировки"`
	FinishedAt        *string                   `json:"finished_at,omitempty" example:"2023-10-05T16:30:00Z" description:"Время окончания тренировки"`
	TotalDuration     *string                   `json:"total_duration,omitempty" example:"1h30m" description:"Общее время тренировки"`
	TotalRestTime     *string                   `json:"total_rest_time,omitempty" example:"30m" description:"Общее время отдыха"`
	TotalExerciseTime *string                   `json:"total_exercise_time,omitempty" example:"1h" description:"Общее время выполнения упражнений"`
	Rating            *int32                    `json:"rating,omitempty" example:"5" description:"Оценка тренировки"`
	Exercises         []TrainedExerciseResponse `json:"exercises,omitempty" description:"Упражнения в тренировке"`
}

// TrainedExerciseResponse представляет ответ с информацией о выполненном упражнении
type TrainedExerciseResponse struct {
	ID         int64    `json:"id" example:"1" description:"ID выполненного упражнения"`
	TrainingID int64    `json:"training_id" example:"1" description:"ID тренировки"`
	ExerciseID int64    `json:"exercise_id" example:"1" description:"ID упражнения"`
	Weight     *float64 `json:"weight,omitempty" example:"50.5" description:"Вес в килограммах"`
	Approaches *int32   `json:"approaches,omitempty" example:"3" description:"Количество подходов"`
	Reps       *int32   `json:"reps,omitempty" example:"10" description:"Количество повторений"`
	Time       *string  `json:"time,omitempty" example:"1h30m" description:"Общее время упражнения"`
	Doing      *string  `json:"doing,omitempty" example:"1h" description:"Время выполнения упражнения"`
	Rest       *string  `json:"rest,omitempty" example:"30m" description:"Время отдыха"`
	Notes      *string  `json:"notes,omitempty" example:"Тяжело далось" description:"Заметки"`
}

// TrainingStatsResponse представляет ответ со статистикой тренировок
type TrainingStatsResponse struct {
	TotalTrainings     int64   `json:"total_trainings" example:"15" description:"Общее количество тренировок"`
	CompletedTrainings int64   `json:"completed_trainings" example:"12" description:"Количество завершенных тренировок"`
	AverageRating      float64 `json:"average_rating" example:"4.5" description:"Средний рейтинг тренировок"`
	TotalDuration      string  `json:"total_duration" example:"45h30m" description:"Общее время тренировок"`
	LastTrainingDate   *string `json:"last_training_date,omitempty" example:"2023-10-05T16:30:00Z" description:"Дата последней тренировки"`
}

// CompleteTrainingRequest представляет запрос на завершение тренировки
type CompleteTrainingRequest struct {
	Rating *int32 `json:"rating,omitempty" example:"5" minimum:"1" maximum:"5" description:"Оценка тренировки от 1 до 5 (опционально)"`
}
