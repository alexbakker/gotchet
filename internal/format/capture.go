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
	Parent   *Test            `json:"parent"`
	Index    int              `json:"index"`
	FullName string           `json:"full_name"`
	Package  string           `json:"package"`
	Output   []*Output        `json:"output"`
	Done     bool             `json:"done"`
	Passed   bool             `json:"passed"`
	Elapsed  time.Duration    `json:"elapsed"`
	Tests    map[string]*Test `json:"tests"`
}

type TestCapture struct {
	*Test
	emulate     bool
	testCount   int
	outputCount int
	ts          Time
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
	test := c.Test
	if test != nil {
		for _, name := range splitTestName(e.Test) {
			if subTest, ok := test.Tests[name]; ok {
				test = subTest
				continue
			}

			break
		}
	}

	switch e.Action {
	case TestActionStart:
		if c.Test != nil {
			return errors.New("received second binary start event")
		}
		c.Test = c.newTest(nil, e)
	case TestActionRun:
		if test == nil {
			return fmt.Errorf("no parent for test: %s", e.Test)
		}
		if test.FullName == e.Test {
			return fmt.Errorf("received second run event for test: %s", e.Test)
		}

		subTest := c.newTest(test, e)
		test.Tests[subTest.Name()] = subTest
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
	case TestActionPass:
		if test == nil {
			return fmt.Errorf("received pass event for unstarted test: %s", e.Test)
		}

		test.Done = true
		test.Passed = true
		test.Elapsed = time.Duration(e.Elapsed)
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
		//Parent:   parent,
		Index:    c.testCount,
		FullName: e.Test,
		Package:  e.Package,
		Tests:    make(map[string]*Test),
	}
	c.testCount++
	return
}

// TODO: When reading from stdin, make this write to disk and seek instead of
// keeping everything in-memory
func Read(r io.Reader, emulate bool) (*TestCapture, error) {
	var c TestCapture

	br := bufio.NewReader(r)
	for {
		event, err := readEvent(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return &c, nil
			}
			return nil, err
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
