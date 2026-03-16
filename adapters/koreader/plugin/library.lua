local json = require("json")
local utils = require("utils")
local hash = require("hash")

local library_path = os.getenv("OPENSHELF_LIBRARY_PATH") or "/path/to/your/openshelf_library/"

local library = {}

function library.readLibraryJson()
    local f, err = io.open(library_path .. "library.json", "r")
    if not f then
    		utils.log("Error opening file: " .. err)
        return nil
    end
    local raw = f:read("*a")
    f:close()
    
    local ok, data = pcall(json.decode, raw)
	if ok then
		return data
	else
		utils.log("JSON parse error: " .. data)
		return nil
	end
end

-- Return the path to the reading state file for a given bookId
function library.get_state_file_path(book_id)
    -- filesystem-safe name: replace ':' with '_'
    local safe_id = book_id:gsub(":", "_")
    return library_path .. "/.state/" .. safe_id .. ".json"
end

-- Load the reading state for a bookId
-- Returns: table (state) or nil
function library.load_reading_state(book_id)
    local state_path = library.get_state_file_path(book_id)

    -- Read the file
    local f, err = io.open(state_path, "r")
    if not f then
        return nil
    end
    local content = f:read("*a")
    f:close()

    -- Decode JSON
    local ok, state = pcall(json.decode, content)
    if not ok then
        utils.log("Failed to parse JSON in state file: " .. state_path)
        return nil
    end

    return state
end

function library.write_reading_state(book_id, location)
    local state_path = library.get_state_file_path(book_id)

    -- Load existing state if it exists
    local state
    local f = io.open(state_path, "r")
    if f then
        local content = f:read("*a")
        f:close()
        local ok
        ok, state = pcall(json.decode, content)
        if not ok or type(state) ~= "table" then
            utils.log("Failed to parse existing state file, starting fresh")
            state = {}
        end
    else
        utils.log("No existing state file, creating new state")
        state = {}
    end

    -- Ensure 'location' table exists
    state.location = state.location or {}

    -- Merge/update location fields
    for k, v in pairs(location) do
        state.location[k] = v
    end

    -- Ensure the .state directory exists
    os.execute('mkdir -p "' .. library_path .. '/.state"')

    -- Write the updated state
    f, err = io.open(state_path, "w")
    if not f then
        utils.log("Failed to open state file for writing: " .. tostring(err))
        return false
    end

    f:write(json.encode(state))
    f:close()

    utils.log("State saved at " .. state_path)
    return true
end

return library
