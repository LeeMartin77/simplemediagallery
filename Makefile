live:
	docker run -it --rm -p 9999:80 -v ./src:/var/www/html docker.io/php:8-apache

downloadvideos:
	./get_test_videos.sh