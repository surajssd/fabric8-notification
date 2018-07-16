package types

type NotificationType string

const (
	WorkitemCreate       NotificationType = "workitem.create"
	WorkitemUpdate       NotificationType = "workitem.update"
	CommentCreate        NotificationType = "comment.create"
	CommentUpdate        NotificationType = "comment.update"
	UserEmailUpdate      NotificationType = "user.email.update"
	InvitationTeamNoorg  NotificationType = "invitation.team.noorg"
	InvitationSpaceNoorg NotificationType = "invitation.space.noorg"
	AnalyticsNotifyCVE   NotificationType = "analytics.notify.cve"
)

var notifiers = map[NotificationType][]string{
	AnalyticsNotifyCVE: {"fabric8-gemini-server"},
	UserEmailUpdate:    {"fabric8-auth"},
}

func (nType NotificationType) Notifiers() []string {
	return notifiers[nType]
}
