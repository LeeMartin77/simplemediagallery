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

function isEmbeddableVideo($filePath) {
  $allowedMimeTypes = [
      'video/mp4',
      'video/webm',
      'video/ogg',
      'video/quicktime',
  ];

  $mime = mime_content_type($filePath);

  return in_array($mime, $allowedMimeTypes);
}

function renderContentViewer($root, $filepath) {
    // minimal gallery view
    echo "is file";
    $raw_file = '/_media' . $root;
    if (isImageFile($filepath)) {
        echo "<img src='$raw_file' />";
    }
    echo mime_content_type($filepath);
    if (isEmbeddableVideo($filepath)) {
        $video_type = mime_content_type($filepath);
        echo "<video><source src='$raw_file' type='$video_type'/>Unsupported Browser</video>";
    }
    echo "<a href='$raw_file'>Full File</a>";
}

?>