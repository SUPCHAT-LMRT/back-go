package di

import (
	"fmt"
	get_last_message2 "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/get_last_message"
	is_first_message2 "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/is_first_message"
	toggle_reaction2 "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/toggle_reaction"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/create_group"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/get_member_by_user"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/group_info"
	kick_member2 "github.com/supchat-lmrt/back-go/internal/group/usecase/kick_member"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/leave_group"
	list_members "github.com/supchat-lmrt/back-go/internal/group/usecase/list_members_users"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mention/usecase/extract_mentions"
	"github.com/supchat-lmrt/back-go/internal/mention/usecase/list_mentionnable_user"
	"github.com/supchat-lmrt/back-go/internal/notification/usecase/list_notifications"
	"github.com/supchat-lmrt/back-go/internal/notification/usecase/mark_as_read"
	send_notification2 "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/send_notification"
	"log"
	"os"
	"reflect"
	"runtime"

	"github.com/supchat-lmrt/back-go/internal/back_identifier/usecase"
	poll "github.com/supchat-lmrt/back-go/internal/bots/poll/repository"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/usecase/create_poll"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/usecase/delete_poll"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/usecase/get_poll_by_id"
	get_polls_listpackage "github.com/supchat-lmrt/back-go/internal/bots/poll/usecase/get_polls_list"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/usecase/unvote_option_poll"
	"github.com/supchat-lmrt/back-go/internal/bots/poll/usecase/vote_option_poll"
	"github.com/supchat-lmrt/back-go/internal/chat/recent/usecase/list_recent_chats"
	"github.com/supchat-lmrt/back-go/internal/dig"
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/gin"
	group_chat_message_repository "github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	list_group_chat_messages "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/list_messages"
	save_group_chat_message "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/save_message"
	group_repository "github.com/supchat-lmrt/back-go/internal/group/repository"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/list_recent_groups"
	"github.com/supchat-lmrt/back-go/internal/logger/zerolog"
	"github.com/supchat-lmrt/back-go/internal/mail"
	"github.com/supchat-lmrt/back-go/internal/meilisearch"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	mongo3 "github.com/supchat-lmrt/back-go/internal/notification/repository/mongo"
	"github.com/supchat-lmrt/back-go/internal/notification/usecase/create_notification"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	"github.com/supchat-lmrt/back-go/internal/search/usecase/search"
	"github.com/supchat-lmrt/back-go/internal/search/user"
	has_job "github.com/supchat-lmrt/back-go/internal/user/app_jobs/gin/middlewares"
	app_jobs "github.com/supchat-lmrt/back-go/internal/user/app_jobs/repository"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/assign_job"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/create_job"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/delete_job"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/get_job_for_user"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/list_jobs"
	permissions2 "github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/permissions"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/unassign_job"
	"github.com/supchat-lmrt/back-go/internal/user/app_jobs/usecase/update_job"
	user_chat_direct_repository "github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/get_last_message"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/is_first_message"
	list_direct_messages "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_messages"
	list_recent_chats_direct "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/list_recent_direct_chats"
	save_direct_message "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/save_message"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/send_notification"
	toggle_chat_direct_reaction "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/toggle_reaction"
	"github.com/supchat-lmrt/back-go/internal/user/gin/middlewares"
	mongo2 "github.com/supchat-lmrt/back-go/internal/user/repository/mongo"
	user_status_repository "github.com/supchat-lmrt/back-go/internal/user/status/repository"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_or_create_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_status"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/save_status"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/delete_user"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists_by_email"
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
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/get_list_invite_link"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/list_all_users"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/login"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/logout"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/oauth"
	user_oauth_repository "github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/repository"
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
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/delete_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_private_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_user_private_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/reoder_channels"
	workspace_middlewares "github.com/supchat-lmrt/back-go/internal/workspace/gin/middlewares"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	add_member2 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_user_by_workspace_member_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_workpace_member"
	repository3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/repository"
	delete3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/delete"
	generate3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/generate"
	get_data_token_invite3 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/get_data_token_invite"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/usecase/join_workspace_invite"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/is_user_in_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/kick_member"
	list_workpace_members2 "github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/list_workspace_members"
	workspace_repository "github.com/supchat-lmrt/back-go/internal/workspace/repository"
	has_permissions "github.com/supchat-lmrt/back-go/internal/workspace/roles/gin/middlewares"
	roles_repository "github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/assign_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/check_permissions"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/create_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/delete_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/dessassign_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/get_list_roles"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/get_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/get_roles_for_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/permissions"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/update_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/create_workspace"
	discovery_list_workspaces "github.com/supchat-lmrt/back-go/internal/workspace/usecase/discover/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_member_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/get_workspace_details"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/list_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_banner"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_icon"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_info_workspaces"
	"github.com/supchat-lmrt/back-go/internal/workspace/usecase/update_type_workspace"
	uberdig "go.uber.org/dig"
)

//nolint:revive
func NewDi() *uberdig.Container {
	isProd := os.Getenv("ENV") == "production"

	di := uberdig.New()
	providers := []dig.Provider{
		// Logger
		dig.NewProvider(zerolog.NewZerologLogger(
			logger.WithMinLevel(utils.IfThenElse(isProd, logger.LogLevelInfo, logger.LogLevelTrace)),
		)),
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
		dig.NewProvider(
			mail.NewMailer(
				os.Getenv("SMTP_HOST"),
				os.Getenv("SMTP_TLS") == "true",
				utils.MustAtoi(os.Getenv("SMTP_PORT")),
				os.Getenv("SMTP_USERNAME"),
				os.Getenv("SMTP_PASSWORD"),
			),
		),
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
		dig.NewProvider(discovery_list_workspaces.NewDiscoverListWorkspacesUseCase),
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
		dig.NewProvider(
			generate.NewSendMailGenerateInviteLinkObserver,
			uberdig.Group("generate_invite_link_observers"),
		),
		dig.NewProvider(get_workspace.NewGetWorkspaceUseCase),
		dig.NewProvider(get_data_token_invite3.NewGetInviteLinkDataUseCase),
		dig.NewProvider(update_info_workspaces.NewUpdateInfoWorkspacesUseCase),
		dig.NewProvider(update_type_workspace.NewUpdateTypeWorkspaceUseCase),
		dig.NewProvider(kick_member.NewKickMemberUseCase),
		dig.NewProvider(get_member_id.NewGetMemberIdUsecase),
		// Workspace handlers
		dig.NewProvider(list_workspaces.NewListWorkspaceHandler),
		dig.NewProvider(discovery_list_workspaces.NewDiscoverListWorkspaceHandler),
		dig.NewProvider(create_workspace.NewCreateWorkspaceHandler),
		dig.NewProvider(get_workspace_details.NewGetWorkspaceDetailsHandler),
		dig.NewProvider(update_icon.NewUpdateWorkspaceIconHandler),
		dig.NewProvider(update_banner.NewUpdateWorkspaceBannerHandler),
		dig.NewProvider(list_workpace_members2.NewListWorkspaceHandler),
		dig.NewProvider(generate3.NewCreateInviteLinkHandler),
		dig.NewProvider(get_workspace.NewGetWorkspaceHandler),
		dig.NewProvider(get_data_token_invite3.NewGetInviteLinkWorkspaceDataHandler),
		dig.NewProvider(update_info_workspaces.NewUpdateInfoWorkspacesHandler),
		dig.NewProvider(update_type_workspace.NewUpdateTypeWorkspaceHandler),
		dig.NewProvider(kick_member.NewKickGroupMemberHandler),
		dig.NewProvider(get_member_id.NewGetMemberIdHandler),
		// Workspace observers
		dig.NewProvider(
			update_info_workspaces.NewUpdateInfoWorkspacesObserver,
			uberdig.Group("update_info_workspaces_observers"),
		),
		dig.NewProvider(
			update_icon.NewUpdateWorkspaceIconObserver,
			uberdig.Group("update_icon_workspace_observers"),
		),
		dig.NewProvider(
			update_type_workspace.NewNotifyUpdateTypeWorkspaceObserver,
			uberdig.Group("update_type_workspace_observers"),
		),
		dig.NewProvider(
			update_banner.NewUpdateWorkspaceBannerObserver,
			uberdig.Group("save_banner_workspace_observers"),
		),
		// Workspace mappers
		dig.NewProvider(repository3.NewRedisInviteLinkMapper),
		// Workspace channels
		// Workspace channels repository
		dig.NewProvider(channel_repository.NewMongoChannelRepository),
		dig.NewProvider(channel_repository.NewMongoChannelMapper),
		// Workspace roles repository
		dig.NewProvider(roles_repository.NewMongoRoleRepository),
		dig.NewProvider(roles_repository.NewMongoRoleMapper),
		// Workspace roles usecases
		dig.NewProvider(create_role.NewCreateRoleHandler),
		dig.NewProvider(get_role.NewGetRoleUseCase),
		dig.NewProvider(get_list_roles.NewGetListRolesUseCase),
		dig.NewProvider(update_role.NewUpdateRoleUseCase),
		dig.NewProvider(delete_role.NewDeleteRoleUseCase),
		// Workspace roles handlers
		dig.NewProvider(create_role.NewCreateRoleUseCase),
		dig.NewProvider(get_role.NewGetRoleHandler),
		dig.NewProvider(get_list_roles.NewGetListRolesHandler),
		dig.NewProvider(update_role.NewUpdateRoleHandler),
		dig.NewProvider(delete_role.NewDeleteRoleHandler),
		dig.NewProvider(assign_role.NewAssignRoleToUserHandler),
		dig.NewProvider(dessassign_role.NewDessassignRoleFromUserHandler),
		dig.NewProvider(get_roles_for_member.NewGetRolesForMemberHandler),
		dig.NewProvider(check_permissions.NewCheckPermissionsHandler),
		dig.NewProvider(list_user_private_channel.NewListPrivateChannelMembersHandler),
		// Workspace channels usecases
		dig.NewProvider(list_channels.NewListChannelsUseCase),
		dig.NewProvider(list_private_channels.NewGetPrivateChannelsUseCase),
		dig.NewProvider(create_channel.NewCreateChannelUseCase),
		dig.NewProvider(get_channel.NewGetChannelUseCase),
		dig.NewProvider(count_channels.NewCountChannelsUseCase),
		dig.NewProvider(reoder_channels.NewReorderChannelsUseCase),
		dig.NewProvider(delete_channels.NewDeleteChannelUseCase),
		dig.NewProvider(assign_role.NewAssignRoleToUserUsecase),
		dig.NewProvider(dessassign_role.NewDessassignRoleFromUserUsecase),
		dig.NewProvider(get_roles_for_member.NewGetRolesForMemberUsecase),
		dig.NewProvider(permissions.NewCheckPermissionUseCase),
		dig.NewProvider(list_user_private_channel.NewListPrivateChannelMembersUseCase),
		// Workspaces channels observers
		dig.NewProvider(
			create_channel.NewCreateChannelObserver,
			uberdig.Group("create_channel_observers"),
		),
		dig.NewProvider(
			reoder_channels.NewUserStatusUpdateObserver,
			uberdig.Group("reorder_channels_observers"),
		),
		dig.NewProvider(
			delete_channels.NewDeleteChannelsObserver,
			uberdig.Group("delete_channels_observers"),
		),
		// Workspace channels handlers
		dig.NewProvider(list_channels.NewListChannelsHandler),
		dig.NewProvider(list_private_channels.NewGetPrivateChannelsHandler),
		dig.NewProvider(create_channel.NewCreateChannelHandler),
		dig.NewProvider(get_channel.NewGetChannelHandler),
		dig.NewProvider(reoder_channels.NewReorderChannelHandler),
		dig.NewProvider(delete_channels.NewDeleteChannelHandler),
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
		dig.NewProvider(
			time_series_message_sent_repository.NewMongoMessageSentTimeSeriesWorkspaceRepository,
		),
		// Workspace time series message sent usecases
		dig.NewProvider(get_minutely.NewGetMinutelyMessageSentUseCase),
		// Workspace time series message sent handlers
		dig.NewProvider(get_minutely.NewGetMinutelyMessageSentHandler),
		// Workspace misc
		dig.NewProvider(workspace_middlewares.NewUserInWorkspaceMiddleware),
		dig.NewProvider(has_permissions.NewHasPermissionsMiddleware),
		// Workspace member usecases
		dig.NewProvider(add_member2.NewAddMemberUseCase),
		dig.NewProvider(add_member2.NewAddMemberHandler),
		dig.NewProvider(delete3.NewDeleteInviteLinkWorkspaceUseCase),
		dig.NewProvider(get_user_by_workspace_member_id.NewGetUserByWorkspaceMemberIdUseCase),
		// Workspace member handlers
		dig.NewProvider(join_workspace_invite.NewJoinWorkspaceInviteUseCase),
		dig.NewProvider(join_workspace_invite.NewJoinWorkspaceInviteHandler),
		// Workspace member repository
		dig.NewProvider(repository.NewMongoWorkspaceMemberRepository),
		// Workspace notification
		dig.NewProvider(send_notification2.NewSendMessageNotificationUseCase),
		dig.NewProvider(send_notification2.NewEmailChannel, uberdig.Group("send_channelmessage_notification_channel")),
		dig.NewProvider(send_notification2.NewPushChannel, uberdig.Group("send_channelmessage_notification_channel")),
		dig.NewProvider(save_message.NewSendNotificationObserver, uberdig.Group("save_channel_message_observers")),
		dig.NewProvider(save_message.NewGetMentionObserver, uberdig.Group("save_channel_message_observers")),
		// User
		dig.NewProvider(mongo2.NewMongoUserRepository),
		dig.NewProvider(mongo2.NewMongoUserMapper),
		// User usecases
		dig.NewProvider(get_by_id.NewGetUserByIdUseCase),
		dig.NewProvider(get_by_email.NewGetUserByEmailUseCase),
		dig.NewProvider(login.NewLoginUserUseCase),
		dig.NewProvider(exists_by_email.NewExistsUserByEmailUseCase),
		dig.NewProvider(register.NewRegisterUserUseCase),
		dig.NewProvider(token.NewRefreshAccessTokenUseCase),
		dig.NewProvider(update_password.NewChangePasswordUseCase),
		dig.NewProvider(update_user.NewUpdateUserUseCase),
		dig.NewProvider(update_user_avatar.NewUpdateUserAvatarUseCase),
		dig.NewProvider(update_user_avatar.NewS3UpdateUserAvatarStrategy),
		dig.NewProvider(public_profile.NewGetPublicUserProfileUseCase),
		dig.NewProvider(delete_user.NewDeleteUserUseCase),
		dig.NewProvider(list_all_users.NewListUserUseCase),
		// User handlers
		dig.NewProvider(update_user_avatar.NewUpdateUserAvatarHandler),
		dig.NewProvider(public_profile.NewGetPublicProfileHandler),
		dig.NewProvider(list_all_users.NewListUserHandler),
		// User forgot password repository
		dig.NewProvider(forgot_password_repository.NewRedisForgotPasswordRepository),
		// User forgot password service
		dig.NewProvider(forgot_password_service.NewDefaultForgotPasswordRequestService),
		// User forgot password usecases
		dig.NewProvider(forgot_password_request_usecase.NewRequestForgotPasswordUseCase),
		dig.NewProvider(forgot_password_validate_usecase.NewValidateForgotPasswordUseCase),
		// User forgot password observers
		dig.NewProvider(
			forgot_password_request_usecase.NewLogRequestForgotPasswordObserver,
			uberdig.Group("forgot_password_request_observers"),
		),
		dig.NewProvider(
			forgot_password_request_usecase.NewSendMailRequestForgotPasswordObserver,
			uberdig.Group("forgot_password_request_observers"),
		),
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
		dig.NewProvider(
			reset_password_request_usecase.NewLogRequestResetPasswordObserver,
			uberdig.Group("reset_password_request_observers"),
		),
		dig.NewProvider(
			reset_password_request_usecase.NewSendMailRequestResetPasswordObserver,
			uberdig.Group("reset_password_request_observers"),
		),
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
		dig.NewProvider(get_list_invite_link.NewGetListInviteLinkUseCase),
		// User invite link handlers
		dig.NewProvider(generate.NewCreateInviteLinkHandler),
		dig.NewProvider(get_data_token_invite.NewGetInviteLinkDataHandler),
		dig.NewProvider(get_list_invite_link.NewGetListInviteLinkHandler),
		dig.NewProvider(delete2.NewDeleteInviteLinkHandler),
		// User handlers
		dig.NewProvider(get_my_account.NewGetMyUserAccountHandler),
		dig.NewProvider(login.NewLoginHandler),
		dig.NewProvider(token.NewRefreshTokenHandler),
		dig.NewProvider(register.NewRegisterHandler),
		dig.NewProvider(logout.NewLogoutHandler),
		dig.NewProvider(update_user.NewUpdateAccountPersonalInformationsHandler),
		dig.NewProvider(delete_user.NewDeleteUserHandler),
		// User misc
		dig.NewProvider(token.NewJwtTokenStrategy(os.Getenv("JWT_SECRET"))),
		dig.NewProvider(crypt.NewBcryptStrategy),
		dig.NewProvider(middlewares.NewAuthMiddleware),
		// User Oauth repositories
		dig.NewProvider(user_oauth_repository.NewMongoOauthConnectionRepository),
		dig.NewProvider(user_oauth_repository.NewMongoOauthConnectionMapper),
		// User chat direct
		// User chat direct repository
		dig.NewProvider(user_chat_direct_repository.NewMongoChatDirectRepository),
		dig.NewProvider(user_chat_direct_repository.NewChatDirectMapper),
		// User chat direct usecases
		dig.NewProvider(list_recent_chats_direct.NewListRecentChatDirectUseCase),
		dig.NewProvider(get_last_message.NewGetLastDirectChatMessageUseCase),
		dig.NewProvider(save_direct_message.NewSaveDirectMessageUseCase),
		dig.NewProvider(is_first_message.NewIsFirstMessageUseCase),
		dig.NewProvider(list_direct_messages.NewListDirectMessagesUseCase),
		dig.NewProvider(toggle_chat_direct_reaction.NewToggleReactionDirectMessageUseCase),
		// USer chat direct handlers
		dig.NewProvider(list_direct_messages.NewListDirectMessagesHandler),
		// User Oauth handler & usecase
		dig.NewProvider(oauth.NewRegisterOAuthHandler),
		dig.NewProvider(oauth.NewLoginOAuthUseCase),
		dig.NewProvider(oauth.NewRegisterOAuthUseCase),
		// User status
		// User status repositories
		dig.NewProvider(user_status_repository.NewMongoUserStatusRepository),
		// User status usecases
		dig.NewProvider(save_status.NewSaveStatusUseCase),
		dig.NewProvider(
			save_status.NewUserStatusUpdateObserver,
			uberdig.Group("save_user_status_observers"),
		),
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
		dig.NewProvider(
			websocket.NewSaveChannelMessageObserver,
			uberdig.Group("send_channel_message_observers"),
		),
		dig.NewProvider(
			websocket.NewSaveDirectMessageObserver,
			uberdig.Group("send_direct_message_observers"),
		),
		dig.NewProvider(
			save_direct_message.NewSyncRecentChatObserver,
			uberdig.Group("save_direct_chat_message_observers"),
		),
		dig.NewProvider(
			save_direct_message.NewSendNotificationObserver,
			uberdig.Group("save_direct_chat_message_observers"),
		),
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
		// Group usecases
		dig.NewProvider(add_member.NewAddMemberToGroupUseCase),
		dig.NewProvider(create_group.NewCreateGroupUseCase),
		dig.NewProvider(list_members.NewListGroupMembersUseCase),
		dig.NewProvider(group_info.NewGetGroupInfoUseCase),
		dig.NewProvider(kick_member2.NewKickMemberUseCase),
		dig.NewProvider(create_group.NewSyncRecentChatObserver, uberdig.Group("group_created_observer")),
		dig.NewProvider(add_member.NewSyncRecentChatObserver, uberdig.Group("add_group_member_observer")),
		dig.NewProvider(kick_member2.NewSyncRecentChatObserver, uberdig.Group("kick_group_member_observer")),
		dig.NewProvider(get_member_by_user.NewGetMemberByUserUseCase),
		// Group handlers
		dig.NewProvider(add_member.NewAddMemberToGroupHandler),
		dig.NewProvider(create_group.NewCreateGroupHandler),
		dig.NewProvider(group_info.NewGetGroupInfoHandler),
		dig.NewProvider(leave_group.NewLeaveGroupHandler),
		dig.NewProvider(kick_member2.NewKickMemberHandler),
		// Group chats
		// Group chats repository
		dig.NewProvider(group_chat_message_repository.NewMongoGroupChatRepository),
		dig.NewProvider(group_chat_message_repository.NewMongoGroupChatMessageMapper),
		// Group chats usecases
		dig.NewProvider(list_recent_groups.NewListRecentGroupsUseCase),
		dig.NewProvider(list_group_chat_messages.NewListGroupChatMessagesUseCase),
		dig.NewProvider(save_group_chat_message.NewSaveGroupChatMessageUseCase),
		dig.NewProvider(get_last_message2.NewGetLastGroupChatMessageUseCase),
		dig.NewProvider(is_first_message2.NewIsFirstGroupChatMessageUseCase),
		dig.NewProvider(toggle_reaction2.NewToggleGroupChatReactionUseCase),
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
		// Jobs
		// Jobs repository
		dig.NewProvider(app_jobs.NewMongoJobRepository),
		dig.NewProvider(app_jobs.NewMongoJobMapper),
		// Jobs middlewares
		dig.NewProvider(has_job.NewHasJobPermissionsMiddleware),
		// Jobs usecases
		dig.NewProvider(create_job.NewCreateJobUseCase),
		dig.NewProvider(delete_job.NewDeleteJobUseCase),
		dig.NewProvider(update_job.NewUpdateJobUseCase),
		dig.NewProvider(list_jobs.NewListJobsUseCase),
		dig.NewProvider(assign_job.NewAssignJobUseCase),
		dig.NewProvider(unassign_job.NewUnassignJobUseCase),
		dig.NewProvider(get_job_for_user.NewGetJobForUserUseCase),
		dig.NewProvider(permissions2.NewCheckPermissionJobUseCase),
		// Jobs handlers
		dig.NewProvider(create_job.NewCreateJobHandler),
		dig.NewProvider(delete_job.NewDeleteJobHandler),
		dig.NewProvider(update_job.NewUpdateJobHandler),
		dig.NewProvider(list_jobs.NewListJobsHandler),
		dig.NewProvider(assign_job.NewAssignJobHandler),
		dig.NewProvider(unassign_job.NewUnassignJobHandler),
		dig.NewProvider(get_job_for_user.NewGetJobForUserHandler),
		// Notifications
		dig.NewProvider(create_notification.NewCreateNotificationUseCase),
		dig.NewProvider(mongo3.NewMongoNotificationRepository),
		dig.NewProvider(mongo3.NewMongoNotificationMapper),
		dig.NewProvider(send_notification.NewEmailChannel, uberdig.Group("send_directmessage_notification_channel")),
		dig.NewProvider(send_notification.NewPushChannel, uberdig.Group("send_directmessage_notification_channel")),
		dig.NewProvider(send_notification.NewSendMessageNotificationUseCase),
		dig.NewProvider(permissions2.NewCheckUserPermissionsHandler),
		dig.NewProvider(list_notifications.NewListNotificationsHandler),
		dig.NewProvider(list_notifications.NewListNotificationsUseCase),
		dig.NewProvider(mark_as_read.NewMarkAsReadUseCase),
		dig.NewProvider(mark_as_read.NewMarkAsReadHandler),

		// bots
		// bots repository
		dig.NewProvider(poll.NewMongoPollMapper),
		dig.NewProvider(poll.NewMongoPollRepository),
		// bots usecases
		dig.NewProvider(create_poll.NewCreatePollUseCase),
		dig.NewProvider(get_poll_by_id.NewGetPollByIdUseCase),
		dig.NewProvider(get_polls_listpackage.NewGetPollsListUseCase),
		dig.NewProvider(delete_poll.NewDeletePollUseCase),
		dig.NewProvider(vote_option_poll.NewVoteOptionPollUseCase),
		dig.NewProvider(unvote_option_poll.NewUnvoteOptionPollUseCase),
		// bots handlers
		dig.NewProvider(create_poll.NewCreatePollHandler),
		dig.NewProvider(get_poll_by_id.NewGetPollByIdHandler),
		dig.NewProvider(get_polls_listpackage.NewGetPollsListHandler),
		dig.NewProvider(delete_poll.NewDeletePollHandler),
		dig.NewProvider(vote_option_poll.NewVoteOptionPollHandler),
		dig.NewProvider(unvote_option_poll.NewUnvoteOptionPollHandler),
		// Mentions
		dig.NewProvider(list_mentionnable_user.NewListMentionnableUserUseCase),
		dig.NewProvider(list_mentionnable_user.NewListMentionnableUserHandler),
		dig.NewProvider(extract_mentions.NewExtractMentionsUseCase),
	}

	for _, provider := range providers {
		if err := di.Provide(provider.Constructor, provider.ProvideOptions...); err != nil {
			//nolint:revive
			log.Fatalf("Unable to provide %s : %s", provider, err.Error())
		}

		funcType := reflect.TypeOf(provider.Constructor)
		firstReturn := funcType.Out(0)
		// If the first return is a pointer, get the name of the type it points to
		if firstReturn.Kind() == reflect.Ptr {
			firstReturn = firstReturn.Elem()
		}
		fnDetails := runtime.FuncForPC(
			reflect.ValueOf(provider.Constructor).Pointer(),
		) // Retrieve function details
		fmt.Printf("[Dig] PROVIDE %20s <= %s\n", firstReturn.Name(), fnDetails.Name())
	}
	fmt.Printf("[Dig] %d providers provided\n", len(providers))

	return di
}
