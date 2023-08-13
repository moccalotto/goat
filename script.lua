---@diagnostic disable: undefined-global

function Keydown(k)
   if k.Escape then
      Log("Goodbye 💙")
      Quit()
   end
end

function Setup()
   Log("Hello From Setup()")
end

function Draw()
   Sleep(200)
end