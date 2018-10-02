package handler

import (
	"net/http"

	"bitbucket.org/boomstarternetwork/mineradmin/store"
	"github.com/labstack/echo"
)

type Handler struct {
	projects store.ProjectsStore
	balances store.BalancesStore
}

func NewHandler(ps store.ProjectsStore, bs store.BalancesStore) Handler {
	return Handler{
		projects: ps,
		balances: bs,
	}
}

func (h Handler) Index(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/projects")
}

// language=HTML
const ProjectsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Projects</title>
	<style>
		.coin {
			display: inline-block;
			padding-left: 0.1em;
		}
	</style>
</head>
<body>
	<h1>Projects</h1>
	{{range .}}
		<div>
			<a href="/project/{{.ProjectID}}/users">{{.ProjectName}}</a>
			{{range .Coins}}
				<span class="coin">{{.Coin}}:{{.Amount}}</span>
			{{end}}
		</div>
	{{end}}
</body>
</html>
`

func (h Handler) Projects(c echo.Context) error {
	balances, err := h.balances.ProjectsBalances()
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Internal server error")
	}
	return c.Render(http.StatusOK, "projects", balances)
}

// language=HTML
const ProjectUsersTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Project :: {{.Project.Name}}</title>
	<style>
		.coin {
			display: inline-block;
			padding-left 0.1em;
		}
	</style>
</head>
<body>
	<h1>Projects :: {{.Project.Name}}</h1>
	{{range .Balances}}
		<div>
			<span>{{.Address}}</sp>
			{{range .Coins}}
				<span class="coin">{{.Coin}}:{{.Amount}}</span>
			{{end}}
		</div>
	{{end}}
</body>
</html>
`

func (h Handler) ProjectUsers(c echo.Context) error {
	project, err := h.projects.Get(c.Param("project-id"))
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Internal server error")
	}

	balances, err := h.balances.ProjectUsersBalances(project.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError,
			"Internal server error")
	}

	return c.Render(http.StatusOK, "project-users", struct {
		Project  store.Project
		Balances []store.UserBalance
	}{project, balances})
}
