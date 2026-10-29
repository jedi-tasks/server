deploy-check-all: go-mod-tidy-all
	cd pkg; make deploy-check
	cd server; make deploy-check

go-mod-tidy-all:
	cd pkg; go mod tidy
	cd server; go mod tidy

submodule-update:
	git submodule update --init --recursive --remote
