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

   local sheet = SpriteSheet("assets/Spritesheets/sheet.xml") -- load if necessesary, otherwise retrieve from cache

   BlueLaserSprite = sheet.GetSprite("blueLaser")  -- image is already in mem, just load subsprite coords

   BlueLaserSprite.Add(
      Accelerator()     -- this object accelerates the sprite in a given direction. But it needs to be told its acceleration
   )

   TGroup("projectiles")
      :Add(BlueLaserSprite)
      :Draw()
   
end
