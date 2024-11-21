package TeamManagement

import (
	"exchangeapp/global"
	"exchangeapp/models/team"
	"exchangeapp/models/user"
	"exchangeapp/rsp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// CreateTeam 新建团队
// @Summary 新建团队
// @Description 创建一个新的团队
// @Tags 团队操作
// @Param name body string true "团队名称"
// @Param description body string false "团队描述"
// @Param owner_id body int true "团队创建者的用户ID"
// @Success 201 {object} string "注册成功，返回 JWT token"
// @Failure 400 {object} rsp.ErrorResponse
// @Router /teams [post]
func CreateTeam(ctx *gin.Context) {
	var team team.Team

	// 绑定请求体
	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, rsp.NewErrorResponse(30001, err.Error()))
		return
	}

	// 自动迁移数据库
	if err := global.Db.AutoMigrate(&team); err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(20004, err.Error()))
		return
	}

	// 创建团队
	if err := global.Db.Create(&team).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(20003, err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, rsp.NewSuccessResponse(11001, team))
}

// AddUserToTeam 将用户添加到团队并设置权限
// @Summary 添加用户到团队
// @Description 将指定用户添加到团队中，并为其设置角色和权限
// @Tags 团队操作
// @Param team_id body int true "团队ID"
// @Param user_id body int true "用户ID"
// @Param role body string true "用户角色（owner/admin/member）"
// @Param permissions body string true "权限（JSON格式）"
// @Success 200 {object} string "注册成功，返回 JWT token"
// @Failure 400 {object} rsp.ErrorResponse
// @Router /teams/add_user [post]
func AddUserToTeam(ctx *gin.Context) {
	var member team.TeamMember

	// 绑定请求体
	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(http.StatusBadRequest, rsp.NewErrorResponse(30001, err.Error()))
		return
	}

	// 检查团队是否存在
	var t team.Team
	if err := global.Db.First(&t, member.TeamID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, rsp.NewErrorResponse(11002, "团队不存在"))
		return
	}

	// 检查用户是否存在
	var u user.User
	if err := global.Db.First(&u, member.UserID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, rsp.NewErrorResponse(10002, "用户不存在"))
		return
	}

	// 自动迁移数据库
	if err := global.Db.AutoMigrate(&member); err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(20004, err.Error()))
		return
	}

	// 添加用户到团队
	if err := global.Db.Create(&member).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(20003, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, rsp.NewSuccessResponse(12001, member))
}

// GetTeamMembers 查询团队的所有成员信息
// @Summary 查询团队成员
// @Description 查询指定团队的所有成员信息
// @Tags 团队操作
// @Param id path int true "团队ID"
// @Success 200 {object} string "注册成功，返回 JWT token"
// @Failure 404 {object} rsp.ErrorResponse
// @Router /teams/{id}/members [get]
func GetTeamMembers(ctx *gin.Context) {
	teamID := ctx.Param("id")

	// 查询团队成员信息并预加载用户
	var members []team.TeamMember
	if err := global.Db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "username", "email") // 只选择需要的字段
	}).Where("team_id = ?", teamID).Find(&members).Error; err != nil {
		ctx.JSON(http.StatusNotFound, rsp.NewErrorResponse(11002, "未找到该团队的成员信息"))
		return
	}

	ctx.JSON(http.StatusOK, rsp.NewSuccessResponse(11002, members))
}
