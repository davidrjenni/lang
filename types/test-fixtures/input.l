{
	assert ~(0 + -1 <= 1 & 1 <= 2) # (3.0 * 4.0 = -12.0) & ("foo" < "foo" + "bar") = true;
	for true => false {
		assert 0 - 5 = -5;
		for true {
			continue;
		}
		break;
	}

	if true {
		assert true;
		for true {
			if true {
				break;
			}
		}
	}

	let a := 123;
	let b := 12.34;
	let c := "abc";
	let d := false;

	if a = 123 {
		let e := 456;
		if e = 123 {
			assert true;
		}
		if c # "def" {
			assert true;
		} else {
			let f := b - 12.04;
		}
	}

	set c <- "def";

	let f := func(i i64) func() bool {
		set i <- i * 2;
		return func() bool {
			return i = 42;
		};
	};
}
