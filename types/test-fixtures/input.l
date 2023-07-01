{
	assert ~(0 + -1 <= 1 & 1 <= 2) # (3.0 * 4.0 = -12.0) & ("foo" < "foo" + "bar") = true;
	for true => false {
		assert 0 - 5 = -5;
	}

	if true {
		assert true;
	}
}
