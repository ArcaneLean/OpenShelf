local ok, json = pcall(require, "json")
if not ok then
    print("KOReader JSON not available, using stub")
    json = require("plugin.json_stub")
end

local ok, lfs = pcall(require, "lfs")
if not ok then
    print("KOReader LFS not available, using stub")
    lfs = require("plugin.lfs_stub")
end

local utils = require("plugin.utils")

local state = {}

-- Config: Path to OpenShelf library root
-- TODO: make dynamic or configurable
local library_root = "/home/arcanelean/Documents/Programming/OpenShelf/examples/openshelf_test/"
local state_dir = library_root .. "/.state/"

-- Ensure the state directory exists
local function ensure_state_dir()
    local attr = lfs.attributes(state_dir)
    if not attr then
        utils.log("Creating .state directory: " .. state_dir)
        local ok, err = lfs.mkdir(state_dir)
        if not ok then
            utils.log("ERROR creating .state directory: " .. err)
        end
    end
end

-- Convert Book ID to filesystem-safe filename
local function book_id_to_filename(book_id)
    -- Replace colon with underscore per spec
    local safe_id = string.gsub(book_id, ":", "_")
    return state_dir .. safe_id .. ".json"
end

-- Read OpenShelf state
function state.read_state(book_id)
    ensure_state_dir()
    local file_path = book_id_to_filename(book_id)
    local file = io.open(file_path, "r")
    if not file then
        utils.log("No state file for book: " .. book_id)
        return nil
    end

    local content = file:read("*all")
    file:close()

    local ok, decoded = pcall(json.decode, content)
    if not ok or type(decoded) ~= "table" then
        utils.log("Invalid JSON in state file: " .. file_path)
        return nil
    end

    -- Optional: normalize location fields into array (per Section 10)
    if decoded.location and type(decoded.location) == "table" then
        local locs = {}
        for k,v in pairs(decoded.location) do
            table.insert(locs, { type = k, value = v })
        end
        decoded.locations = locs
    end

    return decoded
end

-- Write OpenShelf state atomically
function state.write_state(book_id, location_table)
    ensure_state_dir()
    local file_path = book_id_to_filename(book_id)
    local tmp_path = file_path .. ".tmp"

    -- Prepare state object
    local existing_state = state.read_state(book_id) or {}
    local updated_state = {
        specVersion = "0.1.0",
        bookId = book_id,
        updatedAt = os.date("!%Y-%m-%dT%H:%M:%SZ"),
        location = {}
    }

    -- Preserve unknown fields
    for k,v in pairs(existing_state) do
        if k ~= "location" and k ~= "specVersion" and k ~= "bookId" and k ~= "updatedAt" then
            updated_state[k] = v
        end
    end

    -- Merge new locations
    if location_table then
        for _, loc in ipairs(location_table) do
            if loc.type and loc.value then
                updated_state.location[loc.type] = loc.value
            end
        end
    end

    -- Write atomically
    local tmp_file = io.open(tmp_path, "w")
    if not tmp_file then
        utils.log("ERROR: cannot open temp file for writing: " .. tmp_path)
        return false
    end
    tmp_file:write(json.encode(updated_state))
    tmp_file:close()

    -- Replace original
    local ok, err = os.rename(tmp_path, file_path)
    if not ok then
        utils.log("ERROR renaming temp state file: " .. err)
        return false
    end

    utils.log("State updated for book: " .. book_id)
    return true
end

return state

