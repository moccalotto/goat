---@diagnostic disable: undefined-global

local winSize = { 2000, 1500 }
local winTitle = "Rockets"

---------------------------
-- Genetic Rockets
---------------------------


-- Number of rockets per generation
local rocketStartPos = { 1000, 1500 } -- bottom center of window
local rocketCount    = 30
local forceCount     = 255
local forceMagnitude = 1
local rocketSize     = 50
local rocketLifespan = 200
local epsilon        = 1e-10
local goal           = { 900, 0, 1100, 150 }
local centerOfGoal   = { (goal[3] + goal[1]) / 2, (goal[2] + goal[4]) / 2, }



local parents = {}   -- The rockets used to generate the future generations. If empty, we spawn brand new ones.
local rockets = {}


local function createRocket()
   local result = {
      x = rocketStartPos[1],
      y = rocketStartPos[2],
      age = 0,
      done = false,
      score = -1,

      forces = {},

      ResetForces = function(self)
         -- TODO
         -- Parents!!!

         self.forces = {}

         -- Calculate how chance of getting
         -- * a daddy force
         -- * a mommy force
         -- * a random force
         -- 
         -- Example: mommy: 65, daddy: 30, random: 5

         for i = 0, forceCount, 1 do
            -- roll a die
            -- maybe pick from mom or dad
            -- or maybe pick random.
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

      CollidesWithGoal = function(self)
         local minX, minY, maxX, maxY = unpack(goal)
         local x = self.x
         local y = self.y
         return x >= minX and x <= maxX and y >= minY and y <= maxY
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

      CalculateScore = function(self)
         local dist = math.sqrt(
            math.pow(centerOfGoal[1] - self.x, 2)
            +
            math.pow(centerOfGoal[2] - self.y, 2)
         )

         -- in reality we calculate the distance to the center of the goal
         return winSize[2] - dist - self.age
      end,

      Live = function(self)
         if self.score >= 0 then
            -- this rocket is done
            return
         end

         if self.age >= rocketLifespan then
            self.score = self:CalculateScore()
            return
         end

         if self:CollidesWithGoal() then
            -- if you actually hit the goal, you get some extra sugar.
            self.score = 1.3 * self:CalculateScore()
            return
         end

         self.age = self.age + 1 -- moving makes you old.

         local remainingLife = rocketLifespan - self.age
         local penaltyFactor = 1 -- remainingLife / rocketLifespan
         for i = 1, #self.forces, 1 do
            self.x = self.x + self.forces[i][1] * penaltyFactor
            self.y = self.y + self.forces[i][2] * penaltyFactor
         end

         -- ALT MOVE MODE. Fire one rocket at at time!!!!
         -- Not sure if this would work as a genetic level though.
         -- local idx = self.age % #self.forces
         -- local penaltyFactor = 100
         -- self.x = self.x + self.forces[idx][1] * penaltyFactor
         -- self.y = self.y + self.forces[idx][2] * penaltyFactor
      end,
   }

   result:ResetForces()

   return result
end
local function startNewRocketGeneration()
   rockets = {}

   for i = 0, rocketCount, 1 do
      table.insert(rockets, createRocket())
   end
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

   if k.Name == "Space" then
      startNewRocketGeneration()
   end
end

function Draw()
   Sleep(10, true)

   ------------------------------
   -- CREATE GOAL
   --
   -- TODO move to function
   ------------------------------
   Push()
   Color(30, 200, 30, 255)

   Rectangle(unpack(goal))
   Pop()
   ------------------------------

   for i, r in ipairs(rockets) do
      if r.done then
         -- check the rocket's score.
         -- find the two best rockets
         -- and pair them up
         -- add a bit of mutation
      end

      r:Live()
      r:Draw()
   end
end
