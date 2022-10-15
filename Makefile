help:	## Show this help.
	    @sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

init:	## Install dependencies.
	@cp ./scripts/* ./.git/hooks;
	@(cd ./web && npm install);
	@(cd ./backend && go mod download);

run:	## Run the web and backend servers.
	@(cd ./web && make run) & (cd ./backend && make run);

test:	## Run the tests.
	@(cd ./web && npm test);
