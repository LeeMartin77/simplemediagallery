<?php
function renderBreadcrumb($path) {
  $segments = explode("/", $path);
  $cumulativePath = "";
  ?><div class='breadcrumb'>
    <a href='/'>home</a><?php
    foreach ($segments as $segment) {
      if ($segment != "") {
        $cumulativePath = $cumulativePath . "/" . $segment;
        echo " / <a href='$cumulativePath'>$segment</a>";
      }
    }
  ?></div><?php
}
?>