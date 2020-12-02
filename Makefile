build:
		go build -o bin/dokidoki -v
		cd client && yarn build

run: build
		bin/dokidoki