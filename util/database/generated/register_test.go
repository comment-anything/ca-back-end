package generated

import (
	"context"
	"testing"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		Username: "test-user1",
		Password: "1xtdf",
		Email:    "7@7.com",
	}
	var user User
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Error(err)
	}
	if user.Username != "test-user1" || user.Password != "1xtdf" || user.Email != "7@7.com" {
		t.Errorf("User was not populated with values.")
	}
	if user.IsVerified != false {
		t.Errorf("Default is verified should be false.")
	}
	if user.Banned != false {
		t.Errorf("Default banned should be false.")
	}
	arg.Email = "8@8.com"
	_, err = testQueries.CreateUser(context.Background(), arg)
	if err == nil {
		t.Errorf("User should not be created if username already exists, even if email is unique")
	}
	arg.Email = "7@7.com"
	arg.Username = "test-user2"
	_, err = testQueries.CreateUser(context.Background(), arg)
	if err == nil {
		t.Errorf("User should not be created if email already exists, even if username is unique")
	}

	err = testQueries.Tst_DeleteUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetUserByEmail(t *testing.T) {
	arg := CreateUserParams{
		Username: "test-user1",
		Password: "1xtdf",
		Email:    "7@7.com",
	}
	var user User
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Errorf(err.Error())
	}
	retUser, err := testQueries.GetUserByEmail(context.Background(), "7@7.com")
	if err != nil {
		t.Errorf(err.Error())
	}
	if retUser.ID != user.ID {
		t.Errorf("user ids dont match")
	}
	retUser, err = testQueries.GetUserByEmail(context.Background(), "8@8.com")
	if err == nil {
		t.Errorf("getting nonexistant user by email should fail")
	}
	err = testQueries.Tst_DeleteUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

}

func TestGetUserByUsername(t *testing.T) {
	arg := CreateUserParams{
		Username: "test-user1",
		Password: "1xtdf",
		Email:    "7@7.com",
	}
	var user User
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Errorf(err.Error())
	}
	retUser2, err := testQueries.GetUserByUserName(context.Background(), "test-user1")
	if err != nil {
		t.Errorf(err.Error())
	}
	if retUser2.ID != user.ID {
		t.Errorf("user ids dont match")
	}
	retUser2, err = testQueries.GetUserByUserName(context.Background(), "test-user2")
	if err == nil {
		t.Errorf("getting nonexistant user by username should fail")
	}
	err = testQueries.Tst_DeleteUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGetUserByID(t *testing.T) {
	arg := CreateUserParams{
		Username: "test-user1",
		Password: "1xtdf",
		Email:    "7@7.com",
	}
	var user User
	user, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		t.Errorf(err.Error())
	}
	retUser2, err := testQueries.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
	if retUser2.Username != user.Username {
		t.Errorf("usernames dont match")
	}
	err = testQueries.Tst_DeleteUser(context.Background(), user.ID)
	if err != nil {
		t.Errorf(err.Error())
	}

}
