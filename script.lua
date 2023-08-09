---@diagnostic disable: undefined-global

local winSize = { 2000, 1500 }
local winTitle = "Rockets"

---------------------------
-- Genetic Rockets
---------------------------


-- Number of rockets per generation
local rocketStartPos = { 1000, 750 } -- bottom center of window
local rocketCount    = 30
local forceCount     = 255
local forceMagnitude = 1
local rocketSize     = 50
local rocketLifespan = 250
local epsilon        = 1e-10

local function createRocket()
   local result = {
      x = rocketStartPos[1],
      y = rocketStartPos[2],
      age = 0,

      forces = {},

      ResetForces = function(self)
         self.forces = {}
         for i = 0, forceCount, 1 do
            local forceDir = math.random() * math.pi * 2
            local forceStr = math.random() * forceMagnitude

            local forceX = math.cos(forceDir) * forceStr
            local forceY = math.sin(forceDir) * forceStr

            if math.abs(forceX) < epsilon then
               forceX = 0
            end
            if math.abs(forceY) < epsilon then
               forceY = 0
            end

            table.insert(self.forces, { forceX, forceY })
         end
      end,

      Draw = function(self)
         local so = rocketSize / 2.0
         local x = self.x
         local y = self.y

         Rectangle(x - so, y - so, x + so, y + so)
         Dot(x, y)

         ---- TODO ----
         -- Draw a sprite.
         -- Make it have a direction.
      end,

      Move = function(self)
         for i = 1, #self.forces, 1 do
            self.x = self.x + self.forces[i][1]
            self.y = self.y + self.forces[i][2]
         end
      end,
   }

   result:ResetForces()

   return result
end

function Setup()
   WinSize(2000, 1500, true)
   WinTitle("Buffas")
   Background(230)
end

function Keydown(k)
   if k.Escape then
      Quit()
   end
end

local r = createRocket()
function Draw()
   Sleep(10, true)
   r:Move()
   r:Draw()
   print(r.x, r.y)
end
