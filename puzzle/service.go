package puzzle

type Service interface {
	Nick() string
	OnPrivmsg(string, string)
	OnQuit(string, string)
	OnNickChange(string, string)
}