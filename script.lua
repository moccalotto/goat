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

 counter = 0

function __Setup()
    counter = math.random(1, 666666)
    print("Setup::")
end

-- D for draw, S for system
function Draw(D, S)
    -- Diller()
    S:Sleep(100)
    counter = counter + 1
    D:Scale(math.random(1,10) / 3)
    D:Line(
        D.W * 0.3, D.H * 0.3,
        D.W * 0.7, D.H * 0.7
    )

    -- print(counter)
end
