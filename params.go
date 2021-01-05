package main

import (
	"math"
)

// Params represents params for GoThree.js library.
type Params struct {
	Angle          float64 `json:"angle"`
	AngleSecond    float64 `json:"angle2"`
	Caps           bool    `json:"allCaps"`
	Distance       int     `json:"distance"`
	DistanceSecond int     `json:"distance2"`
	AutoAngle      bool    `json:"autoAngle"`
}

func GuessParams(c *Commands) *Params {
	// TODO(justin): This is a bit messy, the idea is that there are a bunch of unused goroutines showing up
	// in the output and this filters them out, but there should be a cleaner way to do this.
	activeNames := map[string]struct{}{}
	for _, cmd := range c.cmds {
		if cmd.Parent != "" || cmd.From != "" || cmd.To != "" || cmd.Channel != "" {
			activeNames[cmd.Name] = struct{}{}
		}
	}

	var cmds []*Command
	for _, cmd := range c.cmds {
		if cmd.Name == "#1" && cmd.Time == 0 {
			// Skip vapid main CMD
			continue
		}
		if cmd.Command == CmdCreate {
			if _, found := activeNames[cmd.Name]; !found && cmd.Name != "#1" {
				// Skip duplicate CMD
				continue
			}
		}
		cmds = append(cmds, cmd)
	}
	c.cmds = cmds

	goroutines := make(map[int]int) // map[depth]quantity
	var totalG int

	// calculate number of goroutines in each depth level
	for _, cmd := range c.cmds {
		if cmd.Command == CmdCreate {
			totalG++
			goroutines[cmd.Depth]++
		}
	}

	// special case for simple programs
	angle := 360.0 / float64(goroutines[1])
	if goroutines[1] < 3 {
		angle = 60.0
	}

	params := &Params{
		Angle:          angle,
		Caps:           totalG < 5, // value from head
		Distance:       80,
		AutoAngle:      false,
		DistanceSecond: 20,
	}

	angle2 := 360.0 / float64(goroutines[2]/goroutines[1])
	if goroutines[2] < goroutines[1] || angle2 == math.Inf(1) {
		angle2 = 60.0
	}
	params.AngleSecond = angle2

	return params
}
