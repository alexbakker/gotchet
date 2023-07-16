package format

import "time"

// Source: https://pkg.go.dev/cmd/test2json

type TestAction string

const (
	// The test binary is about to be executed
	TestActionStart TestAction = "start"
	// The test has started running
	TestActionRun TestAction = "run"
	// The test has been paused
	TestActionPause TestAction = "pause"
	// The test has continued running
	TestActionCont TestAction = "cont"
	// The test passed
	TestActionPass TestAction = "pass"
	// The benchmark printed log output but did not fail
	TestActionBench TestAction = "bench"
	// The test or benchmark failed
	TestActionFail TestAction = "fail"
	// The test printed output
	TestActionOutput TestAction = "output"
	// The test was skipped or the package contained no tests
	TestActionSkip TestAction = "skip"
)

type Time time.Time

type TestEvent struct {
	Time    Time
	Action  TestAction
	Package string
	Test    string
	Elapsed float64
	Output  string
}

func (t *Time) UnmarshalText(text []byte) (err error) {
	parsed, err := time.Parse(time.RFC3339, string(text))
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}

func (e *TestEvent) Duration() time.Duration {
	return time.Duration(e.Elapsed * float64(time.Second))
}
