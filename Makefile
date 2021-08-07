init:
	npm install;
	mkdir src;
	ln -s ./main.js ./src/index.js;

run:
	@(cd ./web && make run) & (cd ./backend && make run);
