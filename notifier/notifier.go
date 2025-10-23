package notifier

// Notifier 通知接口
type Notifier interface {
	Send(title, content string) error
}
