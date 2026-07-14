package admin

import (
	"fmt"
	"net/http"
	"pizza-tracker/internal/shared/util"
	"pizza-tracker/internal/user"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RenderLogin(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Dashboard(c *gin.Context)
}

type AdminDeps struct {
	UserRepo user.UserRepository
}

type handler struct {
	AdminDeps
}

func NewHandler(deps AdminDeps) Handler {
	return &handler{
		AdminDeps: deps,
	}
}

type LoginData struct {
	Error string
}

type DashboardData struct {
	Username string
}

func (h *handler) RenderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", LoginData{})
}

func (h *handler) Login(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required,min=3,max=50"`
		Password string `form:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{Error: "Invalid input: " + err.Error()})
		return
	}

	user, err := h.UserRepo.Authenticate(form.Username, form.Password)
	if err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{Error: "Invalid credentials"})
		return
	}

	util.SetSessionValue(c, "userID", fmt.Sprintf("%d", user.ID))
	util.SetSessionValue(c, "username", user.Username)

	c.Redirect(http.StatusSeeOther, "/admin/dashboard")
}

func (h *handler) Logout(c *gin.Context) {
	if err := util.ClearSession(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (h *handler) Dashboard(c *gin.Context) {
	username := util.GetSessionString(c, "username")

	c.HTML(http.StatusOK, "dashboard.tmpl", DashboardData{Username: username})
}
