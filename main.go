package main

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type (
	User struct {
		Name  string `json:"name" form:"name" query:"name" validate:"required"`
		Email string `json:"email" form:"email" query:"email" validate:"required,email"`
		Role  string `json:"role" form:"role" query:"role"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	glob, err := template.New("form.html").Parse(`<html>
<body>
<main>
    <form id="user-form" data-testid="user-form" method="post" action="/user">
        <label for="name">Name</label>
        <input name="name" id="name" required type="text" data-test-id="name-input"/>
        <label for="email">Email</label>
        <input name="email" id="email" required type="email" data-test-id="email-input"/>
        <label for="role">Role</label>
        <input name="role" id="role" type="text" data-test-id="role-input"/>
        <input type="submit" value="Submit" />
    </form>
</main>
</body>
</html>
`)

	glob, err = glob.New("user.html").Parse(`<html>
<body>
<main>
    <div data-test-id="user-name">{{index . "name"}}</div>
    <div data-test-id="user-email">{{index . "email"}}</div>
    <div data-test-id="user-role">{{index . "role"}}</div>
</main>
</body>
</html>`)

	t := &Template{
		templates: template.Must(glob, err),
	}

	e.Renderer = t

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "form.html", map[string]interface{}{
		})
	})

	e.POST("/user", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return c.Render(
			http.StatusOK,
			"user.html",
			map[string]interface{}{
				"name": u.Name,
				"email": u.Email,
				"role": u.Role,
			},
		)
	})

	e.Logger.Fatal(e.Start(":3002"))
}
