package step

type StepIdentifier uint

const (
	CrowdSourcing = iota
	InternalSourcing
	Manual
	Transformation
	Algorithm
	Bifurcation
	Unification
	Nil
)
