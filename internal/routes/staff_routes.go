package routes

import (
	"github.com/Wizardsmile1412/hospital-middleware/internal/handler"
	"github.com/gin-gonic/gin"
)

func RegisterStaffRoutes(api *gin.RouterGroup, staffHandler *handler.StaffHandler) {
	staff := api.Group("/staff")
	{
		staff.POST("/create", staffHandler.CreateStaff)
		staff.POST("/login", staffHandler.Login)
		staff.POST("/refresh", staffHandler.RefreshToken)
	}
}
