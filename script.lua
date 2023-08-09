---@diagnostic disable: undefined-global

local squareCount = 0 -- number of squares in each direction
local pixelsPerSquare = 0 -- number of pixels per square

local snakeColor     = { 0,    0,   0, 255 }
local foodColor      = { 0,  150,   0, 255 }
local deadBgColor    = { 160,  0,   0, 255 }
local deadSnakeColor = { 20,  10,  10, 255 }
local deadFoodColor  = { 10,  60,  10, 255 }

local snake = {
   alive = true,
   headX = 0,
   headY = 0,
   direction = 0, -- 0 = north, 1 = east, 2 = south, 3 = west
   justAte = false,
   tail = {},
   food = {-1, -1},
   commandQueue = {}
}

snake.shuffleFood = function (self)
   for i=1000, 0, -1 do
      self.food[1] = math.random(0, squareCount)
      self.food[2] = math.random(0, squareCount)

      if self:canPlaceFood(self.food) then
         return true
      end
   end
   error("Could not place food after 1000 tries")
end

snake.getVelocity = function(self)
   -- we do this to make a visually appealing snake death
   if false == self.alive then
      return 0, 0
   end

   if 0 == self.direction then
      return 0, -1
   elseif 1 == self.direction then
      return 1, 0
   elseif 2 == self.direction then
      return 0, 1
   elseif 3 == self.direction then
      return -1, 0
   else
      error("Dir should be between 0 and 3, inclusive")
   end
end

snake.move = function(self)

   -- I've turned this off because i want to see the snake
   -- eat itself slowly.
   -- if not self.alive then
   --    return
   -- end

   -- if we have unprocessed keyboard input commands,
   -- process the oldest one and remove it from the queue.
   if #self.commandQueue > 0 then
      local change = table.remove(self.commandQueue, 0)
      self.direction = (self.direction + change) % 4
   end


   -- first we move the entire tail without moving the head
   -- then we move the head to its new location
   -- we simply do this by only moving the tip of the tail
   -- to where the head is
   table.insert(self.tail, 1, { self.headX, self.headY, })
   if self.justAte then
      -- we just ate, so we dont delete our tail tip
      self.justAte = false
   else
      table.remove(self.tail)
   end

   local dX, dY = self:getVelocity()
   self.headX = (self.headX + dX) % squareCount
   self.headY = (self.headY + dY) % squareCount

   for _, v in ipairs(self.tail) do
      if v[1] == self.headX and v[2] == self.headY then
         self:dead()
      end
   end
end

snake.canPlaceFood = function(self, point)
   if point[1] == self.headX and point[2] == self.headY then
      return false
   end

   local maxX, maxY = CanvasSize()

   if point[1] < 0 or point[1] >= maxX then
      return false
   end
   if point[2] < 0 or point[2] >= maxY then
      return false
   end

   for _, v in ipairs(self.tail) do
      if point[1] == v[1] and point[2] == v[2] then
         return false
      end
   end

   return true
end

snake.draw = function(self)
   if self.alive then
      Color( snakeColor[1], snakeColor[2], snakeColor[3], snakeColor[4])
   else
      Color( deadSnakeColor[1], deadSnakeColor[2], deadSnakeColor[3], deadSnakeColor[4])
   end
   Dot(self.headX, self.headY)
   for _, v in ipairs(self.tail) do
      Dot(v[1], v[2])
   end
end

snake.drawFood = function(self)
   if self.alive  then
      Color(foodColor[1], foodColor[2], foodColor[3], foodColor[4])
   else
      Color(deadFoodColor[1], deadFoodColor[2], deadFoodColor[3], deadFoodColor[4])
   end
   Dot(self.food[1], self.food[2])
end

snake.canEat = function(self)
   if self.headX == self.food[1] and self.headY == self.food[2] then
      return true
   end

   return false
end

snake.eat = function(self)
   if not self:canEat() then
      return
   end
   self:grow()
   self:shuffleFood()
end

snake.grow = function(self)
   table.insert(self.tail, { self.headX, self.headY })
end

snake.dead = function(self)
   Background(deadBgColor[1], deadBgColor[2], deadBgColor[3], deadBgColor[4])
   self.alive = false
end

snake.addToCommandQueue = function(self, num)
   if #self.commandQueue > 2 then
      return
   end

   table.insert(self.commandQueue, num)
end

function Winfo(w) 
   return w
end

function Setup()
   local winSizePx = 1500
   SetWinSize(winSizePx, winSizePx, true)
   SetWinTitle("Snurk")

   squareCount = 30 -- we want 30 squares in the x direction and 30 squares in the y direction
   pixelsPerSquare = winSizePx / squareCount

   Background(220)

   snake:shuffleFood()
end

function Keydown(k)
   if k.Str == "Escape" then
      Quit()
   elseif k.Str == "Space" then
      snake:grow()
   elseif k.Left then
      snake:addToCommandQueue(-1) -- turn left
   elseif k.Right then
      snake:addToCommandQueue(1) -- turn right
   elseif k.Str == "D" then
      snake.alive = false
   end
end

function Draw()

   Background(220)
   Color(20)
   Dot(10, 10)

--[[
   if Counter() > 1 then
      Sleep(250)
   end

   snake:eat()
   snake:move()
   snake:drawFood()
   snake:draw()
   
   ---------------------
   --    DRAW GRID    --
   ---------------------

   Scale(1)

   local maxX, maxY = CanvasSize()

   -- Color(128, 0, 0, 255)

   -- draw al the vertical lines
   for _x = pixelsPerSquare, maxX, pixelsPerSquare do
      Line(_x, 0, _x, maxY)
   end

   -- draw all the horizontal lines
   for _y = pixelsPerSquare, maxY, pixelsPerSquare do
      Line(0, _y, maxX, _y)
   end
   --]]
end
