package format

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type Test struct {
	Parent   *Test
	Index    int
	FullName string
	Package  string
	Output   bytes.Buffer
	Done     bool
	Passed   bool
	Elapsed  time.Duration
	Tests    map[string]*Test
}

func (t *Test) Name() string {
	parts := splitTestName(t.FullName)
	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

type TestCapture struct {
	*Test
	emulate bool
	count   int
	ts      Time
}

func (c *TestCapture) handleEvent(e *TestEvent) {
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
			panic("received second binary start event")
		}
		c.Test = c.newTest(nil, e)
	case TestActionRun:
		if test == nil {
			panic(fmt.Sprintf("no parent for test: %s", e.Test))
		}
		if test.FullName == e.Test {
			panic(fmt.Sprintf("received second run event for test: %s", e.Test))
		}

		subTest := c.newTest(test, e)
		test.Tests[subTest.Name()] = subTest
	case TestActionOutput:
		if test == nil {
			panic(fmt.Sprintf("received output event for unstarted test: %s", e.Test))
		}
		test.Output.WriteString(e.Output)
	case TestActionFail:
		if test == nil {
			panic(fmt.Sprintf("received fail event for unstarted test: %s", e.Test))
		}

		test.Done = true
	case TestActionPass:
		if test == nil {
			panic(fmt.Sprintf("received pass event for unstarted test: %s", e.Test))
		}

		test.Done = true
		test.Passed = true
	}

	c.emulateDuration(e)
	c.ts = e.Time
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
		Index:    c.count,
		FullName: e.Test,
		Package:  e.Package,
		Elapsed:  e.Duration(),
		Tests:    make(map[string]*Test),
	}
	c.count++
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

		c.handleEvent(event)
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
