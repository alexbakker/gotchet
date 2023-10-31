package format

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

type Test struct {
	Parent    *Test         `json:"-"`
	Index     int           `json:"index"`
	StartedAt Time          `json:"started_at"`
	EndedAt   Time          `json:"ended_at"`
	FullName  string        `json:"full_name"`
	Package   string        `json:"package"`
	Output    []*Output     `json:"output"`
	Done      bool          `json:"done"`
	Skipped   bool          `json:"skipped"`
	Passed    bool          `json:"passed"`
	Elapsed   time.Duration `json:"elapsed"`
	Tests     []*Test       `json:"tests"`
}

type TestCapture struct {
	Tests            []*Test `json:"tests"`
	Title            string  `json:"title"`
	StartedAt        Time    `json:"started_at"`
	EndedAt          Time    `json:"ended_at"`
	CaptureStartedAt Time    `json:"capture_started_at"`
	CaptureEndedAt   Time    `json:"capture_ended_at"`
	emulate          bool
	testCount        int
	outputCount      int
	ts               Time
	tests            map[string]*Test
}

type Output struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

func (t *Test) Name() string {
	parts := splitTestName(t.FullName)
	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

func (t *Test) FullOutput() *bytes.Buffer {
	var outputs []*Output
	outputs = append(outputs, t.Output...)

	var merge func(t *Test)
	merge = func(t *Test) {
		for _, test := range t.Tests {
			outputs = append(outputs, test.Output...)
			merge(test)
		}
	}
	merge(t)

	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].Index < outputs[j].Index
	})

	var buf bytes.Buffer
	for _, output := range outputs {
		buf.WriteString(output.Text)
	}

	return &buf
}

func (c *TestCapture) handleEvent(e *TestEvent) error {
	test := c.tests[e.Package]
	if test != nil {
		if subTest, ok := c.tests[e.Test]; ok {
			test = subTest
		} else if subTest, ok := c.tests[e.ParentTest()]; ok {
			test = subTest
		}
	}

	switch e.Action {
	case TestActionStart:
		if test != nil {
			return fmt.Errorf("received second binary start event for: %s", e.Package)
		}
		test = c.newTest(nil, e)
		c.tests[e.Package] = test
		c.Tests = append(c.Tests, test)
	case TestActionRun:
		if test == nil {
			return fmt.Errorf("no parent for test: %s", e.Test)
		}
		if test.FullName == e.Test {
			return fmt.Errorf("received second run event for test: %s", e.Test)
		}

		subTest := c.newTest(test, e)
		c.tests[subTest.FullName] = subTest
		test.Tests = append(test.Tests, subTest)
	case TestActionOutput:
		if test == nil {
			return fmt.Errorf("received output event for unstarted test: %s", e.Test)
		}
		test.Output = append(test.Output, &Output{Index: c.outputCount, Text: e.Output})
		c.outputCount++
	case TestActionFail:
		if test == nil {
			return fmt.Errorf("received fail event for unstarted test: %s", e.Test)
		}

		test.Done = true
		test.Elapsed = time.Duration(e.Elapsed)
		test.EndedAt = e.Time
	case TestActionPass:
		if test == nil {
			return fmt.Errorf("received pass event for unstarted test: %s", e.Test)
		}

		test.Done = true
		test.Passed = true
		test.Elapsed = time.Duration(e.Elapsed)
		test.EndedAt = e.Time
	case TestActionSkip:
		if test == nil {
			return fmt.Errorf("received skip event for unstarted test: %s", e.Test)
		}

		test.Done = true
		test.Skipped = true
		test.Elapsed = time.Duration(e.Elapsed)
		test.EndedAt = e.Time
	}

	c.emulateDuration(e)
	c.ts = e.Time
	return nil
}

func (c *TestCapture) emulateDuration(e *TestEvent) {
	if c.emulate && !time.Time(c.ts).IsZero() && !time.Time(e.Time).IsZero() {
		if dur := time.Time(e.Time).Sub(time.Time(c.ts)); dur > 0 {
			time.Sleep(dur)
		}
	}
}

func (c *TestCapture) newTest(parent *Test, e *TestEvent) (test *Test) {
	test = &Test{
		Parent:    parent,
		Index:     c.testCount,
		StartedAt: e.Time,
		FullName:  e.Test,
		Package:   e.Package,
		// Explicitly initialize the tests slice to that we never get a nil in the JSON output
		Tests: make([]*Test, 0),
	}
	c.testCount++
	return
}

// TODO: When reading from stdin, make this write to disk and seek instead of
// keeping everything in-memory
func Read(r io.Reader, emulate bool) (*TestCapture, error) {
	c := TestCapture{
		CaptureStartedAt: Time(time.Now()),
		// Explicitly initialize the tests slice to that we never get a nil in the JSON output
		Tests: make([]*Test, 0),
		tests: make(map[string]*Test),
	}

	var lastTs Time
	br := bufio.NewReader(r)
	for {
		event, err := readEvent(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				c.EndedAt = lastTs
				c.CaptureEndedAt = Time(time.Now())
				return &c, nil
			}
			return nil, err
		}

		lastTs = event.Time
		if time.Time(c.StartedAt).IsZero() {
			c.StartedAt = lastTs
		}

		if err := c.handleEvent(event); err != nil {
			return nil, fmt.Errorf("handle test event: %w", err)
		}
	}
}

func readEvent(r *bufio.Reader) (*TestEvent, error) {
	var lineBuf bytes.Buffer
	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			return nil, err
		}

		if isPrefix {
			lineBuf.Write(line)
			continue
		}
		if lineBuf.Len() > 0 {
			lineBuf.Write(line)
			line = lineBuf.Bytes()
			lineBuf.Reset()
		}

		var event TestEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, err
		}

		return &event, nil
	}
}

func splitTestName(name string) []string {
	return strings.Split(name, "/")
}
