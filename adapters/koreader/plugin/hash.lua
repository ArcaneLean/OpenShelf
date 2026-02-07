local ok, crypto = pcall(require, "crypto")
if not ok then
    print("KOReader crypto not available, using stub")
    crypto = require("plugin.crypto_stub")  -- your stub
end

local utils = require("plugin.utils")

local hash = {}

function hash.compute_sha256(file_path)
    local file = io.open(file_path, "rb")
    if not file then
        utils.log("ERROR: cannot open file " .. file_path)
        return nil
    end

    local ctx = crypto.sha256_init()
    local chunk_size = 1024 * 1024 -- 1 MB
    while true do
        local chunk = file:read(chunk_size)
        if not chunk then break end
        crypto.sha256_update(ctx, chunk)
    end
    file:close()

    local digest = crypto.sha256_final(ctx)
    return "sha256_" .. digest
end

return hash

