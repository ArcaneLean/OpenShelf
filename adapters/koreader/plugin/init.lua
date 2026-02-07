local adapter = require("adapter")
local state   = require("state")
local hash    = require("hash")
local utils   = require("utils")

-- Hook: plugin load
local function on_plugin_load()
    utils.log("OpenShelf KOReader adapter loaded")
    utils.log("Supported levels: " .. table.concat(adapter.levels, ", "))
end

-- Helper: extract KOReader location
local function extract_ko_location(book)
    -- TODO: replace with actual KOReader API calls
    -- Returns a table of {type="percentage|epubcfi|page", value=<...>}
    local locations = {}

    -- Example placeholders:
    table.insert(locations, { type = "percentage", value = book:get_percentage() or 0 })
    table.insert(locations, { type = "epubcfi", value = book:get_cfi() or "" })
    table.insert(locations, { type = "page", value = book:get_page() or 0 })

    return locations
end

-- Hook: book open
local function on_book_open(book_path, book)
    utils.log("Opening book: " .. book_path)
    local book_id = hash.compute_sha256(book_path)
    if not book_id then
        utils.log("ERROR: cannot compute Book ID")
        return
    end

    local reading_state = state.read_state(book_id)
    if reading_state then
        utils.log("Found OpenShelf state for book: " .. book_id)

        -- Choose best location to resume
        local resume_location
        if reading_state.location.epubcfi and reading_state.location.epubcfi ~= "" then
            resume_location = { type="epubcfi", value=reading_state.location.epubcfi }
        elseif reading_state.location.percentage then
            resume_location = { type="percentage", value=reading_state.location.percentage }
        elseif reading_state.location.page then
            resume_location = { type="page", value=reading_state.location.page }
        end

        if resume_location then
            utils.log("Resuming at " .. resume_location.type .. ": " .. tostring(resume_location.value))
            -- TODO: apply to KOReader book object
            book:set_location(resume_location)
        else
            utils.log("No usable location found; starting fresh")
        end
    else
        utils.log("No OpenShelf state found; starting fresh")
    end
end

-- Hook: book close
local function on_book_close(book_path, book)
    utils.log("Closing book: " .. book_path)
    local book_id = hash.compute_sha256(book_path)
    if not book_id then
        utils.log("ERROR: cannot compute Book ID")
        return
    end

    -- Extract current KOReader location
    local locations = extract_ko_location(book)
    if #locations == 0 then
        utils.log("No location info; skipping state write")
        return
    end

    -- Write to OpenShelf state
    local ok = state.write_state(book_id, locations)
    if not ok then
        utils.log("Failed to write OpenShelf state for book: " .. book_id)
    end
end

-- Register plugin hooks (KOReader API)
koreader.register_plugin({
    name       = adapter.name,
    on_load    = on_plugin_load,
    on_book_open  = on_book_open,
    on_book_close = on_book_close,
})

