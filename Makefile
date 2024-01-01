build: 


.PHONY: all

all: build

build:
	go build -o dist/smg src/main.go

run:
	SMG_MEDIA_DIRECTORY='example' go run src/main.go

test:
	go test ./src

downloadvideos:
	./get_test_videos.sh

docker:
	docker build . -t simplemediagallery:testimage

docker-run:
	docker run -it --rm -p 3333:3333 -v ./example:/_media simplemediagallery:testimage
