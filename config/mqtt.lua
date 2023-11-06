 function
    local str = ...
	if string.find(str, "^rrpc,") or string.find(str, "^config,") then

    	return str
	end
  		local lat,lng= create.getRealLocation()
 		local dt = os.date("*t", os.time())
 		local timestr = string.format("%04d%02d%02d%02d%02d00", dt.year, dt.month, dt.day, dt.hour, dt.min)
        -- Create a new table in desired format
        local res = {}
        res.data = {}

          res.data.clientid = _G.model .. mobile.imei()
          res.data["@timestamp"] = os.date("!%Y-%m-%dT%H:%M:%S", os.time()) .. "+08:00"
          res.data.timestamp = os.time() -- current timestamp in milliseconds
          res.data.Node = mobile.eci()
        --res.data.id = "aaa" -- replace with your actual ID
        res.data.index = "testindex15"
        res.data.payload = {} -- Create payload table
        res.data.payload.data = {}
        res.data.payload.data.factor = {}

        res.data.payload.info = {}
        res.data.payload.info.model = _G.model
        res.data.payload.info.csq = mobile.csq()
        res.data.payload.info.iccid = mobile.iccid()
        res.data.payload.info.number = mobile.number(0)
        res.data.payload.info.imei = mobile.imei()
        res.data.payload.info.ip = ip
        res.data.payload.info.lat = lat
        res.data.payload.info.lng = lng

        res.data.payload.data.mn = _G.model .. mobile.imei()
        res.data.payload.data.datatime = timestr
        res.data.payload.data.datatype = "hj212"

--   local factorArray = {} -- 创建一个新的数组来存放因子
--
-- --创建第一个因子
-- local factor1 = {}
-- factor1.id = "a01001"
-- factor1.unit = "℃"
-- factor1.name = "温度"
-- factor1.value = tonumber(_G.t1)
--
-- --插入到数组
-- table.insert(factorArray, factor1)
--
-- --然后同样的方式处理其他因子...
-- local factor2 = {}
-- factor2.id = "a01002"
-- factor2.unit = "%"
-- factor2.name = "湿度"
-- factor2.value = tonumber(_G.h1)
--
-- table.insert(factorArray, factor2)
--
-- local factor3 = {}
-- factor3.id = "a05001"
-- factor3.unit = "ppm"
-- factor3.name = "二氧化碳"
-- factor3.value = tonumber(_G.c1)
--
-- table.insert(factorArray, factor3)
--   res.data.payload.data.factor = factorArray
       res.data.payload.data.factor.a01001 = {}
		res.data.payload.data.factor.a01001.id =  "a01001"
  		res.data.payload.data.factor.a01001.unit =  "℃"
  		res.data.payload.data.factor.a01001.name =  "温度"
		res.data.payload.data.factor.a01001.value =  tonumber(_G.t1)
       res.data.payload.data.factor.a01002 = {}
		res.data.payload.data.factor.a01002.id =  "a01002"
		res.data.payload.data.factor.a01002.unit =  "%"
  		res.data.payload.data.factor.a01002.name =  "湿度"
		res.data.payload.data.factor.a01002.value =  tonumber(_G.h1)
       res.data.payload.data.factor.a05001 = {}
		res.data.payload.data.factor.a05001.id =  "a05001"
		res.data.payload.data.factor.a05001.unit =  "ppm"
  		res.data.payload.data.factor.a05001.name =  "二氧化碳"
		res.data.payload.data.factor.a05001.value =  tonumber(_G.c1)



            -- Modify the payload values based on your actual logic
            --res.data.payload.value = 100
            return json.encode(res)


    --return str -- return original string if not processed
end
