// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

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
