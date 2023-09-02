<?php 
include("_includes/ContentViewer.php"); 
include("_includes/BreadCrumb.php");
include("_includes/Thumbnail.php");
$root = "";

if ($_SERVER['REQUEST_URI'] != '/') {
    $root = $_SERVER['REQUEST_URI'] ;
}

?>
<!DOCTYPE html>
<html>

  <head>
    <meta charset="UTF-8" />
    <link href="/styles.css" rel="stylesheet"/>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Simple Media Gallery</title>
  </head>
<body>
<?php
if ($root != "") {
    renderBreadcrumb($root);
}


$mediaroot = __DIR__ . '/_media';
$directory = $mediaroot . $root;

function omitDots($string) {
    return $string != "." && $string != "..";
}

if (is_file($directory)) {
    renderContentViewer($root, $directory);
} else if (is_dir($directory)) {
    $files = scandir($directory);
    $files = array_filter($files, "omitDots");
    echo "<div class='gallery'>";
    foreach ($files as $file) {
        renderThumbnail($mediaroot . $root, $file);
    }
    echo "</div>";
} else {
    http_response_code(404);
    echo "Nothing here friend";
    die();
}
?>

</body>
</html>