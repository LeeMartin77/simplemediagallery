#!/bin/bash

for x in *; do 
  filename=$(echo $x | sed 's/\.tiff$//');
  mkdir $filename
  mv $x "$filename"
  magick "$filename/$x" "$filename/$filename.jpeg"
  magick "$filename/$x" "$filename/$filename.png"
  magick "$filename/$x" "$filename/$filename.bmp"
done
