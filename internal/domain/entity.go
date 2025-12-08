package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Training struct {
	ID                int64             `db:"id" json:"id"`
	UserID            uuid.UUID         `db:"user_id" json:"user_id"`
	IsDone            bool              `db:"is_done" json:"is_done"`
	PlannedDate       time.Time         `db:"planned_date" json:"planned_date"`
	ActualDate        *time.Time        `db:"actual_date" json:"actual_date"`
	StartedAt         *time.Time        `db:"started_at" json:"started_at"`
	FinishedAt        *time.Time        `db:"finished_at" json:"finished_at"`
	TotalDuration     *time.Duration    `db:"total_duration" json:"total_duration"`
	TotalRestTime     *time.Duration    `db:"total_rest_time" json:"total_rest_time"`
	TotalExerciseTime *time.Duration    `db:"total_exercise_time" json:"total_exercise_time"`
	Rating            *int32            `db:"rating" json:"rating"`
	Exercises         []TrainedExercise `db:"exercises" json:"exercises"`
}

type TrainedExercise struct {
	ID         int64            `db:"id" json:"id"`
	TrainingID int64            `db:"training_id" json:"training_id"`
	ExerciseID int64            `db:"exercise_id" json:"exercise_id"`
	Weight     *decimal.Decimal `db:"weight" json:"weight"`
	Approaches *int32           `db:"approaches" json:"approaches"`
	Reps       *int32           `db:"reps" json:"reps"`
	Time       *time.Duration   `db:"time" json:"time"`
	Doing      *time.Duration   `db:"doing" json:"doing"`
	Rest       *time.Duration   `db:"rest" json:"rest"`
	Notes      *string          `db:"notes" json:"notes"`
}

type Exercise struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
	Href        string `db:"href" json:"href"`
	Tags        []Tag  `db:"tags" json:"tags"`
}

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Type string `db:"type" json:"type"`
}

type ExerciseFilter struct {
	TagID  *int64
	Search *string
}
