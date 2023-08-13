---@diagnostic disable: undefined-global

---------------------------
-- Configuration constants
---------------------------
local winSize                = Vector(2000, 1500)
local winTitle               = "Rockets"

local rocketStartPos         = { 1000, 1200 } -- Starting position of rockets
local rocketCount            = 45             -- Number of rockets per generation
local rocketLifespan         = 200            -- rocketLifespan must not be higher than forceCount
local rocketSize             = 90             -- the size (in pixels) of the rocket

local forceCount             = rocketLifespan -- number of forces that can affect a rocket
local forceMagnitude         = 2.5            -- how strong the forces are
local chanceOfUsingParentDNA = 0.95           -- This is the chance that a gene is taken from a paren (otherwise its random generated)
local mutationMagnitude      = 0.05           -- when a gene is taken from a parent, its still mutated a bit. This is the magniude
local parentCount            = 4              -- number of parents per generation. Rockets are awesome. They can have many parents.

local obstacleCenter         = winSize:Scale(0.5):Sub(Vector(0, 200))
local obstacleSize           = Vector(winSize:X() / 3, 50)

local goalSize               = 200
local goalPos                = Vector(winSize:X() / 2, goalSize / 2)


---------------------------
-- Simulation variables
-- don't touch
---------------------------
local parents = {} -- The rockets used to generate the future generations. If empty, we spawn brand new ones.
local rockets = {}
local doneCount = 0
local obstacleRect = {
   obstacleCenter:X() - obstacleSize:X() / 2,
   obstacleCenter:Y() - obstacleSize:Y() / 2,
   obstacleCenter:X() + obstacleSize:X() / 2,
   obstacleCenter:Y() + obstacleSize:Y() / 2,
}


---------------------------
-- Library imports.
-- don't touch
---------------------------
local dumpTable = require("dumpTable").dumpTable


local function createRocket()
   local result = {
      pos = Vector(unpack(rocketStartPos)),
      age = 0,
      done = false,
      score = -1,
      velocity = Vector(0, 0),

      forces = {},

      HasParents = function(self)
         return #parents > 0
      end,

      ShouldUseParentDNA = function(self)
         if self:HasParents() == false then
            return false -- we didn't use parent dna cuz there were no parents
         end

         local die_roll = math.random()

         -- if there are parents chance that we use the dna
         -- from one of them
         return die_roll < chanceOfUsingParentDNA
      end,

      GenerateRandomForce = function(self)
         return PolarVector(
            math.random() * math.pi * 2,
            math.random() * forceMagnitude
         )
      end,

      InitializeForces = function(self)
         self.forces = {}

         -- Calculate how chance of getting
         -- * a daddy force
         -- * a mommy force
         -- * a random force
         --
         -- Example: mommy: 65, daddy: 30, random: 5

         for i = 1, forceCount, 1 do
            if self:ShouldUseParentDNA() then
               -- pick a random parent and use their applicable gene
               local parentIdx = math.random(1, #parents)

               -- apply some mutation
               local force = parents[parentIdx].forces[i]:Clone()
               local mutationVector = Vector(
                  math.random() - 0.5,
                  math.random() - 0.5
               )

               mutationVector:Scale(mutationMagnitude * forceMagnitude)

               force = force:Add(mutationVector)

               table.insert(self.forces, force)
            else
               -- Generate random force
               local force = self:GenerateRandomForce()
               table.insert(self.forces, force)
            end
         end
      end,

      CollidesWithObstacle = function(self)
         local X = self.pos:X()
         local Y = self.pos:Y()

         if X < obstacleRect[1] then
            return false
         end
         if X > obstacleRect[3] then
            return false
         end
         if Y < obstacleRect[2] then
            return false
         end
         if Y > obstacleRect[4] then
            return false
         end

         return true
      end,

      CollidesWithGoal = function(self)
         local dist = self.pos:Sub(goalPos):Len()

         if dist < (goalSize + rocketSize) / 2 then
            return true
         end

         return false
      end,

      IsParent = function(self)
         for i = 1, #parents, 1 do
            if self == parents[i] then
               return true
            end
         end

         return false
      end,

      Draw = function(self)
         Push()
         if self:IsParent() then
            Color(30, 30, 255, 230)
         else
            Color(10, 10, 30, 150)
         end

         local angle = self.velocity:Angle() - math.pi / 6

         Polygon(
            self.pos:X(),     -- center x position
            self.pos:Y(),     -- center y position
            rocketSize / 2.0, -- "radius" of polygon
            angle,            -- rotation
            3                 -- number of sides
         )
         Pop()

         ---- TODO ----
         -- Draw a sprite.
         -- Make it have a direction.
      end,

      SetScore = function(self, scoreFactor)
         if self.done then
            return
         end

         if scoreFactor == nil then
            scoreFactor = 1
         end

         local x = self.pos:X()
         local y = self.pos:Y()

         local worstPossibleScore = math.sqrt(
            math.pow(winSize:X(), 2) + math.pow(winSize:Y(), 2)
         )

         local distanceToGoal = self.pos:Sub(goalPos):Len()

         self.score = scoreFactor * (worstPossibleScore - distanceToGoal) / 1000

         self.done = true
         doneCount = doneCount + 1
      end,

      Fly = function(self)
         if self.done then
            return
         end

         if self.age >= rocketLifespan then
            self:SetScore(1)
            return
         end

         self.age = self.age + 1 -- moving makes you old.

         if self:CollidesWithGoal() then
            -- if you actually hit the goal, you get some extra sugar.
            self:SetScore(1.3)
            return
         end

         if self:CollidesWithObstacle() then
            self:SetScore(-200)
            return
         end

         local force = self.forces[self.age]

         self.velocity = self.velocity:Add(force)
         self.pos = self.pos:Add(self.velocity)
      end,
   }

   result:InitializeForces()

   return result
end
local function startNewRocketGeneration()
   rockets = {}
   doneCount = 0

   for i = 1, rocketCount, 1 do
      table.insert(rockets, createRocket())
   end
end
function Setup()
   WinSize(winSize:X(), winSize:Y(), true)
   WinTitle(winTitle)
   Background(230)
   FrameRateCap(15)
end

function Keydown(k)
   if k.Escape then
      Quit()
   end

   if k.Name == "Space" then
      startNewRocketGeneration()
      FrameRateCap(-1)
   end
end

function Draw()
   if doneCount == rocketCount then
      -- return
   end
   ------------------------------
   -- CREATE GOAL
   --
   -- TODO move to function
   ------------------------------
   Push()
   Color(30, 200, 30, 255)

   Polygon(
      goalPos:X(),  -- pos X
      goalPos:Y(),  -- pos Y
      goalSize / 2, -- size
      0,            -- angle
      6             -- sides
   )
   Pop()

   Push()
   Color(150, 75, 0)
   Rectangle(unpack(obstacleRect))
   Pop()


   ------------------------------
   -- FLY AND DRAW THE ROCKETS
   --
   ------------------------------
   for i, r in ipairs(rockets) do
      r:Fly()
      r:Draw()
   end

   if doneCount < rocketCount then
      return
   end

   ------------------------------
   -- CALCULATE THE TOTAL SCORE
   --
   ------------------------------

   table.sort(rockets, function(a, b)
      return a.score > b.score
   end)


   for i = 1, math.min(parentCount, #rockets), 1 do
      parents[i] = rockets[i]
      parents[i]:Draw()
   end
   FrameRateCap(15)
end
