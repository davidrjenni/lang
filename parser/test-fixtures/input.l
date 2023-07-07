{
	assert 42 * 2 <= 84 <= 85 * (1 + 1);
	assert 5.0*5.000_000 = 25.00;
	assert --1 = 1;
	assert ~~true;

	// This is a line comment.
	for true {}
	for 2*3 < 5 {
		assert -1 = 1;
		break;
		continue; // This is a line comment.
	}

	if true {
		assert true;
	}

	if false {
		assert false;
	} else {
		assert true;
	}

	if false {
		assert false;
	} else if false {
		assert true;
	} else if false {
		assert true;
	} else {
		assert true;
	}

	let x := 123;
	if x < 123 {
		assert true;
	}
	set x <- 456;

	let f := func() i64 { };
	let g := func(a i64, b bool) i64 { };
	let h := func(a i64, b bool) func(bool) func(func(string) bool) f64 { };
	return 42;
}
