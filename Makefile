build: 
	docker build . -t simplemediagallery:testimage

live:
	docker run -it --rm -p 9999:80 -v ./src:/var/www/html -v ./example:/var/www/html/_media simplemediagallery:testimage

downloadvideos:
	./get_test_videos.sh