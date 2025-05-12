package entity

const (
	PermissionManager uint64 = 1 << iota // Générer un lien d’invitation
	PermissionAdmin                      // Gérer les utilisateurs (ajout/suppression, etc.)
)

type (
	JobsId string
)

type Job struct {
	Id                 JobsId
	Name               string
	Permissions        uint64
	OrganizationalOnly bool // True si le rôle n’a aucune permission fonctionnelle
	IsAssigned         bool // True si le rôle est assigné à un utilisateur
}

func (r Job) HasPermission(permission uint64) bool {
	return r.Permissions&permission != 0
}

func (id JobsId) String() string {
	return string(id)
}
