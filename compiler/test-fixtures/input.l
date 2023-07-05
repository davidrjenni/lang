{
	assert ~(~(true));
	assert ~(~(false));
	assert 27 = 3 + 5*5 - 1;
	assert false = true | true;
	assert -1 = 0 - 1 => true;
	let x := 2 * 3;
	let y := x * 3;
	assert x = 6;
	let z := true & x = 6;
	assert z;
}
