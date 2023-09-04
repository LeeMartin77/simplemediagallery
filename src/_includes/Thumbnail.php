<?php 
function renderThumbnail($mediaroot, $file) {
  echo "<div class='thumbnail'>";
  $mediafile = $mediaroot . "/" . $file;
  $link = $_SERVER['REQUEST_URI'] . '/' . $file;
  if ($_SERVER['REQUEST_URI'] === '/') {
    $link = '/' . $file;
  }
  echo "<a href='$link'><img src='/thumbnail.php?file=$link' /></a>";
  if (is_file($mediafile)) {
    echo "<a href='$link'>$file</a>";
  } else {
    echo "<a href='$link'>$file</a>";
  }
  echo "</div>";
}
?>