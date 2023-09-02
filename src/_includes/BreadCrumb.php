<?php
function renderBreadcrumb($path) {
  $segments = explode("/", $path);
  $cumulativePath = "";
  ?><div><ul>
    <li><a href='/'>Home</a></li><?php
    foreach ($segments as $segment) {
      if ($segment != "") {
        $cumulativePath = $cumulativePath . "/" . $segment;
        echo "<li><a href='$cumulativePath'>$segment</a></li>";
      }
    }
  ?></ul></div><?php
}
?>