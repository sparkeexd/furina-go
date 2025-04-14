package action

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
)

// Cron job structure holding the job definition and the function task to run by the bot.
type CronJob struct {
	Definition gocron.JobDefinition
	Task       gocron.Task
	Option     gocron.JobOption
	CronTab    string
}

// Service's cron jobs to be registered.
type JobService interface {
	Jobs(session *discordgo.Session) []CronJob
}

// Create a new cron job model.
func NewCronJob(definition gocron.JobDefinition, task gocron.Task, option gocron.JobOption, cronTab string) CronJob {
	return CronJob{definition, task, option, cronTab}
}
