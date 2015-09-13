package main

import (
	"fmt"
	// "github.com/desmondhume/furrytemple/job"
	// videoFetch "github.com/desmondhume/furrytemple/job/videos/fetch"
	"github.com/desmondhume/furrytemple/pageBuilder"
)

func main() {
	// output := make(chan map[string]interface{})
	// jobsReports := make(chan job.JobReport)

	// videoFetch.Run(output, jobsReports)
	err := pageBuilder.BuildHomepage()
	if err != nil {
		fmt.Println(err)
	}
	return
}
