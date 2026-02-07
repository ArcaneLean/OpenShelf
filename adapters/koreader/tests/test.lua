local hash = require("plugin.hash")
local state = require("plugin.state")
local utils = require("plugin.utils")

local dummy_file = "/home/arcanelean/Documents/Programming/OpenShelf/examples/openshelf_test/books/sample.epub"
print("Book ID:", hash.compute_sha256(dummy_file))

state.write_state("sha256_dummyhash", {
    { type="percentage", value=42.3 },
    { type="epubcfi", value="/6/2[chapter1]!/4/2/14" }
})
    
local s = state.read_state("sha256_dummyhash")
if s and s.updatedAt then
    print("Updated:", s.updatedAt)
else
    print("No updatedAt found")
end
if s and s.location then
    for k,v in pairs(s.location) do
        print(k,v)
    end
else
    print("No state or location found")
end

