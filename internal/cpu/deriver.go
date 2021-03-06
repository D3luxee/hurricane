/*-
 * Copyright © 2017, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved.
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package cpu // import "github.com/solnx/hurricane/internal/cpu"

import (
	"github.com/mjolnir42/erebos"
	wall "github.com/solnx/eye/lib/eye.wall"
	"github.com/solnx/hurricane/internal/intf"
	"github.com/solnx/legacy"
)

// Implementation of the intf.Deriver interface

// NewDeriver ...
func NewDeriver(conf *erebos.Config) *Deriver {
	d := &Deriver{}
	d.data = make(map[int64]*CPU)
	d.lookup = wall.NewLookup(conf, `hurricane`)
	return d
}

// Deriver ...
type Deriver struct {
	data   map[int64]*CPU
	lookup *wall.Lookup
}

// Start ...
func (d *Deriver) Start() error {
	return d.lookup.Start()
}

// Close ...
func (d *Deriver) Close() {
	d.lookup.Close()
}

// Register ...
func (d *Deriver) Register(m map[string]intf.Deriver) {
	for _, s := range []string{
		`/sys/cpu/count/idle`,
		`/sys/cpu/count/iowait`,
		`/sys/cpu/count/irq`,
		`/sys/cpu/count/nice`,
		`/sys/cpu/count/softirq`,
		`/sys/cpu/count/system`,
		`/sys/cpu/count/user`,
	} {
		m[s] = d
	}
}

// Update ...
func (d *Deriver) Update(m *legacy.MetricSplit, t *erebos.Transport) ([]*legacy.MetricSplit, []*erebos.Transport, bool, error) {
	if _, ok := d.data[m.AssetID]; !ok {
		d.data[m.AssetID] = &CPU{
			lookup: d.lookup,
		}
	}

	return d.data[m.AssetID].update(m, t)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
