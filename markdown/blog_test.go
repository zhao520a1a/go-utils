package markdown

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// 读取文件夹，在文件头部加"title"标签
func TestAddTitle(t *testing.T) {
	// 读取命令行参数
	//if len(os.Args) < 3 {
	//	fmt.Println("请输入两个文件夹路径")
	//	return
	//}
	//args := os.Args[1:]
	//if len(args) != 2 {
	//	fmt.Println("请输入两个文件夹路径")
	//	return
	//}
	//srcDir, destDir := args[0], args[1]
	srcDir := "/Users/Golden/Documents/Blog/source/_posts"
	destDir := "/Users/Golden/Desktop/tmp"
	// 遍历源文件夹下的Markdown文件
	err := filepath.Walk(srcDir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断是否为Markdown文件
		if strings.ToLower(filepath.Ext(path)) == ".md" {
			fmt.Println("正在处理文件：", path)
			// 读取文件内容
			contentBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(contentBytes)

			// 如果内容中没有“title”的关键词，则将文件名称以"title:xxx"拼接到内容的首行
			if !strings.Contains(content, "title:") {
				// 添加title字段到文件头部
				title := strings.TrimSpace(strings.TrimRight(file.Name(), ".md"))
				content = strings.Replace(content, "---", fmt.Sprintf("---\ntitle: %s", title), 1)
				// 将修改后的内容写入第二个文件夹路径下的同名Markdown文档中
				newFilePath := filepath.Join(destDir, file.Name())
				if err := writeStringToFile(newFilePath, content); err != nil {
					err = fmt.Errorf("写入文件失败：%v", err)
					return err
				}
				fmt.Println("处理完成：", file.Name())
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

// 将字符串写入文件中
func writeStringToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

// 转化文件中的图片标签
/*
用 Go 语言编码实现如下：正则表达式匹配 Markdown 语法中的图片标签 ![xxx](xxx)，将匹配到的内容替换成 <img src="xxx" alt="img" />
*/

func TestConvertImgTag(t *testing.T) {
	srcDir := "/Users/golden/Downloads/zhao520a1a.github.io-hexo/source/_posts/内部原理"
	destDir := "/Users/golden/Desktop/tmp"
	// 遍历源文件夹下的Markdown文件
	err := filepath.Walk(srcDir, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 判断是否为Markdown文件
		if strings.ToLower(filepath.Ext(path)) == ".md" {
			fmt.Println("正在处理文件：", path)
			// 读取文件内容
			contentBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content := ConvertMarkdownImg(string(contentBytes))
			newFilePath := filepath.Join(destDir, file.Name())
			if err := writeStringToFile(newFilePath, content); err != nil {
				err = fmt.Errorf("写入文件失败：%v", err)
				return err
			}
			fmt.Println("处理完成：", file.Name())
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestMatchMarkdownImg(t *testing.T) {
	text := "This is an ![example](https://example.com/image.png) image.  ![xxx](https://xxx.com/image.png) "
	result := ConvertMarkdownImg(text)
	fmt.Println(text)
	fmt.Println(result)
}

func ConvertMarkdownImg(text string) string {
	pattern := `!\[.*?\]\(.*\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllString(text, -1)
	result := text
	for _, match := range matches {
		prefix := strings.Split(match, "(")[0]
		temp := strings.TrimPrefix(match, prefix+"(")
		imgHref := strings.TrimSuffix(temp, ")")
		// 特殊处理
		imgHref = strings.ReplaceAll(imgHref, ".resources", "")
		imgStr := fmt.Sprintf(`<img src="%s" alt="img" />`, imgHref)
		result = strings.ReplaceAll(result, match, imgStr)
		fmt.Println(match)
		fmt.Println(imgStr)

	}
	return result
}
