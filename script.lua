---@diagnostic disable: undefined-global
function Setup()
   Width(2000)
   Height(600)
   Title("Helmuth")
   CanResize()
end

-- D for draw, S for system
local bgRed = 20
local bgGreen = 20
local bgBlue = 20

local fgRed = 200
local fgGreen = 200
local fgBlue = 200

local scale = 1

function Draw()
   RandomizeColors()

   Background(bgRed, bgGreen, bgBlue, 255)

   Color(fgRed, fgGreen, fgBlue, 255)

   Scale(scale)

   Line(0, 0, Width, Height)
end

function RandomizeColors()
   if math.random(2, 30) > scale then
      scale = scale + 1
   else
      scale = scale - 1
   end
   if scale <= 1 then
      scale = 2
   end

   if math.random(1, 255) > bgRed then
      bgRed = bgRed + 1
   else
      bgRed = bgRed - 1
   end
   if math.random(1, 255) > bgGreen then
      bgGreen = bgGreen + 1
   else
      bgGreen = bgGreen - 1
   end
   if math.random(1, 255) > bgBlue then
      bgBlue = bgBlue + 1
   else
      bgBlue = bgBlue - 1
   end
   if math.random(1, 255) > fgRed then
      fgRed = fgRed + 1
   else
      fgRed = fgRed - 1
   end
   if math.random(1, 255) > fgGreen then
      fgGreen = fgGreen + 1
   else
      fgGreen = fgGreen - 1
   end
   if math.random(1, 255) > fgBlue then
      fgBlue = fgBlue + 1
   else
      fgBlue = fgBlue - 1
   end
end