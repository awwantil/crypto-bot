https://tproger.ru/translations/deploy-a-secure-golang-rest-api

1. POST - http://localhost:8000/api/user/new
   {
   "email": "pomo@gmail.com",
   "password": "rrr444@@"
   }
2. POST - http://localhost:8000/api/user/login
   {
   "email": "pomo@gmail.com",
   "password": "rrr444@@"
   }
3. POST - http://localhost:8000/api/contacts/new
   {
   "name": "Иван",
   "phone": "981-082-099"
   }
4. GET - http://localhost:8000/api/me/contacts

   	resp1, resp2, err := app.GetApi().Isolated.GetAccount("USDT")
   	logger.Info("err: ", err)
   	logger.Info("resp1: ", resp1)
   	logger.Info("resp2: ", string(resp2))
INFO[2024-04-07T14:21:25+03:00] resp1: map[USDT:{USDT 4980.455107212083 4920.07177387875 60.383333333333326}]  app=okx-bot component=app.main-rest
INFO[2024-04-07T14:21:25+03:00] resp2: {"code":"0","data":[{"adjEq":"","borrowFroz":"","details":[{"availBal":"4920.07177387875","availEq":"4920.07177387875","borrowFroz":"","cashB
al":"4920.07177387875","ccy":"USDT","clSpotInUseAmt":"","crossLiab":"","disEq":"4981.5508073356705","eq":"4980.455107212083","eqUsd":"4981.5508073356705","fixedBal":"0","frozenBal"
:"60.383333333333326","imr":"0","interest":"","isoEq":"60.383333333333326","isoLiab":"","isoUpl":"0.3799999999999955","liab":"","maxLoan":"","maxSpotInUse":"","mgnRatio":"","mmr":"
0","notionalLever":"0","ordFrozen":"0","rewardBal":"0","smtSyncEq":"0","spotInUseAmt":"","spotIsoBal":"0","stgyEq":"0","twap":"0","uTime":"1712488777025","upl":"0.3799999999999955"
,"uplLiab":""}],"imr":"","isoEq":"60.396617666666664","mgnRatio":"","mmr":"","notionalUsd":"","ordFroz":"","totalEq":"83579.94080733568","uTime":"1712488885622","upl":""}],"msg":""}  app=okx-bot component=app.main-rest


INFO[2024-04-08T17:25:49+03:00] resp1: 
map[USDT:{USDT 4980.919718892807 4920.786385559473 60.13333333333333}]  app=okx-bot component=app.main-rest

INFO[2024-04-08T17:25:49+03:00] resp2: {"code":"0","data":[{"adjEq":"","borrowFroz":"","details":
[{"availBal":"4920.786385559473","availEq":"4920.786385559473","borrowFroz":"","cashBal":"4920.786385559473","ccy":"USDT",
"crossLiab":"","disEq":"4981.26838327313","eq":"4980.919718892807","eqUsd":"4981.26838327313","fixedBal":"0",
"frozenBal":"60.13333333333333",
"imr":"0","interest":"","isoEq":"60.13333333333333","isoLiab":"","isoUpl":"0.0600000000000023","liab":"","maxLoan":"","mgnRatio":"","mmr":"0","notionalLever":"0","ordFrozen":"0","r
ewardBal":"0","smtSyncEq":"0","spotInUseAmt":"","spotIsoBal":"0","stgyEq":"0","twap":"0","uTime":"1712586286269","upl":"0.0600000000000023","uplLiab":""}],"imr":"","isoEq":"60.13754266666667","mgnRatio":"","mmr":"","notionalUsd":"","ordFroz":"","totalEq":"86262.16838327312","uTime":"1712586350746","upl":""}],"msg":""}  app=okx-bot component=app.main-rest    
INFO[2024-04-08T17:25:50+03:00] err: <nil>                                    app=okx-bot component=app.main-rest

INFO[2024-04-08T17:25:50+03:00] resp1: &{SWAP SOL-USDT-SWAP   1344992036067729408  86d4a3bf87bcBCDE     10 0 market sell net isolated 10 180.22 68435015 10 1712586286268 180.22 filled  cancel_maker 3        []              2024-04-08 17:24:46 +0300 MSK 2024-04-08 17:24:46 +0300 MSK}  app=okx-bot component=app.main-rest

INFO[2024-04-08T17:25:50+03:00] resp2: {"code":"0","data":[
{"accFillSz":"10","algoClOrdId":"","algoId":"","attachAlgoClOrdId":"","attachAlgoOrds":[],"avgPx":"180.22","cTime":"1712586286267",
"cancelSource":"","cancelSourceReason":"","category":"normal","ccy":"","clOrdId":"","fee":"-0.09011","feeCcy":"USDT",
"fillPx":"180.22","fillSz":"10","fillTime":"1712586286268",
"instId":"SOL-USDT-SWAP","instType":"SWAP","isTpLimit":"false","lever":"3",
"linkedAlgoOrd":{"algoId":""},"ordId":"1344992036067729408","ordType":"market","pnl":"0","posSide":"net",
"px":"","pxType":"","pxUsd":"","pxVol":"","quickMgnType":"","rebate":"0",
"rebateCcy":"USDT","reduceOnly":"false","side":"sell","slOrdPx":"","slTriggerPx":"","slTriggerPxType"
:"","source":"","state":"filled","stpId":"","stpMode":"cancel_maker","sz":"10",
"tag":"86d4a3bf87bcBCDE","tdMode":"isolated","tgtCcy":"","tpOrdPx":"","tpTriggerPx":"","tpTriggerPxType":"",
"tradeId":"68435015","uTime":"1712586286269"}],"msg":""}  app=okx-bot component=app.main-rest

INFO[2024-04-08T17:25:50+03:00] resp1: 
&{SWAP SOL-USDT-SWAP   1344992036067729408  86d4a3bf87bcBCDE     10 0 market sell net isolated 
10 180.22 68435015 10 1712586286268 180.22 filled  cancel_maker 3        
[]              2024-04-08 17:24:46 +0300 MSK 2024-04-08 17:24:46 +0300 MSK}  app=okx-bot component=app.main-rest


INFO[2024-04-08T17:47:32+03:00] resp4: {"code":"0","data":[{"accFillSz":"10","algoClOrdId":"","algoId":"",
"attachAlgoClOrdId":"","attachAlgoOrds":[],"avgPx":"180.22","cTime":"1712586286267",
"cancelSource":"","cancelSourceReason":"","category":"normal","ccy":"","clOrdId":"","fee":"-0.09011","feeCcy":"USDT",
"fillPx":"180.22","fillSz":"10","fillTime":"1712586286268",
"instId":"SOL-USDT-SWAP","instType":"SWAP","isTpLimit":"false","lever":"3","linkedAlgoOrd":{"algoId":""},
"ordId":"1344992036067729408","ordType":"market","pnl":"0","posSide":"net","px":"","pxType":"","pxUsd":"","pxVol":"",
"quickMgnType":"","rebate":"0","rebateCcy":"USDT","reduceOnly":"false","side":"sell","slOrdPx":"","slTriggerPx":"","slTriggerPxType"
:"","source":"","state":"filled","stpId":"","stpMode":"cancel_maker","sz":"10","tag":"86d4a3bf87bcBCDE",
"tdMode":"isolated","tgtCcy":"","tpOrdPx":"","tpTriggerPx":"","tpTriggerPxType":"","tradeId":"68435015",
"uTime":"1712586286269"}],"msg":""}  app=okx-bot component=app.main-rest
