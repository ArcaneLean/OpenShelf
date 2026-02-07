local logger = require("logger")

local utils = {}

function utils.log(msg)
    logger.info("[OpenShelf] " .. msg)
end

function utils.dump(tbl, indent, visited)
    indent = indent or 0
    visited = visited or {}

    if visited[tbl] then
        utils.log(string.rep("  ", indent) .. "*cycle*")
        return
    end

    visited[tbl] = true

    for k, v in pairs(tbl) do
        local prefix = string.rep("  ", indent) .. tostring(k) .. ": "
        if type(v) == "table" then
            utils.log(prefix .. "{")
            utils.dump(v, indent + 1, visited)
            utils.log(string.rep("  ", indent) .. "}")
        else
            utils.log(prefix .. tostring(v))
        end
    end
end

return utils

