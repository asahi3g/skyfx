##

all: release

debug:
	C:\\msys64\\mingw64\\bin/go build -gcflags=all="-N -l" -o ~/go/bin/skycad cmd/skycad/*.go

release:
	C:\\msys64\\mingw64\\bin/go build -o ~/go/bin/skycad cmd/skycad/*.go

cleanbin:
	rm -rf ~/go/bin/skycad