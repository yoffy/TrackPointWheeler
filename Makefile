TrackPointWheeler.exe: *.go
	GOOS=windows go build

.PHONY: format
format:
	go fmt
