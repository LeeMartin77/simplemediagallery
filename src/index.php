<?php 
$root = "";
if ($_SERVER['REQUEST_URI'] != '/') {
    $root = $_SERVER['REQUEST_URI'] ;
}

$directory = __DIR__ . '/_media' . $root;

// breadcrumb

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

if (is_file($directory)) {
    // minimal gallery view
    echo "is file";
    $raw_file = '/_media' . $root;
    if (isImageFile($directory)) {
        echo "<img src='$raw_file' />";
    }
    echo mime_content_type($directory);
    if (isEmbeddableVideo($directory)) {
        $video_type = mime_content_type($directory);
        echo "<video><source src='$raw_file' type='$video_type'/>Unsupported Browser</video>";
    }
    echo "<a href='$raw_file'>Full File</a>";
    
} else if (is_dir($directory)) {
    $files = scandir($directory);

    foreach ($files as $file) {
        // this becomes thumbnails
        ?><ul><?php
        if ($file != '.' && $file != '..') {
            echo "<li><a href='$root/$file'>$file</a></li>";
        }
        ?></ul><?php
    }
} else {
    http_response_code(404);
    echo "Nothing here friend";
    die();
}
?>