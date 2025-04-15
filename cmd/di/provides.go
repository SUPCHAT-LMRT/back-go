package di

import (
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/back_identifier/usecase"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/usecase/list_recent_chats"
	"github.com/supchat-lmrt/back-go/internal/dig"
	"github.com/supchat-lmrt/back-go/internal/event"
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
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/search/usecase/search"
	"github.com/supchat-lmrt/back-go/internal/search/user"
	user_chat_direct_repository "github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/is_first_message"
	list_direct_messages "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_messages"
	list_recent_chats_direct "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_recent_direct_chats"
	save_direct_message "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/save_message"
	toggle_chat_direct_reaction "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/toggle_reaction"
	"github.com/supchat-lmrt/back-go/internal/user/gin/middlewares"
	mongo2 "github.com/supchat-lmrt/back-go/internal/user/repository/mongo"
	user_status_repository "github.com/supchat-lmrt/back-go/internal/user/status/repository"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_or_create_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/save_status"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists"
	forgot_password_repository "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/repository"
	forgot_password_service "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/service"
	forgot_password_request_usecase "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/usecase/request"
	forgot_password_validate_usecase "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/usecase/validate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_my_account"
	repository2 "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
	delete2 "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/delete"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/generate"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/get_data_token_invite"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login_oauth"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/logout"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/public_profile"
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
	"github.com/supchat-lmrt/back-go/internal/utils"
	"github.com/supchat-lmrt/back-go/internal/websocket"
	chat_message_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	time_series_message_sent_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/usecase/get_minutely"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/count_messages_by_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/list_messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/save_message"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/toggle_reaction"
	channel_repository "github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/count_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/create_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/reoder_channels"
	workspace_middlewares "github.com/supchat-lmrt/back-go/internal/workspace/gin/middlewares"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	add_member2 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	repository3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/repository"
	delete3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/delete"
	generate3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/generate"
	get_data_token_invite3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/get_data_token_invite"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/join_workspace_invite"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/is_user_in_workspace"
	list_workpace_members2 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/list_workpace_members"
	workspace_repository "github.com/supchat-lmrt/back-go/internal/workspace/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/create_workspace"
	discovery_list_workspaces "github.com/supchat-lmrt/back-go/internal/workspace/usecase/discover/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace_details"
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
		// Meilisearch
		dig.NewProvider(meilisearch.NewClient),
		// Mailer
		dig.NewProvider(mail.NewMailer(os.Getenv("SMTP_HOST"), os.Getenv("SMTP_TLS") == "true", utils.MustAtoi(os.Getenv("SMTP_PORT")), os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))),
		// Identifier workspace
		dig.NewProvider(usecase.NewGetBackIdentifierUseCase),
		dig.NewProvider(usecase.NewHostnameBackIdentifierStrategy),
		// Workspaces
		// Workspaces repository
		dig.NewProvider(workspace_repository.NewMongoWorkspaceRepository),
		dig.NewProvider(workspace_repository.NewMongoWorkspaceMapper),
		dig.NewProvider(repository.NewMongoWorkspaceMemberMapper),
		dig.NewProvider(repository3.NewRedisInviteLinkRepository),
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
		dig.NewProvider(list_workpace_members2.NewListWorkspaceMembersUseCase),
		dig.NewProvider(generate3.NewInviteLinkUseCase),
		dig.NewProvider(get_workspace.NewGetWorkspaceUseCase),
		dig.NewProvider(get_data_token_invite3.NewGetInviteLinkDataUseCase),
		// Workspace handlers
		dig.NewProvider(list_workspaces.NewListWorkspaceHandler),
		dig.NewProvider(discovery_list_workspaces.NewDiscoverListWorkspaceHandler),
		dig.NewProvider(create_workspace.NewCreateWorkspaceHandler),
		dig.NewProvider(get_workspace_details.NewGetWorkspaceDetailsHandler),
		dig.NewProvider(update_icon.NewUpdateWorkspaceIconHandler),
		dig.NewProvider(update_banner.NewUpdateWorkspaceBannerHandler),
		dig.NewProvider(list_workpace_members2.NewListWorkspaceHandler),
		dig.NewProvider(generate3.NewCreateInviteLinkHandler),
		dig.NewProvider(get_data_token_invite3.NewGetInviteLinkWorkspaceDataHandler),
		// Workspace mappers
		dig.NewProvider(repository3.NewRedisInviteLinkMapper),
		// Workspace channels
		// Workspace channels repository
		dig.NewProvider(channel_repository.NewMongoChannelRepository),
		dig.NewProvider(channel_repository.NewMongoChannelMapper),
		// Workspace channels usecases
		dig.NewProvider(list_channels.NewListChannelsUseCase),
		dig.NewProvider(create_channel.NewCreateChannelUseCase),
		dig.NewProvider(get_channel.NewGetChannelUseCase),
		dig.NewProvider(count_channels.NewCountChannelsUseCase),
		dig.NewProvider(reoder_channels.NewReorderChannelsUseCase),
		// Workspaces channels observers
		dig.NewProvider(create_channel.NewNotifyWebSocketObserver, uberdig.Group("create_channel_observers")),
		dig.NewProvider(reoder_channels.NewUserStatusUpdateObserver, uberdig.Group("reorder_channels_observers")),
		// Workspace channels handlers
		dig.NewProvider(list_channels.NewListChannelsHandler),
		dig.NewProvider(create_channel.NewCreateChannelHandler),
		dig.NewProvider(get_channel.NewGetChannelHandler),
		dig.NewProvider(reoder_channels.NewReorderChannelHandler),
		// Workspace channels chat
		// Workspace channels chat repository
		dig.NewProvider(chat_message_repository.NewMongoChannelMessageRepository),
		dig.NewProvider(chat_message_repository.NewChannelMessageMapper),
		// Workspace channels chat usecases
		dig.NewProvider(list_messages.NewListMessageUseCase),
		dig.NewProvider(toggle_reaction.NewToggleReactionChannelMessageUseCase),
		dig.NewProvider(count_messages_by_workspace.NewCountMessagesUseCase),
		// Workspace channels chat handlers
		dig.NewProvider(list_messages.NewListChannelMessagesHandler),
		// Workspace channels chat messages usecases
		dig.NewProvider(save_message.NewSaveChannelMessageUseCase),
		// Workspace time series
		// Workspace time series message sent
		// Workspace time series message sent repository
		dig.NewProvider(time_series_message_sent_repository.NewMongoMessageSentTimeSeriesWorkspaceRepository),
		// Workspace time series message sent usecases
		dig.NewProvider(get_minutely.NewGetMinutelyMessageSentUseCase),
		// Workspace time series message sent handlers
		dig.NewProvider(get_minutely.NewGetMinutelyMessageSentHandler),
		// Workspace misc
		dig.NewProvider(workspace_middlewares.NewUserInWorkspaceMiddleware),
		// Workspace member usecases
		dig.NewProvider(add_member2.NewAddMemberUseCase),
		dig.NewProvider(delete3.NewDeleteInviteLinkWorkspaceUseCase),
		// Workspace member handlers
		dig.NewProvider(join_workspace_invite.NewJoinWorkspaceInviteUseCase),
		dig.NewProvider(join_workspace_invite.NewJoinWorkspaceInviteHandler),
		// Workspace member repository
		dig.NewProvider(repository.NewMongoWorkspaceMemberRepository),
		// User
		dig.NewProvider(mongo2.NewMongoUserRepository),
		dig.NewProvider(mongo2.NewMongoUserMapper),
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
		dig.NewProvider(public_profile.NewGetPublicUserProfileUseCase),
		// User handlers
		dig.NewProvider(update_user_avatar.NewUpdateUserAvatarHandler),
		dig.NewProvider(public_profile.NewGetPublicProfileHandler),
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
		// User invite link repository
		dig.NewProvider(repository2.NewRedisInviteLinkRepository),
		dig.NewProvider(repository2.NewRedisInviteLinkMapper),
		// User invite link usecases
		dig.NewProvider(generate.NewInviteLinkUseCase),
		dig.NewProvider(get_data_token_invite.NewGetInviteLinkDataUseCase),
		dig.NewProvider(delete2.NewDeleteInviteLinkUseCase),
		// User invite link handlers
		dig.NewProvider(generate.NewCreateInviteLinkHandler),
		dig.NewProvider(get_data_token_invite.NewGetInviteLinkDataHandler),
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
		dig.NewProvider(save_direct_message.NewSaveDirectMessageUseCase),
		dig.NewProvider(is_first_message.NewIsFirstMessageUseCase),
		dig.NewProvider(list_direct_messages.NewListDirectMessagesUseCase),
		dig.NewProvider(toggle_chat_direct_reaction.NewToggleReactionDirectMessageUseCase),
		// USer chat direct handlers
		dig.NewProvider(list_direct_messages.NewListDirectMessagesHandler),
		// User Oauth handler & usecase
		dig.NewProvider(login_oauth.NewOAuthHandler),
		dig.NewProvider(login_oauth.NewOAuthUseCase),
		// User status
		// User status repositories
		dig.NewProvider(user_status_repository.NewMongoUserStatusRepository),
		// User status usecases
		dig.NewProvider(save_status.NewSaveStatusUseCase),
		dig.NewProvider(save_status.NewUserStatusUpdateObserver, uberdig.Group("save_user_status_observers")),
		dig.NewProvider(get_status.NewGetStatusUseCase),
		dig.NewProvider(get_or_create_status.NewGetOrCreateStatusUseCase),
		dig.NewProvider(get_public_status.NewGetPublicStatusUseCase),
		// User status handlers
		dig.NewProvider(save_status.NewSaveStatusHandler),
		// Mail usecases
		dig.NewProvider(sendmail.NewSendMailUseCase),
		// Event bus
		dig.NewProvider(event.NewEventBus),
		// Ws
		dig.NewProvider(websocket.NewWsServer),
		dig.NewProvider(websocket.NewSaveChannelMessageObserver, uberdig.Group("send_channel_message_observers")),
		dig.NewProvider(websocket.NewSaveDirectMessageObserver, uberdig.Group("send_direct_message_observers")),
		dig.NewProvider(save_direct_message.NewSyncRecentChatObserver, uberdig.Group("save_direct_chat_message_observers")),
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
		// Search
		// Search usecases
		dig.NewProvider(search.NewSearchTermUseCase),
		// Search handlers
		dig.NewProvider(search.NewSearchTermHandler),
		// Search sync managers
		dig.NewProvider(message.NewMeilisearchSearchMessageSyncManager),
		dig.NewProvider(channel.NewMeilisearchSearchChannelSyncManager),
		dig.NewProvider(user.NewMeilisearchSearchUserSyncManager),
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
