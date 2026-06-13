package notification

type Provider interface {
	SendNotification(subject, message string) error
}
