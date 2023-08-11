local function dumpTable(table, depth)
   if type(depth) ~= "number" then
      depth = 0
   end
   if depth > 200 then
      print("Error: Depth > 200 in dumpTable()")
      return
   end
   for k, v in pairs(table) do
      if (type(v) == "table") then
         print(string.rep("  ", depth) .. k .. ":")
         dumpTable(v, depth + 1)
      else
         print(string.rep("  ", depth) .. k .. ": ", v)
      end
   end
end

return {
   ["dumpTable"] = dumpTable,
}