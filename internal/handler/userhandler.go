package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ty16akin/ConfirmNG/internal/database"
	"github.com/ty16akin/ConfirmNG/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{}

func (u *User) GetUsers(w http.ResponseWriter, r *http.Request) {
	// cursorStr := r.URL.Query().Get("cursor")
	// if cursorStr == ""{
	// 	cursorStr = "0"
	// }

	// const decimal = 10
	// const bitSize = 64
	// cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	cursor, err := database.Users.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "err: unable to get users", http.StatusInternalServerError)
		return
	}

	var users []model.User
	if err = cursor.All(context.Background(), &users); err != nil {
		http.Error(w, "err: unable to get users", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)

}

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body model.CreateUserRequest
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// user.ID = primitive.NewObjectID()
	now := time.Now().UTC()

	user = model.User{
		ID:       primitive.NewObjectID(),
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
		Created:  now,
	}

	_, err := database.Users.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("failed to create user")
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
	fmt.Println("User Created")
}

func (u *User) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "err: unable to get users", http.StatusBadRequest)
		return
	}

	result := database.Users.FindOne(context.Background(), primitive.M{"_id": _id})
	user := model.User{}
	err = result.Decode(&user)
	if err != nil {
		http.Error(w, "err: unable to find user", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (u *User) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string    `bson:"username" json:"username"`
		Email    string    `bson:"email" json:"email"`
		Updated  time.Time `bson:"updated" json:"updated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "err: unable to find user", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	now := time.Now().UTC()

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "err: unable to get users", http.StatusBadRequest)
		return
	}

	_, err = database.Users.UpdateOne(context.Background(), bson.M{"_id": _id}, bson.M{"$set": bson.M{
		"username": body.Username,
		"email":    body.Email,
		"updated":  now}})
	if err != nil {
		fmt.Println("Failed to update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Success: User Updated"))

}

func (u *User) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "err: unable to get users", http.StatusBadRequest)
		return
	}

	result, err := database.Users.DeleteOne(context.Background(), bson.M{"_id": _id})

	if err != nil {
		fmt.Println("Failed to update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "err: unable to find user", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Success: User deleted"))
	fmt.Println("Deletes user by id")
}
