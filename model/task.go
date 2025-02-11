package model

import (
    "context"
    "errors"
    "fmt"
    "task-manager/database"
    "time"

    "github.com/gookit/color"
    "github.com/jedib0t/go-pretty/v6/table"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
    ID        primitive.ObjectID `bson:"_id"`
    Title     string             `bson:"title"`
    Completed bool               `bson:"completed"`
    CreatedAt time.Time          `bson:"created_at"`
    UpdatedAt time.Time          `bson:"updated_at"`
}

// Function to get all tasks
func GetAllTasks(filter interface{}) ([]bson.M, error) {
    var tasks []bson.M
    ctx := context.TODO()
    collection := database.GetDBCollection()
    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return tasks, err
    }

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

    cursor.Close(ctx)

    if len(tasks) == 0 {
        return tasks, mongo.ErrNoDocuments
    }
    return tasks, nil
}

// Function to filter tasks
func FilterTasks(status bool) ([]bson.M, error) {
	// filter completed tasks
	filter := bson.D{primitive.E{Key: "completed", Value: status}}
	results, err := GetAllTasks(filter)
	if err != nil {
		return results, err
	}
	return results, err
}

// Function to add a task
func AddTask(task Task) error {
    collection := database.GetDBCollection()
    _, err := collection.InsertOne(context.TODO(), &task)
    if err != nil {
        return err
    }
    return nil
}

// Function to complete a task
func CompleteTask(taskID string) error {
    objectID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        return err
    }

    collection := database.GetDBCollection()
    filter := bson.D{primitive.E{Key: "_id", Value: objectID}}
    update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "completed", Value: true}}}}

    _, err = collection.UpdateOne(context.TODO(), filter, update)
    return err
}

// Function to delete a task
func DeleteTask(taskID string) error {
    objectID, err := primitive.ObjectIDFromHex(taskID)
    if err != nil {
        return err
    }

    collection := database.GetDBCollection()
    filter := bson.D{primitive.E{Key: "_id", Value: objectID}}
    result, err := collection.DeleteOne(context.TODO(), filter)
    if err != nil {
        return err
    } else if result.DeletedCount == 0 {
        return errors.New("task not found")
    }
    return nil
}

// Function to print tasks
func PrintTasks(tasks []bson.M) {
    t := table.NewWriter()
    t.SetTitle("All Tasks")

    t.AppendHeader(table.Row{"#ID", "Title", "Completed"})
    for _, task := range tasks {
        id := task["_id"].(primitive.ObjectID).Hex()
        title := task["title"]
        completed := task["completed"]

        if task["completed"] == true {
            id = color.Green.Sprint(id)
            title = color.Green.Sprint(title)
            completed = color.Green.Sprint(completed)
        } else {
            id = color.Yellow.Sprint(id)
            title = color.Yellow.Sprint(title)
            completed = color.Yellow.Sprint(completed)
        }

        t.AppendRow(table.Row{id, title, completed})
    }

    fmt.Println(t.Render())
}
