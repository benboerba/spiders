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
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	formula    string = `https://[^"{]+\.jpg`
	cookie     string = `_xsrf=pCztpO8EN6FSWNC7ih4f444fxSVvkITd; _zap=ef68d577-531e-4c5d-b38e-ddb7eaf49dc1; d_c0="AKAj7XI2zg-PTuQsnObF_rvWLm3o2AyNtqY=|1564316737"; capsion_ticket="2|1:0|10:1564316750|14:capsion_ticket|44:NGRlZjRjNzNlNDliNDRlMGJlYjg2Y2UxMmQ4MWY4NDk=|f64af661b96d24449ec403d7ef22745564a92c1be9fea7fdd05818ccbf154344"; z_c0="2|1:0|10:1564316784|4:z_c0|92:Mi4xQjY1bUJBQUFBQUFBb0NQdGNqYk9EeVlBQUFCZ0FsVk5jT0lxWGdCVGRQMWxDRXRZMlNPWjJwMDd2XzhYS1F5UnNB|3896b5f1763304a98c2001ae5d449cefb84d5d953db0e5b3a43779ddabf8075c"; tst=r; q_c1=2619259abaad49f897c5d4e7f4e5fb26|1564326458000|1564326458000; tgw_l7_route=f2979fdd289e2265b2f12e4f4a478330`
	prefix     string = "https://www.zhihu.com/api/v4/questions/"
	postfix    string = "/answers?include=data%5B%2A%5D.is_normal%2Cadmin_closed_comment%2Creward_info%2Cis_collapsed%2Cannotation_action%2Cannotation_detail%2Ccollapse_reason%2Cis_sticky%2Ccollapsed_by%2Csuggest_edit%2Ccomment_count%2Ccan_comment%2Ccontent%2Ceditable_content%2Cvoteup_count%2Creshipment_settings%2Ccomment_permission%2Ccreated_time%2Cupdated_time%2Creview_info%2Crelevant_info%2Cquestion%2Cexcerpt%2Crelationship.is_authorized%2Cis_author%2Cvoting%2Cis_thanked%2Cis_nothelp%2Cis_labeled%2Cis_recognized%2Cpaid_info%3Bdata%5B%2A%5D.mark_infos%5B%2A%5D.url%3Bdata%5B%2A%5D.author.follower_count%2Cbadge%5B%2A%5D.topics&limit=20&offset="
	ppos       string = "&sort_by=default"
	myfilepath string = `C:\Users\Administrator\Desktop\picture100\`
	//keyword    string = `_hd`
)

func main() {
	t := time.Now()
	//name := "zhihuurl.txt"
	offset := 0
	questionid := "318927654"
	content := geturlRespHtml(prefix + questionid + postfix + strconv.Itoa(offset) + ppos)
	tmpindex := strings.LastIndex(content, "totals") //totals出现的下标
	//total处理      开始
	tmp := []byte{}
	tmpindex += 8
	for content[tmpindex] >= 48 && content[tmpindex] <= 57 {
		tmp = append(tmp, content[tmpindex])
		tmpindex++
	}
	total, err := strconv.Atoi(string(tmp))
	if err != nil {
		return
	}
	//total处理      结束

	//创建文件夹
	err = os.Mkdir(myfilepath, os.ModePerm)
	if err != nil {
		fmt.Printf("mkdir failed![%v]\n", err)
	} else {
		fmt.Printf("mkdir success!\n")
	}

	pagenum := total/20 + 1
	ch := make(chan int, pagenum)

	fmt.Println(total)

	fmt.Println("pagenum : ", pagenum)
	WriteWithIoutil("1.txt", content)

	for offset < total {
		go downloadgoroutine(prefix+questionid+postfix+strconv.Itoa(offset)+ppos, ch)

		// content = geturlRespHtml(prefix + questionid + postfix + strconv.Itoa(offset) + ppos)
		// str += analyze(content)
		offset += 20
		fmt.Println("offset : ", offset)
	}

	for j := 0; j < pagenum; j++ {
		tmp := <-ch
		fmt.Println(tmp)
	}

	elapsed := time.Since(t)
	fmt.Println("program elapsed:", elapsed)
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
