package handlers

import (
	"database/sql"
	"gofr.dev/pkg/gofr"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUserHandlers(app *gofr.App, db *sql.DB) {
	// Create
	app.POST("/signup", func(c *gofr.Context) (interface{}, error) {
		var user User
		err := c.Bind(&user)
		if err != nil {
			return nil, gofr.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}

		if user.Username == "" || user.Password == "" {
			return nil, gofr.NewHTTPError(http.StatusBadRequest, "Username or password cannot be empty")
		}

		_, err = db.Exec("INSERT INTO users VALUES (?, ?)", user.Username, user.Password)
		if err != nil {
			return nil, gofr.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return map[string]string{"status": "User signed up successfully"}, nil
	})

	// Read
	app.GET("/users/:username", func(c *gofr.Context) (interface{}, error) {
		username := c.Param("username")
		var user User
		err := db.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user.Username, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, gofr.NewHTTPError(http.StatusNotFound, "User not found")
			} else {
				return nil, gofr.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}

		return user, nil
	})

	// Update
	app.PUT("/users/:username", func(c *gofr.Context) (interface{}, error) {
		username := c.Param("username")
		var user User
		err := c.Bind(&user)
		if err != nil {
			return nil, gofr.NewHTTPError(http.StatusBadRequest, "Bad Request")
		}

		_, err = db.Exec("UPDATE users SET password = ? WHERE username = ?", user.Password, username)
		if err != nil {
			return nil, gofr.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return map[string]string{"status": "User updated successfully"}, nil
	})

	// Delete
	app.DELETE("/users/:username", func(c *gofr.Context) (interface{}, error) {
		username := c.Param("username")
		_, err := db.Exec("DELETE FROM users WHERE username = ?", username)
		if err != nil {
			return nil, gofr.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return map[string]string{"status": "User deleted successfully"}, nil
	})
}
