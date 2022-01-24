package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"user-crud/internal/models"
	"user-crud/pkg/fielderror"
	"user-crud/pkg/filter"
	"user-crud/pkg/paginator"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

//Client is the DB client struct
type Client struct {
	db *sqlx.DB
}

func (c *Client) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

//Get gets a user's info from the DB
func (c *Client) Get(ctx context.Context, id string) (*models.User, error) {
	qb := pq().Select(
		"id",
		"first_name",
		"last_name",
		"nickname",
		"email",
		"country",
		"created_at",
		"updated_at",
	).From(usersTable).Where(squirrel.Eq{"id": id})

	query, args, err := qb.ToSql()

	if err != nil {
		log.Printf("postgres: failed to construct sql: %s", err)
		return nil, err
	}

	user := &models.User{}
	err = c.db.GetContext(ctx, user, query, args...)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fielderror.NonFieldErrorf("no such user exists")
	}

	if err != nil {
		log.Printf("postgres: failed to fetch user: %s", err)
		return nil, err
	}

	return user, nil
}

//Insert inserts a user in the DB
func (c *Client) Insert(ctx context.Context, user models.User) (*models.User, error) {
	qb := pq().Insert(usersTable).Suffix("RETURNING *")

	values := make(map[string]interface{})

	if user.FirstName != "" {
		values["first_name"] = user.FirstName
	}

	if user.LastName != "" {
		values["last_name"] = user.LastName
	}

	if user.Nickname != "" {
		values["nickname"] = user.Nickname
	}

	if user.Password != "" {
		values["password"] = user.Password
	}

	if user.Email != "" {
		values["email"] = user.Email
	}

	if user.Country != "" {
		values["country"] = user.Country
	}

	qb = qb.SetMap(values)

	query, args, err := qb.ToSql()

	if err != nil {
		log.Printf("postgres: failed to construct sql: %s", err)
		return nil, err
	}

	createdUser := &models.User{}
	err = c.db.GetContext(ctx, createdUser, query, args...)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return nil, fielderror.FieldErrorf("email", "email %s already exists", user.Email)
		}

		switch pgErr.ColumnName {
		case "password":
			return nil, fielderror.FieldErrorf("password", "password is required")
		case "email":
			return nil, fielderror.FieldErrorf("email", "email is required")
		}
	}

	if err != nil {
		log.Printf("postgres: failed to insert user: %s", err)
		return nil, err
	}

	return createdUser, nil
}

//Update updates user partially meaning all set fields will be updated
func (c *Client) Update(ctx context.Context, user models.User) (*models.User, bool, error) {
	qb := pq().Update(usersTable).Where(squirrel.Eq{"id": user.ID}).Suffix("RETURNING *")

	if user.FirstName != "" {
		qb = qb.Set("first_name", user.FirstName)
	}

	if user.LastName != "" {
		qb = qb.Set("last_name", user.LastName)
	}

	if user.Nickname != "" {
		qb = qb.Set("nickname", user.Nickname)
	}

	if user.Password != "" {
		qb = qb.Set("password", user.Password)
	}

	if user.Email != "" {
		qb = qb.Set("email", user.Email)
	}

	if user.Country != "" {
		qb = qb.Set("country", user.Country)
	}

	qb = qb.Set("updated_at", time.Now())

	updatedUser := &models.User{}

	query, args, err := qb.ToSql()

	if err != nil {
		log.Printf("postgres: failed to construct sql: %s", err)
		return nil, false, err
	}

	if len(args) <= 2 {
		fetchedUser, fetchErr := c.Get(ctx, user.ID)
		return fetchedUser, false, fetchErr
	}

	err = c.db.GetContext(ctx, updatedUser, query, args...)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return nil, false, fielderror.FieldErrorf("email", "email %s already exists", user.Email)
		}
	}

	if err != nil {
		log.Printf("postgres: failed to update user: %s", err)
		return nil, false, err
	}

	return updatedUser, true, nil
}

//List is the function that lists the users
func (c *Client) List(ctx context.Context, pg paginator.Paginator, filters ...filter.Filter) (*models.UsersList, error) {
	qb := pq().Select("id",
		"email",
		"first_name",
		"last_name",
		"nickname",
		"email",
		"country",
		"created_at",
		"updated_at",
		"count(id) OVER() as total_count",
	).
		From(usersTable).
		Limit(pg.Size()).
		Offset((pg.Page() - 1) * pg.Size())

	for _, f := range filters {
		qb = applyFilter(qb, f)
	}

	query, args, err := qb.ToSql()

	if err != nil {
		log.Printf("postgres: failed to construct sql: %s", err)
		return nil, err
	}

	var usersWithCount []struct {
		models.User
		TotalCount uint64 `db:"total_count"`
	}

	if err = c.db.SelectContext(ctx, &usersWithCount, query, args...); err != nil {
		log.Printf("postgres: failed to list users: %s", err)
		return nil, err
	}

	users := models.UsersList{Users: make([]*models.User, 0, len(usersWithCount))}

	for _, u := range usersWithCount {
		user := u
		users.Users = append(users.Users, &user.User)
		users.TotalCount = user.TotalCount
	}

	return &users, nil
}

// Delete deletes user from db
func (c *Client) Delete(ctx context.Context, id string) error {
	qb := pq().Delete(usersTable).Where(squirrel.Eq{"id": id})

	query, args, err := qb.ToSql()

	if err != nil {
		log.Printf("postgres: failed to construct sql: %s", err)
		return err
	}

	_, err = c.db.ExecContext(ctx, query, args...)

	if err != nil {
		log.Printf("postgres: failed to delete user: %s", err)
		return err
	}

	return nil
}

func (c *Client) Database() *sqlx.DB {
	return c.db
}
