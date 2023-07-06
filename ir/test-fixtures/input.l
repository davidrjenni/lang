{
	assert ~(~(true));
	assert ~(~(false));
	assert 27 = 3 + 5*5 - 1;
	assert false = true | true;
	assert -1 = 0 - 1;

	for true {
		assert -1 = 0 - 1;
	}

	if true {
		assert true;
	}

	for true {
		if true {
			break;
		}
	}

	for true {
		if true {
			for true {
				break;
				assert true;
			}
			continue;
		}
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

	let x := 2 * 3;
	let y := x * 3;
	assert x = 6;

	let z := true & y = 18;
	assert z;

	set x <- y * 2;
	assert x = 36;
}
