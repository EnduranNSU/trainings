package domain

import (
	"time"

	"github.com/google/uuid"
)

type Training struct {
	ID        int64         `db:"id" json:"id"`
	UserID    uuid.UUID     `db:"user_id" json:"user_id"`
	IsDone    bool          `db:"is_done" json:"is_done"`
	Planned   time.Time     `db:"planned" json:"planned"`
	Done      *time.Time    `db:"done" json:"done"`
	TotalTime *time.Duration `db:"total_time" json:"total_time"`
	Rating    *int32          `db:"rating" json:"rating"`
	Exercises []TrainedExercise `db:"exercises" json:"exercises"`
}

type TrainedExercise struct {
	ID         int64   `db:"id" json:"id"`
	TrainingID int64   `db:"training_id" json:"training_id"`
	ExerciseID int64   `db:"exercise_id" json:"exercise_id"`
	Weight     *float64 `db:"weight" json:"weight"`
	Approaches *int64   `db:"approaches" json:"approaches"`
	Reps       *int64   `db:"reps" json:"reps"`
	Time       *time.Time  `db:"time" json:"time"`
	Notes      *string  `db:"notes" json:"notes"`
}

type Exercise struct {
	ID          int64     `db:"id" json:"id"`
	Description string    `db:"description" json:"description"`
	Href        string    `db:"href" json:"href"`
	Tags        []Tag     `db:"tags" json:"tags"`
}

type Tag struct {
	ID   int64  `db:"id" json:"id"`
	Type string `db:"type" json:"type"`
}

type ExerciseFilter struct {
	TagID *int64
	Search *string
}