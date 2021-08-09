init:
	@cp ./scripts/* ./.git/hooks;
	@(cd ./web && npm install);
	@(cd ./backend && go mod download);

run:
	@(cd ./web && make run) & (cd ./backend && make run);
