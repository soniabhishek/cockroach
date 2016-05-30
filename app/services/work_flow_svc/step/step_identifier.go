package step

type StepIdentifier uint

const (
	CrowdSourcing = iota + 1
	InternalSourcing
	Transformation
	Algorithm
	Bifurcation
	Unification
	Manual
	Nil
)
