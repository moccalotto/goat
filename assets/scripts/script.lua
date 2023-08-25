---@diagnostic disable: undefined-global

function Keydown(k)
   if k.Escape then
      Log("Goodbye 💙")
      Quit()
   end
end

local el1, el2

function Setup()
   Log("Hello From Setup()")
   -- WinSize(1000, 1000)
   -- FrameRateCap(30)
end

local i = 0.0
function Draw()
   i = i + 0.05
   local rfactor = math.min(255, 128.0 + math.cos(i) * 128)
   Background(254, rfactor, 128, 255)
end
