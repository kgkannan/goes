// autogenerated: do not edit!
// generated from gentemplate [gentemplate -d Package=elib -id Float64 -d VecType=Float64Vec -d Type=float64 vec.tmpl]

// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package elib

type Float64Vec []float64

func (p *Float64Vec) Resize(n uint) {
	c := uint(cap(*p))
	l := uint(len(*p)) + n
	if l > c {
		c = NextResizeCap(l)
		q := make([]float64, l, c)
		copy(q, *p)
		*p = q
	}
	*p = (*p)[:l]
}

func (p *Float64Vec) validate(new_len uint, zero float64) *float64 {
	c := uint(cap(*p))
	lʹ := uint(len(*p))
	l := new_len
	if l <= c {
		// Need to reslice to larger length?
		if l > lʹ {
			*p = (*p)[:l]
			for i := lʹ; i < l; i++ {
				(*p)[i] = zero
			}
		}
		return &(*p)[l-1]
	}
	return p.validateSlowPath(zero, c, l, lʹ)
}

func (p *Float64Vec) validateSlowPath(zero float64, c, l, lʹ uint) *float64 {
	if l > c {
		cNext := NextResizeCap(l)
		q := make([]float64, cNext, cNext)
		copy(q, *p)
		for i := c; i < cNext; i++ {
			q[i] = zero
		}
		*p = q[:l]
	}
	if l > lʹ {
		*p = (*p)[:l]
	}
	return &(*p)[l-1]
}

func (p *Float64Vec) Validate(i uint) *float64 {
	var zero float64
	return p.validate(i+1, zero)
}

func (p *Float64Vec) ValidateInit(i uint, zero float64) *float64 {
	return p.validate(i+1, zero)
}

func (p *Float64Vec) ValidateLen(l uint) (v *float64) {
	if l > 0 {
		var zero float64
		v = p.validate(l, zero)
	}
	return
}

func (p *Float64Vec) ValidateLenInit(l uint, zero float64) (v *float64) {
	if l > 0 {
		v = p.validate(l, zero)
	}
	return
}

func (p *Float64Vec) ResetLen() {
	if *p != nil {
		*p = (*p)[:0]
	}
}

func (p Float64Vec) Len() uint { return uint(len(p)) }