FROM php:8-apache
RUN a2enmod rewrite
COPY src /var/www/html/