debug:
	go test -c -gcflags "-N -l"
	gdb SimpleRTree.test -d $GOROOT
