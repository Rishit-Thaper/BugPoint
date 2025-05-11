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

func GetBugs(c *fiber.Ctx) error {
	bugsCollection := db.GetCollection("bugs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := bugsCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	defer cursor.Close(ctx)
	var bugs []models.Bug
	if err := cursor.All(ctx, &bugs); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": bugs})
}

func CreateBug(c *fiber.Ctx) error {
	var bug models.Bug
	if err := c.BodyParser(&bug); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Input"})
	}
	if bug.Title == "" || bug.Description == "" || bug.Status != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Required fields are missing"})
	}
	bugsCollection := db.GetCollection("bugs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bug.ID = primitive.NewObjectID()
	bug.CreatedAt = time.Now()
	bug.UpdatedAt = bug.CreatedAt
	_, err := bugsCollection.InsertOne(ctx, bug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Some error occured while adding bug."})

	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Bug added successfully"})
}

func UpdateBug(c *fiber.Ctx) error {
	bugId, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var updatedBug struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	bugsCollection := db.GetCollection("bugs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.BodyParser(&updatedBug); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Input"})
	}

	update := bson.M{
		"$set": bson.M{
			"title":       updatedBug.Title,
			"description": updatedBug.Description,
			"updated_at":  time.Now(),
		},
	}
	result, err := bugsCollection.UpdateOne(ctx, bson.M{"_id": bugId}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update bug"})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Bug don't exist"})
	}
	var bug models.Bug
	err = bugsCollection.FindOne(ctx, bson.M{"_id": bugId}).Decode(&bug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update bug"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": bug, "message": "Bug Updated successfully!"})
}

func DeleteBug(c *fiber.Ctx) error {
	bugId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	bugsCollection := db.GetCollection("bugs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := bugsCollection.DeleteOne(ctx, bson.M{"_id": bugId})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete bug"})
	}
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bug don't exist"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Bug deleted successfully!"})
}

func GetSingleBug(c *fiber.Ctx) error {
	bugId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	bugsCollection := db.GetCollection("bugs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var bug models.Bug
	err = bugsCollection.FindOne(ctx, bson.M{"_id": bugId}).Decode(&bug)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Bug don't exist"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Bug deleted successfully!"})
}
