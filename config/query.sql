-- name: GetExercisesWithTags :many
SELECT 
    e.id,
    e.description,
    e.href,
    e.tags as exercise_tags,  -- Rename this column
    COALESCE(
        json_agg(
            json_build_object(
                'id', t.id,
                'type', t.type
            )
        ) FILTER (WHERE t.id IS NOT NULL),
        '[]'
    ) as tags  -- This stays as 'tags' for the JSON result
FROM exercise e
LEFT JOIN exercise_to_tag et ON e.id = et.execrise_id
LEFT JOIN tag t ON et.tag_id = t.id
GROUP BY e.id, e.description, e.href, e.tags
ORDER BY e.id;

-- name: GetExerciseByID :one
SELECT 
    e.id,
    e.description,
    e.href,
    e.tags as exercise_tags,  -- Rename this column
    COALESCE(
        json_agg(
            json_build_object(
                'id', t.id,
                'type', t.type
            )
        ) FILTER (WHERE t.id IS NOT NULL),
        '[]'
    ) as tags  -- This stays as 'tags' for the JSON result
FROM exercise e
LEFT JOIN exercise_to_tag et ON e.id = et.execrise_id
LEFT JOIN tag t ON et.tag_id = t.id
WHERE e.id = $1
GROUP BY e.id, e.description, e.href, e.tags;

-- name: GetAllTags :many
SELECT id, type FROM tag ORDER BY id;

-- name: GetExercisesByTag :many
SELECT 
    e.id,
    e.description,
    e.href,
    e.tags
FROM exercise e
INNER JOIN exercise_to_tag et ON e.id = et.execrise_id
WHERE et.tag_id = $1
ORDER BY e.id;

-- name: GetTrainingsByUser :many
SELECT 
    t.id,
    t.user_id,
    t.isDone,
    t.planned,
    t.done,
    t.total_time,
    t.rating,
    COALESCE(
        json_agg(
            json_build_object(
                'id', te.id,
                'exercise_id', te.exercise_id,
                'weight', te.weight,
                'approaches', te.approaches,
                'reps', te.reps,
                'time', te.time,
                'notes', te.notes
            )
        ) FILTER (WHERE te.id IS NOT NULL),
        '[]'
    ) as exercises
FROM training t
LEFT JOIN trained_exercise te ON t.id = te.training_id
WHERE t.user_id = $1
GROUP BY t.id
ORDER BY t.planned DESC;

-- name: CreateTraining :one
INSERT INTO training (
    user_id,
    isDone,
    planned,
    done,
    total_time,
    rating
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: AddExerciseToTraining :one
INSERT INTO trained_exercise (
    training_id,
    exercise_id,
    weight,
    approaches,
    reps,
    time,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateTrainedExercise :one
UPDATE trained_exercise
SET 
    weight = COALESCE($1, weight),
    approaches = COALESCE($2, approaches),
    reps = COALESCE($3, reps),
    time = COALESCE($4, time),
    notes = COALESCE($5, notes)
WHERE id = $6
RETURNING *;

-- name: UpdateTraining :one
UPDATE training
SET 
    isDone = COALESCE($1, isDone),
    planned = COALESCE($2, planned),
    done = COALESCE($3, done),
    total_time = COALESCE($4, total_time),
    rating = COALESCE($5, rating)
WHERE id = $6
RETURNING *;

-- name: GetTrainingWithExercises :one
SELECT 
    t.*,
    COALESCE(
        json_agg(
            json_build_object(
                'id', te.id,
                'exercise_id', te.exercise_id,
                'weight', te.weight,
                'approaches', te.approaches,
                'reps', te.reps,
                'time', te.time,
                'notes', te.notes
            )
        ) FILTER (WHERE te.id IS NOT NULL),
        '[]'
    ) as exercises
FROM training t
LEFT JOIN trained_exercise te ON t.id = te.training_id
WHERE t.id = $1
GROUP BY t.id;

-- name: DeleteExerciseFromTraining :exec
DELETE FROM trained_exercise 
WHERE id = $1 AND training_id = $2;

-- name: DeleteTrainingAndExercises :exec
WITH deleted_exercises AS (
    DELETE FROM trained_exercise WHERE training_id = $1
)
DELETE FROM training WHERE training.id = $1;