box.cfg { listen = 3301 }
box.schema.user.grant('guest', 'super', nil, nil, { if_not_exists = true })

function test_func(param)
    return 'it works: ' .. param
end
