kill -9 `ps aux | grep canal | grep -v grep | awk '{print $2}'`
