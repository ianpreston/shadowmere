package puzzle

type Service interface {
	Nick() string
	OnPrivmsg(string, string)
}