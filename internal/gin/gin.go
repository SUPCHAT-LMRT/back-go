package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/usecase/list_recent_chats"
	validator2 "github.com/supchat-lmrt/back-go/internal/gin/validator"
	list_group_chat_messages "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/list_messages"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/user/gin/middlewares"
	request_forgot_password "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/usecase/request"
	validate_forgot_password "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_my_account"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/logout"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/register"
	request_reset_password "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/usecase/request"
	validate_reset_password "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_user"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_user_avatar"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/usecase/request"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/list_messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/create_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/create_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_icon"
	uberdig "go.uber.org/dig"
	"os"
)

type GinRouter interface {
	RegisterRoutes()
	AddCorsHeaders()
	Run() error
}

type DefaultGinRouter struct {
	Router *gin.Engine
	deps   GinRouterDeps
}

type GinRouterDeps struct {
	uberdig.In
	// Middlewares
	AuthMiddleware *middlewares.AuthMiddleware
	// Handlers
	// Workspace
	ListWorkspaceHandler       *list_workspaces.ListWorkspaceHandler
	CreateWorkspaceHandler     *create_workspace.CreateWorkspaceHandler
	UpdateWorkspaceIconHandler *update_icon.UpdateWorkspaceIconHandler
	// Workspaces channels
	ListChannelsHandler        *list_channels.ListChannelsHandler
	CreateChannelHandler       *create_channel.CreateChannelHandler
	ListChannelMessagesHandler *list_messages.ListChannelMessagesHandler
	// User
	GetMyAccountHandler                      *get_my_account.GetMyAccountHandler
	LoginHandler                             *login.LoginHandler
	RegisterHandler                          *register.RegisterHandler
	RefreshTokenHandler                      *token.RefreshTokenHandler
	LogoutHandler                            *logout.LogoutHandler
	UpdateAccountPersonalInformationsHandler *update_user.UpdateAccountPersonalInformationsHandler
	UpdateUserAvatarHandler                  *update_user_avatar.UpdateUserAvatarHandler
	// User validate
	RequestAccountValidationHandler *request.RequestAccountValidationHandler
	ValidateAccountHandler          *validate.ValidateAccountHandler
	// User forgot password
	RequestForgotPasswordHandler  *request_forgot_password.RequestForgotPasswordHandler
	ValidateForgotPasswordHandler *validate_forgot_password.ValidateForgotPasswordHandler
	// User reset password
	RequestResetPasswordHandler  *request_reset_password.RequestResetPasswordHandler
	ValidateResetPasswordHandler *validate_reset_password.ValidateResetPasswordHandler
	// Ws
	WebsocketHandler *websocket.WebsocketHandler
	// Chat
	ListRecentChatsHandler *list_recent_chats.ListRecentChatsHandler
	// Group
	AddMemberToGroupHandler *add_member.AddMemberToGroupHandler
	// Group chat
	ListGroupChatMessagesHandler *list_group_chat_messages.ListGroupChatMessagesHandler
}

func NewGinRouter(deps GinRouterDeps) GinRouter {
	router := gin.Default()
	return &DefaultGinRouter{Router: router, deps: deps}
}

func (d *DefaultGinRouter) RegisterRoutes() {
	authMiddleware := d.deps.AuthMiddleware.Execute

	apiGroup := d.Router.Group("/api")
	apiGroup.GET("/ws", authMiddleware, d.deps.WebsocketHandler.Handle)

	accountGroup := apiGroup.Group("/account")
	{
		accountGroup.GET("/me", authMiddleware, d.deps.GetMyAccountHandler.Handle)
		accountGroup.PUT("/personal-informations", authMiddleware, d.deps.UpdateAccountPersonalInformationsHandler.Handle)
		accountGroup.PUT("/avatar", authMiddleware, d.deps.UpdateUserAvatarHandler.Handle)

		authGroup := accountGroup.Group("/auth")
		{
			authGroup.POST("/login", d.deps.LoginHandler.Handle)
			authGroup.POST("/register", d.deps.RegisterHandler.Handle)
			authGroup.POST("/token/access/renew", d.deps.RefreshTokenHandler.Handle)
			authGroup.POST("/logout", authMiddleware, authMiddleware, d.deps.LogoutHandler.Handle)
		}

		validationGroup := accountGroup.Group("/validation")
		{
			validationGroup.POST("/request", d.deps.RequestAccountValidationHandler.Handle)
			validationGroup.POST("/validate", d.deps.ValidateAccountHandler.Handle)
		}

		forgotPasswordGroup := accountGroup.Group("/forgot-password")
		{
			forgotPasswordGroup.POST("/request", d.deps.RequestForgotPasswordHandler.Handle)
			forgotPasswordGroup.POST("/validate", d.deps.ValidateForgotPasswordHandler.Handle)
		}

		resetPasswordGroup := accountGroup.Group("/reset-password")
		{
			resetPasswordGroup.POST("/request", authMiddleware, d.deps.RequestResetPasswordHandler.Handle)
			resetPasswordGroup.POST("/validate", d.deps.ValidateResetPasswordHandler.Handle)
		}
	}

	chatGroup := apiGroup.Group("/chats")
	{
		chatGroup.GET("/recents", authMiddleware, d.deps.ListRecentChatsHandler.Handle)
	}

	groupGroup := apiGroup.Group("/groups")
	{
		groupGroup.POST("/members", authMiddleware, d.deps.AddMemberToGroupHandler.Handle)
		groupGroup.GET("/:groupId/messages", authMiddleware, d.deps.ListGroupChatMessagesHandler.Handle)
	}

	workspacesGroup := apiGroup.Group("/workspaces")
	{
		workspacesGroup.Use(authMiddleware)
		workspacesGroup.GET("", d.deps.ListWorkspaceHandler.Handle)
		workspacesGroup.POST("", d.deps.CreateWorkspaceHandler.Handle)
		workspacesGroup.PUT("/:workspaceId/icon", d.deps.UpdateWorkspaceIconHandler.Handle)

		channelGroup := workspacesGroup.Group("/:workspaceId/channels")
		{
			channelGroup.GET("", d.deps.ListChannelsHandler.Handle)
			channelGroup.POST("", d.deps.CreateChannelHandler.Handle)
			channelGroup.GET("/:channelId/messages", d.deps.ListChannelMessagesHandler.Handle)
		}
	}
}

func (d *DefaultGinRouter) AddCorsHeaders() {
	d.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

func (d *DefaultGinRouter) Run() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("ISO8601date", validator2.IsISO8601Date)
		if err != nil {
			return err
		}
	}

	return d.Router.Run(":" + os.Getenv("HTTP_SERVER_PORT"))
}
