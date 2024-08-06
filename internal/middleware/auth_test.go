package middleware

import (
	"encoding/hex"
	"fmt"
	"github.com/gophero/goal/ciphers"
	"github.com/gophero/goal/conv"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/testx"
	"github.com/gophero/goal/uuid"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

var userId int64 = 1001
var reqParam = "p1=hello&p2=world&p3=11"

func mockToken(userId uint) string {
	prefix := ciphers.MD5(conv.Int64ToStr(int64(userId)))
	middle := ciphers.MD5(uuid.UUID32())
	suffix := ciphers.MD5(conv.Int64ToStr(time.Now().UnixMilli())) // 时间戳
	return prefix + "-" + middle + "-" + suffix
}

func TestCreateToken(t *testing.T) {
	tk := mockToken(0)
	fmt.Println(tk)
}

func TestRequest(t *testing.T) {
	tl := testx.Wrap(t)
	tl.Case("mock request encrypting")

	token := mockToken(uint(userId))
	tl.Log("token: ", token)
	key := requestKey(token)
	iv := []byte(key) // 随机iv，每次的密文不同
	data, err := ciphers.AES.Encrypt([]byte(reqParam), []byte(key), ciphers.CBC, iv)
	encstr := fmt.Sprintf("%x\n", data)
	tl.Log("encrypt string: ", encstr)
	tl.Require(err == nil, "encrypt should have no error")

	key = requestKey(token)
	iv = []byte(key) // 随机iv，每次的密文不同
	encdata, err := hex.DecodeString(encstr)
	data, err = ciphers.AES.Decrypt(encdata, []byte(key), ciphers.CBC, iv)
	tl.Log("encrypt string:", string(data))
	tl.Require(err == nil, "decrypt should have not error")
	tl.Require(string(data) == reqParam, "decrypt result and request params should be match")
}

func TestEnc(t *testing.T) {
	s := "{\"id\":1,\"uuid\":1680158346,\"url\":\"https://rr1---sn-oguesn6k.googlevideo.com/videoplayback?expire=1683167417&ei=WcRSZMScEZ6L1d8Pw8Ox4AY&ip=193.38.139.241&id=o-AJOCIAuuHj9gShPCMydr8iIw_g5MtMdK4WAxAzk6bnd2&itag=22&source=youtube&requiressl=yes&mh=-5&mm=31%2C29&mn=sn-oguesn6k%2Csn-oguelnzz&ms=au%2Crdu&mv=m&mvi=1&pl=24&initcwndbps=1201250&spc=qEK7Bx-CwqJ4JWOMoFDmTV_yB7e9PY_CaJJD1D1OEA&vprv=1&mime=video%2Fmp4&ns=CuA0OFtJqdx8GhS7S71cIqkN&cnr=14&ratebypass=yes&dur=33.204&lmt=1681468054707870&mt=1683145516&fvip=1&fexp=24007246&c=WEB&txp=5432434&n=GpVL3pVn0wexxA&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cspc%2Cvprv%2Cmime%2Cns%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AOq0QJ8wRQIhALhL9S3hDxCPBiLYX6bX3wc2cmlxpxELRKYFKBtPWSk2AiAJr-PqcL5OzBwzS-sGeS2R8r93cvasxwOZ79dpyIbrbQ%3D%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRAIgXicHNLxlopZSD_Vy-fRjoCQwTRoIXoTTnjCiEnH5-FQCICcZ-nbl131xEe5bZHUXBLz4Q7Oow3mR5JtBe6Ls7R3-\"}"
	token := "c4ca4238a0b923820dcc509a6f75849b-75cd8f9ed159f30143eab39564442755"

	key := requestKey(token)
	fmt.Println(key)
	iv := []byte(key)                                                                       // 随机iv，每次的密文不同
	rawBytes := errorx.Throwv(ciphers.AES.Encrypt([]byte(s), []byte(key), ciphers.CBC, iv)) // aes解密
	encstr := fmt.Sprintf("%x\n", rawBytes)
	fmt.Println("encrypt string: ", encstr)
}

func TestDec(t *testing.T) {
	s := "afc846a469794230e523c5a0660c0a2f3309b14451a2b80cf834153eca32828ae916fa91fa0729e4ad2c18f1a4a87999318c04ee54e20fc9017c9db4937e060bfa10c6e5a29d795f2607e15c54ab9eb34e9294437a0ea895d42be21e785afce104f6555ef51e7e88921ba789fecbbeb723859feee621e7483922b2b65cde44832dd28eace72b6f9066f6ab218c9504385a80978a35eb0c84865bad3c5734814c8457629d52121e1cfb9193241c7c9628d1f889d64d05dafb07b61bd6bdf3011ee2f461fe601b5b2d0d986bed015d534b42c9bac747547c93ee6816455b916ef861fb672775f243c30f04ad89637ba0df5ad8394f51c86a121f24c9b2c18a1e8dbbc8f98fcf1a3e041a8e866154ff39e906677a52d8477ead1fe432d54498d7e7eacfe984e528c778c994954eb319b8da637ef46c829195162f9eaaaaafe77273fc7590e692c6ead56f840ed59e34a6d1291c1636ce10bf70f7cebb14fde724acca0721c471db2502d0b2c7a18a7a66fc72c5c7344e9d9ee9c8f5eea477fba9f2fb93c77175ac047ac3dbf1af3049fa20d229c242f2d59ac9cb3dded5bcbbe05da11a689b320bdf81dd1897ed319875a336628aaaae9568f22cc733ea3d4d0178dfdce909fe569f76e8d79dd6927b5b53b721b6789514bc4e9e9fb077e4dab4ff744593bf446d0294f6b482d5da3a0fc52de0060008f36a9aace8059d1091126df778f84859d0e711459e9321ca1c6242d38a550178dbfa140c87cf5e0952a137b7ae893d84b387ea9776b1a07d6b610dd37bc838015225c5f57c145c5df8902585b60148dd425b90fa797c25d41b2492cbbf87ee4864591d2c952f5f4db363cc5c15adc23cb51f5dfdf39df273429cdc16964423caf81712e8af37b3425aa0c57d72bb6beb5e96e14fae155921e873b3e41a75c5d5e09deee550da7b7c235438674e21ed98d53cd7c57d54763105c09c4fe31053b83a38d8b5059fdf40f7d54f6774b88739b4e7660235a37666a65debc52861183aed43ace8e50d75cb1bc120c72317214b0bfab94c6d5c1d6b6734f491d264efa666dd2e3e29a12b58a13b5c31612ba0fa551d1d6a331cabd9ce9edaf9afef79e4a5d7236114f4f48fed779f28da64ea776b0b2a80f476c069150f15eed9280bcdda70533288eeb9dc5737cd0b518385cf20c69f591d46ae96a483f5b547b257adfb00d9ee8aad768dabac94699a3e11ebc0482719344f6933e3313ac35a06fecb62e5d5ce4149430cac0996cd3f67fc15651d33b1108059cdfd5daf4d87b43bcb605f128104eab07fd4eeed813a3a87e3abab16742e27423dae6bfccf3161bc69b7210b4106c5d6ef68c275"
	token := "c4ca4238a0b923820dcc509a6f75849b-75cd8f9ed159f30143eab39564442755"

	key := requestKey(token)
	fmt.Println(key)
	iv := []byte(key)                                                                      // 随机iv，每次的密文不同
	encBytes := errorx.Throwv(hex.DecodeString(s))                                         // 十六进制解码
	rawBytes := errorx.Throwv(ciphers.AES.Decrypt(encBytes, []byte(key), ciphers.CBC, iv)) // aes解密
	fmt.Println("decode string: ", string(rawBytes))
}

// 测试重复请求防止 SecCoin
func TestRepeatLogin(t *testing.T) {
	url := "http://127.0.0.1:8080/login/sign_in"
	token := "c4ca4238a0b923820dcc509a6f75849b-e052676e596a4662918a51e1efbe4c89"
	param := "{\"email\":\"t1@gmail.com\",\"password\":\"112233\",\"uuid\":\"123\"}"

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			key := requestKey(token)
			iv := []byte(key)
			fmt.Println(key)
			data, err := ciphers.AES.Encrypt([]byte(param), []byte(key), ciphers.CBC, iv)
			encstr := fmt.Sprintf("%x", data)
			fmt.Println("encrypt string: ", encstr)
			r := strings.NewReader("a=" + encstr)
			req, err := http.NewRequest(http.MethodPost, url, r)
			if err != nil {
				fmt.Printf("err: %v\n", err)
			} else {
				header := http.Header{}
				header.Add("Content-Type", "application/x-www-form-urlencoded")
				header.Add("X-TOKEN", token)
				header.Add("Accept-Language", "en")
				req.Header = header
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Printf("err: %v\n", err)
				} else {
					fmt.Println(resp.Status)
				}
			}
		}()
	}
	wg.Wait()
}
