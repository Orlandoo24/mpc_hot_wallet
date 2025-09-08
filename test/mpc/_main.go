package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// generateTestPartyIDs is a helper function to create party IDs for the test.
func generateTestPartyIDs(count, start int) tss.SortedPartyIDs {
	partyIDs := make(tss.UnSortedPartyIDs, 0, count)
	for i := start; i < start+count; i++ {
		pMoniker := fmt.Sprintf("party-%d", i)
		pID := tss.NewPartyID(fmt.Sprintf("%d", i), pMoniker, big.NewInt(int64(i)))
		partyIDs = append(partyIDs, pID)
	}
	return tss.SortPartyIDs(partyIDs)
}

func main() {
	// 1. 设置常量
	const (
		partyCount = 3
		threshold  = 1 // threshold + 1 = 2 signers
	)

	// 2. 生成各参与方的排序 ID
	partyIDs := generateTestPartyIDs(partyCount, 0)

	// 3. 运行密钥生成协议
	fmt.Println("--- 阶段 1: 密钥生成 ---")
	savedData, err := runKeygen(partyIDs, threshold)
	if err != nil {
		fmt.Printf("密钥生成失败: %v\n", err)
		return
	}
	fmt.Println("密钥生成成功完成。")

	// 从任意一方的保存数据中获取聚合公钥
	pubKeyPoint := savedData[0].ECDSAPub
	// 派生以太坊地址
	ecdsaPubKey := &ecdsa.PublicKey{
		Curve: tss.S256(),
		X:     pubKeyPoint.X(),
		Y:     pubKeyPoint.Y(),
	}
	address := PubKeyToEthAddress(ecdsaPubKey)
	fmt.Println("生成的新钱包地址:", address)

	// 4. 运行签名协议
	fmt.Println("\n--- 阶段 2: 签名 ---")
	// 定义待签名的消息 (通常是交易哈希)
	msgToSign := []byte("hashlink-mpc-demo")
	msgHash := ethcrypto.Keccak256(msgToSign)

	// We need threshold + 1 signers
	signingPartyIDs := partyIDs[:threshold+1]
	signingSaveData := make([]keygen.LocalPartySaveData, 0, threshold+1)

	// Find the signing parties' save data
	partyIDMap := make(map[string]struct{})
	for _, pID := range signingPartyIDs {
		partyIDMap[pID.Id] = struct{}{}
	}
	for _, data := range savedData {
		if _, ok := partyIDMap[data.ShareID.String()]; ok {
			signingSaveData = append(signingSaveData, data)
		}
	}

	signature, err := runSigning(signingPartyIDs, signingSaveData, msgHash)
	if err != nil {
		fmt.Printf("签名失败: %v\n", err)
		return
	}
	fmt.Println("签名成功完成。")

	// 5. 验证签名
	fmt.Printf("签名 R: %s\n", hex.EncodeToString(signature.R))
	fmt.Printf("签名 S: %s\n", hex.EncodeToString(signature.S))

	// 将签名数据转换为以太坊格式
	ethSig := append(signature.R, signature.S...)
	ethSig = append(ethSig, signature.SignatureRecovery[0]) // 添加 V

	// 使用 go-ethereum 的库来验证签名
	isValid := ethcrypto.VerifySignature(
		ethcrypto.FromECDSAPub(ecdsaPubKey),
		msgHash,
		ethSig[:64], // 去掉 V
	)

	if isValid {
		fmt.Println("ECDSA 签名验证成功!")
	} else {
		fmt.Println("ECDSA 签名验证失败!")
	}
}

// runKeygen 模拟并运行密钥生成协议
func runKeygen(partyIDs tss.SortedPartyIDs, threshold int) ([]keygen.LocalPartySaveData, *tss.Error) {
	// 初始化各方的通信通道
	outCh := make(chan tss.Message, len(partyIDs))
	endCh := make(chan *keygen.LocalPartySaveData, len(partyIDs))
	errCh := make(chan *tss.Error, len(partyIDs))

	parties := make(map[*tss.PartyID]tss.Party, len(partyIDs))
	var wg sync.WaitGroup
	wg.Add(len(partyIDs))

	// 提前生成安全素数等参数，可以加速密钥生成过程
	preParams, _ := keygen.GeneratePreParams(1 * time.Minute)

	// 为每个参与方创建和启动一个 keygen.LocalParty 实例
	for _, pID := range partyIDs {
		params := tss.NewParameters(tss.S256(), tss.NewPeerContext(partyIDs), pID, len(partyIDs), threshold)
		P := keygen.NewLocalParty(params, outCh, endCh, *preParams)
		parties[pID] = P
		go func(P tss.Party) {
			defer wg.Done()
			if err := P.Start(); err != nil {
				errCh <- err
			}
		}(P)
	}

	// 模拟网络，路由消息
	go routeMessages(outCh, parties, errCh)

	// 等待所有参与方完成
	wg.Wait()
	close(endCh)
	close(errCh)

	// 收集结果
	var savedData []keygen.LocalPartySaveData
	for data := range endCh {
		savedData = append(savedData, *data)
	}
	for err := range errCh {
		return nil, err
	}
	return savedData, nil
}

// runSigning 模拟并运行签名协议
func runSigning(partyIDs tss.SortedPartyIDs, savedData []keygen.LocalPartySaveData, msgHash []byte) (*common.SignatureData, *tss.Error) {
	outCh := make(chan tss.Message, len(partyIDs))
	endCh := make(chan *common.SignatureData, len(partyIDs))
	errCh := make(chan *tss.Error, len(partyIDs))

	parties := make(map[*tss.PartyID]tss.Party, len(partyIDs))
	var wg sync.WaitGroup
	wg.Add(len(partyIDs))

	msgInt := new(big.Int).SetBytes(msgHash)

	ctx := tss.NewPeerContext(partyIDs)
	for i, pID := range partyIDs {
		params := tss.NewParameters(tss.S256(), ctx, pID, len(partyIDs), len(partyIDs)-1)
		P := signing.NewLocalParty(msgInt, params, savedData[i], outCh, endCh)
		parties[pID] = P
		go func(P tss.Party) {
			defer wg.Done()
			if err := P.Start(); err != nil {
				errCh <- err
			}
		}(P)
	}

	go routeMessages(outCh, parties, errCh)

	wg.Wait()
	close(endCh)
	close(errCh)

	var signature *common.SignatureData
	for sig := range endCh {
		// All parties will have the same signature, so we can just take the first one
		signature = sig
		break
	}
	for err := range errCh {
		return nil, err
	}
	return signature, nil
}

// routeMessages is a simple message router that simulates network communication between parties.
func routeMessages(outCh <-chan tss.Message, parties map[*tss.PartyID]tss.Party, errCh chan<- *tss.Error) {
	for msg := range outCh {
		// Use a temporary variable for the goroutine
		msgToRoute := msg

		// Verifying the message is ParsedMessage
		parsedMsg, ok := msgToRoute.(tss.ParsedMessage)
		if !ok {
			// This should not happen with the current tss-lib implementation
			errCh <- tss.NewError(fmt.Errorf("could not cast message to ParsedMessage"), "routing", -1, msgToRoute.GetFrom())
			return
		}

		dest := parsedMsg.GetTo()
		if dest == nil { // broadcast
			for pID, p := range parties {
				if pID == parsedMsg.GetFrom() {
					continue
				}
				go func(p tss.Party, parsedMsg tss.ParsedMessage) {
					if _, err := p.Update(parsedMsg); err != nil {
						errCh <- err
					}
				}(p, parsedMsg)
			}
		} else { // point-to-point
			for _, pID := range dest {
				if pID == parsedMsg.GetFrom() {
					continue
				}
				p, ok := parties[pID]
				if !ok {
					errCh <- tss.NewError(fmt.Errorf("party %s not found", pID), "routing", -1, parsedMsg.GetFrom())
					continue
				}
				go func(p tss.Party, parsedMsg tss.ParsedMessage) {
					if _, err := p.Update(parsedMsg); err != nil {
						errCh <- err
					}
				}(p, parsedMsg)
			}
		}
	}
}

// PubKeyToEthAddress converts a TSS public key to an Ethereum address.
func PubKeyToEthAddress(pubKey *ecdsa.PublicKey) string {
	addressBytes := ethcrypto.PubkeyToAddress(*pubKey)
	return "0x" + hex.EncodeToString(addressBytes[:])
}
