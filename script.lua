function Dump(o)
    if type(o) == 'table' then
       local s = '{ '
       for k,v in pairs(o) do
          if type(k) ~= 'number' then k = '"'..k..'"' end
          s = s .. '['..k..'] = ' .. Dump(v) .. ','
       end
       return s .. '} '
    else
       return tostring(o)
    end
 end

function Setup()
    print("Setupa")
end

function Draw(dm)
    print("Draw::", dm.stuff)
    dm:Line(
        dm.W * 0.3, dm.H * 0.3,
        dm.W * 0.7, dm.H * 0.7
    )
end

function Update(param)
    print("Update:: ", param)
end
