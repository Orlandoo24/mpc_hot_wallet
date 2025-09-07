curl --noproxy localhost -X POST http://localhost:8888/api/transaction/send -H "Content-Type: application/json" -d '{"from_address": "0xDC881Fe66502a439dF3bCaE13e4068b5163AfD37", "to_address": "0x5c5670d9EaeC91AdcDfC4424a8C2490FCC6049b2", "chain": "BSC", "from_token": "0x0000000000000000000000000000000000000000", "to_token": "0x0000000000000000000000000000000000000000", "amount": "100000000000000"}'


{"tx_hash":"0x98776d6257e84f7a4c9183fdaaaec2b07df67070caef44ce7170e75210178919",
"message":"✅ BSC 主网 原生代币转账已提交！交易正在异步处理中，请通过区块浏览器查询最终状态。","explorer_url":"https:/
/bscscan.com/tx/0x98776d6257e84f7a4c9183fdaaaec2b07df67070caef44ce7170e752101789
19","chain":"BSC","status":"pending"}