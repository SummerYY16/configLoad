package reconf

// 通知应用程序文件改变

// Notifyer 通知观察者
type Notifyer  interface {
	Callback(*Config)
}
