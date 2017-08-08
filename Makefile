test:
	go test -v

bench:
	go test -run none -v  -bench=.

mem:
	go test -bench=. -memprofile=mem.out

cpu:
	go test -bench=. -cpuprofile=cpu.out