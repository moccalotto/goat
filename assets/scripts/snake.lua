---@diagnostic disable: undefined-global

local score          = 250
local winSizePx      = 1500
local tileCount      = 30 -- number of squares in each direction
local pxPerTile      = winSizePx / tileCount -- number of pixels per square

local snakeColor     = { 0,     0,   0, 255 }
local foodColor      = { 0,   150,   0, 255 }
local bgColor        = { 220, 220, 220, 220 }
local deadBgColor    = { 160,   0,   0, 255 }
local deadSnakeColor = { 20,   10,  10, 255 }
local deadFoodColor  = { 10,   60,  10, 255 }
local gridColor      = { 128, 128, 128, 255 }

local snake = {
   alive = true,
   headX = math.random(1, tileCount - 2),
   headY = math.random(1, tileCount - 2),
   direction = 0, -- 0 = north, 1 = east, 2 = south, 3 = west
   justAte = false,
   foodCountdown = 0,
   tail = {},
   food = {-1, -1},
   commandQueue = {}
}

snake.placeFood = function (self)
   for i=1000, 0, -1 do
      self.food[1] = math.random(0, tileCount - 1)
      self.food[2] = math.random(0, tileCount - 1)

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
   self.headX = (self.headX + dX) % tileCount
   self.headY = (self.headY + dY) % tileCount

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

   if point[1] < 0 or point[1] >= tileCount then
      return false
   end
   if point[2] < 0 or point[2] >= tileCount then
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
      Color(unpack(snakeColor))
   else
      Color(unpack(deadSnakeColor))
   end
   Dot(self.headX, self.headY)
   for _, v in ipairs(self.tail) do
      Dot(v[1], v[2])
   end
end

snake.drawFood = function(self)
   if self.alive  then
      Color(unpack(foodColor))
   else
      Color(unpack(deadFoodColor))
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
   self:placeFood()
   score = score + 1
   WinTitle(string.format("Snurk - %d points (%d, %d)", score, snake.food[1], snake.food[2]))
end

snake.grow = function(self)
   table.insert(self.tail, { self.headX, self.headY })
end

snake.dead = function(self)
   Background(unpack(deadBgColor))
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
   WinSize(winSizePx, winSizePx, true)
   WinTitle("Snurk")

   Background(unpack(bgColor))

   snake:placeFood()
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
   Sleep(math.max(150 + #snake.tail), true)

   if snake.foodCountdown <= 0 then
      snake.foodCountdown = 100 + math.random(0, #snake.tail)
      snake:placeFood()
   else
      snake.foodCountdown = snake.foodCountdown - 1
   end

   Scale(pxPerTile)

   Rectangle(0, 0, 1500, 1500)

   snake:eat()
   snake:move()
   snake:drawFood()
   snake:draw()

   
   ---------------------
   --    DRAW GRID    --
   ---------------------

   Scale(1)

   local maxX = winSizePx
   local maxY = winSizePx

   Color(unpack(gridColor))

   -- draw al the vertical lines
   for _x = pxPerTile, maxX, pxPerTile do
      Line(_x, 0, _x, maxY)
   end

   -- draw all the horizontal lines
   for _y = pxPerTile, maxY, pxPerTile do
      Line(0, _y, maxX, _y)
   end
end