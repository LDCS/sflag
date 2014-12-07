#!/bin/bash

$GOROOT/bin/go run sflag1.go
$GOROOT/bin/go run sflag2.go
$GOROOT/bin/go run sflag2.go --Iq=60
$GOROOT/bin/go run sflag2.go --iq=1400
$GOROOT/bin/go run sflag2.go --iq=140 --Iq=60 
$GOROOT/bin/go run sflag3.go hello world
$GOROOT/bin/go run sflag4.go --foo
$GOROOT/bin/go run sflag4.go world hello
$GOROOT/bin/go run sflag5.go world hello
$GOROOT/bin/go run sflag6.go --Bar=2 world hello
$GOROOT/bin/go run sflag6.go world hello
$GOROOT/bin/go run sflag7.go --SomeCommand "foo | bar | baz" world hello
$GOROOT/bin/go run sflag8.go world hello
$GOROOT/bin/go run sflag9.go --Age 10 --Bar=7 --baz=2 --GDP=2 hello world


