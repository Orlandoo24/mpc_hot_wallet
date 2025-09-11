curl --noproxy localhost -X POST http://localhost:8888/api/bridge/wrap \
-H "Content-Type: application/json" \
-d '{
"from_address": "0xDC881Fe66502a439dF3bCaE13e4068b5163AfD37",
"from_chain": 56,
"to_chain": 137,
"from_token": "0x55d398326f99059fF775485246999027B3197955",
"to_token": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
"amount": "1000000000000000",
"to_address": "0xDC881Fe66502a439dF3bCaE13e4068b5163AfD37",
"order": "CHEAPEST",
"slippage": "0.005"
}'

{"tx_hash":"0x73354232250f2956bdcd4c7e071f8ce282369a08e1beef2d94faefb33d08a562","message":"✅ 跨链转账已提交！从链 56 到链 137，交易哈希: 0x73354232250f2956bdcd4c7e071f8ce282369a08e1beef2d94faefb33d08a562。请使用 /bridge/status 查询进度。","explorer_url":"https://bscscan.com/tx/0x73354232250f2956bdcd4c7e071f8ce282369a08e1beef2d94faefb33d08a562","from_chain":56,"to_chain":137,"status":"pending"}