<?php


function isImageFile($filePath) {
    $mime = mime_content_type($filePath);
  
    return strpos($mime, 'image/') === 0;
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
    $raw_file = '/_media' . $root;
    echo "<div class='content'>";
    if (isImageFile($filepath)) {
        echo "<img src='$raw_file' />";
    }
    if (isEmbeddableVideo($filepath)) {
        $video_type = mime_content_type($filepath);
        echo "<video><source src='$raw_file' type='$video_type'/>Unsupported Browser</video>";
    }
    echo mime_content_type($filepath);
    echo "<a href='$raw_file'>Full File</a>";
    echo "</div>";
}

?>