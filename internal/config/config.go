package config

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Configuration contains all types of configurations
//easyjson:json
type Configuration struct {
	Server     ServerConfig      `json:"server"`
	Cors       CORSConfig        `json:"cors"`
	DataBase   DatabaseConfig    `json:"dataBase"`
	Storage    FileStorageConfig `json:"storage"`
	AWS        AwsPublicConfig   `json:"aws"`
	Game       GameConfig        `json:"game"`
	Session    SessionConfig     `json:"session"`
	WebSocket  WebSocketConfig   `json:"websocket"`
	AuthClient AuthClient        `json:"authClient"`
}

// ServerConfig set host, post and buffers sizes
//easyjson:json
type ServerConfig struct {
	Host      string `json:"host"`
	PortURL   string `json:"portUrl"`
	PortValue string `json:"portValue"`
}

// CORSConfig set allowable origins, headers and methods
//easyjson:json
type CORSConfig struct {
	Origins     []string `json:"origins"`
	Headers     []string `json:"headers"`
	Credentials string   `json:"credentials"`
	Methods     []string `json:"methods"`
}

// DatabaseConfig set type of database management system
//   the url of connection string, max amount of
//   connections, tables, sizes of page  of gamers
//   and users
//easyjson:json
type DatabaseConfig struct {
	DriverName       string `json:"driverName"`
	URL              string `json:"url"`
	ConnectionString string `json:"connectionString"`
	MaxOpenConns     int    `json:"maxOpenConns"`
	PageGames        int    `json:"pageGames"`
	PageUsers        int    `json:"pageUsers"`
}

// FileStorageConfig set, where avatars store and
//    what mode set to files/directories
//easyjson:json
type FileStorageConfig struct {
	PlayersAvatarsStorage string `json:"playersAvatarsStorage"`
	DefaultAvatar         string `json:"defaultAvatar"`
	Region                string `json:"region"`
	Endpoint              string `json:"endpoint"`
}

// AwsPublicConfig public aws information as region and endpoint
//easyjson:json
type AwsPublicConfig struct {
	AwsConfig *aws.Config `json:"-"`
	Region    string      `json:"region"`
	Endpoint  string      `json:"endpoint"`
}

// AwsPrivateConfig private aws information. Need another json.
//easyjson:json
type AwsPrivateConfig struct {
	AccessURL string `json:"accessUrl"`
	AccessKey string `json:"accessKey"`
	SecretURL string `json:"secretUrl"`
	SecretKey string `json:"secretKey"`
}

//easyjson:json
type FieldConfig struct {
	MinAreaSize    int `json:"minAreaSize"`
	MaxAreaSize    int `json:"maxAreaSize"`
	MinProbability int `json:"minProbability"`
	MaxProbability int `json:"maxProbability"`
}

// GameConfig set, how much rooms server can create and
// how much connections can join. Also there are flags:
// can server close rooms or not(for history mode),
// metrics should be recorded or not
//easyjson:json
type GameConfig struct {
	RoomsCapacity      int          `json:"roomsCapacity"`
	ConnectionCapacity int          `json:"connectionCapacity"`
	Location           string       `json:"location"`
	CanClose           bool         `json:"closeRoom"`
	Metrics            bool         `json:"metrics"`
	Field              *FieldConfig `json:"field"`
}

// AuthClient client of auth microservice
//easyjson:json
type AuthClient struct {
	URL     string `json:"url"`
	Address string `json:"address"`
}

// SessionConfig set cookie name, path, length, expiration time
// and HTTPonly flag
//easyjson:json
type SessionConfig struct {
	Name            string `json:"name"`
	Path            string `json:"path"`
	Length          int    `json:"length"`
	LifetimeSeconds int    `json:"lifetime"`
	HTTPOnly        bool   `json:"httpOnly"`
}

// WebSocketConfig set timeouts
//easyjson:json
type WebSocketConfig struct {
	WriteWait       int   `json:"writeWait"`
	PongWait        int   `json:"pongWait"`
	PingPeriod      int   `json:"pingPeriod"`
	MaxMessageSize  int64 `json:"maxMessageSize"`
	ReadBufferSize  int   `json:"readBufferSize"`
	WriteBufferSize int   `json:"writeBufferSize"`
}

// WebSocketSettings set timeouts
//easyjson:json
type WebSocketSettings struct {
	WriteWait       time.Duration `json:"writeWait"`
	PongWait        time.Duration `json:"pongWait"`
	PingPeriod      time.Duration `json:"pingPeriod"`
	MaxMessageSize  int64         `json:"maxMessageSize"`
	ReadBufferSize  int           `json:"readBufferSize"`
	WriteBufferSize int           `json:"writeBufferSize"`
}

func set(URL, value string) {
	if URL != "" && os.Getenv(URL) == "" {
		os.Setenv(URL, value)
	}
	utils.Debug(false, "environment -", URL, " :", value)
}

// InitEnvironment set environmental variables
func InitEnvironment(c *Configuration) {

	set(c.DataBase.URL, c.DataBase.ConnectionString)
	set(c.Server.PortURL, c.Server.PortValue)
	set(c.AuthClient.URL, c.AuthClient.Address)
}

// InitPublic load configuration file
func InitPublic(publicConfigPath string) (conf *Configuration, err error) {
	conf = &Configuration{}
	var data []byte
	if data, err = ioutil.ReadFile(publicConfigPath); err != nil {
		return
	}
	if err = conf.UnmarshalJSON(data); err != nil {
		return
	}
	conf.AWS.AwsConfig = &aws.Config{
		Region:   aws.String(conf.AWS.Region),
		Endpoint: aws.String(conf.AWS.Endpoint)}
	InitEnvironment(conf)
	return
}

// InitPrivate load configuration file and set private environment
func InitPrivate(privateConfigPath string) (err error) {
	var (
		data  []byte
		place = "secret json -"
	)

	if data, err = ioutil.ReadFile(privateConfigPath); err != nil {

		utils.Debug(false, place, "not found:", err.Error())
		return
	}
	var apc = &AwsPrivateConfig{}
	if err = apc.UnmarshalJSON(data); err != nil {
		utils.Debug(false, place, "wrong:", err.Error())
		return
	}

	os.Setenv(apc.AccessURL, apc.AccessKey)
	os.Setenv(apc.SecretURL, apc.SecretKey)
	return
}
