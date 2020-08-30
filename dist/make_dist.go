package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ver     = kingpin.Arg(`version`, `version to build`).Required().String()
	develop = kingpin.Flag(`dev`, `set development mode`).Bool()
	appName = `tinygo-autocmpl`
	test    = kingpin.Flag(`test`, `set test mode`).Bool()
	now     = time.Now()
)

var targets = []string{
	// ここに成果物にコピーしたいファイル、フォルダを記載する
	// 例)
	// `dist/make_dist.go`,
}

func main() {
	kingpin.Parse()

	version := *ver
	version = regexp.MustCompile(`^v`).ReplaceAllString(version, ``)
	versionInfo := version

	if *test {
		doTests()
		if !checkChanges(version) {
			fmt.Fprintf(os.Stderr, "No mention of version `%s` in changelog file `Changes.md`\n", version)
			os.Exit(1)
		}
		return
	}

	if *develop {
		versionInfo += `_dev`
		versionInfo += fmt.Sprintf("_%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		if _, err := os.Stat(`.git`); os.IsNotExist(err) {
		} else {
			gitHash, err := exec.Command(`git`, `rev-parse`, `HEAD`).Output()
			if err != nil {
				panic(err)
			}
			versionInfo += fmt.Sprintf("_%s", strings.TrimSpace(string(gitHash)))

			gitDiff, err := exec.Command(`git`, `diff`, `--name-only`, `HEAD`).Output()
			if err != nil {
				panic(err)
			}
			if strings.TrimSpace(string(gitDiff)) != "" {
				versionInfo += `_differ`
			}
		}
	}

	fmt.Println("make_dist :", versionInfo)

	// リリースファイルを集約するフォルダを削除(Cleanup)
	os.RemoveAll(`dist\files`)

	for {
		if _, err := os.Stat(`dist\files`); !os.IsNotExist(err) {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	os.MkdirAll(`dist\files`, 0666)
	os.MkdirAll(`dist\release`, 0666)

	if !*develop {
		// Changesが記載できていることを確認する
		if !checkChanges(version) {
			fmt.Fprintf(os.Stderr, "No mention of version `%s` in changelog file `Changes.md`\n", version)
			os.Exit(1)
		}
	}

	if true {
		// buildを実施し、実行体を dist/files に移動
		doTests()
		doBuild(versionInfo)
	}

	if true {
		createReadme()
		createGitInfo()
		copyToFiles(targets)
	}

	if true {
		createZip(versionInfo)
	}
}

func checkChanges(version string) bool {
	rfp, err := os.Open(`Changes.md`)
	if err != nil {
		panic(err)
	}
	defer rfp.Close()

	checkChangesOk := false

	scanner := bufio.NewScanner(rfp)

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), version) {
			checkChangesOk = true
			break
		}
	}

	return checkChangesOk
}

func createZip(version string) {
	datestr := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	wfp, err := os.Create(`dist/release/` + datestr + `_` + appName + `_` + version + `.zip`)
	if err != nil {
		panic(err)
	}
	defer wfp.Close()

	fileList := []compressFile{}
	err = filepath.Walk(`dist/files`, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			distRelPath, err := filepath.Rel(`dist/files`, path)
			if err != nil {
				return err
			}
			fileList = append(fileList, compressFile{path, distRelPath})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	b := doCompress(fileList)
	wfp.Write(b.Bytes())
}

type compressFile struct {
	src, distRelPath string
}

func doCompress(files []compressFile) *bytes.Buffer {
	b := new(bytes.Buffer)
	w := zip.NewWriter(b)
	defer w.Close()

	for _, file := range files {
		info, _ := os.Stat(file.src)

		hdr, _ := zip.FileInfoHeader(info)
		hdr.Name = file.distRelPath
		hdr.Method = zip.Deflate
		f, err := w.CreateHeader(hdr)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadFile(file.src)
		if err != nil {
			panic(err)
		}
		f.Write(body)
	}

	return b
}

func createReadme() {
	rfp, err := os.Open(`README.md`)
	if err != nil {
		panic(err)
	}
	defer rfp.Close()

	rfpCl, err := os.Open(`Changes.md`)
	if err != nil {
		panic(err)
	}
	defer rfpCl.Close()

	wfp, err := os.Create(`dist/files/readme.txt`)
	if err != nil {
		panic(err)
	}
	defer wfp.Close()

	wfpSjis := transform.NewWriter(wfp, japanese.ShiftJIS.NewEncoder())
	defer wfpSjis.Close()

	scanner := bufio.NewScanner(rfp)
	for scanner.Scan() {
		fmt.Fprintf(wfpSjis, "%s\r\n", scanner.Text())
	}

	fmt.Fprintf(wfpSjis, "\r\n")

	scanner = bufio.NewScanner(rfpCl)
	for scanner.Scan() {
		fmt.Fprintf(wfpSjis, "%s\r\n", scanner.Text())
	}

}

func createGitInfo() {
	if _, err := os.Stat(`.git`); os.IsNotExist(err) {
		return
	}

	wfp, err := os.Create(`dist/files/git_info.txt`)
	if err != nil {
		panic(err)
	}
	defer wfp.Close()

	gitHash, err := exec.Command(`git`, `rev-parse`, `HEAD`).Output()
	if err != nil {
		panic(err)
	}

	gitDiff, err := exec.Command(`git`, `diff`, `--name-only`, `HEAD`).Output()
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(wfp, "$ git rev-parse HEAD\r\n")
	fmt.Fprintf(wfp, "%s\r\n", string(gitHash))
	fmt.Fprintf(wfp, "$ git diff --name-only HEAD\r\n")
	fmt.Fprintf(wfp, "%s\r\n", string(gitDiff))
}

func extractFromZip(version string) {
	r, err := zip.OpenReader(`./dist/release/` + version + `/` + appName + `_` + version + `_windows_386.zip`)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.Compare(f.Name, appName+`.exe`) == 0 {
			extractToFile(f, `./dist/files/`+appName+`.exe`)
		}
	}
}

func extractToFile(from *zip.File, to string) error {
	f, err := from.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := os.Create(to)
	if err != nil {
		return err
	}
	defer t.Close()

	_, err = io.Copy(t, f)
	if err != nil {
		return err
	}

	return nil
}

func doTests() {
	items, err := ioutil.ReadDir(`.`)
	if err != nil {
		panic(err)
	}

	testTarget := []string{`.`}
	for _, x := range items {
		if x.IsDir() {
			if strings.HasPrefix(x.Name(), `.`) {
				// skip
				// ex) .git
			} else if x.Name() == `vendor` {
			} else if x.Name() == `dist` {
			} else {
				testTarget = append(testTarget, fmt.Sprintf(`./%s/...`, x.Name()))
			}
		}
	}
	//fmt.Println(testTarget)

	{
		cmd := exec.Command(`go`, `test`)
		cmd.Args = append(cmd.Args, testTarget...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			panic(err)
		}

		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "go test failed\n")
			os.Exit(1)
		}
	}

	{
		cmd := exec.Command(`golangci-lint`, `run`, `--disable-all`, `--enable=vet`, `--enable=vetshadow`, `--enable=golint`, `--enable=ineffassign`, `--enable=goconst`, `--enable=goimports`, `--tests`)

		cmd.Args = append(cmd.Args, testTarget...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			panic(err)
		}

		if err := cmd.Wait(); err != nil {
			fmt.Fprintf(os.Stderr, "gometalinter failed\n")
			os.Exit(1)
		}
	}

}

func doBuild(version string) {
	os.Setenv(`GOOS`, `windows`)
	os.Setenv(`GOARCH`, `386`)
	buildDate := time.Now()

	opt := fmt.Sprintf(`-ldflags=-X main.VERSION=%s -X "main.BUILDDATE=%s"`, version, buildDate.Format(`2006/01/02 15:04:05 -0700 MST`))
	fmt.Println(`go`, `build`, `-o`, appName+`.exe`, opt)
	cmd := exec.Command(`go`, `build`, `-o`, appName+`.exe`, opt)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	if err := fileCopy(appName+`.exe`, `dist/files/`+appName+`.exe`); err != nil {
		panic(err)
	}
}

func fileCopy(srcName, dstName string) error {
	fi, err := os.Stat(srcName)
	if err != nil {
		panic(err)
	}
	if fi.IsDir() {
		os.MkdirAll(dstName, 0666)
		files, err2 := ioutil.ReadDir(srcName)
		if err2 != nil {
			panic(err2)
		}
		for _, f := range files {
			fileCopy(filepath.Join(srcName, f.Name()), filepath.Join(dstName, f.Name()))
		}
	} else {
		d := filepath.Dir(dstName)
		if d != `.` {
			os.MkdirAll(d, 0666)
		}
	}

	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func zipFile(writer *zip.Writer, zipPath string, targetFile string) error {
	f, err := os.Open(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	info, err := os.Stat(targetFile)
	if err != nil {
		return err
	}

	header, _ := zip.FileInfoHeader(info)
	//zip用のパスを設定
	//これを設定しないと、zip内でディレクトリの中に作られない。
	header.Name = zipPath

	zf, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	if _, err := zf.Write(body); err != nil {
		return err
	}
	return nil
}

func copyToFiles(targets []string) {
	for _, t := range targets {
		fi, err := os.Stat(t)
		if err != nil {
			panic(err)
		}
		if fi.IsDir() {
			fileCopy(t, filepath.Join(`dist/files`, t))
		} else {
			fileCopy(t, filepath.Join(`dist/files`, t))
		}
	}
}
