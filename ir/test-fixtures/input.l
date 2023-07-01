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
}
