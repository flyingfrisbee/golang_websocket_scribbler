package environment

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	SQLFilePath string
	DBPath      string
)

func GenerateEnvVar() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	SQLFilePath = os.Getenv("SQL_FILE_PATH")
	DBPath = os.Getenv("DB_PATH")
}
