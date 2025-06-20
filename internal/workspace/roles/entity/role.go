package entity

import (
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

// Définition des permissions sous forme de bits
const (
	PermissionManageChannels          uint64 = 1 << iota // Gérer les canaux
	PermissionManageRoles                                // Gérer les rôles
	PermissionManageMessages                             // Gérer les messages
	PermissionSendMessages                               // Envoyer un message
	PermissionAttachFiles                                // Joindre des fichiers
	PermissionMentionEveryone                            // Mentionner @everyone
	PermissionKickMembers                                // Kick des membres
	PermissionInviteMembers                              // Inviter des membres
	PermissionManageWorkspaceSettings                    // Gérer les paramètres de l'espace de travail
)

type (
	RoleId      string
	WorkspaceId string
)

type Role struct {
	Id          RoleId
	Name        string
	WorkspaceId workspace_entity.WorkspaceId
	Permissions uint64
	Color       string
	IsAssigned  bool
}

type UserRole struct {
	UserId      string
	RoleId      RoleId
	WorkspaceId WorkspaceId
}

// Vérifie si un rôle possède une permissions spécifique
func (r Role) HasPermission(permission uint64) bool {
	return r.Permissions&permission != 0 //nolint:revive
}

func (id RoleId) String() string {
	return string(id)
}
