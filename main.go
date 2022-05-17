package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/quantfall/holmes"
)

const (
	AGENTMANAGERPROTO        = "https"
	AGENTMANAGERPORT         = 9000
	REGISTRATIONENDPOINT     = "/api/v1/agent"
	GETIDANDKEYENDPOINT      = "/api/v1/agent-id-key-by-name"
	GETCOMMANDSENDPOINT      = "/api/v1/incident-commands"
	COMMANDSRESPONSEENDPOINT = "/api/v1/incident-command/result"
	TLSCA                    = "ca.crt"
	TLSCRT                   = "client.crt"
	TLSKEY                   = "client.key"
)

var h = holmes.New("debug", "UTMStack")
var hostname string

type agentDetails struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type jobResult struct {
	JobId  int64  `json:"jobId"`
	Result string `json:"result"`
}

func main() {
	fmt.Println("loading")
	var error error
	hostname, error = os.Hostname()
	if error != nil {
		panic(error)
	}

	createDB()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	//router.SetTrustedProxies([]string{"192.168.1.2"})
	router.LoadHTMLGlob("templates/**/*.tmpl")
	router.Static("/assets", "./assets")

	store := cookie.NewStore([]byte("secret"))
	// Set session expiration time
	log.Println(store)
	store.Options(sessions.Options{MaxAge: 3600 * 24}) // 24hr
	router.Use(sessions.Sessions("auth", store))

	router.POST("/log-in", login)
	router.POST("/create-password", createPassword)

	router.GET("/set-password", func(c *gin.Context) {
		if isAuthenticated(c) {
			log.Println("auth exist")
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.HTML(http.StatusOK, "auth/set-password.tmpl", gin.H{
			"title": "Set Password",
		})
	})

	router.GET("/sign-in", func(c *gin.Context) {
		if isAuthenticated(c) {
			log.Println("auth exist")
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.HTML(http.StatusOK, "auth/sign-in.tmpl", gin.H{
			"title": "Sign In To UTMStack Agent",
		})
	})

	auth := router.Group("/")
	auth.Use(Authentication())
	{
		auth.GET("/", func(c *gin.Context) {
			modules := getModules()
			settings := getSettings()

			c.HTML(http.StatusOK, "app/index.tmpl", gin.H{
				"title":    "Home",
				"modules":  modules,
				"settings": settings,
			})
		})

		auth.GET("/logs-info", func(c *gin.Context) {
			data := getLogs()

			//log.Println(getLogs())
			c.HTML(http.StatusOK, "app/logs-info.tmpl", gin.H{
				"title": "Logs",
				"logs":  data,
			})
		})

		auth.GET("/settings", func(c *gin.Context) {
			data := getSettings()

			c.HTML(http.StatusOK, "app/settings.tmpl", gin.H{
				"title":    "Settings",
				"settings": data,
			})
		})

		auth.GET("/change-password", func(c *gin.Context) {
			c.HTML(http.StatusOK, "auth/change-password.tmpl", gin.H{
				"title": "",
			})
		})

		auth.POST("/log-out", logout)

		router.POST("/change-password", changePassword)
		router.POST("/save-settings", saveSettings)
		router.POST("/edit-module", editModule)
	}

	router.Run(":8080")

	/*if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "run":
			incidentResponse()
			startBeat()
			startWazuh()
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
			<-signals
			stopWazuh()

		case "install":
			var ip string
			var utmKey string
			var skip string

			fmt.Println("Manager IP or FQDN:")
			if _, err := fmt.Scanln(&ip); err != nil {
				h.Error("can't get the manager IP or FQDN: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			fmt.Println("Registration Key:")
			if _, err := fmt.Scanln(&utmKey); err != nil {
				h.Error("can't get the registration key: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			fmt.Println("Skip certificate validation (yes or no):")
			if _, err := fmt.Scanln(&skip); err != nil {
				h.Error("can't get certificate validation response: %v", err)
				time.Sleep(10 * time.Second)
				os.Exit(1)
			}

			install(ip, utmKey, skip)

		case "silent-install":
			ip := os.Args[2]
			utmKey := os.Args[3]
			skip := os.Args[4]

			install(ip, utmKey, skip)

		default:
			fmt.Println("unknown option")
		}
	} else {
		err := uninstall()
		if err != nil {
			h.Error("can't remove agent dependencies or configurations: %v", err)
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}

		os.Exit(0)
	}*/
}
