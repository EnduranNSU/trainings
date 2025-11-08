package httpin

import (
	"github.com/gin-gonic/gin"

	_ "github.com/EnduranNSU/trainings/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewGinRouter(training *TrainingHandler, exercise *ExerciseHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.StaticFile("/openapi.yaml", "docs/swagger.yaml")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Training routes
		trainings := api.Group("/trainings/training")
		{
			trainings.GET("", training.GetTrainingsByUser)
			trainings.POST("", training.CreateTraining)
			trainings.GET("/stats", training.GetUserTrainingStats)
			trainings.GET("/:id", training.GetTrainingWithExercises)
			trainings.PUT("/:id", training.UpdateTraining)
			trainings.DELETE("/:id", training.DeleteTraining)
			trainings.PATCH("/:id/complete", training.CompleteTraining)
		}

		// Training exercises routes
		trainingExercises := api.Group("/trainings/training-exercises")
		{
			trainingExercises.POST("", training.AddExerciseToTraining)
			trainingExercises.PUT("/:id", training.UpdateTrainedExercise)
			trainingExercises.DELETE("", training.RemoveExerciseFromTraining)
		}

		// Exercise routes
		exercises := api.Group("/trainings/exercises")
		{
			exercises.GET("", exercise.GetAllExercises)
			exercises.GET("/search", exercise.SearchExercises)
			exercises.POST("/by-tags", exercise.GetExercisesByMultipleTags)
			exercises.GET("/:id", exercise.GetExerciseByID)
			exercises.GET("/tag/:tag_id", exercise.GetExercisesByTag)
			exercises.GET("/:exercise_id/tags", exercise.GetExerciseTags)
		}

		// Tag routes
		tags := api.Group("/trainings/tags")
		{
			tags.GET("", exercise.GetAllTags)
			tags.GET("/popular", exercise.GetPopularTags)
			tags.GET("/:id", exercise.GetTagByID)
		}
	}

	return r
}