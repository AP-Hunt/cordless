package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Bios-Marcel/cordless/authentication"
	"github.com/Bios-Marcel/cordless/messagebus"
	"github.com/Bios-Marcel/cordless/tview"
	"github.com/Bios-Marcel/cordless/version"
	"github.com/Bios-Marcel/cordless/windowman"
	"github.com/Bios-Marcel/discordgo"

	"github.com/Bios-Marcel/cordless/commands/commandimpls"
	"github.com/Bios-Marcel/cordless/config"
	"github.com/Bios-Marcel/cordless/readstate"
	"github.com/Bios-Marcel/cordless/shortcuts"
	"github.com/Bios-Marcel/cordless/ui"
	"github.com/Bios-Marcel/cordless/ui/login"
	"github.com/Bios-Marcel/cordless/ui/updatedialog"
)

// StartApplication is the application's composition root. It sets up all known
// windows, and handles all of the wiring between them.
func StartApplication(windowManager windowman.WindowManagerInterface, accountToUse string) error {
	configuration := config.Current
	//App that will be reused throughout the process runtime.
	tviewApp := windowManager.GetUnderlyingApp()

	configDir, configErr := config.GetConfigDirectory()
	if configErr != nil {
		log.Fatalf("Unable to determine configuration directory (%s)\n", configErr.Error())
	}

	if accountToUse != nil && strings.TrimSpace(accountToUse) != "" {
		configuration.Token = configuration.GetAccountToken(account)
	}

	authenticator := authentication.Authenticator{}

	loginWindow := login.NewLoginWindow(configDir, authenticator)
	windowManager.RegisterWindow("login", &loginWindow)

	sess, ready, err := authenticator.AuthenticateFromLastSession(configuration, &accountToUse)
	if err != nil {
		log.Fatalf("Something went wrong trying to authenticate from the last session: %s", err)
	}

	// Can't auth from last session
	if sess == nil {
		windowManager.GetMessageBus().Subscribe(messagebus.TopicDiscordLoginSuccess, func(sess discordgo.Session, ready discordgo.Ready) {
			panic("I need to start up the message publishing for discord events and then show the main window")
			windowManager.ShowWindow("main-screen")
		})
		windowManager.ShowWindow("login")

		var firstWindow windowman.Window
		if accountToUse != "" {
			firstWindow = SetupApplicationWithAccount(tviewApp, accountToUse)
		} else {
			firstWindow = SetupApplication(tviewApp)
		}

		windowManager.RegisterWindow("root-screen", firstWindow)
		windowManager.ShowWindow("root-screen")
	} else {
		panic("I need to show the main window here")
	}

	configureUpdatesDialog(windowManager, *configuration)
	configureConfigMessageHandling(windowManager.GetMessageBus(), *configuration)

	return windowManager.Run()
}

func configureUpdatesDialog(windowManager windowman.WindowManagerInterface, configuration config.Config) {
	updatesDialog := updatedialog.NewUpdateDialog(configuration)
	updateChannel := version.CheckForUpdate(configuration.DontShowUpdateNotificationFor)
	version.EmitVersionCheckingToMessageBus(windowManager.GetMessageBus(), updateChannel)

	windowManager.GetMessageBus().Subscribe(version.MessageUpdateAvailable, func(available bool) {
		if available {
			windowManager.Dialog(&updatesDialog)
		}
	})
}

//configureConfigMessageHandling subscribes to events that result in config changes
func configureConfigMessageHandling(messageBus messagebus.MessageBus, configuration config.Config) {
	messageBus.Subscribe(messagebus.TopicDiscordLoginSuccess, func(sess discordgo.Session, ready discordgo.Ready) {
		configuration.Token = sess.Token
		config.PersistConfig()
	})
}

// SetupApplicationWithAccount launches the whole application and might
// abort in case it encounters an error. The login will attempt
// using the account specified, unless the argument is empty.
// If the account can't be found, the login page will be shown.
func SetupApplicationWithAccount(app *tview.Application, account string) windowman.Window {
	configuration := config.Current

	configDir, configErr := config.GetConfigDirectory()
	if configErr != nil {
		log.Fatalf("Unable to determine configuration directory (%s)\n", configErr.Error())
	}

	go func() {
		shortcutsLoadError := shortcuts.Load()
		if shortcutsLoadError != nil {
			panic(shortcutsLoadError)
		}

		config.Current.Token = discord.Token

		persistError := config.PersistConfig()
		if persistError != nil {
			app.Stop()
			log.Fatalf("Error persisting configuration (%s).\n", persistError.Error())
		}

		discord.State.MaxMessageCount = 100

		readstate.Load(discord.State)

		app.QueueUpdateDraw(func() {
			window, createError := ui.NewWindow(app, discord, readyEvent)

			if createError != nil {
				app.Stop()
				//Otherwise the logger output can't be seen, since we are stopping the TUI either way.
				log.SetOutput(os.Stdout)
				log.Fatalf("Error constructing window (%s).\n", createError.Error())
			}

			window.RegisterCommand(commandimpls.NewVersionCommand())
			statusGetCmd := commandimpls.NewStatusGetCommand(discord)
			statusSetCmd := commandimpls.NewStatusSetCommand(discord)
			statusSetCustomCmd := commandimpls.NewStatusSetCustomCommand(discord)
			window.RegisterCommand(statusSetCmd)
			window.RegisterCommand(statusGetCmd)
			window.RegisterCommand(statusSetCustomCmd)
			window.RegisterCommand(commandimpls.NewStatusCommand(statusGetCmd, statusSetCmd, statusSetCustomCmd))
			window.RegisterCommand(commandimpls.NewFileSendCommand(discord, window))
			accountLogout := commandimpls.NewAccountLogout(func() { SetupApplication(app) }, window)
			window.RegisterCommand(accountLogout)
			window.RegisterCommand(commandimpls.NewAccount(accountLogout, window))
			window.RegisterCommand(commandimpls.NewManualCommand(window))
			window.RegisterCommand(commandimpls.NewFixLayoutCommand(window))
			window.RegisterCommand(commandimpls.NewFriendsCommand(discord))
			userSetCmd := commandimpls.NewUserSetCommand(window, discord)
			userGetCmd := commandimpls.NewUserGetCommand(window, discord)
			window.RegisterCommand(userSetCmd)
			window.RegisterCommand(userGetCmd)
			window.RegisterCommand(commandimpls.NewUserCommand(userSetCmd, userGetCmd))
			serverJoinCmd := commandimpls.NewServerJoinCommand(window, discord)
			serverLeaveCmd := commandimpls.NewServerLeaveCommand(window, discord)
			serverCreateCmd := commandimpls.NewServerCreateCommand(discord)
			window.RegisterCommand(serverJoinCmd)
			window.RegisterCommand(serverLeaveCmd)
			window.RegisterCommand(serverCreateCmd)
			window.RegisterCommand(commandimpls.NewServerCommand(serverJoinCmd, serverLeaveCmd, serverCreateCmd))
			window.RegisterCommand(commandimpls.NewNickSetCmd(discord, window))
			tfaEnableCmd := commandimpls.NewTFAEnableCommand(window, discord)
			tfaDisableCmd := commandimpls.NewTFADisableCommand(discord)
			tfaBackupGetCmd := commandimpls.NewTFABackupGetCmd(discord, window)
			tfaBackupResetCmd := commandimpls.NewTFABackupResetCmd(discord, window)
			window.RegisterCommand(commandimpls.NewTFACommand(tfaEnableCmd, tfaDisableCmd, tfaBackupGetCmd, tfaBackupResetCmd))
			window.RegisterCommand(tfaEnableCmd)
			window.RegisterCommand(tfaDisableCmd)
			window.RegisterCommand(tfaBackupGetCmd)
			window.RegisterCommand(tfaBackupResetCmd)
			window.RegisterCommand(commandimpls.NewDMOpenCmd(discord, window))
		})
	}()

	return &loginWindow
}

// SetupApplication launches the whole application and might abort in case
// it encounters an error.
func SetupApplication(app *tview.Application) windowman.Window {
	return SetupApplicationWithAccount(app, "")
}

func attemptLogin(loginScreen *login.Login, loginMessage string, configuration *config.Config) (*discordgo.Session, *discordgo.Ready) {
	var (
		session      *discordgo.Session
		readyEvent   *discordgo.Ready
		discordError error
	)

	if configuration.Token == "" {
		session, discordError = loginScreen.RequestLogin(loginMessage)
	} else {
		session, discordError = discordgo.NewWithToken(configuration.Token)
	}

	if discordError != nil {
		configuration.Token = ""
		return attemptLogin(loginScreen, fmt.Sprintf("Error during last login attempt:\n\n[red]%s", discordError), configuration)
	}

	if session == nil {
		configuration.Token = ""
		return attemptLogin(loginScreen, "Error during last login attempt:\n\n[red]Received session is nil", configuration)
	}

	readyChan := make(chan *discordgo.Ready, 1)
	session.AddHandlerOnce(func(s *discordgo.Session, event *discordgo.Ready) {
		readyChan <- event
	})

	discordError = session.Open()

	if discordError != nil {
		configuration.Token = ""
		return attemptLogin(loginScreen, fmt.Sprintf("Error during last login attempt:\n\n[red]%s", discordError), configuration)
	}

	readyEvent = <-readyChan
	close(readyChan)

	return session, readyEvent
}
