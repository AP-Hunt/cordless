package authentication

import (
	"fmt"

	"github.com/Bios-Marcel/cordless/config"
	"github.com/Bios-Marcel/discordgo"
)

type Authenticator struct {
}

func (auth *Authenticator) AuthenticateWithToken(token string) (discordgo.Session, *discordgo.Ready, error) {
	return completeAuth(discordgo.NewWithToken(Token))

}

func (auth *Authenticator) AuthenicateWithCredentials(username string, password string, mfaToken string) (discordgo.Session, *discordgo.Ready, error) {
	return completeAuth(discordgo.NewWithPasswordAndMFA(username, password, mfaToken))
}

func (auth *Authenticator) AuthenticateFromLastSession(configuration *config.Config, accountToUse *string) (*discordgo.Session, *discordgo.Ready, error) {
	var token string
	if accountToUse != nil {
		token = configuration.GetAccountToken(*accountToUse)
	} else {
		token = configuration.Token
	}

	if token == "" {
		return nil, nil, nil
	}

	return completeAuth(discordgo.NewWithToken(token))
}

// completeAuth handles the shared logic between different Discord authentication scenarios
func completeAuth(session discordgo.Session, discordError error) (discordgo.Session, *discordgo.Ready, error) {
	if discordError != nil {
		return nil, nil, discordError
	}

	if session == nil {
		return nil, nil, fmt.Errorf("received session in nil")
	}

	readyChan := make(chan *discordgo.Ready, 1)
	session.AddHandlerOnce(func(s *discordgo.Session, event *discordgo.Ready) {
		readyChan <- event
	})

	discordError = session.Open()

	if discordError != nil {
		return nil, nil, fmt.Errorf("error opening discord session")
	}

	readyEvent = <-readyChan
	close(readyChan)

	return session, readyEvent, nil
}
