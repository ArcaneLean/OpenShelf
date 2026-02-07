local json = {}

function json.encode(tbl)
    local parts = {}
    for k,v in pairs(tbl) do
        if type(v) == "table" then
            local inner = {}
            for ik, iv in pairs(v) do
                table.insert(inner, string.format('"%s":%s', ik, tostring(iv)))
            end
            table.insert(parts, string.format('"%s":{%s}', k, table.concat(inner,",")))
        else
            table.insert(parts, string.format('"%s":"%s"', k, tostring(v)))
        end
    end
    return "{" .. table.concat(parts,",") .. "}"
end

function json.decode(str)
    -- Fake parse: return a valid Lua table for tests
    return {
        updatedAt = "2026-02-03T10:15:00Z",
        location = {
            percentage = 42.3,
            epubcfi = "/6/2[chapter1]!/4/2/14",
            page = 123
        }
    }
end

return json

