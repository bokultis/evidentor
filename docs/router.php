<?php

$path = pathinfo($_SERVER["SCRIPT_FILENAME"]);
if ($path["extension"] != "yaml") {
    return FALSE;
}

header("Access-Control-Allow-Origin: *");
header("Access-Control-Allow-Credentials: true");
header('Access-Control-Allow-Methods: GET, POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Accept');

header("Content-Type: text/yaml");
readfile($_SERVER["SCRIPT_FILENAME"]);