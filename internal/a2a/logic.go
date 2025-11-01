package a2a

import (
	"log"
	"time"
)

func Remember() {
	log.Println("Starting daily birthday check...")

	checkTomorrowBirthdays()

	checkTodayBirthdays()
}
func checkTomorrowBirthdays() {
	tomorrow := time.Now().AddDate(0, 0, 1)

	log.Printf("🔔 Checking for reminders - tomorrow is %s %d",
		tomorrow.Month().String(), tomorrow.Day())
}

// checkTodayBirthdays sends birthday wishes for birthdays happening today
func checkTodayBirthdays() {
	today := time.Now()

	log.Printf("🎂 Checking for birthdays - today is %s %d",
		today.Month().String(), today.Day())
}
