#!/bin/bash -eu

cd `dirname $0`

if [ $BIRDSEYEAPI_EXECUTION_MODE = 'PRODUCTION' ] ; then
    echo 'running PRODUCTION mode...'
    ./bin/server
else 
    echo 'running DEBUG mode...'
    # bashを対話モードで動かしてdockerが終了しないようにする
    /bin/bash
fi
