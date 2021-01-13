package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/PECHIVKO/task-manager/models"
)

type User struct {
	Username string `boil:"username" json:"username" toml:"username" yaml:"username"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(dbConn *sql.DB) *UserRepository {
	var repo UserRepository
	repo.DB = dbConn
	return &repo
}

func (r UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	model := toUser(user)
	insertUserQuery := "INSERT INTO users (username) VALUES ($1);"

	_, err := r.DB.Exec(insertUserQuery, model.Username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", err)
		// os.Exit(1)
	}
	return nil
}

func (r UserRepository) GetUser(ctx context.Context, username string) (*models.User, error) {
	user := new(User)

	getUserQuery := "select * from users where id = $1"

	err := r.DB.QueryRow(getUserQuery, username).Scan(&user.Username)

	if err != nil {
		return nil, err
	}

	return toModel(user), nil
}

func (r UserRepository) FetchProjectsUsers(ctx context.Context) ([]*models.User, error) {
	FetchProjectsUsersQuery := "select * from users"

	rows, err := r.DB.Query(FetchProjectsUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)

	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Username)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}

	return toUsers(users), nil
}

func toUser(u *models.User) *User {
	return &User{
		Username: u.Username,
	}
}

func toUsers(us []*User) []*models.User {
	out := make([]*models.User, len(us))

	for i, u := range us {
		out[i] = toModel(u)
	}

	return out
}

func toModel(u *User) *models.User {
	return &models.User{
		Username: u.Username,
	}
}
