package bulk_processor

type Job struct {
	Do func()
}

func NewJob(do func()) Job {
	return Job{
		Do: do,
	}
}