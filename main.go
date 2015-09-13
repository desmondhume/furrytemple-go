package main

import (
	"github.com/desmondhume/furrytemple/job"
	videoFetch "github.com/desmondhume/furrytemple/job/videos/fetch"
)

// Structure of parser and normalizer

// Parser
// 	-> find video source (youtube, reddit, ...)
// 	-> pass video to correct normalizer
// --> CHANNEL -->
// Normalizer
// 	-> normalize video
// 	-> return unified video struct to the factory
//  --> CHANNEL -->
// Carpenter
// 	-> save the video inside the database

func main() {
	output := make(chan map[string]interface{})
	jobsReports := make(chan job.JobReport)

	videoFetch.Run(output, jobsReports)
	return
}
