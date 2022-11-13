package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

const Domain = "http://localhost:3333/"

type urlList struct {
	List      []string
	Listshort []string
}

// simple GET
func Get(surl string) string {
	client := &http.Client{}
	body := bytes.NewBuffer([]byte(surl))
	request, err := http.NewRequest("GET", Domain, body)
	if err != nil {
		log.Println(err)
	}
	// params := request.URL.Query()
	// params.Add("url", "")
	// params.Add("short_url", surl)
	//request.URL.RawQuery = params.Encode()
	// send
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	// response
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(result))
	return string(result)
}

// simple POST
func Post(url string) string {
	client := &http.Client{}
	body := bytes.NewBuffer([]byte(url))
	request, err := http.NewRequest("POST", Domain, body)
	if err != nil {
		log.Println(err)
	}
	// params := request.URL.Query()
	// params.Add("url", url)
	// params.Add("short_url", "")
	// request.URL.RawQuery = params.Encode()
	// send
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	// response
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(result))
	return string(result)
}

// POST as JSON
func PostJson(url string) {
	var std map[string]string = map[string]string{"url": url, "short_url": ""}
	data, err := json.Marshal(std)
	if err != nil {
		log.Println(err)
	}
	body := bytes.NewBuffer([]byte(data))
	req, err := http.NewRequest("POST", Domain, body)
	if err != nil {
		log.Println(err)
	}
	// header
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(result))
}

func main() {
	ok := "OK"
	// POST already exitst
	for i, lurl := range Lurls {
		if Post(lurl) != Surls[i] {
			fmt.Printf("error with lurl %s\n", lurl)
			ok = "NOT OK"
		}
	}
	// ==================
	// POST random JSON
	// lurl := ""
	// for i := 0; i < 3; i++ {
	// 	lurl = Randomer(20)
	// 	fmt.Println(Post(fmt.Sprintf("json_%d_%s", i, lurl)))
	// }
	//========

	// POST random
	// lurl := ""
	// for i := 0; i < 10; i++ {
	// 	lurl = Randomer(20)
	// 	fmt.Println(Post(fmt.Sprintf("lurl_%d_%s", i, lurl)))

	// }
	//=====

	// GET
	for i, surl := range Surls {
		if Get(surl) != Lurls[i] {
			fmt.Printf("error with surl %s\n", surl)
			ok = "NOT OK"
		}
	}
	if Get("FakeSurl") != "" {
		fmt.Printf("error with surl %s\n", "fakeUrl")
		ok = "NOT OK"
	}
	//====
	fmt.Println(ok)
}

func Randomer(n int) string {
	letters := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	if n <= 1 {
		n = 10
	}
	b := make([]byte, n)
	for i := 0; i < 10; i++ {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

var Surls = []string{
	"VNm4jAoXB3",
	"80fpMNcH0n",
	"RcL9SgeGKK",
	"4ZktNZWUrh",
	"AR_jjCuc5X",
	"VicglzquU2",
	"RcJwzK2XME",
	"0yOPKKLqsG",
	"cF3wglyj6d",
	"oAc215bVm5",
	"0qjNHedeRf",
	"hvOrBFGl1f",
	"ABU2durCAz",
	"TyJM5uDXcv",
	"yIS5GSWXvE",
	"mhoMjlmA9X",
	"7NzmA7D18X",
	"UkQjkfICNe",
	"wyY1aK2gUw",
	"kbbSwoU7jk",
	"x2GqZGruOE",
	"Y5dopYYyfZ",
	"1IVHAHgQAT",
}
var Lurls = []string{
	"surl_0_5p0_XqiFzz",
	"surl_1_9MyP_t8swo",
	"surl_2__traBocgTn",
	"surl_3_DRVOtj1hxM",
	"surl_4_QsJQGGV4jD",
	"surl_5_piWvfxI4l7",
	"surl_6_7DjDu4EUAm",
	"surl_7_spi5g8Ywn7",
	"surl_8_RFO2ckycvS",
	"surl_9_eMMuWmty0h",
	"json_0_5p0_XqiFzz",
	"json_1_9MyP_t8swo",
	"json_2__traBocgTn",
	"lurl_0_surl_0_5p0_XqiFzz",
	"lurl_1_surl_1_9MyP_t8swo",
	"lurl_2_surl_2__traBocgTn",
	"lurl_3_surl_3_DRVOtj1hxM",
	"lurl_4_surl_4_QsJQGGV4jD",
	"lurl_5_surl_5_piWvfxI4l7",
	"lurl_6_surl_6_7DjDu4EUAm",
	"lurl_7_surl_7_spi5g8Ywn7",
	"lurl_8_surl_8_RFO2ckycvS",
	"lurl_9_surl_9_eMMuWmty0h",
}
