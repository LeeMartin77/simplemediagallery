<?php


function isImageFile($filePath) {
  $allowedMimeTypes = [
      'image/jpeg',
      'image/png',
      'image/gif',
      'image/bmp',
      'image/webp',
      // Add more image MIME types as needed
  ];

  $mime = mime_content_type($filePath);

  return in_array($mime, $allowedMimeTypes);
}

// Set the target dimensions
$targetWidth = 250;
$targetHeight = 250;

$mediaroot = __DIR__ . '/_media';
$sourceImagePath = $mediaroot . $_GET['file'];

if (isImageFile($sourceImagePath)) {

  $sourceImage = imagecreatefromjpeg($sourceImagePath);

  // Get the source image dimensions
  $sourceWidth = imagesx($sourceImage);
  $sourceHeight = imagesy($sourceImage);

  // Calculate the aspect ratio of the source image
  $aspectRatio = $sourceWidth / $sourceHeight;

  // Calculate the dimensions to fit within the target area while maintaining aspect ratio
  if ($aspectRatio > 1) {
      $newWidth = $targetWidth;
      $newHeight = $targetWidth / $aspectRatio;
  } else {
      $newHeight = $targetHeight;
      $newWidth = $targetHeight * $aspectRatio;
  }

  // Create a new image with the calculated dimensions
  $resizedImage = imagecreatetruecolor($newWidth, $newHeight);

  // Resize and copy the source image to the new image
  imagecopyresampled($resizedImage, $sourceImage, 0, 0, 0, 0, $newWidth, $newHeight, $sourceWidth, $sourceHeight);

  // Set the header to indicate it's a JPEG image
  header('Content-Type: image/jpeg');

  // Output the resized image as JPEG
  imagejpeg($resizedImage);

  // Clean up resources
  imagedestroy($sourceImage);
  imagedestroy($resizedImage);
} else {
  $image = imagecreate($targetWidth, $targetHeight);
  
  // Allocate the white color
  $whiteColor = imagecolorallocate($image, 255, 255, 255);
  
  // Fill the image with white color
  imagefill($image, 0, 0, $whiteColor);
  
  // Set the header to indicate it's a JPEG image
  header('Content-Type: image/jpeg');
  
  // Output the image as JPEG
  imagejpeg($image);
  
  // Clean up resources
  imagedestroy($image);
}
?>
