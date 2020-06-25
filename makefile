build: go

go: prepare leaderboard-summary name-from-key

prepare:
	mkdir -p functions

leaderboard-summary:
	cp config.go leaderboard_summary/config.go && \
	cd leaderboard_summary && \
	go get ./... && \
	GO111MODULE=on go build -o "../functions/leaderboard-summary"

name-from-key:
	cp config.go name_from_key/config.go && \
	cd name_from_key && \
	go get ./... && \
	GO111MODULE=on go build -o "../functions/name-from-key"
