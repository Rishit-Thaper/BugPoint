package handlers

import (
	"bug-point-be/db"
	"bug-point-be/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var TasksCollection = db.GetCollection("tasks")

func GetTasks(c *fiber.Ctx) error {
	laneId, err := primitive.ObjectIDFromHex(c.Params("laneId"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid Lane"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := TasksCollection.Find(ctx, bson.M{"laneId": laneId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	defer cursor.Close(ctx)

	var tasks []models.Task

	if err := cursor.All(ctx, &tasks); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": tasks})
}

func CreateTask(c *fiber.Ctx) error {
	laneId, err := primitive.ObjectIDFromHex(c.Params("laneId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Lane"})
	}

	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if task.Title == "" || task.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Required fields are missing"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	task.ID = primitive.NewObjectID()
	task.LaneId = laneId
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt

	_, insertErr := TasksCollection.InsertOne(ctx, task)

	if insertErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Some error occured while adding task."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Task added successfully"})
}

func UpdateTask(c *fiber.Ctx) error {
	taskId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var updatedTask struct {
		Title       string             `json:"title"`
		Description string             `json:"description"`
		LaneId      primitive.ObjectID `json:"laneId"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.BodyParser(&updatedTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Input"})
	}

	updateFields := bson.M{}
	if updatedTask.Title != "" {
		updateFields["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		updateFields["description"] = updatedTask.Description
	}
	if !updatedTask.LaneId.IsZero() {
		updateFields["laneId"] = updatedTask.LaneId
	}
	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := TasksCollection.UpdateOne(ctx, bson.M{"_id": taskId}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Task don't exist"})
	}

	var task models.Task

	err = TasksCollection.FindOne(ctx, bson.M{"_id": taskId}).Decode(&task)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": task, "message": "Task Updated successfully!"})
}

func DeleteTask(c *fiber.Ctx) error {
	taskId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := TasksCollection.DeleteOne(ctx, bson.M{"_id": taskId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
	}
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task don't exist"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task deleted successfully!"})
}

func GetSingleTask(c *fiber.Ctx) error {
	taskId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task models.Task
	err = TasksCollection.FindOne(ctx, bson.M{"_id": taskId}).Decode(&task)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task don't exist"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": task, "message": "Task fetched successfully!"})
}
