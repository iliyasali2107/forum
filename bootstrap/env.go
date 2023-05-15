package bootstrap

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Env struct {
	ServerAddress  string
	DBName         string
	DBPath         string
	DBDriver       string
	Port           int
	ContextTimeout int
}

func NewEnv() *Env {
	env := Env{}
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

	env.ServerAddress = envRows["SERVER_ADDRESS"].(string)
	env.DBName = envRows["DB_NAME"].(string)
	env.DBPath = envRows["DB_PATH"].(string)
	env.DBDriver = envRows["DB_DRIVER"].(string)
	env.Port = envRows["PORT"].(int)
	env.ContextTimeout = envRows["CONTEXT_TIMEOUT"].(int)

	return &env
}
