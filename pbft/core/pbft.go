	package core

	import (
		"bufio"
		"crypto"
		"crypto/rand"
		"crypto/rsa"
		"crypto/sha256"
		"crypto/x509"
		"encoding/hex"
		"encoding/json"
		"encoding/pem"
		"fmt"
		"io"
		"log"
		"os"
		"strconv"
		"sync"
		"time"
	)

	type State int

	var nodeCount int

	const (
		cPrePrepare State = iota
		cPrepare
		cCommit
	)

	type Message string

	// 本地消息池（模拟持久化层），只有确认提交成功后才会存入此池
	var localMessagePool = []Message{}

	type node struct {
		//节点ID
		nodeID string
		//节点监听地址
		addr string
		//RSA私钥
		rsaPrivKey []byte
		//RSA公钥
		rsaPubKey []byte
	}

	type Request struct {
		Message Message
	}

	type pbft struct {
		//节点信息
		node node
		//每笔请求自增序号
		sequenceID int
		//锁
		lock sync.Mutex
		//临时消息池，消息摘要对应消息本体
		messagePool map[string]Request
		//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
		prePareConfirmCount map[string]map[string]bool
		//存放收到的commit数量(至少需要收到并确认2f+1个)，根据摘要来对应
		commitConfirmCount map[string]map[string]bool
		//该笔消息是否已进行Commit广播
		isCommitBordcast map[string]bool
		//该笔消息是否已对客户端进行Reply
		isReply map[string]bool
	}

	var cnt int

	func (p *pbft) sequenceIDAdd() {
		p.sequenceID = cnt
		cnt++
	}

	func getDigest(r Request) string {
		data, _ := json.Marshal(r)
		tmp := sha256.Sum256(data)
		slice := tmp[:]
		return string(slice)
	}

	func (p *pbft) RsaSignWithSha256(digestByte []byte, rsaPrivKey []byte) []byte {
		// 解码PEM格式的私钥
		block, _ := pem.Decode(rsaPrivKey)
		if block == nil {
			return nil
		}

		// 解析RSA私钥
		priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil
		}

		signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, digestByte)
		if err != nil {
			return nil
		}

		return signature
	}

	// request
	func (p *pbft) handleClientRequest(content []byte) {
		fmt.Println("主节点已接收到客户端发来的request ...")
		//使用json解析出Request结构体
		r := new(Request)
		err := json.Unmarshal(content, r)
		if err != nil {
			log.Panic(err)
		}
		//添加信息序号
		p.sequenceIDAdd()
		//获取消息摘要
		digest := getDigest(*r)
		fmt.Println("收到的request消息为: ", r.Message)
		fmt.Println("已将request存入临时消息池")
		//存入临时消息池
		p.messagePool[digest] = *r
		//主节点对消息摘要进行签名
		digestByte, _ := hex.DecodeString(digest)
		signInfo := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
		//拼接成PrePrepare，准备发往follower节点
		pp := PrePrepare{*r, digest, p.sequenceID, signInfo}
		b, err := json.Marshal(pp)
		if err != nil {
			log.Panic(err)
		}
		pause()
		fmt.Println("正在向其他节点进行进行PrePrepare广播 ...")
		fmt.Println("PrePrepare消息内容为: ", pp)
		//进行PrePrepare广播
		p.broadcast(cPrePrepare, b)
		fmt.Println("PrePrepare广播完成")
		pause()
	}

	func (p *pbft) broadcast(state State, jsonData []byte) {
		fifoName := "pbftFiFo"
		file, err := os.OpenFile(fifoName, os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Error opening FIFO: %v", err)
		}
		defer file.Close()
		// 将State编码为JSON
		stateJSON, err := json.Marshal(state)
		if err != nil {
			log.Fatal("error marshaling state to JSON: %v", err)
		}

		// 构建组合消息
		compositeMessage := map[string]interface{}{
			"state":    stateJSON,
			"jsonData": string(jsonData),
		}
		combinedJSON, err := json.Marshal(compositeMessage)
		if err != nil {
			log.Fatal("error marshaling composite message to JSON: %v", err)
		}

		// 写入FIFO
		_, err = file.Write(combinedJSON)
		if err != nil {
			log.Fatal("error writing JSON to FIFO: %v", err)
		}
		fmt.Println("Broadcasted State and JSON Data to FIFO.")
	}

	// 在开启时，使用一个协程运行它
	func (p *pbft) received() {
		fifoName := "pbftFiFo"
		file, err := os.Open(fifoName)
		if err != nil {
			log.Fatalf("Error opening FIFO: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		for {
			var combinedMsg map[string]interface{}
			err := decoder.Decode(&combinedMsg)
			if err != nil {
				if err == io.EOF {
					fmt.Println("No more data.")
					time.Sleep(500 * time.Millisecond)
				} else {
					log.Fatalf("Error decoding JSON from FIFO: %v", err)
				}
			}

			// 解析State
			var state State
			err = json.Unmarshal([]byte(combinedMsg["state"].(string)), &state)
			if err != nil {
				log.Fatalf("Error unmarshaling state: %v", err)
			}
			fmt.Printf("Received State: %+v\n", state)

			rawJSONData := combinedMsg["jsonData"].([]byte)
			switch state {
			case cPrePrepare:
				p.handlePrePrepare(rawJSONData)
			case cPrepare:
				p.handlePrepare(rawJSONData)
			case cCommit:
				p.handleCommit(rawJSONData)
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}

	func pause() {
		time.Sleep(time.Millisecond)
	}

	// pre-prepare
	type PrePrepare struct {
		RequestMessage Request
		Digest         string
		SequenceID     int
		Sign           []byte
	}

	func (p *pbft) getPubKey(str string) []byte {
		return nil
	}

	func (p *pbft) RsaVerySignWithSha256(digest []byte, sign []byte, pubKey []byte) bool {
		// 解码PEM格式的公钥
		block, _ := pem.Decode(pubKey)
		if block == nil {
			return false
		}

		// 解析RSA公钥
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return false
		}
		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return false
		}

		// 验证签名
		err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest, sign)
		if err != nil {
			return false
		}

		return true
	}

	type Prepare struct {
		Digest     string
		SequenceID int
		NodeID     string
		Sign       []byte
	}

	func (p *pbft) handlePrePrepare(content []byte) {
		fmt.Println("本节点已接收到主节点发来的PrePrepare ...")
		pause()
		// //使用json解析出PrePrepare结构体
		pp := new(PrePrepare)
		err := json.Unmarshal(content, pp)
		if err != nil {
			log.Panic(err)
		}
		//获取主节点的公钥，用于数字签名验证
		primaryNodePubKey := p.getPubKey("N0")
		digestByte, _ := hex.DecodeString(pp.Digest)
		if digest := getDigest(pp.RequestMessage); digest != pp.Digest {
			fmt.Println("信息摘要对不上，拒绝进行prepare广播")
		} else if p.sequenceID+1 != pp.SequenceID {
			fmt.Println("消息序号对不上，拒绝进行prepare广播")
		} else if !p.RsaVerySignWithSha256(digestByte, pp.Sign, primaryNodePubKey) {
			fmt.Println("主节点签名验证失败！,拒绝进行prepare广播")
		} else {
			//序号赋值
			p.sequenceID = pp.SequenceID
			//将信息存入临时消息池
			fmt.Println("已将消息存入临时节点池")
			p.messagePool[pp.Digest] = pp.RequestMessage
			//节点使用私钥对其签名
			sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
			//拼接成Prepare
			pre := Prepare{pp.Digest, pp.SequenceID, p.node.nodeID, sign}
			bPre, err := json.Marshal(pre)
			if err != nil {
				log.Panic(err)
			}
			//进行准备阶段的广播
			fmt.Println("正在进行Prepare广播 ...")
			fmt.Println("广播的Prepare消息内容为: ", pre)
			p.broadcast(cPrepare, bPre)
			fmt.Println("Prepare广播完成")
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Press enter to continue...")
			_, _ = reader.ReadString('\n')
		}
	}

	// prepare
	func (p *pbft) setPrePareConfirmMap(digest string, nodeID string, flag bool) {

	}

	func (p *pbft) handlePrepare(content []byte) {
		//使用json解析出Prepare结构体
		pre := new(Prepare)
		err := json.Unmarshal(content, pre)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("本节点已接收到%s节点发来的Prepare ... \n", pre.NodeID)
		//获取消息源节点的公钥，用于数字签名验证
		MessageNodePubKey := p.getPubKey(pre.NodeID)
		digestByte, _ := hex.DecodeString(pre.Digest)
		if _, ok := p.messagePool[pre.Digest]; !ok {
			fmt.Println("当前临时消息池无此摘要，拒绝执行commit广播")
		} else if p.sequenceID != pre.SequenceID {
			fmt.Println("消息序号对不上，拒绝执行commit广播")
		} else if !p.RsaVerySignWithSha256(digestByte, pre.Sign, MessageNodePubKey) {
			fmt.Println("节点签名验证失败！,拒绝执行commit广播")
		} else {
			p.setPrePareConfirmMap(pre.Digest, pre.NodeID, true)
			count := 0
			for range p.prePareConfirmCount[pre.Digest] {
				count++
			}
			//因为主节点不会发送Prepare，所以不包含自己
			specifiedCount := 0
			if p.node.nodeID == "N0" {
				specifiedCount = nodeCount / 3 * 2
			} else {
				specifiedCount = (nodeCount / 3 * 2) - 1
			}
			//如果节点至少收到了2f个prepare的消息（包括自己）,并且没有进行过commit广播，则进行commit广播
			p.lock.Lock()
			//获取消息源节点的公钥，用于数字签名验证
			if count >= specifiedCount && !p.isCommitBordcast[pre.Digest] {
				pause()
				fmt.Println("本节点已收到至少2f个节点(包括本地节点)发来的Prepare信息，内容为： ", pre)
				//节点使用私钥对其签名
				sign := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
				c := Commit{pre.Digest, pre.SequenceID, p.node.nodeID, sign}
				bc, err := json.Marshal(c)
				if err != nil {
					log.Panic(err)
				}
				//进行提交信息的广播
				fmt.Println("正在进行commit广播 ...")
				fmt.Println("广播的commit消息内容为: ", bc)
				p.broadcast(cCommit, bc)
				p.isCommitBordcast[pre.Digest] = true
				fmt.Println("commit广播完成")
			}
			p.lock.Unlock()
			pause()
		}
	}

	// commit
	type Commit struct {
		Digest     string
		SequenceID int
		NodeID     string
		Sign       []byte
	}

	func (p *pbft) setCommitConfirmMap(digest string, nodeID string, flag bool) {

	}
	func (p *pbft) handleCommit(content []byte) {
		//使用json解析出Commit结构体
		c := new(Commit)
		err := json.Unmarshal(content, c)
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("本节点已接收到%s节点发来的Commit ... \n", c.NodeID)
		//获取消息源节点的公钥，用于数字签名验证
		MessageNodePubKey := p.getPubKey(c.NodeID)
		digestByte, _ := hex.DecodeString(c.Digest)
		if _, ok := p.prePareConfirmCount[c.Digest]; !ok {
			fmt.Println("当前prepare池无此摘要，拒绝将信息持久化到本地消息池")
		} else if p.sequenceID != c.SequenceID {
			fmt.Println("消息序号对不上，拒绝将信息持久化到本地消息池")
		} else if !p.RsaVerySignWithSha256(digestByte, c.Sign, MessageNodePubKey) {
			fmt.Println("节点签名验证失败！,拒绝将信息持久化到本地消息池")
		} else {
			p.setCommitConfirmMap(c.Digest, c.NodeID, true)
			count := 0
			for range p.commitConfirmCount[c.Digest] {
				count++
			}
			//如果节点至少收到了2f+1个commit消息（包括自己）,并且节点没有回复过,并且已进行过commit广播，则提交信息至本地消息池，并reply成功标志至客户端！
			p.lock.Lock()
			if count >= nodeCount/3*2 && !p.isReply[c.Digest] && p.isCommitBordcast[c.Digest] {
				fmt.Println("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
				//将消息信息，提交到本地消息池中！
				localMessagePool = append(localMessagePool, p.messagePool[c.Digest].Message)
				info := ""
				if p.node.nodeID != "N0" {
					info = p.node.nodeID + "节点已将msgid:" + strconv.Itoa(p.messagePool[c.Digest].ID) + "存入本地消息池中,消息内容为：" + p.messagePool[c.Digest].Content
				} else {
					info = "主节点已将msgid:" + strconv.Itoa(p.messagePool[c.Digest].ID) + "存入本地消息池中,消息内容为：" + p.messagePool[c.Digest].Content
				}
				//reply
				fmt.Println(info)
				fmt.Println("正在reply客户端 ...")
				tcpDial([]byte(info), p.messagePool[c.Digest].ClientAddr)
				p.isReply[c.Digest] = true
				fmt.Println("reply完毕")
			}
			p.lock.Unlock()
		}
	}
