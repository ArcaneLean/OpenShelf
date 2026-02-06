local adapter = {}

adapter.name = "openshelf-koreader"
adapter.levels = {1, 2}  -- read + write
adapter.supportedLocations = {"epubcfi", "percentage", "page"}
adapter.specVersion = "0.1.0"

return adapter

