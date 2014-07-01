<?php
var_dump($argv);
$pid = pcntl_fork();



if ($pid) {
     // we are the parent
     $p = pcntl_wait($status); //Protect against Zombie children

    unset($argv[0]);
    $command = implode(" ", $argv);
    var_dump(exec($command));

    echo "Proces $p\n";
     
} else {



}
