local _ = require("gettext")

return {
    name = "openshelf",
    fullname = _("OpenShelf Adapter"),
    description = _("Provides OpenShelf library integration for KOReader."),

    adapter = {
        levels = {1, 2},
        supportedLocations = {"percentage", "epubcfi", "page"},
    },
}
