package tests

import (
	"log"

	"api-tests-template/internal/utils"
)

func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

func TearDownSuite() {
	log.Println("Tear down suite")
}

func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}
