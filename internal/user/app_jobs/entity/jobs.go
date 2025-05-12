package entity

const (
	CREATE_INVITATION      uint64 = 1 << iota // Générer un lien d’invitation
	DELETE_INVITATION                         // Supprimer un lien d’invitation
	ASSIGN_JOB                                // Assigner un job à un utilisateur
	UNASSIGN_JOB                              // Désassigner un job d’un utilisateur
	DELETE_JOB                                // Supprimer un job
	UPDATE_JOB                                // Mettre à jour un job
	UPDATE_JOB_PERMISSIONS                    // Mettre à jour les permissions d’un job
)

type (
	JobsId string
)

type Job struct {
	Id          JobsId
	Name        string
	Permissions uint64
	IsAssigned  bool
}

func (r Job) HasPermission(permission uint64) bool {
	return r.Permissions&permission != 0
}

func (id JobsId) String() string {
	return string(id)
}
