-- WRK: https://github.com/wg/wrk
-- ./wrk -s benchmark.lua -d 10s -t 2 http://localhost:8084/

math.randomseed(os.time())
wrk.method = "POST"

request = function()
    wrk.body   = string.format("https://kfd.me/%s", tostring(math.random(1,999999)))
    return wrk.format(nil, "/")
 end
