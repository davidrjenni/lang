// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

type Pass func(Seq) Seq

var Loads = Pass(loads)

func flatten(seq Seq) (tseq Seq) {
	for _, n := range seq {
		switch n := n.(type) {
		case *BinaryExpr:
			if seqx, ok := n.LHS.(*seqExpr); ok {
				s := flatten(seqx.Seq)
				tseq = append(tseq, s...)
				tseq = append(tseq, &BinaryExpr{
					RHS: n.RHS,
					Op:  n.Op,
					LHS: seqx.Dst,
					pos: n.Pos(),
				})
				continue
			}
		case *Load:
			if seqx, ok := n.Src.(*seqExpr); ok {
				s := flatten(seqx.Seq)
				tseq = append(tseq, s...)
				tseq = append(tseq, &Load{
					Src: seqx.Dst,
					Dst: n.Dst,
					pos: n.Pos(),
				})
				continue
			}
		case Seq:
			s := flatten(n)
			tseq = append(tseq, s...)
			continue
		}
		tseq = append(tseq, n)
	}
	return tseq
}

func loads(seq Seq) (tseq Seq) {
	for _, n := range seq {
		if n, ok := n.(*Load); ok {
			src, ok := n.Src.(*Reg)
			if ok && src.Type == n.Dst.Type && src.Second == n.Dst.Second {
				continue
			}
		}
		tseq = append(tseq, n)
	}
	return tseq
}
