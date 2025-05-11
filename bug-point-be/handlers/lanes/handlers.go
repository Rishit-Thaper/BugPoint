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

var LanesCollection = db.GetCollection("lane")

func GetLanes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := LanesCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	defer cursor.Close(ctx)
	var lanes []models.Lane
	if err := cursor.All(ctx, &lanes); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": lanes})
}

func CreateLane(c *fiber.Ctx) error {
	var lane models.Lane
	if err := c.BodyParser(&lane); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if lane.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Required fields are missing"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lane.ID = primitive.NewObjectID()
	lane.CreatedAt = time.Now()
	lane.UpdatedAt = lane.CreatedAt
	_, err := LanesCollection.InsertOne(ctx, lane)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Some error occured while adding Lane."})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Lane added successfully"})
}

func Updatelane(c *fiber.Ctx) error {
	laneId, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var updatedBug struct {
		Title string `json:"title"`
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.BodyParser(&updatedBug); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Input"})
	}

	updateFields := bson.M{}
	if updatedBug.Title != "" {
		updateFields["title"] = updatedBug.Title
	}
	updateFields["updated_at"] = time.Now()

	if len(updateFields) == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No fields to update"})
	}

	update := bson.M{
		"$set": updateFields,
	}
	result, err := LanesCollection.UpdateOne(ctx, bson.M{"_id": laneId}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Lane"})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Lane don't exist"})
	}
	var lane models.Lane
	err = LanesCollection.FindOne(ctx, bson.M{"_id": laneId}).Decode(&lane)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update Lane"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": lane, "message": "Lane Updated successfully!"})
}

func DeleteLane(c *fiber.Ctx) error {
	laneId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := LanesCollection.DeleteOne(ctx, bson.M{"_id": laneId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete lane"})
	}
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Lane don't exist"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Lane deleted successfully!"})
}
