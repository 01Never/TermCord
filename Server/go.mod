module example/TermCord

go 1.26.1

require (
	example/TermCord/shared v0.0.0
	github.com/coder/websocket v1.8.14
	golang.org/x/time v0.15.0
)

replace example/TermCord/shared => ../shared
