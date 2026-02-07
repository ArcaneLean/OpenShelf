local lfs = {}

function lfs.attributes(path)
    local f = io.open(path, "r")
    if f then
        f:close()
        return { mode = "file" }
    else
        return nil
    end
end

function lfs.mkdir(path)
    local ok, err = os.execute("mkdir -p " .. path)
    if ok then return true end
    return nil, err
end

return lfs

