package blockchaincommon

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserRegit_Message struct {
	Sha256Value []byte `json:"sha256value"`
	Appid       []byte `json:"appid"`
	Time        []byte `json:"emit"`
	Token       []byte `json:"token"`
	Data        []byte `json:"data"`
}
type AdminCreateNFT_Message struct {
	Sha256Value []byte `json:"sha256value"`
	Appid       []byte `json:"appid"`
	Time        []byte `json:"emit"`
	Token       []byte `json:"token"`
	Nonce       []byte `json:"nonce"`
	LifeTime    []byte `json:"lifeTime"`
	Password    []byte `json:"password"`
	From        []byte `json:"From"`
	To          []byte `json:"To"`
	ChainType   []byte `json:"chainType"`
}
type AdminCreateNFTBatch_Message struct {
	Sha256Value []byte   `json:"sha256value"`
	Appid       []byte   `json:"appid"`
	Time        []byte   `json:"emit"`
	Token       []byte   `json:"token"`
	Nonce       []byte   `json:"nonce"`
	LifeTime    []byte   `json:"lifeTime"`
	Password    []byte   `json:"password"`
	From        []byte   `json:"From"`
	Tos         []string `json:"Tos"`
	ChainType   []byte   `json:"chainType"`
}
type AdminTransferNFTBatch_Message struct {
	Sha256Value []byte   `json:"sha256value"`
	Appid       []byte   `json:"appid"`
	Time        []byte   `json:"emit"`
	Token       []byte   `json:"token"`
	Nonce       []byte   `json:"nonce"`
	LifeTime    []byte   `json:"lifeTime"`
	Password    []byte   `json:"password"`
	From        []byte   `json:"From"`
	Tos         []string `json:"Tos"`
	Ids         []string `json:"ids"`
	ChainType   []byte   `json:"chainType"`
}
type TransferFrom_Message struct {
	Sha256Value []byte `json:"sha256value"`
	Appid       []byte `json:"appid"`
	Time        []byte `json:"emit"`
	Token       []byte `json:"token"`
	Nonce       []byte `json:"nonce"`
	LifeTime    []byte `json:"lifeTime"`
	Password    []byte `json:"password"`
	From        []byte `json:"From"`
	To          string `json:"to"`
	Id          string `json:"id"`
	ChainType   []byte `json:"chainType"`
}
type FreeGasMint_Message struct {
	Sha256Value []byte `json:"sha256value"`
	Appid       []byte `json:"appid"`
	Time        []byte `json:"emit"`
	Token       []byte `json:"token"`
	Nonce       []byte `json:"nonce"`
	LifeTime    []byte `json:"lifeTime"`
	Password    []byte `json:"password"`
	From        []byte `json:"From"`
	ChainType   []byte `json:"chainType"`
}
type UserNFTs_Message struct {
	Sha256Value []byte `json:"sha256value"`
	Appid       []byte `json:"appid"`
	Time        []byte `json:"emit"`
	Token       []byte `json:"token"`
	From        []byte `json:"From"`
	ChainType   []byte `json:"chainType"`
}

type UserRegitRes_Message struct {
	Confluxaddress string `json:"ConfluxAddress"`
	ETHaddress     string `json:"ETHAddress"`
}

var body []byte

const TestAPPID_CFXName string = "0xd67c8aed16df25b21055993449222fa895c67eb87bb1d7130c38cc469d8625b5" //测试专用 APPID就是对应项目ID，不同项目会有不同合约,APPID也会不同
const TestAPPID_ETHName string = "0xd8350f2533aa38e5ed0b99b1b4af1a134ba9854bbcf66435239b9c18f60276d7" //测试专用 APPID就是对应项目ID，不同项目会有不同合约,APPID也会不同
//密钥托管服务器Post请求url
const TestIPandPort string = "https://1.116.87.151:13149"                                       //测试专用
const TestAdministratorPassword string = "dx123456"                                             //测试专用 管理员密码，cfx和eth都是这个
const TestCFXAdministratorAddress string = "cfxtest:aakmdj7tutgdy3h558rr5621mhrrx75kfyw3e3sfz0" //测试专用
const TestETHAdministratorAddress string = "0xfec36af44b8be6AB6ba97aF3b71940D3f3B8B539"         //测试专用
var publickey []byte

/**
 * @name:Reg
 * @test: test font
 * @msg:用户注册Verse密钥托管系统账户
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} APPID 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {string} RegPassword 用户注册时填写的二级支付密码：建议大小写+数字0-9
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{Reg}
 * @return {*}conflux地址，ETH地址
 */
func Reg(IPandPort string, APPID string, RegPassword string, flag string) (string, string, error) {
	body, err := regitPost(IPandPort, "UserRegit", APPID, RegPassword, flag)
	if err != nil {
		return string(body), string(body), err
	}
	fmt.Println(string(body))
	res := &UserRegitRes_Message{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err.Error(), err.Error(), err
	}
	fmt.Println("Confluxaddress:", res.Confluxaddress)
	fmt.Println("ETHaddress:", res.ETHaddress)
	return res.Confluxaddress, res.ETHaddress, nil
}

/**
 * @name:InitRSAPuk
 * @test: test font
 * @msg:初始化与密钥系统的加密通信RSA2018公钥
 * @param {string} filename 加密通信公钥文件路径 public.pem  也可以自己重命名名称
 * @return {*}
 */
func InitRSAPuk(filename string) error {
	//1. 读取公钥信息 放到data变量中
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	stat, _ := file.Stat() //得到文件属性信息
	data := make([]byte, stat.Size())
	file.Read(data)
	file.Close()
	publickey = data
	return nil
}
func regitPost(IPandPort string, actionName string, myappid string, Password string, flag string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	//post请求提交json数据
	messages := UserRegit_Message{sha256Value, src_appid, src_mytime, src_token, src_Password}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte(""), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return body, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}

/////////////////////////////////////CONFLUX AND ETH//////////////////////////////////////////
/**
 * @name:TotalSupplyPost
 * @test: 本函数为【查询类函数】，只读取区块链信息，不写入和改变区块链信息
 * @msg:获取指定区块链，指定NFT合约的NFT总发行数量
 * @param {string} IPandPort  密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_TotalSupply代表conflux区块链上的合约，ETH_TotalSupply代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{TotalSupplyPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*}返回值[]byte统一string()处理即可，本函数返回值为一个big.Int的字符串表示
 */
func TotalSupplyPost(IPandPort string, actionName string, myappid string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	// src_Password := publicEncode([]byte(Password), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := UserRegit_Message{sha256Value, src_appid, src_mytime, src_token, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:UserNFTsPost
 * @test: 本函数为【查询类函数】，只读取区块链信息，不写入和改变区块链信息
 * @msg:查询指定用户区块链地址的所有NFT持有情况
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_UserNFTs代表conflux区块链上的合约，ETH_UserNFTs代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {string} From 被查询用户的区块链地址
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{UserNFTsPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*}返回NFT的id,例如：1,2,3,6,33,666,7543   以逗号隔开，请使用字符串split即可c拆分为数组
 */
func UserNFTsPost(IPandPort string, actionName string, myappid string, From string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)

	src_From := publicEncode([]byte(From), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := UserNFTs_Message{sha256Value, src_appid, src_mytime, src_token, src_From, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:AdminCreateNFTPost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。通常生成批量NFT请调用函数AdminCreateNFTBatchPost()
 * @msg:合约管理员创建单个NFT
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_AdminCreateNFT代表conflux区块链上的合约，ETH_AdminCreateNFT代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password 合约管理员的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From 合约管理员的区块链地址
 * @param {string} to NFT创建出来拥有者的地址，这个地址可以是管理员地址，也可以是其他用户的地址。例如：统一分发型的NFT就是先给管理员创建NFT，最后再转移给用户，那么这里传入的就是管理员地址。
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{AdminCreateNFTPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func AdminCreateNFTPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, to string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_To := publicEncode([]byte(to), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := AdminCreateNFT_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, src_To, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:AdminCreateNFTBatchPost
 * @test:本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。本函数采用Rollup模式打包交易，可以大幅加快区块链确认时间，减少Gas
 * @msg:合约管理员创建批量NFT，数量必须小于240.建议【200为佳】，因为区块链一个区块大小300K，如果超过240可能会交易失败
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_AdminCreateNFTBatch代表conflux区块链上的合约，ETH_AdminCreateNFTBatch代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password 合约管理员的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From 合约管理员的区块链地址
 * @param {[]string} tos NFT创建出来拥有者的地址，【是一个数组】，这些地址可以都是管理员地址，也可以是其他用户的地址。例如：统一分发型的NFT就是先给管理员创建NFT，最后再转移给用户，那么这里数组传入的都是管理员地址。
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{AdminCreateNFTBatchPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func AdminCreateNFTBatchPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, tos []string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_Tos := tos
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := AdminCreateNFTBatch_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, src_Tos, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:AdminTransferNFTBatchPost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。本函数采用Rollup模式打包交易，可以大幅加快区块链确认时间，减少Gas
 * @msg:管理员批量转移NFT至指定地址，地址数量和ID数量必须相等，且都必须小于240.建议【200为佳】，因为区块链一个区块大小300K，如果超过240可能会交易失败
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_AdminTransferNFTBatch代表conflux区块链上的合约，ETH_AdminTransferNFTBatch代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password 合约管理员的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From  合约管理员的区块链地址
 * @param {[]string} tos 转移的目的地址,也就是用户地址，与ids数组一一对应
 * @param {[]string} ids 转移的ID，与tos数组一一对应
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{AdminTransferNFTBatchPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func AdminTransferNFTBatchPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, tos []string, ids []string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_Tos := tos
	src_ids := ids
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := AdminTransferNFTBatch_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, src_Tos, src_ids, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:TransferFromPost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。
 * @msg:NFT的通用转移函数，本函数实现NFTID为id的NFT从From地址转移至to地址。From地址必须是id的NFT拥有者
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_TransferFrom代表conflux区块链上的合约，ETH_TransferFrom代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password From地址用户的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From  NFT拥有者地址，转移源地址
 * @param {string} to  转移NFT的目的地址
 * @param {string} id 转移NFT的ID
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{TransferFromPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func TransferFromPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, to string, id string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := TransferFrom_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, to, id, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:BurnPost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。
 * @msg:From地址用户销毁属于自己NFTID为id的NFT，销毁后NFT打入零地址，即黑洞，【再也无法找回】
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_TransferFrom代表conflux区块链上的合约，ETH_TransferFrom代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password From地址用户的的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From  NFT拥有者地址
 * @param {string} id 需要销毁的NFT的ID，该ID的NFT必须属于From地址
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{BurnPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func BurnPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, id string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := TransferFrom_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, "", id, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:ApprovePost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。
 * @msg:授权From地址内的指定id的NFT给to合约。to通常是NFTmarket即NFT市场合约地址，目前业务【暂时不会用到】
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_TransferFrom代表conflux区块链上的合约，ETH_TransferFrom代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password From地址用户的的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From  NFT拥有者地址
 * @param {string} to NFTmarket即NFT市场合约地址
 * @param {string} id 需要授权的NFT的ID，该ID的NFT必须属于From地址
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{BurnPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func ApprovePost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, to string, id string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := TransferFrom_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, to, id, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:BurnBatchPost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。本函数采用Rollup模式打包交易，可以大幅加快区块链确认时间，减少Gas
 * @msg:From地址用户批量销毁属于自己NFTIDs为id的NFT，销毁后对应的NFT打入零地址，即黑洞，【再也无法找回】
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_TransferFrom代表conflux区块链上的合约，ETH_TransferFrom代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password From地址用户的的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From  NFT拥有者地址
 * @param {[]string} ids NFTID数组
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{BurnPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63     可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func BurnBatchPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, ids []string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	var src_Tos []string
	src_ids := ids
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := AdminTransferNFTBatch_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, src_Tos, src_ids, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

/**
 * @name:UserFreeMintPost
 * @test: 本函数为【写入类函数】谨慎调用，会参数区块链交易信息并通过密钥系统签名，会改变区块链信息。
 * @msg:用户自主创建NFT【目前暂不开放该功能】
 * @param {string} IPandPort 密钥系统请求链接 例如 https://127.0.0.1:13149
 * @param {string} actionName 请求名称，同样名称的含义也是请求的功能，注意区分ETH和Conflux。本参数传入：CFX_FreeGasMint代表conflux区块链上的合约，ETH_FreeGasMint代表以太坊及其侧链、L2的合约
 * @param {string} myappid 项目认证的APPID 例如 0xd67c9aed16df25b21055993449229fa895c67eb87bb1d7130c38cc469d8625b5
 * @param {int64} Nonce 随机数，该参数用于同一地址的并发区块链交易执行，如果管理员执行一次或者串行执行N次，传参-1  如果需要并发连续执行20次，那么需要一次传入0-19，作为区块链随机数以避免交易冲突
 * @param {uint64} LifeTime 私钥的生命周期，单位毫秒，设置该参数可以加快处理速度，减少区块链密钥系统解密重复运算，如果执行一次建议值2000，如果并发执行N次，建议1000*N
 * @param {string} Password From地址用户的的区块链密钥解密密码，也可以理解为支付密码
 * @param {string} From  NFT创建者地址
 * @param {string} flag 标记，用于同一地址区块链并发交易使用，通常就填写本函数名称{BurnPost}
 * @param {string} ChainType 区块链类型，参数：cfx代表conflux  eth代表以太坊  bsc代表币安链  arb代表以太坊L2 Arbitrum，注意全部为小写字母哦
 * @return {*} 返回值为交易hash代表成功。例如：0xcc07051ca530dbb1982b25438ca1a0d5c874a3c4c104256b7d7981e78bb02e63    可以通过判断err!=nil；其他返回他信息错误原因err.Error()在[]byte内
 */
func UserFreeMintPost(IPandPort string, actionName string, myappid string, Nonce int64, LifeTime uint64, Password string, From string, flag string, ChainType string) ([]byte, error) {
	now := uint64(time.Now().Unix())    //获取当前时间
	by := make([]byte, 8)               //建立数组
	binary.BigEndian.PutUint64(by, now) //uint64转数组
	//加密数据
	sha256Value := []byte(CalculateHashcode(myappid)) //APPID的sha256
	src_appid := publicEncode([]byte(myappid), publickey)
	src_mytime := publicEncode([]byte(by), publickey)
	src_token := publicEncode([]byte(fmt.Sprint(time.Now().UnixNano())+myappid+flag), publickey)
	binary.BigEndian.PutUint64(by, uint64(Nonce)) //uint64转数组
	src_Nonce := publicEncode([]byte(by), publickey)
	binary.BigEndian.PutUint64(by, uint64(LifeTime)) //uint64转数组
	src_LifeTime := publicEncode([]byte(by), publickey)
	src_Password := publicEncode([]byte(Password), publickey)
	src_From := publicEncode([]byte(From), publickey)
	src_ChainType := publicEncode([]byte(ChainType), publickey)
	//post请求提交json数据
	messages := FreeGasMint_Message{sha256Value, src_appid, src_mytime, src_token, src_Nonce, src_LifeTime, src_Password, src_From, src_ChainType}
	ba, err := json.Marshal(messages)
	if err != nil {
		return []byte("json.Marshal error"), err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Post(IPandPort+"/"+actionName+"", "application/json", bytes.NewBuffer([]byte(ba)))
	if err != nil {
		body, err := ioutil.ReadAll(resp.Body)
		return []byte("http error:" + fmt.Sprint(err) + "internel:" + string(body)), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte("ReadAll error"), err
	}
	return body, nil
}

//使用rsa公钥加密文件
func publicEncode(plainText []byte, data []byte) []byte {
	//1. 读取公钥信息 放到data变量中
	//2. 将得到的字符串pem解码
	//1. 读取公钥信息 放到data变量中

	//2. 将得到的字符串pem解码
	block, _ := pem.Decode(data)
	//3. 使用x509将编码之后的公钥解析出来
	pubInterface, err2 := x509.ParsePKIXPublicKey(block.Bytes)
	if err2 != nil {
		panic(err2)
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	//4. 使用公钥加密
	cipherText, err3 := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainText)
	if err3 != nil {
		panic(err3)
	}
	return cipherText
}

//使用rsa公钥加密文件
func publicEncodeLong(plainText []byte, data []byte) ([]byte, error) {

	//2. 将得到的字符串pem解码
	block, _ := pem.Decode(data)

	//3. 使用x509将编码之后的公钥解析出来
	pubInterface, err2 := x509.ParsePKIXPublicKey(block.Bytes)
	if err2 != nil {
		panic(err2)
	}
	pubKey := pubInterface.(*rsa.PublicKey)
	partLen := pubKey.N.BitLen()/8 - 11
	chunks := split(plainText, partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

// 、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、、
func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:])
	}
	return chunks
}

//sha256运算
/**
 * @name: CalculateHashcode
 * @test: test font
 * @msg: 计算字符串的sh256值即hash值
 * @param {string} data 待计算的字符串
 * @return {*} 返回hash
 */
func CalculateHashcode(data string) string {
	nonce := 0
	var str string
	var check string
	pass := false
	var dif int = 4
	for nonce = 0; ; nonce++ {
		str = ""
		check = ""
		check = data + strconv.Itoa(nonce)
		h := sha256.New()
		h.Write([]byte(check))
		hashed := h.Sum(nil)
		str = hex.EncodeToString(hashed)
		for i := 0; i < dif; i++ {
			if str[i] != '0' {
				break
			}
			if i == dif-1 {
				pass = true
			}
		}
		if pass == true {
			return str
		}
	}
}

//跨域
