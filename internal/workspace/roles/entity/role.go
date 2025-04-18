package entity

import (
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

// Définition des permissions sous forme de bits
const (
	PermissionManageChannels  uint64 = 1 << iota // Gérer les canaux
	PermissionManageRoles                        // Gérer les rôles
	PermissionManageMessages                     // Gérer les messages
	PermissionManageInvites                      // Gérer les invitations
	PermissionSendMessages                       // Envoyer un message
	PermissionAttachFiles                        // Joindre des fichiers
	PermissionPinMessages                        // Épingler des messages
	PermissionMentionEveryone                    // Mentionner @everyone
	PermissionKickMembers                        // Kick des membres
	PermissionInviteMembers                      // Inviter des membres
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
}

// Vérifie si un rôle possède une permission spécifique
func (r Role) HasPermission(permission uint64) bool {
	return r.Permissions&permission != 0
}

func (id RoleId) String() string {
	return string(id)
}
