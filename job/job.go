package job

type JobReport struct {
	Status  string
	Message string
}

type Reports chan JobReport

type Job interface {
	Run(chan map[string]interface{}, Reports)
}
