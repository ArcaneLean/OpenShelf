local sha = require("hash.sha2")

local hash = {}

-- Compute SHA-256 of a file and return as "sha256_<hex>"
function hash.sha256(file_path)
    local f, err = io.open(file_path, "rb")
    if not f then
        return nil, err
    end

    -- get the append function for streaming SHA-256
    local append = sha.sha256()

    -- feed file in chunks
    while true do
        local chunk = f:read(8192)
        if not chunk then break end
        append(chunk)
    end

    f:close()

    -- get final digest
    local digest = append()
    return "sha256_" .. digest
end

return hash

