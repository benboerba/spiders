//------------------------知乎爬虫工具使用说明------------------------------------
//      1.更新cookie
//      2.准备待下载文件夹
//      3.等
//      author by nxb
//------------------------------------------------------------------------------

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	formula    string = `https://[^"{]+\.jpg`
	cookie     string = `_zap=8f3061cf-1d28-4f66-ac6e-5bc9b8ec2fc3; _xsrf=C10wcdvVeQa9tEd1DLJHxzsgWZiuMzGM; d_c0="AEDiBEcrvQ-PTjkYe2trS7dcTNmL2x0Ky0w=|1563172957"; capsion_ticket="2|1:0|10:1563172970|14:capsion_ticket|44:ZTQ0ZjQ5YzkxNmJiNGNhZWJjNzExMGYxNTdmY2QzNzI=|c9f4c633ec22ae7328495bb76860f53a012b700e73fc22e0540709c374703b0b"; z_c0="2|1:0|10:1563173003|4:z_c0|92:Mi4xQjY1bUJBQUFBQUFBUU9JRVJ5dTlEeVlBQUFCZ0FsVk5pMjRaWGdBQzNwSnRMSVFMM29OUXl1T1UwOHRxSHl1aUN3|3dda9b7b871a7c8c0edfc159f5e50a1017102c2b02cc5aadafa1229a43beaa4b"; tst=r; q_c1=095bc3f3f18a429baec8651cb3d78b30|1563173071000|1563173071000; tgw_l7_route=060f637cd101836814f6c53316f73463`
	prefix     string = "https://www.zhihu.com/api/v4/questions/"
	postfix    string = "/answers?include=data%5B%2A%5D.is_normal%2Cadmin_closed_comment%2Creward_info%2Cis_collapsed%2Cannotation_action%2Cannotation_detail%2Ccollapse_reason%2Cis_sticky%2Ccollapsed_by%2Csuggest_edit%2Ccomment_count%2Ccan_comment%2Ccontent%2Ceditable_content%2Cvoteup_count%2Creshipment_settings%2Ccomment_permission%2Ccreated_time%2Cupdated_time%2Creview_info%2Crelevant_info%2Cquestion%2Cexcerpt%2Crelationship.is_authorized%2Cis_author%2Cvoting%2Cis_thanked%2Cis_nothelp%2Cis_labeled%2Cis_recognized%2Cpaid_info%3Bdata%5B%2A%5D.mark_infos%5B%2A%5D.url%3Bdata%5B%2A%5D.author.follower_count%2Cbadge%5B%2A%5D.topics&limit=20&offset="
	ppos       string = "&sort_by=default"
	myfilepath string = `E:\Resources\photo\`
	//keyword    string = `_hd`
)

func main() {
	if err := InitMgo(); err != nil {
		fmt.Println("init mgo ", err)
		return
	}

	fmt.Println("init finish")

	sess := GetMgoS()
	defer sess.Close()

	c := GetMgoC(sess, "medex", "student")

	var data interface{}
	err := c.Find(bson.M{"name": "hhh"}).One(&data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GetDataViaName() student ", data)

	// t := time.Now()
	// //name := "zhihuurl.txt"
	// offset := 0
	// questionid := "26620889"
	// content := geturlRespHtml(prefix + questionid + postfix + strconv.Itoa(offset) + ppos)
	// tmpindex := strings.LastIndex(content, "totals") //totals出现的下标
	// //total处理      开始
	// tmp := []byte{}
	// tmpindex += 8
	// for content[tmpindex] >= 48 && content[tmpindex] <= 57 {
	// 	tmp = append(tmp, content[tmpindex])
	// 	tmpindex++
	// }
	// total, err := strconv.Atoi(string(tmp))
	// if err != nil {
	// 	return
	// }
	// //total处理      结束

	// pagenum := total/20 + 1
	// ch := make(chan int, pagenum)

	// fmt.Println(total)

	// fmt.Println("pagenum : ", pagenum)
	// WriteWithIoutil("1.txt", content)

	// for offset < total {
	// 	go downloadgoroutine(prefix+questionid+postfix+strconv.Itoa(offset)+ppos, ch)

	// 	// content = geturlRespHtml(prefix + questionid + postfix + strconv.Itoa(offset) + ppos)
	// 	// str += analyze(content)
	// 	offset += 20
	// 	fmt.Println("offset : ", offset)
	// }

	// for j := 0; j < pagenum; j++ {
	// 	tmp := <-ch
	// 	fmt.Println(tmp)
	// }

	// elapsed := time.Since(t)
	// fmt.Println("program elapsed:", elapsed)
}

func downloadgoroutine(url string, ch chan int) {
	content := geturlRespHtml(url)
	analyze(content)
	ch <- 0
}

func geturlRespHtml(url string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("addr error")
	}
	req.Header.Set("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("login error")
	}
	resp_byte, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	respHtml := string(resp_byte)
	return respHtml
}

func WriteWithIoutil(name, content string) {
	data := []byte(content)
	if ioutil.WriteFile(name, data, 0644) == nil {
	}
}

func downloadpicture(imgurl string) {
	imgpath := myfilepath
	res, err := http.Get(imgurl)
	if err != nil {
		fmt.Println("get error")
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	//reader := bufio.NewReaderSize(res.Body, 512*1024)
	FileName := filepath.Base(imgurl) //创建一个文件名

	file, err := os.Create(imgpath + FileName)
	if err != nil {
		fmt.Println("create error")
	}
	defer file.Close()
	// 获得文件的writer对象
	//writer := bufio.NewWriter(file)

	io.Copy(file, res.Body)
	//fmt.Printf("Total length: %d", written)
}

func analyze(str string) string {
	rp1 := regexp.MustCompile(formula)
	heads := rp1.FindAllStringSubmatch(str, -1)
	var strtmp string
	index := 0
	for _, value := range heads {
		for _, val := range value {
			if !strings.Contains(val, "<") && !strings.Contains(val, "_hd") && !strings.Contains(val, "_is") {
				if index%1000 == 0 {
					time.Sleep(time.Duration(1) * time.Second)
				}
				index++
				downloadpicture(val)
				strtmp += val
				strtmp += "\n"
			}
		}
	}
	return strtmp
}

func analyzenodownload(str string) string {
	rp1 := regexp.MustCompile(formula)
	heads := rp1.FindAllStringSubmatch(str, -1)
	var strtmp string
	index := 0
	for _, value := range heads {
		for _, val := range value {
			if !strings.Contains(val, "<") && !strings.Contains(val, "_hd") && !strings.Contains(val, "_is") {
				if index%50 == 0 {
					time.Sleep(time.Duration(2))
				}
				index++
				strtmp += val
				strtmp += "\n"
			}

		}
	}
	return strtmp
}

func deleterepetition(filepath, keyword string) {
	rd, err := ioutil.ReadDir(filepath)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		} else {
			if strings.Contains(fi.Name(), keyword) {
				filename := filepath + "\\" + fi.Name()
				err := os.Remove(filename)
				if err != nil {
					fmt.Println("remove error info:", filename)
				}
			}
		}
	}
}
