local crypto = {}

-- simple SHA-256 using Lua + luarocks library, or a fake value
-- for offline development we can return a fixed hex for testing
function crypto.sha256(data)
    -- TODO: optional: use an actual Lua SHA-256 lib
    return string.rep("a", 64)  -- 64-character hex string
end

function crypto.sha256_init()
    return { buffer = "" }
end

function crypto.sha256_update(ctx, chunk)
    ctx.buffer = ctx.buffer .. chunk
end

function crypto.sha256_final(ctx)
    -- simple fake hash
    return string.rep("b", 64)
end

return crypto

