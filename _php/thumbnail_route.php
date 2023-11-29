<?php


function isImageFile($filePath) {
  $mime = mime_content_type($filePath);

  return strpos($mime, 'image/') === 0;
}

function extractSourceImage($filePath) {
  $mime = mime_content_type($filePath);
  if ($mime === 'image/jpeg') {
    return imagecreatefromjpeg($filePath);
  }
  if ($mime === 'image/bmp' || $mime === 'image/x-ms-bmp') {
    return imagecreatefrombmp($filePath);
  }
  if ($mime === 'image/png') {
    return imagecreatefrompng($filePath);
  }
  return imagecreatefrompng(__DIR__ . '/static/picture.png');
}


function isVideoFile($filePath) {
  $mime = mime_content_type($filePath);

  return strpos($mime, 'video/') === 0;
}

// Set the target dimensions
$targetWidth = 250;
$targetHeight = 250;

if (isset($_GET['size'])) {
  $targetWidth = intval($_GET['size']);
  $targetHeight = intval($_GET['size']);
}

$mediaroot = __DIR__ . '/_media';
$sourceImagePath = $mediaroot . $_GET['file'];

if (isImageFile($sourceImagePath)) {
  $sourceImage = extractSourceImage($sourceImagePath);

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
  $resizedImage = imagecreatetruecolor(floor($newWidth), floor($newHeight));

  // Resize and copy the source image to the new image
  imagecopyresampled($resizedImage, $sourceImage, 0, 0, 0, 0, floor($newWidth), floor($newHeight), $sourceWidth, $sourceHeight);

  // Set the header to indicate it's a JPEG image
  header('Content-Type: image/jpeg');

  // Output the resized image as JPEG
  imagejpeg($resizedImage);

  // Clean up resources
  imagedestroy($sourceImage);
  imagedestroy($resizedImage);
} else if (isVideoFile($sourceImagePath)) {
  $sourceImage = imagecreatefrompng(__DIR__ . '/static/play.png');

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
  // probably a directory
  $sourceImage = imagecreatefrompng(__DIR__ . '/static/folder.png');

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
}
?>
