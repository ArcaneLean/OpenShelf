local json = require("json") -- KOReader has built-in JSON support
local utils = require("utils")

local state = {}

function state.read_state(book_id)
    utils.log("Reading state for: " .. book_id)
    -- TODO: open /.state/sha256_<book_id>.json
    return nil
end

function state.write_state(book_id, location)
    utils.log("Writing state for: " .. book_id)
    -- TODO: write /.state/sha256_<book_id>.json
end

return state

