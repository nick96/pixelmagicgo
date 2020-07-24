#!/usr/bin/env bash

# Compare pixelmatch-go and pixelmatch-js (https://github.com/mapbox/pixelmatch)
# using the images provided in the pixelmatch-js tets.
#
# This script is intended for use with the docker image built by
# ./docker/Dockerfile.compare so it makes a few assumptions about the locations
# of things:
#
# - Anything in the JS pixelmatch is expected to be in the pixelmatch-js
#   sub-directory.
# - Anything in the Go pixelmatch is expected to be in the pixelmatch-go
#   sub-directory.

js_output_tmp_fmt="/tmp/js.XXXXXX"
go_output_tmp_fmt="/tmp/go.XXXXXX"

function diffTest {
    local expected=$1
    local actual=$2
    local js_diff=$(mktemp $js_output_tmp_fmt)
    local go_diff=$(mktemp $go_output_tmp_fmt)

    ./pixelmatch-js/bin/pixelmatch $expected $actual $js_diff >/dev/null
    ./pixelmatch-go/pixelmatch -output $go_diff $expected $actual >/dev/null

    diff <(xxd $js_diff) <(xxd $go_diff) &>/dev/null
}

fixture_dir="./pixelmatch-js/test/fixtures"
failed=0

diffTest "$fixture_dir/1a.png" "$fixture_dir/1b.png"
if [ $? ]
then
    echo "Fixture 1 failed" >/dev/stderr
    failed=1
fi

diffTest "$fixture_dir/2a.png" "$fixture_dir/2b.png"
if [ $? ]
then
    echo "Fixture 2 failed" >/dev/stderr
    failed=1
fi

diffTest "$fixture_dir/3a.png" "$fixture_dir/3b.png"
if [ $? ]
then
    echo "Fixture 3 failed" >/dev/stderr
    failed=1
fi

diffTest "$fixture_dir/4a.png" "$fixture_dir/4b.png"
if [ $? ]
then
    echo "Fixture 4 failed" >/dev/stderr
    failed=1
fi

diffTest "$fixture_dir/5a.png" "$fixture_dir/5b.png"
if [ $? ]
then
    echo "Fixture 5 failed" >/dev/stderr
    failed=1
fi

diffTest "$fixture_dir/6a.png" "$fixture_dir/6b.png"
if [ $? ]
then
    echo "Fixture 6 failed" >/dev/stderr
    failed=1
fi

diffTest "$fixture_dir/7a.png" "$fixture_dir/7b.png"
if [ $? ]
then
    echo "Fixture 7 failed" >/dev/stderr
    failed=1
fi
