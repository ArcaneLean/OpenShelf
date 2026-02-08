local WidgetContainer = require("ui/widget/container/widgetcontainer")
local ReaderUI = require("apps/reader/readerui")
local utils = require("utils")
local hash = require("hash")
local library = require("library")
local meta = require("_meta")
local Event = require("ui/event")

local OpenShelf = WidgetContainer:extend{
    name = "openshelf",
    is_doc_only = false,
}

local session_initialized = false

function OpenShelf:init()
    if session_initialized then return end
    session_initialized = true

    utils.log("OpenShelf KOReader adapter initialized:")
    utils.log("\tSupported levels: " .. table.concat(meta.adapter.levels, ", "))
    utils.log("\tSupported locations: " .. table.concat(meta.adapter.supportedLocations, ", "))
    local library_json = library.readLibraryJson()
    utils.log("OpenShelf Library found:")
    utils.log("\tVersion: " .. library_json.spec.version)
    utils.log("\tCapabilities: " .. table.concat(library_json.capabilities, ", "))
end

function OpenShelf:onReaderReady(doc_settings)
	self.doc_data = doc_settings and doc_settings.data
	local doc_path = self.doc_data
		and self.doc_data.doc_path
	local filename = doc_path:match("^.+[/\\](.+)$") or doc_path
	local book_id = hash.sha256(doc_path)
	local reading_state = library.load_reading_state(book_id)
	if reading_state then
		utils.log("Found reading state for " .. filename)
	else
		utils.log("No reading state found for " .. filename)
		return
	end
	
	self.ui:handleEvent(Event:new("GotoPercent", reading_state.location.percentage))
end

function OpenShelf:onCloseDocument()
    local doc_path = self.doc_data
        and self.doc_data.doc_path
	local filename = doc_path:match("^.+[/\\](.+)$") or doc_path
	utils.log("Closing " .. filename)
	local book_id = hash.sha256(doc_path)
	local loc = {
		percentage = self.doc_data.percent_finished * 100,
		page = math.floor(self.doc_data.percent_finished * self.doc_data.doc_pages),
		epubcfi = self.doc_data.last_xpointer
	}
	library.write_reading_state(book_id, loc)
end

return OpenShelf
