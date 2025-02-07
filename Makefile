GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
env =

build:
	mkdir -p target/runtime_log/
	cp cmd/conf/$(env)/logic.toml target/logic.toml
	$(GOBUILD) -o target/go_project_template_logic cmd/logic/main.go

run:
	nohup target/go_project_template_logic -conf=target/logic.toml >> target/runtime_log/logic.log 2>&1 &

stop:
	pkill -f target/go_project_template_logic