package auth

import (
	v1 "spiritFruit/app/http/controllers/admin/v1"
	requests "spiritFruit/app/requests"
	adminRequest "spiritFruit/app/requests/admin"
	"spiritFruit/app/models/admins"
	"spiritFruit/pkg/auth"
	"spiritFruit/pkg/hash"
	"spiritFruit/pkg/jwt"
	"spiritFruit/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginController 登录控制器
type LoginController struct {
	v1.BaseADMINController
}

// LoginByPassword 多种方法登录，支持手机号和用户名
func (lc *LoginController) LoginByPassword(c *gin.Context) {
	// 1. 验证表单
	request := adminRequest.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, adminRequest.LoginByPassword); !ok {
		return
	}

	// 2. 尝试登录
	admin, err := auth.AdminAttempt(request.LoginID, request.Password)
	if err != nil {
		// 失败，显示错误提示
        response.JSON(c, gin.H{
            "code":    1001,
            "data":    nil,
            "message": "账号不存在或密码错误",
        })

	} else {
		token := jwt.NewJWT().IssueAdminToken(admin.GetStringID(), *admin.Username)
		admin.Password = nil
		response.JSON(c, gin.H{
        	"code": 0,
            "data": map[string]interface{}{
                "token":     token,
                "user_info": admin,
            },
        })
	}
}

// CurrentAdmin 获取当前登录代理信息
func (lc *LoginController) CurrentAdmin(c *gin.Context) {
	adminModel := auth.CurrentAdmin(c)
	response.Data(c, adminModel)
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {

	token, err := jwt.NewJWT().RefreshToken(c)

	if err != nil {
		response.Error(c, err, "令牌刷新失败")
	} else {
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}

// Register 注册
func (lc *LoginController) Register(c *gin.Context) {
	// 1. 验证表单
	request := adminRequest.RegisterRequest{}
	if ok := requests.Validate(c, &request, adminRequest.Register); !ok {
		return
	}

	// 2. 创建用户
	adminModel := admins.Admins{
		Username: &request.Username,
		Mobile:   &request.Mobile,
		Email:    &request.Email,
		Password: &request.Password,
	}

	// 密码加密
	passwordHash := hash.BcryptHash(request.Password)
	adminModel.Password = &passwordHash

	adminModel.Create()

	if adminModel.ID > 0 {
		response.JSON(c, gin.H{
			"code": 0,
			"data": adminModel,
			"message": "注册成功",
		})
	} else {
		response.Abort500(c, "更新失败，请稍后尝试~")
	}
}
