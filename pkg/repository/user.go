package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fastbyt3/diy-mssql-test/pkg/models"
	"github.com/fastbyt3/diy-mssql-test/pkg/utils"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetUser(ctx context.Context, id int32) (*models.User, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := "INSERT INTO users(username, email) VALUES(@username, @email)"
	fmt.Println(user.Username, user.Email)
	_, err := r.db.ExecContext(ctx, query,
		sql.Named("username", user.Username),
		sql.Named("email", user.Email),
	)
	return err
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User

	rows, err := r.db.Query("SELECT Id, Username, Email FROM dbo.users;")
	if err != nil {
		return nil, fmt.Errorf("getAllUsers failed", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email); err != nil {
			return nil, fmt.Errorf("failed to parse result into User struct in getAllUsers", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (r *UserRepo) GetUser(ctx context.Context, id int32) (*models.User, error) {
	query := "SELECT id, username, email FROM dbo.users WHERE id = @p1"
	row := r.db.QueryRowContext(ctx, query, id)

	user := &models.User{}
	err := row.Scan(&user.Id, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("No such user exists")
	}

	if err != nil {
		utils.Logger.Error("Unexpected error fetching user",
			zap.String("error", err.Error()),
			zap.Int32("user id", id),
		)

		return nil, fmt.Errorf("Unexpected error")
	}

	return user, nil
}
