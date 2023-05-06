module main

go 1.20

require geecache  v0.0.0
require lru       v0.0.0

replace geecache => ./geecache
replace lru      => ./geecache/lru