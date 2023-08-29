<?php 
$root = "";
if ($_SERVER['REQUEST_URI'] != '/') {
    $root = $_SERVER['REQUEST_URI'] ;
}

$directory = __DIR__ . '/_media' . $root;

// breadcrumb

if (is_file($directory)) {
    // minimal gallery view
    echo "is file";
    $raw_file = '/_media' . $root;
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