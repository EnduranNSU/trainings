package dto

// GetTrainingsByUserRequest представляет запрос на получение тренировок пользователя
type GetTrainingsByUserRequest struct {
	UserID string `json:"user_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid" description:"UUID пользователя"`
}

// CreateTrainingRequest представляет запрос на создание тренировки
type CreateTrainingRequest struct {
	UserID    string  `json:"user_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid" description:"UUID пользователя"`
	IsDone    bool    `json:"is_done" example:"false" description:"Завершена ли тренировка"`
	Planned   string  `json:"planned" binding:"required" example:"2023-10-05T15:00:00Z" pattern:"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$" description:"Запланированная дата и время тренировки"`
	Done      *string `json:"done,omitempty" example:"2023-10-05T16:30:00Z" pattern:"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$" description:"Фактическая дата и время завершения тренировки (опционально)"`
	TotalTime *string `json:"total_time,omitempty" example:"1h30m" description:"Общее время тренировки в формате duration (опционально)"`
	Rating    *int32  `json:"rating,omitempty" example:"5" minimum:"1" maximum:"5" description:"Оценка тренировки от 1 до 5 (опционально)"`
}

// UpdateTrainingRequest представляет запрос на обновление тренировки
type UpdateTrainingRequest struct {
	ID        int64   `json:"id" binding:"required" example:"1" minimum:"1" description:"ID тренировки"`
	IsDone    *bool   `json:"is_done,omitempty" example:"true" description:"Завершена ли тренировка (опционально)"`
	Planned   string  `json:"planned" binding:"required" example:"2023-10-05T15:00:00Z" pattern:"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$" description:"Запланированная дата и время тренировки"`
	Done      *string `json:"done,omitempty" example:"2023-10-05T16:30:00Z" pattern:"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$" description:"Фактическая дата и время завершения тренировки (опционально)"`
	TotalTime *string `json:"total_time,omitempty" example:"1h30m" description:"Общее время тренировки в формате duration (опционально)"`
	Rating    *int32  `json:"rating,omitempty" example:"5" minimum:"1" maximum:"5" description:"Оценка тренировки от 1 до 5 (опционально)"`
}

// AddExerciseToTrainingRequest представляет запрос на добавление упражнения к тренировке
type AddExerciseToTrainingRequest struct {
	TrainingID int64   `json:"training_id" binding:"required" example:"1" minimum:"1" description:"ID тренировки"`
	ExerciseID int64   `json:"exercise_id" binding:"required" example:"1" minimum:"1" description:"ID упражнения"`
	Weight     *float64 `json:"weight,omitempty" example:"50.5" minimum:"0" maximum:"1000" description:"Вес в килограммах (опционально)"`
	Approaches *int64   `json:"approaches,omitempty" example:"3" minimum:"1" maximum:"20" description:"Количество подходов (опционально)"`
	Reps       *int64   `json:"reps,omitempty" example:"10" minimum:"1" maximum:"100" description:"Количество повторений (опционально)"`
	Time       *string  `json:"time,omitempty" example:"2023-10-05T15:30:00Z" pattern:"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$" description:"Время выполнения упражнения (опционально)"`
	Notes      *string  `json:"notes,omitempty" example:"Тяжело далось" description:"Заметки к упражнению (опционально)"`
}

// UpdateTrainedExerciseRequest представляет запрос на обновление выполненного упражнения
type UpdateTrainedExerciseRequest struct {
	ID         int64   `json:"id" binding:"required" example:"1" minimum:"1" description:"ID выполненного упражнения"`
	Weight     *float64 `json:"weight,omitempty" example:"55.0" minimum:"0" maximum:"1000" description:"Вес в килограммах (опционально)"`
	Approaches *int64   `json:"approaches,omitempty" example:"4" minimum:"1" maximum:"20" description:"Количество подходов (опционально)"`
	Reps       *int64   `json:"reps,omitempty" example:"12" minimum:"1" maximum:"100" description:"Количество повторений (опционально)"`
	Time       *string  `json:"time,omitempty" example:"2023-10-05T15:45:00Z" pattern:"^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$" description:"Время выполнения упражнения (опционально)"`
	Notes      *string  `json:"notes,omitempty" example:"Стало легче" description:"Заметки к упражнению (опционально)"`
}

// CompleteTrainingRequest представляет запрос на завершение тренировки
type CompleteTrainingRequest struct {
	Rating *int32 `json:"rating,omitempty" example:"5" minimum:"1" maximum:"5" description:"Оценка тренировки от 1 до 5 (опционально)"`
}

// TrainingResponse представляет ответ с информацией о тренировке
type TrainingResponse struct {
	ID        int64      `json:"id" example:"1" description:"ID тренировки"`
	UserID    string     `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" description:"UUID пользователя"`
	IsDone    bool       `json:"is_done" example:"true" description:"Завершена ли тренировка"`
	Planned   string     `json:"planned" example:"2023-10-05T15:00:00Z" description:"Запланированная дата и время"`
	Done      *string    `json:"done,omitempty" example:"2023-10-05T16:30:00Z" description:"Фактическая дата и время завершения"`
	TotalTime *string    `json:"total_time,omitempty" example:"1h30m" description:"Общее время тренировки"`
	Rating    *int32     `json:"rating,omitempty" example:"5" description:"Оценка тренировки"`
	Exercises []TrainedExerciseResponse `json:"exercises,omitempty" description:"Упражнения в тренировке"`
}

// TrainedExerciseResponse представляет ответ с информацией о выполненном упражнении
type TrainedExerciseResponse struct {
	ID         int64    `json:"id" example:"1" description:"ID выполненного упражнения"`
	TrainingID int64    `json:"training_id" example:"1" description:"ID тренировки"`
	ExerciseID int64    `json:"exercise_id" example:"1" description:"ID упражнения"`
	Weight     *float64 `json:"weight,omitempty" example:"50.5" description:"Вес в килограммах"`
	Approaches *int64   `json:"approaches,omitempty" example:"3" description:"Количество подходов"`
	Reps       *int64   `json:"reps,omitempty" example:"10" description:"Количество повторений"`
	Time       *string  `json:"time,omitempty" example:"2023-10-05T15:30:00Z" description:"Время выполнения"`
	Notes      *string  `json:"notes,omitempty" example:"Тяжело далось" description:"Заметки"`
}

// TrainingStatsResponse представляет ответ со статистикой тренировок
type TrainingStatsResponse struct {
	TotalTrainings   int64   `json:"total_trainings" example:"15" description:"Общее количество тренировок"`
	CompletedTrainings int64 `json:"completed_trainings" example:"12" description:"Количество завершенных тренировок"`
	AverageRating    float64 `json:"average_rating" example:"4.5" description:"Средний рейтинг тренировок"`
	TotalDuration    string  `json:"total_duration" example:"45h30m" description:"Общее время тренировок"`
	LastTrainingDate *string `json:"last_training_date,omitempty" example:"2023-10-05T16:30:00Z" description:"Дата последней тренировки"`
}