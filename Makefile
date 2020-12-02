build:
		go build -o bin/dokidoki -v
		cd client && yarn && yarn build

run: build
		bin/dokidoki