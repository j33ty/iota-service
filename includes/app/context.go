package app

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	iota "github.com/iotaledger/iota.go/api"
	pow "github.com/iotaledger/iota.go/pow"
	app "github.com/j33ty/iota-service/models/app"
	misc "gitlab.cdmarvels.com/cdtargaryens/libs/misc"
)

// Context maintains the main app context
type Context struct {
	AppName string
	Config  *app.Config
	IOTA    *iota.API
}

var instance *Context
var once sync.Once

// GetInstance - Gets the instance of the context
func GetInstance() *Context {
	once.Do(func() {
		instance = &Context{}
	})
	return instance
}

// Initialize prepares and returns the app context
func (c *Context) Initialize(appName string) error {

	c.AppName = appName
	log.Println("Initializing app: " + appName)

	err := InitConfig(c)
	if err != nil {
		return err
	}

	err = InitIOTA(c)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)
	go func() {
		<-sigChan
		c.Shutdown()
		os.Exit(0)
	}()

	return nil
}

// Shutdown - Gracefully terminate all the variables in App Context
func (c *Context) Shutdown() error {
	log.Println("Shutting down app...")
	return nil
}

// InitConfig - Reads dev.json or production.json from filesystem as TIER env
func InitConfig(c *Context) error {
	err := misc.ReadConfigFile(&c.Config)
	if c.Config == nil || err != nil {
		return err
	}
	return nil
}

// InitIOTA - Reads dev.json or production.json from filesystem as TIER env
func InitIOTA(c *Context) error {

	var err error

	_, powFunc := pow.GetBestPoW()
	c.IOTA, err = iota.ComposeAPI(iota.HttpClientSettings{
		URI:          c.Config.Endpoint,
		LocalPowFunc: powFunc,
	})
	if err != nil {
		return err
	}
	return nil
}
