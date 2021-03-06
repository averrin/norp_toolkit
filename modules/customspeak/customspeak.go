package customspeak

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	_ "log"

	yaml "gopkg.in/yaml.v1"

	"path"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/configor"
	"github.com/skratchdot/open-golang/open"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/qml"
)

var root string
var config = ConfigStruct{}
var bridge *MocBridge

//go:generate qtmoc
type MocBridge struct {
	core.QObject

	_ func(email string, password string) `signal:"serve"`
	_ func()                              `signal:"disconnect"`
	_ func()                              `signal:"offline"`
	_ func(err string)                    `signal:"error"`
	_ func(guild string)                  `signal:"setguild"`
	_ func(channel string)                `signal:"setchannel"`
}

func StartUI(view *qml.QQmlApplicationEngine) {
	bridge = NewMocBridge(nil)
	bridge.ConnectServe(func(email string, password string) {
		err := Serve(email, password)
		if err != nil {
			bridge.Error(fmt.Sprintf("Connection error: %v", err))
			bridge.Offline()
		} else {
			bridge.Error("")
		}
	})
	bridge.ConnectDisconnect(func() {
		sc <- nil
	})
	view.RootContext().SetContextProperty("customSpeak", bridge)
	view.RootContext().SetContextProperty2("initEmail", core.NewQVariant14(config.Email))
	view.RootContext().SetContextProperty2("initPassword", core.NewQVariant14(config.Password))
}

func Close() {
	sc <- nil
}

func Init() error {
	root, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	root = path.Join(root, "modules/customspeak")
	return configor.Load(&config, path.Join(root, "config.yml"))
}

var buffer = make([][]byte, 0)
var dg *discordgo.Session
var botSession *discordgo.Session
var me *discordgo.User
var bot *discordgo.User
var sc chan os.Signal
var handlerInstalled bool
var voiceConnection *discordgo.VoiceConnection
var mutex *sync.Mutex
var currentChannel string

func saveConfig() {
	d, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Writing file: %v\n", path.Join(root, "config.yml"))
	err = ioutil.WriteFile(path.Join(root, "config.yml"), d, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func createBot() string {
	// Create a new application.
	ap := &discordgo.Application{}
	ap.Name = fmt.Sprintf("%v's spy", me.Username)
	ap.Description = "NoRoleplaying Toolkit Discord spy"
	ap, err := dg.ApplicationCreate(ap)
	if err != nil {
		fmt.Println("error creating new applicaiton,", err)
		return ""
	}

	fmt.Printf("Application created successfully:\n")
	b, _ := json.MarshalIndent(ap, "", " ")
	fmt.Println(string(b))

	// Create the bot account under the application we just created
	bot, err := dg.ApplicationBotCreate(ap.ID)
	if err != nil {
		fmt.Println("error creating bot account,", err)
		return ""
	}

	fmt.Printf("Bot account created successfully.\n")
	b, _ = json.MarshalIndent(bot, "", " ")
	fmt.Println(string(b))
	return bot.Token
}

func Serve(email string, password string) error {
	handlerInstalled = false
	var err error
	dg, err = discordgo.New(email, password)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return err
	}
	config.Email = email
	config.Password = password
	saveConfig()

	me, _ = dg.User("@me")
	// Open the websocket and begin listening.
	go func() {
		err = dg.Open()
		if err != nil {
			fmt.Println("Error opening Discord session: ", err)
			bridge.Error(fmt.Sprintf("Connection error: %v", err))
			bridge.Offline()
			return
		}

		if config.Token == "" {
			config.Token = createBot()
		}
		if config.Token == "" {
			bridge.Error(fmt.Sprintf("Error Bot creation: %v", err))
			bridge.Offline()
			return
		}

		saveConfig()

		botSession, err = discordgo.New("Bot " + config.Token)
		if err != nil {
			bridge.Error(fmt.Sprintf("Error creating Bot session: %v", err))
			bridge.Offline()
			return
		}
		bot, _ = botSession.User("@me")
		if bot == nil {
			bridge.Error("Error creating Bot session: Bot user not found")
			bridge.Offline()
			return
		}
		dg.Close()

		img, err := ioutil.ReadFile(path.Join(root, "icon.png"))
		if err != nil {
			fmt.Println(err)
		}

		contentType := http.DetectContentType(img)
		base64img := base64.StdEncoding.EncodeToString(img)

		// Now lets format our base64 image into the proper format Discord wants
		// and then call UserUpdate to set it as our user's Avatar.
		avatar := fmt.Sprintf("data:%s;base64,%s", contentType, base64img)
		_, err = botSession.UserUpdate("", "", "", avatar, "")

		err = botSession.Open()
		if err != nil {
			fmt.Println("Error opening bot session: ", err)
			bridge.Error(fmt.Sprintf("Bot Connection error: %v", err))
			bridge.Offline()
			return
		}

		botSession.UpdateStatus(0, "NoRP Toolkit")

		botSession.AddHandler(messageCreate)
		botSession.AddHandler(ready)

		// Wait here until CTRL-C or other term signal is received.
		fmt.Println("CustomSpeak is now running.  Press CTRL-C to exit.")
		sc = make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		bridge.Offline()

		bridge.Setguild("")
		bridge.Setchannel("")

		// Cleanly close down the Discord session.
		botSession.Close()
		handlerInstalled = false
		fmt.Println("Disconnected")
	}()
	return nil
}

var userStates map[string]bool
var userHist map[string]bool
var users map[string]string
var channel *discordgo.Channel

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	userStates = map[string]bool{}
	userHist = map[string]bool{}
	if users == nil {
		users = map[string]string{}
	}

	fmt.Println("Ready handler")
	for _, g := range event.Guilds {
		for _, vs := range g.VoiceStates {
			if vs.UserID == me.ID {
				installHandler(botSession, vs)
				return
			}
		}
	}
	botSession.AddHandler(voiceStateUpdate)

}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message handler")
	if strings.HasPrefix(strings.ToLower(m.Content), "слава укр") {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Героям Слава!")
	}
}

func addMemeber(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	users[event.User.ID] = event.User.Username
}

func fetchUsers(guild *discordgo.Guild) {
	for _, m := range guild.Members {
		// u, _ := s.User(vs.UserID)
		users[m.User.ID] = m.User.Username
		fmt.Println(m.User.ID, "  ", m.User.Username)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func voiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if voiceConnection != nil && voiceConnection.ChannelID != "" {
		if vs.ChannelID == "" {
			fmt.Println("Disconnect from channel")
			voiceConnection.Disconnect()
			voiceConnection.Close()

			bridge.Setguild("")
			bridge.Setchannel("")
			currentChannel = ""
			handlerInstalled = false
			return
		}
	}

	if vs.UserID != me.ID {
		fmt.Println("Not me")
		return
	}

	fmt.Println("Voice handler")

	if handlerInstalled {
		fmt.Println("Already installed")
		return
	}
	fmt.Println("install from voice")
	installHandler(botSession, vs.VoiceState)
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	fmt.Println("Guild handler")

	for _, vs := range event.Guild.VoiceStates {
		if vs.UserID == me.ID {
			installHandler(botSession, vs)
			return
		}
	}
}

func installHandler(s *discordgo.Session, u *discordgo.VoiceState) {
	if handlerInstalled {
		return
	}

	fmt.Println("Installing handler")
	bridge.Error("")
	handlerInstalled = true
	mutex = &sync.Mutex{}
	c, err := s.State.Channel(u.ChannelID)
	if err != nil {
		// Could not find channel.
		bridge.Error(fmt.Sprintf("Channel fetch error: %v", err))
		bridge.Offline()
		sc <- nil
		return
	}
	guild, err := s.State.Guild(u.GuildID)

	// Find the guild for that channel.
	if err != nil {
		// Could not find guild.
		bridge.Error(fmt.Sprintf("Guild fetch error: %v", err))
		bridge.Offline()
		sc <- nil
		return
	}

	_, err = botSession.Guild(u.GuildID)
	if err != nil {
		bridge.Error("Waiting for invitation... Invitation page will be opened in your browser.")
		open.Run(fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%v&scope=bot", bot.ID))

		handlerInstalled = false
		botSession.AddHandler(guildCreate)
		return
	}

	folderName := fmt.Sprintf("%v_%v", guild.Name, c.Name)
	folderName = path.Join(root, folderName)
	bridge.Setguild(guild.Name)
	bridge.Setchannel(c.Name)

	// Look for the message sender in that guild's current voice states.
	fmt.Println("Fetching users")
	fetchUsers(guild)
	ch := make(chan Event, 50)
	os.Mkdir(folderName, 0774)
	vc, err := botSession.ChannelVoiceJoin(guild.ID, c.ID, false, false)
	currentChannel = c.ID
	voiceConnection = vc
	if err != nil {
		fmt.Printf("Voice channel connection error: %v\n", err)
		bridge.Error(fmt.Sprintf("Voice channel connection error: %v", err))
		bridge.Offline()
		sc <- nil
		return
	}
	vc.AddHandler(func(conn *discordgo.VoiceConnection, event *discordgo.VoiceSpeakingUpdate) {
		ch <- Event{event.UserID, event.Speaking}
	})

	go func(ch chan Event) {
		for handlerInstalled {
			e := <-ch
			mutex.Lock()
			userStates[e.Username] = e.Speak
			mutex.Unlock()
		}
		fmt.Println("Channel listener finished")
	}(ch)
	go func() {
		for handlerInstalled {
			mutex.Lock()
			states := map[string]bool{}
			for k, v := range userStates {
				states[k] = v
			}
			mutex.Unlock()
			for u, speak := range states {
				if userHist[u] == speak {
					continue
				}
				userHist[u] = speak

				userName, ok := users[u]
				if !ok {
					fmt.Print(users)
					panic(u)
				}
				fmt.Printf("[%v] %v speaks: %v\n", time.Now(), userName, speak)
				var cp string
				var src string
				if speak {
					src = "on.png"
					cp = path.Join(root, "custom", userName, src)
					if _, err := os.Stat(cp); err == nil {
						src = cp
					} else {
						src = path.Join(root, src)
					}
				} else {
					src = "off.png"
					cp = path.Join(root, "custom", userName, src)
					if _, err := os.Stat(cp); err == nil {
						src = cp
					} else {
						src = path.Join(root, src)
					}
				}
				dest := path.Join(folderName, userName+".png")
				os.Remove(dest)
				os.Link(src, dest)
			}
			// time.Sleep(200)
		}
		fmt.Println("Image handler finished")
	}()
}
