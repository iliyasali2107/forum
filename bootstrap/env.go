package bootstrap

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Env struct {
	ServerAddress   string
	DBName          string
	DBPath          string
	DBDriver        string
	Port            int
	ContextTimeout  int
	LoginExpireTime time.Duration
}

func NewEnv() *Env {
	env := Env{
		ServerAddress:   ":8080",
		DBName:          "forum",
		DBPath:          "./db/forum.db",
		DBDriver:        "sqlite3",
		Port:            8080,
		ContextTimeout:  2,
		LoginExpireTime: 30 * time.Minute,
	}
	envRows := make(map[string]any)

	file, err := os.Open("./.env")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "=")
		envRows[row[0]] = row[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	serverAddress, ok := envRows["SERVER_ADDRESS"]
	if ok {
		env.ServerAddress = serverAddress.(string)
	}

	dbName, ok := envRows["DB_NAME"]
	if ok {
		env.DBName = dbName.(string)
	}

	dbPath, ok := envRows["DB_PATH"]
	if ok {
		env.DBPath = dbPath.(string)
	}

	dbDriver, ok := envRows["DB_DRIVER"]
	if ok {
		env.DBDriver = dbDriver.(string)
	}

	env.Port = 8080
	portStr, ok := envRows["PORT"].(string)
	portInt, err := strconv.Atoi(portStr)
	if ok && err == nil {
		env.Port = portInt
	}

	contextTimeoutStr, ok := envRows["CONTEXT_TIMEOUT"].(string)
	contextTimeoutInt, err := strconv.Atoi(contextTimeoutStr)
	if ok && err == nil {
		env.ContextTimeout = contextTimeoutInt
	}

	expireTimeStr, ok := envRows["LOGIN_EXPIRE_TIME"].(string)
	expireTimeDuration, err := time.ParseDuration(expireTimeStr)
	if ok && err == nil {
		env.LoginExpireTime = expireTimeDuration
	}

	return &env
}
