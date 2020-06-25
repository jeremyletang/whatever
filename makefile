build: go

go: leaderboard-summary

leaderboard-summary:
	mkdir -p functions
	cd leaderboard_summary && \
	go get ./... && \
	GO111MODULE=on go build -o "../functions/leaderboard-summary"
