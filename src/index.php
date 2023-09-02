<?php 
include("_includes/ContentViewer.php"); 
include("_includes/BreadCrumb.php");
$root = "";
if ($_SERVER['REQUEST_URI'] != '/') {
    $root = $_SERVER['REQUEST_URI'] ;
    renderBreadcrumb($root);
}

$directory = __DIR__ . '/_media' . $root;

// breadcrumb
if (is_file($directory)) {
    renderContentViewer($root, $directory);
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