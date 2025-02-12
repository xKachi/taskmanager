package model

import (
	"context"
	"errors"
	"fmt"
	"task-manager/database"
	"time"

	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Task struct {
	ID        bson.ObjectID `bson:"_id"`
	Title     string        `bson:"title"`
	Completed bool          `bson:"completed"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func GetAllTasks(filter interface{}) ([]bson.M, error) {
	var tasks []bson.M
	ctx := context.TODO()
	collection := database.GetDBCollection()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	// iterate over cursor
	for cursor.Next(ctx) {
		task := bson.M{}
		err := cursor.Decode(&task)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err := cursor.Err(); err != nil {
		return tasks, err
	}

	// once exhausted, close the cursor
	cursor.Close(ctx)

	if len(tasks) == 0 {
		return tasks, mongo.ErrNoDocuments
	}
	return tasks, nil
}

func FilterTasks(status bool) ([]bson.M, error) {
	// filter completed tasks
	filter := bson.D{{Key: "completed", Value: status}}
	results, err := GetAllTasks(filter)
	if err != nil {
		return results, err
	}
	return results, err
}

func AddTask(task Task) error {
	collection := database.GetDBCollection()
	_, err := collection.InsertOne(context.TODO(), &task)
	if err != nil {
		return err
	}
	return nil
}

func CompleteTask(taskID string) error {

	objectID, err := bson.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	// get collection
	collection := database.GetDBCollection()

	filter := bson.D{{Key: "_id", Value: objectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "completed", Value: true}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func DeleteTask(taskID string) error {
	objectID, err := bson.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	// get database collection
	collection := database.GetDBCollection()

	filter := bson.D{{Key: "_id", Value: objectID}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	} else if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil

}

func PrintTasks(tasks []bson.M) {
	t := table.NewWriter()
	t.SetTitle("All Tasks")

	t.AppendHeader(table.Row{"#ID", "Title", "Completed"})
	for _, tasks := range tasks {
		id := tasks["_id"].(bson.ObjectID).Hex()
		title := tasks["title"]
		completed := tasks["completed"]

		if tasks["completed"] == true {
			id = color.Green.Sprint(id)
			title = color.Green.Sprint(title)
			completed = color.Green.Sprint(completed)
		} else {
			id = color.Yellow.Sprint(id)
			title = color.Yellow.Sprint(title)
			completed = color.Yellow.Sprint(completed)
		}

		// add to table row
		t.AppendRow(table.Row{id, title, completed})
	}

	// render table
	fmt.Println(t.Render())
}