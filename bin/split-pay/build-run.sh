set -exu
HERE=$(dirname $(realpath $BASH_SOURCE))
cd $HERE

go build -o split-pay.exe split-pay.go
./split-pay.exe