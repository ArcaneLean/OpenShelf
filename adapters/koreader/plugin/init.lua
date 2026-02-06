local adapter = require("adapter")
local state = require("state")
local hash = require("hash")
local utils = require("utils")

-- Called when KOReader plugin loads
local function on_plugin_load()
    utils.log("OpenShelf KOReader adapter loaded")
    utils.log("Supported levels: " .. table.concat(adapter.levels, ", "))
end

-- Called when a book is opened
local function on_book_open(book_path)
    utils.log("Opening book: " .. book_path)
    local book_id = hash.compute_sha256(book_path)
    local reading_state = state.read_state(book_id)
    if reading_state then
        -- TODO: apply location to KOReader
        utils.log("Found OpenShelf state for book: " .. book_id)
    else
        utils.log("No OpenShelf state found for book: " .. book_id)
    end
end

-- Called when a book is closed
local function on_book_close(book_path)
    utils.log("Closing book: " .. book_path)
    local book_id = hash.compute_sha256(book_path)
    local ko_state = {}  -- TODO: extract KOReader location
    state.write_state(book_id, ko_state)
end

-- Register plugin hooks (KOReader API)
koreader.register_plugin({
    name = adapter.name,
    on_load = on_plugin_load,
    on_book_open = on_book_open,
    on_book_close = on_book_close,
})

