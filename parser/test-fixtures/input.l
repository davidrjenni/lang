{
	assert 42 * 2 <= 84 <= 85 * (1 + 1);
	assert 5.0*5.000_000 = 25.00;
	assert --1 = 1;
	assert ~~true;

	for true {}
	for 2*3 < 5 {
		assert -1 = 1;
		break;
		continue;
	}

	if true {
		assert true;
	}
}
