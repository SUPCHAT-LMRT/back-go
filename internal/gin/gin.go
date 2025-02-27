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
	user_invite_link_generate "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/generate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/get_data_token_invite"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login_oauth"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/logout"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/register"
	request_reset_password "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/usecase/request"
	validate_reset_password "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_user"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_user_avatar"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/usecase/get_minutely"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/list_messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/create_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_channels"
	workspace_middlewares "github.com/supchat-lmrt/back-go/internal/workspace/gin/middlewares"
	workspace_invite_link_generate "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/generate"
	get_data_token_invite2 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/get_data_token_invite"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/join_workspace_invite"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/list_workpace_members"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/create_workspace"
	discovery_list_workspaces "github.com/supchat-lmrt/back-go/internal/workspace/usecase/discover/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace_details"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_banner"
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
	AuthMiddleware            *middlewares.AuthMiddleware
	UserInWorkspaceMiddleware *workspace_middlewares.UserInWorkspaceMiddleware
	// Handlers
	// Workspace
	ListWorkspaceHandler              *list_workspaces.ListWorkspaceHandler
	CreateWorkspaceHandler            *create_workspace.CreateWorkspaceHandler
	UpdateWorkspaceIconHandler        *update_icon.UpdateWorkspaceIconHandler
	UpdateWorkspaceBannerHandler      *update_banner.UpdateWorkspaceBannerHandler
	ListWorkspaceMembersHandler       *list_workpace_members.ListWorkspaceMembersHandler
	DiscoverListWorkspaceHandler      *discovery_list_workspaces.DiscoverListWorkspaceHandler
	GetWorkspaceDetailsHandler        *get_workspace_details.GetWorkspaceDetailsHandler
	GetMinutelyMessageSentHandler     *get_minutely.GetMinutelyMessageSentHandler
	CreateInviteLinkWorkspaceHandler  *workspace_invite_link_generate.CreateInviteLinkHandler
	GetInviteLinkWorkspaceDataHandler *get_data_token_invite2.GetInviteLinkWorkspaceDataHandler
	// Workspaces channels
	ListChannelsHandler        *list_channels.ListChannelsHandler
	CreateChannelHandler       *create_channel.CreateChannelHandler
	ListChannelMessagesHandler *list_messages.ListChannelMessagesHandler
	// User
	GetMyAccountHandler                      *get_my_account.GetMyUserAccountHandler
	LoginHandler                             *login.LoginHandler
	RegisterHandler                          *register.RegisterHandler
	RefreshTokenHandler                      *token.RefreshTokenHandler
	LogoutHandler                            *logout.LogoutHandler
	UpdateAccountPersonalInformationsHandler *update_user.UpdateAccountPersonalInformationsHandler
	UpdateUserAvatarHandler                  *update_user_avatar.UpdateUserAvatarHandler
	// User forgot password
	RequestForgotPasswordHandler  *request_forgot_password.RequestForgotPasswordHandler
	ValidateForgotPasswordHandler *validate_forgot_password.ValidateForgotPasswordHandler
	// User reset password
	RequestResetPasswordHandler  *request_reset_password.RequestResetPasswordHandler
	ValidateResetPasswordHandler *validate_reset_password.ValidateResetPasswordHandler
	// Invite link
	CreateInviteLinkHandler    *user_invite_link_generate.CreateInviteLinkHandler
	GetInviteLinkDataHandler   *get_data_token_invite.GetInviteLinkDataHandler
	JoinWorkspaceInviteHandler *join_workspace_invite.JoinWorkspaceInviteHandler
	// User OAuth connection
	OAuthHandler *login_oauth.OAuthHandler
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
	userInWorkspaceMiddleware := d.deps.UserInWorkspaceMiddleware.Execute

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
			authGroup.POST("/logout", authMiddleware, d.deps.LogoutHandler.Handle)
			authGroup.GET("/oauth/:provider", d.deps.OAuthHandler.AuthProvider)
			authGroup.GET("/oauth/:provider/callback", d.deps.OAuthHandler.AuthCallback)
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

		inviteLinkGroup := accountGroup.Group("/invite-link")
		{
			inviteLinkGroup.POST("", d.deps.CreateInviteLinkHandler.Handle)
			inviteLinkGroup.GET("/:token", d.deps.GetInviteLinkDataHandler.Handle)

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
		workspacesGroup.GET("/discover", d.deps.DiscoverListWorkspaceHandler.Handle)

		specificWorkspaceGroup := workspacesGroup.Group("/:workspaceId")
		{
			specificWorkspaceGroup.Use(userInWorkspaceMiddleware)
			specificWorkspaceGroup.PUT("/icon", d.deps.UpdateWorkspaceIconHandler.Handle)
			specificWorkspaceGroup.PUT("/banner", d.deps.UpdateWorkspaceBannerHandler.Handle)
			specificWorkspaceGroup.GET("/members", d.deps.ListWorkspaceMembersHandler.Handle)
			specificWorkspaceGroup.GET("/details", d.deps.GetWorkspaceDetailsHandler.Handle)
			specificWorkspaceGroup.GET("/time-series/messages", d.deps.GetMinutelyMessageSentHandler.Handle)

			channelGroup := specificWorkspaceGroup.Group("/channels")
			{
				channelGroup.GET("", d.deps.ListChannelsHandler.Handle)
				channelGroup.POST("", d.deps.CreateChannelHandler.Handle)
				channelGroup.GET("/:channelId/messages", d.deps.ListChannelMessagesHandler.Handle)
			}
		}
	}

	inviteLinkGroup := apiGroup.Group("/workspace-invite-link")
	{
		inviteLinkGroup.POST("", d.deps.CreateInviteLinkWorkspaceHandler.Handle)
		inviteLinkGroup.GET("/:token", d.deps.GetInviteLinkWorkspaceDataHandler.Handle)
		inviteLinkGroup.POST("/:token/join", authMiddleware, d.deps.JoinWorkspaceInviteHandler.Handle)
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
