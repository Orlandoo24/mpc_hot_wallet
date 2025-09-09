curl --noproxy localhost -X POST http://localhost:8888/api/transaction/swap \
-H "Content-Type: application/json" \
-d '{
"from_address": "0xDC881Fe66502a439dF3bCaE13e4068b5163AfD37",
"to_address": "0xDC881Fe66502a439dF3bCaE13e4068b5163AfD37",
"chain": "BSC",
"from_token": "0x0000000000000000000000000000000000000000",
"to_token": "0x55d398326f99059fF775485246999027B3197955",
"amount": "10000000000000"
}'

{"tx_hash":"0x343dd7ed585a1def11ce51629bc9469b07cf80244f27e4e707ef5ca10039b4bb","message":"✅ Swap 交易已提交！使用 kyberswap 工具，交易哈希: 0x343dd7ed585a1def11ce51629bc9469b07cf80244f27e4e707ef5ca10039b4bb","explorer_url":"https://bscscan.com/tx/0x343dd7ed585a1def11ce51629bc9469b07cf80244f27e4e707ef5ca10039b4bb","chain":"BSC","status":"pending"}