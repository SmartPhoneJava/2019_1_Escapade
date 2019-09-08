package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

type VaultInitPayload struct {
	SecretShares    int `json:"secret_shares"`
	SecretThreshold int `json:"secret_threshold"`
}

type VaultErrors struct {
	Errors []string `json:"errors"`
}

type VaultInitResult struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base64"`
	RootToken  string   `json:"root_token"`
}

type VaultUnsealPayload struct {
	Key string `json:"key"`
}

type VaultUnsealResult struct {
	Sealed          bool   `json:"sealed"`
	SecretThreshold int    `json:"t"`
	SecretShares    int    `json:"n"`
	Progress        int    `json:"progress"`
	Version         string `json:"version"`
	ClusterName     string `json:"cluster_name"`
	ClusterID       string `json:"cluster_id"`
}

var KEYS_DIVIDE = " | "

func main() {
	var (
		consul    *consulapi.Client
		sessionID string

		publicPath  = flag.String("config", "config.json", "public config")
		privatePath = flag.String("secret", "secret.json", "private config")
		prefix      = flag.String("consul_prefix", "config/", "prefix for consul kv")

		consulAddr     = flag.String("consul_addr", "consul:8500", "address of consul server")
		consulCheckKey = flag.String("consul_check_key", "config/loaded",
			"the key that is used to check whether the data should be filled in(if key's value nil - yes)")

		vaultAddr            = flag.String("vault_addr", "vault:8200", "address of vault server")
		vaultSecretShares    = flag.Int("vault_secret_shares", 3, "amount of unsealed keys")
		vaultSecretThreshold = flag.Int("vault_secret_threshold", 2, "minimum of keys to get master key")
	)

	/*
		curl -s --request PUT -d '{"secret_shares": 3,"secret_threshold": 2}' http://127.0.0.1:8200/v1/sys/init

	*/

	log.Println("vaultInit")
	err := vaultInit(*vaultAddr, *vaultSecretShares, *vaultSecretThreshold)
	if err != nil {
		log.Fatal("Error", err.Error())
	}

	log.Println("vaultUnseal")
	err = vaultUnseal(*vaultAddr, *vaultSecretShares, *vaultSecretThreshold)
	if err != nil {
		log.Fatal("Error", err.Error())
	}

	vaultPush(*privatePath, *prefix, *vaultAddr)

	consul, sessionID = check(*consulAddr, *consulCheckKey)
	if consul == nil {
		return
	}
	consulPush(consul, *publicPath, *prefix, sessionID)

}

// vaultInit send init request to vault api
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
func vaultInit(vaultAddr string, shares, threshold int) error {

	var (
		result VaultInitResult
		vi     = VaultInitPayload{
			SecretShares:    shares,
			SecretThreshold: threshold,
		}
	)
	data, err := json.Marshal(vi)
	if err != nil {
		return err
	}
	body := bytes.NewReader(data)

	log.Println("body l", data)

	req, err := http.NewRequest("PUT", "http://"+vaultAddr+"/v1/sys/init", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var vaultErr VaultErrors
		err = json.Unmarshal(bodyBytes, &vaultErr)
		if err != nil {
			log.Println("Undefined answer:", string(bodyBytes))
		} else {
			log.Println("Vault return errors:", vaultErr)
		}
		return err
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return err
	}

	os.Setenv("VaultSecretShares", strconv.Itoa(shares))
	os.Setenv("VaultSecretThreshold", strconv.Itoa(threshold))
	os.Setenv("VaultRootToken", result.RootToken)
	os.Setenv("VaultUnsealKeys", result.Keys[0])
	os.Setenv("VaultUnsealKeysBase64", result.KeysBase64[0])

	keysLen := len(result.Keys)
	for i := 1; i < keysLen; i++ {
		os.ExpandEnv("$VaultUnsealKeys" + KEYS_DIVIDE + result.Keys[i])
		os.ExpandEnv("$VaultUnsealKeysBase64" + KEYS_DIVIDE + result.KeysBase64[i])
	}

	return nil
}

// vaultUnseal send unseal request to vault api
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
func vaultUnseal(vaultAddr string, threshold, shares int) error {

	var (
		keys      = os.Getenv("VaultUnsealKeys")
		keysSlice []string
	)
	if shares > 1 {
		keysSlice = strings.Split(keys, KEYS_DIVIDE)
	}

	for i, key := range keysSlice {
		if i == threshold {
			break
		}
		log.Println("key:", key)
		vi := VaultUnsealPayload{
			Key: key,
		}
		data, err := json.Marshal(vi)
		if err != nil {
			return err
		}
		log.Println("data:", data)
		body := bytes.NewReader(data)

		req, err := http.NewRequest("PUT", "http://"+vaultAddr+"/v1/sys/unseal", body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			var vaultErr VaultErrors
			err = json.Unmarshal(bodyBytes, &vaultErr)
			if err != nil {
				log.Println("Undefined answer:", string(bodyBytes))
			} else {
				log.Println("Vault return errors:", vaultErr)
			}
			return err
		}

		var result VaultUnsealResult
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return err
		}

	}

	return nil
}

// vaultPush send write secret request to vault api
// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
func vaultPush(privatePath, prefix, vaultAddr string) {
	data, err := ioutil.ReadFile(privatePath)
	if err != nil {
		log.Println("Read file error:", err.Error())
		return
	}

	body := bytes.NewReader(data)

	req, err := http.NewRequest("POST", "http://"+vaultAddr+"/v1/secret/"+prefix+"values", body)
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}
	req.Header.Set("X-Vault-Token", os.Getenv("VaultRootToken"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}
	defer resp.Body.Close()
	return
}

func consulPush(consul *consulapi.Client, publicPath, prefix, sessionID string) {
	var m map[string]interface{}

	data, err := ioutil.ReadFile(publicPath)
	if err != nil {
		log.Println("Read file error:", err.Error())
		return
	}

	if err = json.Unmarshal([]byte(data), &m); err != nil {
		log.Println("Unmarhal error:", err.Error())
		return
	}

	err = addDataToConsulKV(m, prefix, consul, sessionID)
	if err != nil {
		log.Println("Fail:", err.Error())
	} else {
		log.Println("Success")
	}
	return
}

func check(consulAddr, checkKey string) (*consulapi.Client, string) {
	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Println("Consul new client error:", err.Error())
		return nil, ""
	}

	sessionID, _, err := consul.Session().Create(nil, nil)
	if err != nil {
		log.Println("Consul session error:", err.Error())
		return nil, ""
	}

	kvp, _, err := consul.KV().Get(checkKey, nil)
	if err != nil {
		log.Println("Cant check", err.Error())
		return nil, ""
	} else if kvp != nil {
		log.Println("Check key is not null. It is:", string(kvp.Value))
		return nil, ""
	}
	return consul, sessionID
}

func putToConsulKV(key string, value interface{},
	consul *consulapi.Client, sessionID string) error {
	bytes := []byte(fmt.Sprintf("%v", value))

	_, _, err := consul.KV().Acquire(
		&consulapi.KVPair{
			Key:     key,
			Value:   bytes,
			Session: sessionID},
		nil)

	return err
}

func addDataToConsulKV(data map[string]interface{}, prefix string,
	consul *consulapi.Client, sessionID string) error {
	var err error
	for key, value := range data {
		if deeper, yes := value.(map[string]interface{}); yes {
			err = addDataToConsulKV(deeper, prefix+key+"/", consul, sessionID)
		} else {
			err = putToConsulKV(prefix+key, value, consul, sessionID)
		}
		if err != nil {
			break
		}
	}
	return err
}