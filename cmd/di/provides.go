package di

import (
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/usecase/list_recent_chats"
	"github.com/supchat-lmrt/back-go/internal/dig"
	"github.com/supchat-lmrt/back-go/internal/gin"
	group_chat_message_repository "github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	list_group_chat_messages "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/list_messages"
	save_group_chat_message "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/save_message"
	group_repository "github.com/supchat-lmrt/back-go/internal/group/repository"
	"github.com/supchat-lmrt/back-go/internal/group/strategies"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/list_recent_groups"
	logger "github.com/supchat-lmrt/back-go/internal/logger/zerolog"
	"github.com/supchat-lmrt/back-go/internal/mail"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/s3"
	user_chat_direct_repository "github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	list_recent_chats_direct "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_recent_direct_chats"
	"github.com/supchat-lmrt/back-go/internal/user/gin/middlewares"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists"
	forgot_password_repository "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/repository"
	forgot_password_service "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/service"
	forgot_password_request_usecase "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/usecase/request"
	forgot_password_validate_usecase "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_my_account"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/logout"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/register"
	reset_password_repository "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/repository"
	reset_password_service "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/service"
	reset_password_request_usecase "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/usecase/request"
	reset_password_validate_usecase "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/sendmail"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/token"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_password"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_user"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_user_avatar"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/service"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/usecase/request"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/utils"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	chat_message_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/list_messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/save_message"
	channel_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/create_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_channels"
	workspace_middlewares "github.com/supchat-lmrt/back-go/internal/workspace/gin/middlewares"
	workspace_repository "github.com/supchat-lmrt/back-go/internal/workspace/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/create_workspace"
	discovery_list_workspaces "github.com/supchat-lmrt/back-go/internal/workspace/usecase/discover/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workpace_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace_details"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/is_user_in_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/list_workpace_members"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_banner"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_icon"
	uberdig "go.uber.org/dig"
	"log"
	"os"
	"reflect"
	"runtime"
)

func NewDi() *uberdig.Container {
	di := uberdig.New()
	providers := []dig.Provider{
		// Logger
		dig.NewProvider(logger.NewZerologLogger),
		// Gin
		dig.NewProvider(gin.NewGinRouter),
		// Mongo
		dig.NewProvider(mongo.NewClient),
		// Redis
		dig.NewProvider(redis.NewClient),
		// S3
		dig.NewProvider(s3.NewS3Client),
		// Mailer
		dig.NewProvider(mail.NewMailer(os.Getenv("SMTP_HOST"), os.Getenv("SMTP_TLS") == "true", utils.MustAtoi(os.Getenv("SMTP_PORT")), os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))),
		// Workspaces
		// Workspaces repository
		dig.NewProvider(workspace_repository.NewMongoWorkspaceRepository),
		dig.NewProvider(workspace_repository.NewMongoWorkspaceMapper),
		dig.NewProvider(workspace_repository.NewMongoWorkspaceMemberMapper),
		// Workspace usecases
		dig.NewProvider(list_workspaces.NewListWorkspacesUseCase),
		dig.NewProvider(discovery_list_workspaces.NewDiscoveryListWorkspacesUseCase),
		dig.NewProvider(create_workspace.NewCreateWorkspaceUseCase),
		dig.NewProvider(get_workpace_member.NewGetWorkspaceMemberUseCase),
		dig.NewProvider(is_user_in_workspace.NewIsUserInWorkspaceUseCase),
		dig.NewProvider(get_workspace_details.NewGetWorkspaceDetailsUseCase),
		dig.NewProvider(update_icon.NewUpdateWorkspaceIconUseCase),
		dig.NewProvider(update_icon.NewS3UpdateWorkspaceIconStrategy),
		dig.NewProvider(update_banner.NewUpdateWorkspaceBannerUseCase),
		dig.NewProvider(update_banner.NewS3UpdateWorkspaceBannerStrategy),
		dig.NewProvider(list_workpace_members.NewListWorkspaceMembersUseCase),
		// Workspace handlers
		dig.NewProvider(list_workspaces.NewListWorkspaceHandler),
		dig.NewProvider(discovery_list_workspaces.NewDiscoverListWorkspaceHandler),
		dig.NewProvider(create_workspace.NewCreateWorkspaceHandler),
		dig.NewProvider(get_workspace_details.NewGetWorkspaceDetailsHandler),
		dig.NewProvider(update_icon.NewUpdateWorkspaceIconHandler),
		dig.NewProvider(update_banner.NewUpdateWorkspaceBannerHandler),
		dig.NewProvider(list_workpace_members.NewListWorkspaceHandler),
		// Workspace channels
		// Workspace channels repository
		dig.NewProvider(channel_repository.NewMongoChannelRepository),
		dig.NewProvider(channel_repository.NewMongoChannelMapper),
		// Workspace channels usecases
		dig.NewProvider(list_channels.NewListChannelsUseCase),
		dig.NewProvider(create_channel.NewCreateChannelUseCase),
		dig.NewProvider(get_channel.NewGetChannelUseCase),
		// Workspaces channels observers
		dig.NewProvider(create_channel.NewNotifyWebSocketObserver, uberdig.Group("create_channel_observers")),
		// Workspace channels handlers
		dig.NewProvider(list_channels.NewListChannelsHandler),
		dig.NewProvider(create_channel.NewCreateChannelHandler),
		// Workspace channels chat
		// Workspace channels chat repository
		dig.NewProvider(chat_message_repository.NewMongoChannelMessageRepository),
		dig.NewProvider(chat_message_repository.NewChannelMessageMapper),
		// Workspace channels chat usecases
		dig.NewProvider(list_messages.NewListMessageUseCase),
		// Workspace channels chat handlers
		dig.NewProvider(list_messages.NewListChannelMessagesHandler),
		// Workspace channels chat messages usecases
		dig.NewProvider(save_message.NewSaveChannelMessageUseCase),
		// Workspace misc
		dig.NewProvider(workspace_middlewares.NewUserInWorkspaceMiddleware),
		// User
		dig.NewProvider(user_repository.NewMongoUserRepository),
		dig.NewProvider(user_repository.NewMongoUserMapper),
		// User usecases
		dig.NewProvider(get_by_id.NewGetUserByIdUseCase),
		dig.NewProvider(get_by_email.NewGetUserByEmailUseCase),
		dig.NewProvider(login.NewLoginUserUseCase),
		dig.NewProvider(exists.NewExistsUserUseCase),
		dig.NewProvider(register.NewRegisterUserUseCase),
		dig.NewProvider(token.NewRefreshAccessTokenUseCase),
		dig.NewProvider(update_password.NewChangePasswordUseCase),
		dig.NewProvider(update_user.NewUpdateUserUseCase),
		dig.NewProvider(update_user_avatar.NewUpdateUserAvatarUseCase),
		dig.NewProvider(update_user_avatar.NewS3UpdateUserAvatarStrategy),
		// User handlers
		dig.NewProvider(update_user_avatar.NewUpdateUserAvatarHandler),
		// User register observers
		dig.NewProvider(register.NewRequestValidationObserver, uberdig.Group("register_user_observers")),
		// User validate repository
		dig.NewProvider(repository.NewRedisValidationRepository),
		dig.NewProvider(repository.NewRedisValidationRequestMapper),
		// User validate service
		dig.NewProvider(service.NewDefaultValidationRequestService),
		// User validate usecases
		dig.NewProvider(request.NewRequestAccountValidationUseCase),
		dig.NewProvider(validate.NewValidateAccountUseCase),
		// User validate observers
		dig.NewProvider(request.NewLogRequestValidationObserver, uberdig.Group("validation_request_observers")),
		dig.NewProvider(request.NewSendMailRequestValidationObserver, uberdig.Group("validation_request_observers")),
		// User validate handlers
		dig.NewProvider(request.NewRequestAccountValidationHandler),
		dig.NewProvider(validate.NewValidateAccountHandler),
		// User forgot password repository
		dig.NewProvider(forgot_password_repository.NewRedisForgotPasswordRepository),
		// User forgot password service
		dig.NewProvider(forgot_password_service.NewDefaultForgotPasswordRequestService),
		// User forgot password usecases
		dig.NewProvider(forgot_password_request_usecase.NewRequestForgotPasswordUseCase),
		dig.NewProvider(forgot_password_validate_usecase.NewValidateForgotPasswordUseCase),
		// User forgot password observers
		dig.NewProvider(forgot_password_request_usecase.NewLogRequestForgotPasswordObserver, uberdig.Group("forgot_password_request_observers")),
		dig.NewProvider(forgot_password_request_usecase.NewSendMailRequestForgotPasswordObserver, uberdig.Group("forgot_password_request_observers")),
		// User forgot password handlers
		dig.NewProvider(forgot_password_request_usecase.NewRequestForgotPasswordHandler),
		dig.NewProvider(forgot_password_validate_usecase.NewValidateForgotPasswordHandler),
		// User reset password repository
		dig.NewProvider(reset_password_repository.NewRedisResetPasswordRepository),
		// User reset password service
		dig.NewProvider(reset_password_service.NewDefaultResetPasswordService),
		// User reset password usecases
		dig.NewProvider(reset_password_request_usecase.NewRequestResetPasswordUseCase),
		dig.NewProvider(reset_password_validate_usecase.NewValidateResetPasswordUseCase),
		// User reset password observers
		dig.NewProvider(reset_password_request_usecase.NewLogRequestResetPasswordObserver, uberdig.Group("reset_password_request_observers")),
		dig.NewProvider(reset_password_request_usecase.NewSendMailRequestResetPasswordObserver, uberdig.Group("reset_password_request_observers")),
		// User reset password handlers
		dig.NewProvider(reset_password_request_usecase.NewRequestResetPasswordHandler),
		dig.NewProvider(reset_password_validate_usecase.NewValidateResetPasswordHandler),
		// User handlers
		dig.NewProvider(get_my_account.NewGetMyUserAccountHandler),
		dig.NewProvider(login.NewLoginHandler),
		dig.NewProvider(token.NewRefreshTokenHandler),
		dig.NewProvider(register.NewRegisterHandler),
		dig.NewProvider(logout.NewLogoutHandler),
		dig.NewProvider(update_user.NewUpdateAccountPersonalInformationsHandler),
		// User misc
		dig.NewProvider(token.NewJwtTokenStrategy(os.Getenv("JWT_SECRET"))),
		dig.NewProvider(crypt.NewBcryptStrategy),
		dig.NewProvider(middlewares.NewAuthMiddleware),
		// User chat direct
		// User chat direct repository
		dig.NewProvider(user_chat_direct_repository.NewMongoChatDirectRepository),
		dig.NewProvider(user_chat_direct_repository.NewChatDirectMapper),
		// User chat direct usecases
		dig.NewProvider(list_recent_chats_direct.NewListRecentChatDirectUseCase),
		// Mail usecases
		dig.NewProvider(sendmail.NewSendMailUseCase),
		// Ws
		dig.NewProvider(websocket.NewWsServer),
		dig.NewProvider(websocket.NewSaveMessageObserver, uberdig.Group("send_message_observers")),
		// Ws handlers
		dig.NewProvider(websocket.NewWebsocketHandler),
		// Chat
		// Chat recent
		// Chat recent usecases
		dig.NewProvider(list_recent_chats.NewListRecentChatsUseCase),
		dig.NewProvider(list_recent_chats.NewGroupMapper),
		dig.NewProvider(list_recent_chats.NewDirectChatMapper),
		// Chat recent handlers
		dig.NewProvider(list_recent_chats.NewListRecentChatsHandler),
		dig.NewProvider(list_recent_chats.NewResponseMapper),
		// Group
		// Group repository
		dig.NewProvider(group_repository.NewMongoGroupRepository),
		dig.NewProvider(group_repository.NewMongoGroupMapper),
		dig.NewProvider(group_repository.NewMongoGroupMemberMapper),
		// Group strategies
		dig.NewProvider(strategies.NewMembersNamesGroupNameStrategy),
		// Group usecases
		dig.NewProvider(add_member.NewAddMemberToGroupUseCase),
		// Group handlers
		dig.NewProvider(add_member.NewAddMemberToGroupHandler),
		// Group chats
		// Group chats repository
		dig.NewProvider(group_chat_message_repository.NewMongoGroupChatMessageRepository),
		dig.NewProvider(group_chat_message_repository.NewGroupChatMessageMapper),
		// Group chats usecases
		dig.NewProvider(list_recent_groups.NewListRecentGroupsUseCase),
		dig.NewProvider(list_group_chat_messages.NewListGroupChatMessagesUseCase),
		dig.NewProvider(save_group_chat_message.NewSaveGroupChatMessageUseCase),
		// Group chats handlers
		dig.NewProvider(list_group_chat_messages.NewListGroupChatMessagesHandler),
	}

	for _, provider := range providers {
		if err := di.Provide(provider.Constructor, provider.ProvideOptions...); err != nil {
			log.Fatalf("Unable to provide %s : %s", provider, err.Error())
		}

		funcType := reflect.TypeOf(provider.Constructor)
		firstReturn := funcType.Out(0)
		// If the first return is a pointer, get the name of the type it points to
		if firstReturn.Kind() == reflect.Ptr {
			firstReturn = firstReturn.Elem()
		}
		fnDetails := runtime.FuncForPC(reflect.ValueOf(provider.Constructor).Pointer()) // Retrieve function details
		fmt.Printf("[Dig] PROVIDE %20s <= %s\n", firstReturn.Name(), fnDetails.Name())
	}
	fmt.Printf("[Dig] %d providers provided\n", len(providers))

	return di
}
