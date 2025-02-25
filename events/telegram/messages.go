package telegram

const (
	msgHelp = `I can save and keep your links. Also I can offer you them to read.
	
Available commands:
/help - show this help message
/start - show welcome message
/rnd - get a random link from your list. It Will be deleted from your reading list after that!`

	msgHello          = "Hi there! 👋\n\n" + msgHelp
	msgNoSavedPages   = "You have no saved pages 🙊"
	msgSaved          = "Url saved! 👌"
	msgAlreadyExists  = "You already have this page in your list 🤗"
	msgCommandUnknown = "Sorry, i don't understant this command"
)
