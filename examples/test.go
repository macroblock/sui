package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"github.com/veandco/go-sdl2/sdl"
	"golang.org/x/crypto/ssh"

	"github.com/macroblock/sui"
)

const (
	minThreads = 1
	maxThreads = 50
	tempExt    = ".part"
)

var (
	fileMutex   sync.Mutex
	db          map[string]struct{}
	fdb         *os.File
	ftpMode     = "ftp"
	ftpHost     = ""
	ftpPort     = -1
	ftpUser     = ""
	ftpPassword = ""
	ftpPath     = "/master" //"/for_ott" //"/master" // "/temp"

	files []string

	// ftpConn    *ftp.ServerConn
	numThreads = 1

	root    *sui.RootWindow
	lbFiles *ListBox
)

func ftpInit() (IFtp, error) { //(*ftp.ServerConn, error) {
	switch ftpMode {
	case "ftp":
		c, err := ftp.DialTimeout(ftpHost+":"+strconv.Itoa(ftpPort), 5*time.Second)
		if err != nil {
			return nil, err
		}
		err = c.Login(ftpUser, ftpPassword)
		if err != nil {
			return nil, err
		}
		err = c.ChangeDir(ftpPath)
		if err != nil {
			return nil, err
		}
		return c, nil
	case "sftp":
		addr := ftpHost + ":" + strconv.Itoa(ftpPort)
		config := &ssh.ClientConfig{
			User: ftpUser,
			Auth: []ssh.AuthMethod{
				ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
					// Just send the password back for all questions
					answers := make([]string, len(questions))
					for i := range answers {
						answers[i] = ftpPassword
					}
					return answers, nil
				}),
				ssh.PasswordCallback(func() (string, error) { return ftpPassword, nil }),
				ssh.Password(ftpPassword),
			},
			// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,
			// Config: ssh.Config{
			// 	//Ciphers: []string{"aes128-cbc"},
			// 	//Ciphers: []string{"3des-cbc", "aes256-cbc", "aes192-cbc", "aes128-cbc"},
			// 	Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com",
			// 		"arcfour256", "arcfour128", "aes128-cbc", "aes192-cbc", "aes256-cbc", "3des-cbc", "des-cbc",
			// 	},
			// },
		}
		conn, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			return nil, fmt.Errorf("Failed to dial: " + err.Error())
		}
		client, err := sftp.NewClient(conn, sftp.MaxPacketUnchecked(261120))
		if err != nil {
			return nil, fmt.Errorf("Failed to create client: " + err.Error())
		}
		c := &TSftp{}
		c.client = client
		// c.client.
		return c, nil
	}
	return nil, fmt.Errorf("unknown create ftp client mode '%v'", ftpMode)
}

func ftpClose(c *ftp.ServerConn) {
	c.Quit()
}

type ftpItem struct {
	log       *os.File
	c         IFtp //*ftp.ServerConn
	filename  string
	stopped   bool
	working   bool
	done      bool
	fileSize  int64
	bytesSent int64
	file      *os.File
	//fileSize int64
	err        error
	oldErr     error
	bps        [10]int
	currentBps int
	prevTime   time.Time //int64
	bpsIndex   int
	numSent    int64
	started    time.Time
	completed  time.Time
}

func getTime() int64 {
	return time.Now().UnixNano()
}

func NewFtpItem(name string) *ftpItem {
	item := &ftpItem{}
	item.prevTime = time.Now() //getTime()
	item.filename = name
	return item
}

func (o *ftpItem) Close() {
	o.stopped = true
	if o.file != nil {
		o.file.Close()
		o.file = nil
	}
	if o.c != nil {
		o.c.Quit()
		o.c = nil
	}
}

func (o *ftpItem) Stop() {
	o.err = fmt.Errorf("was stopped")
	o.stopped = true
	o.Clean()
}

func (o *ftpItem) Play() {
	if !o.working {
		o.stopped = false
	}
}

func (o *ftpItem) NextIndex() {
	o.bpsIndex++
	o.bpsIndex %= len(o.bps)
}

func (o *ftpItem) Bps() int {
	count := 0
	bps := 0
	for _, v := range o.bps {
		if v != 0 {
			bps += v
			count++
		}
	}
	if count == 0 {
		o.currentBps = 0
	} else {
		o.currentBps = bps / count
	}
	return o.currentBps
}

func (o *ftpItem) Read(p []byte) (int, error) {
	n, err := o.file.Read(p)
	if err == nil {
		o.numSent += int64(n)
		now := time.Now()            //getTime()
		delta := now.Sub(o.prevTime) //now - o.prevTime
		//o.numSent += int64(n)
		if delta >= time.Second {
			//o.bps[o.bpsIndex] = int(o.numSent * int64(time.Second) / int64(delta))
			o.bps[o.bpsIndex] = int(float64(o.numSent) / (float64(delta) / float64(time.Second))) // 10
			//fmt.Println("bps: ", o.bps[o.bpsIndex], n, delta)
			o.NextIndex()
			o.numSent = 0
			o.prevTime = now
		} else {
			//fmt.Println("Delta is Zero", n)
		}

		o.bytesSent += int64(n)
		//fmt.Println("Read", n, "bytes for a total of", pt.total)
	} else {
		// TODO: should it be handled ?
		if err != io.EOF {
			// fmt.Println("!!! read err: ", err, n)
		}
	}

	return n, err
}

func (o *ftpItem) InitLogFile() {
	err := error(nil)
	o.log, err = os.OpenFile(filepath.Join("log", filepath.Base(o.filename)+".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
}

func (o *ftpItem) InitLocalFile() {
	if o.err != nil {
		// fmt.Println("init local: was error")
		return
	}
	writeLog(o.log, "opening local file: %v", o.filename)
	file, err := os.Open(o.filename)
	if err != nil {
		o.err = err
		return
	}
	// fmt.Println("open ok")
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		file = nil
		o.err = err
		return
	}
	// fmt.Println("stat ok")
	o.fileSize = stat.Size()
	o.file = file
}

func (o *ftpItem) InitRemoteFile() bool {
	if o.err != nil {
		// fmt.Println("init remote: was error")
		return false
	}
	o.bytesSent = 0
	remoteFile := ftpPath + "/" + filepath.Base(o.filename)
	writeLog(o.log, "checking remote file %q", remoteFile)
	if size, err := o.c.FileSize(remoteFile); err == nil {
		// file already uploaded
		if o.fileSize == size {
			return true
		}
	}

	remoteFile = ftpPath + "/" + filepath.Base(o.filename) + tempExt
	if size, err := o.c.FileSize(remoteFile); err == nil && size <= o.fileSize {
		o.bytesSent = size
		writeLog(o.log, "upload will be continued from %v bytes", o.bytesSent)
	} else {
		writeLog(o.log, "deleting file %q because of inconvinient size (%v > %v) if it exists", remoteFile, size, o.fileSize)
		o.c.Delete(remoteFile)
	}
	return false
}

func (o *ftpItem) ReadyToWork() bool {
	if !o.working && !o.stopped && !o.done {
		return true
	}
	return false
}

func (o *ftpItem) Stor() {
	if o.err != nil {
		// fmt.Println("stor: was an error")
		return
	}
	// fmt.Println("before stor is ok")
	writeLog(o.log, "skipping already uploaded data (%v bytes)", o.bytesSent)
	if pos, err := o.file.Seek(o.bytesSent, 0); pos != o.bytesSent || err != nil {
		o.bytesSent = 0
		if pos, err = o.file.Seek(0, 0); pos != o.bytesSent || err != nil {
			o.err = errors.New("File.Seek() ERROR: something criticaly wrong")
			o.err = err
			return
		}
	}

	o.started = time.Now()
	// fmt.Println("store started: ", o.bytesSent, filepath.Base(o.filename)+tempExt)
	writeLog(o.log, "upload begun")
	o.err = o.c.StorFrom(ftpPath+"/"+filepath.Base(o.filename)+tempExt, o, uint64(o.bytesSent))

	o.completed = time.Now()
	// fmt.Println("delta: ", o.completed.Sub(o.started))
	writeLog(o.log, "upload ended")

}

func (o *ftpItem) PostProcess() {
	if o.err != nil {
		// fmt.Println("PostProcess was error: " + fmt.Sprintf("%v", o.err))
		return
	}
	if o.c == nil || o.stopped {
		// fmt.Println("PostProcess o.c is nil or stopped")
		writeLog(o.log, "upload stopped by user (or o.c is nil)")
		return
	}
	src := path.Join(ftpPath, filepath.Base(o.filename)+tempExt)
	dst := path.Join(ftpPath, filepath.Base(o.filename))

	writeLog(o.log, "deleting old file %q if it exists:", dst)
	o.c.Delete(dst)

	writeLog(o.log, "renaming %q -> %q:", src, dst)
	// fmt.Println("rename to  :", dst)
	o.err = o.c.Rename(src, dst)
	if o.err != nil {
		// fmt.Println("post process ERROR: " + fmt.Sprintf("%v", o.err))
		writeLog(o.log, "post process ERROR")
	}
}

func (o *ftpItem) Clean() {
	if o.file != nil {
		o.file.Close()
		o.file = nil
	}
	if o.c != nil {
		o.c.Quit()
		o.c = nil
	}
}

func (o *ftpItem) StartJob() {
	o.currentBps = 0
	o.numSent = 0
	o.bpsIndex = 0
	for i := range o.bps {
		o.bps[i] = 0
	}
	o.bytesSent = 0
	o.fileSize = 1 // hack
	o.working = true
	go o.job()
}

func (o *ftpItem) job() {
	//o.stopped = true
	o.oldErr = o.err
	o.InitLogFile()
	defer o.log.Close()

	writeLog(o.log, "------------------------------")
	writeLog(o.log, "job started")
	o.c, o.err = ftpInit()
	if o.c != nil {
		defer o.c.Quit()
	}
	o.InitLocalFile()
	if o.InitRemoteFile() {
		// file already uploaded
		fmt.Printf("WARNING: file has been already uploaded but wasn't found in db: %v\n", filepath.Base(o.filename))
		writeLog(o.log, "WARNING: file has been already uploaded but wasn't found in db!")
		writeLog(o.log, "job closed. (1)")
		o.Clean()
		o.done = true
		o.working = false
		return
	}
	o.Stor()
	o.PostProcess()
	if o.err == nil {
		if !o.stopped {
			writeLog(o.log, "upload successfully done")
			appendName(fdb, &db, filepath.Base(o.filename))
			o.done = true
		} else {
			writeLog(o.log, "upload interrupted by user")
		}
	} else {
		writeLog(o.log, "upload interrupted by ERROR: %v", o.err)
	}
	writeLog(o.log, "job closed. (2)")
	o.Clean()
	o.working = false
}

func loop() {
	sui.PostUpdate()
	items := lbFiles.Items()

	// delete what is done
	for i := 0; i < len(items); {
		item := items[i].Data.(*ftpItem)
		if item.done {
			if lbFiles.itemIndex == i {
				lbFiles.itemIndex = -1
			} else if lbFiles.itemIndex > i {
				lbFiles.itemIndex--
			}
			items = deleteFtpItem(items, i)
			lbFiles.items = items
		} else {
			i++
		}
	}

	// count active workers
	numWorkers := 0
	for i := range items {
		item := items[i].Data.(*ftpItem)
		//percent := 0
		if item.working {
			//percent = int(item.bytesSent * 100 / item.fileSize)
			numWorkers++
		}

		//items[i].Name = fmt.Sprint(percent, item.stopped, item.working, item.bytesSent, " ", filepath.Base(item.filename))
	}
	if numWorkers >= numThreads {
		return
	}
	// run new workers
	for i := range items {
		item := items[i].Data.(*ftpItem)
		if item.ReadyToWork() {
			item.StartJob()
			return
		}
	}
}

func deleteFtpItem(items []listBoxItem, i int) []listBoxItem {
	if i < 0 && i >= len(items) {
		panic("trying delete wrong index")
	}
	//fmt.Println("\nin: ", i, "\n", items)
	if items[i].Data != nil {
		item := items[i].Data.(*ftpItem)
		item.Close()
		items[i].Data = nil
	}
	if i == 0 {
		items = items[1:]
	} else if i == len(items)-1 {
		items = items[:i]
	} else {
		items = append(items[:i], items[i+1:]...)
	}
	//fmt.Println("\nout: ", i, "\n", items)
	return items
}

func onKeyPress() {
	switch sui.KeySym() {
	case sdl.K_DELETE:
		items := lbFiles.Items()
		for i := 0; i < len(items); {
			if items[i].Selected || i == lbFiles.itemIndex {
				items = deleteFtpItem(items, i)
				lbFiles.itemIndex = -1
				lbFiles.items = items
				sui.PostUpdate()
			} else {
				i++
			}
		}
		writeFtpItems()

	case sdl.K_HOME, sdl.K_LEFT:
		lbFiles.itemIndex = sui.MinInt(0, len(lbFiles.Items())-1)
		lbFiles.offset = 0
		sui.PostUpdate()
	case sdl.K_END, sdl.K_RIGHT:
		lbFiles.itemIndex = len(lbFiles.Items()) - 1
		lbFiles.offset = sui.MaxInt(0, lbFiles.itemIndex-lbFiles.Size().Y/itemHeight+1)
		sui.PostUpdate()
	case sdl.K_UP:
		lbFiles.itemIndex--
		lbFiles.itemIndex = sui.MaxInt(0, lbFiles.itemIndex)
		lbFiles.itemIndex = sui.MinInt(len(lbFiles.Items())-1, lbFiles.itemIndex)
		lbFiles.CalcOffset()
		sui.PostUpdate()
	case sdl.K_DOWN:
		//fmt.Println("key DOWN")
		lbFiles.itemIndex++
		lbFiles.itemIndex = sui.MaxInt(0, lbFiles.itemIndex)
		lbFiles.itemIndex = sui.MinInt(len(lbFiles.Items())-1, lbFiles.itemIndex)
		lbFiles.CalcOffset()
		sui.PostUpdate()
	}
}

func moveTo(toTop bool) {
	items := lbFiles.Items()
	newItems := []listBoxItem{}
	for i := range items {
		if items[i].Selected || i == lbFiles.itemIndex {
			p := items[i]
			items[i].Data = nil
			//p.Selected = false
			newItems = append(newItems, p)
			lbFiles.itemIndex = -1
			lbFiles.items = items
			sui.PostUpdate()
		}
	}
	for i := 0; i < len(items); {
		if items[i].Data == nil {
			items = deleteFtpItem(items, i)
		} else {
			i++
		}
	}

	if toTop {
		items = append(newItems, items...)
		lbFiles.itemIndex = sui.MinInt(0, len(items)-1)
	} else {
		items = append(items, newItems...)
		lbFiles.itemIndex = len(items) - 1
	}
	lbFiles.items = items
	lbFiles.CalcOffset()
	writeFtpItems()
}

func onDropFile() {
	dropFileName := sui.DropFile()

	for _, v := range lbFiles.items {
		item := v.Data.(*ftpItem)
		if item.filename == dropFileName {
			fmt.Println("duplicated file in queue: ", dropFileName)
			return
		}
	}

	if existsName(db, filepath.Base(dropFileName)) {
		fmt.Println("file has been already uploaded (found in db): ", filepath.Base(dropFileName))
		return
	}

	item := NewFtpItem(dropFileName)
	lbFiles.AddItem(fmt.Sprint(item.stopped, " ", item.filename), item)
	lbFiles.itemIndex = len(lbFiles.items) - 1
	lbFiles.CalcOffset()
	sui.PostUpdate()
	writeFtpItems()
	//files = append(files, sui.DropFile())
}

func onMouseClick() {
	//o := sui.Sender()
	//fmt.Println("!!!!!!!! MouseClick: ", o)
	//o.SetClearColor(sui.Palette.BackgroundLo)
}

func onMouseOver() {
	x := sui.MouseOver()
	if x != nil && x != root {
		x.SetClearColor(sui.Palette.Active.Hi())
	}
	if sui.PrevMouseOver() != nil && sui.PrevMouseOver() != root {
		if sui.PrevMouseOver().ClearColor() != sui.Palette.Passive {
			sui.PrevMouseOver().SetClearColor(sui.Palette.Active)
		}
	}
}

func onMouseOverPassive() {
	x := sui.MouseOver()
	if x != nil && x != root {
		x.SetClearColor(sui.Palette.Passive)
	}
	if sui.PrevMouseOver() != nil && sui.PrevMouseOver() != root {
		if sui.PrevMouseOver().ClearColor() != sui.Palette.Passive {
			sui.PrevMouseOver().SetClearColor(sui.Palette.Active)
		}

	}
}

func onPressMouseDown() {
	//o := sui.Sender()
	//fmt.Println("MousePressDown: ", o)
}

func onPressMouseUp() {
	//o := sui.Sender()
	//fmt.Println("MousePressUp: ", o)
}

func readFtpItems() error {
	file, err := os.Open("queue.lst")
	if err != nil {
		return err
	}
	defer file.Close()

	//lines := []ftpItem{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		name = strings.TrimSpace(name)
		if existsName(db, name) {
			continue
		}
		if existsName(db, filepath.Base(name)) {
			fmt.Println("load queue error: file has been already uploaded (found in db): ", filepath.Base(name))
			continue
		}
		item := NewFtpItem(scanner.Text())
		lbFiles.AddItem(fmt.Sprint(item.stopped, " ", item.filename), item)
		//lines = append(lines, scanner.Text())
	}
	return scanner.Err()
}

// writeLines writes the lines to the given file.
func writeFtpItems() error {
	file, err := os.Create("queue.lst")
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, item := range lbFiles.items {
		fmt.Fprintln(w, item.Data.(*ftpItem).filename)
	}
	return w.Flush()
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readDb(path string) (map[string]struct{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// var lines []string
	db := map[string]struct{}{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		// fmt.Println("line: ", s)
		x := strings.Split(s, string('\t'))
		// fmt.Println("slice: ", x[1])
		if len(x) != 2 {
			panic("wrong db format")
		}
		s = strings.TrimSpace(x[1])
		if s == "" {
			continue
		}
		// lines = append(lines, s)
		db[s] = struct{}{}
	}
	return db, scanner.Err()
}

func existsName(db map[string]struct{}, name string) bool {
	// for _, s := range db {
	// 	// fmt.Printf("%v <-> %v\n", name, s)
	// 	if s == name {
	// 		return true
	// 	}
	// }
	// return false
	_, ok := db[name]
	return ok
}

func appendName(f *os.File, db *map[string]struct{}, name string) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if _, err := f.WriteString(fmt.Sprintf("%v\t%v\n", time.Now().Format("2006-01-02 15:04:05"), name)); err != nil {
		panic(err)
	}
	(*db)[name] = struct{}{}
	// *db = append(*db, name)
}

func writeLog(f *os.File, format string, args ...interface{}) {
	if _, err := f.WriteString(fmt.Sprintf("%v\t%v\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, args...))); err != nil {
		panic(err)
	}
}

func main() {
	//fmt.Println(ftpUser + ":" + ftpPassword + "@ftp://" + ftpHost + ":" + strconv.Itoa(ftpPort))
	//ftpTest()

	err := error(nil)
	db, err = readDb("db.list")
	if err != nil {
		panic("somethin wrong while reading db")
	}
	// fmt.Println("db: ", db)
	// fmt.Println("err: ", err)

	fdb, err = os.OpenFile("db.list", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fdb.Close()

	// if _, err = f.WriteString(text); err != nil {
	// 	panic(err)
	// }

	err = sui.Init()
	defer sui.Close()
	if err != nil {
		panic(err)
	}
	root = sui.NewRootWindow("("+ftpMode+")://"+ftpHost+ftpPath, 800, 600)
	//root.SetClearColor(sui.Palette.BackgroundLo)
	//root.SetClearColor(sui.Color32(0x00000000))
	root.OnDropFile = onDropFile
	//root.OnDraw = onDraw
	//root.OnEnter = onEnter
	//root.OnLeave = onLeave
	root.OnMouseOver = onMouseOver
	root.OnMouseButtonDown = onPressMouseDown
	root.OnMouseButtonUp = onPressMouseUp
	root.OnKeyPress = onKeyPress
	//root.OnMouseClick = onMouseClick

	btnInc := sui.NewBox(40, 35)
	btnInc.Move(5, 5)
	btnInc.OnMouseOver = onMouseOver
	btnInc.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Inc")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnInc.OnMouseClick = func() {
		numThreads = sui.MinInt(numThreads+1, maxThreads)
		sui.PostUpdate()
	}

	btnDec := sui.NewBox(40, 35)
	btnDec.Move(50, 5)
	btnDec.OnMouseOver = onMouseOver
	btnDec.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Dec")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnDec.OnMouseClick = func() {
		numThreads = sui.MaxInt(numThreads-1, minThreads)
		sui.PostUpdate()
	}

	lblNumThreads := sui.NewBox(40, 35)
	lblNumThreads.Move(95, 5)
	lblNumThreads.SetClearColor(sui.Palette.Passive)
	lblNumThreads.OnMouseOver = onMouseOverPassive
	lblNumThreads.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(10, 5), strconv.Itoa(numThreads))
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}

	btnStop := sui.NewBox(50, 35)
	btnStop.Move(140, 5)
	btnStop.OnMouseOver = onMouseOver
	btnStop.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Stop")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnStop.OnMouseClick = func() {
		items := lbFiles.items
		for i := range items {
			if items[i].Selected || i == lbFiles.itemIndex {
				item := items[i].Data.(*ftpItem)
				item.Stop()
			}
		}
		sui.PostUpdate()
	}

	btnPlay := sui.NewBox(50, 35)
	btnPlay.Move(195, 5)
	btnPlay.OnMouseOver = onMouseOver
	btnPlay.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Start")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnPlay.OnMouseClick = func() {
		items := lbFiles.items
		for i := range items {
			if items[i].Selected || i == lbFiles.itemIndex {
				item := items[i].Data.(*ftpItem)
				item.Play()
			}
		}
		sui.PostUpdate()
	}

	btnMoveToTop := sui.NewBox(95, 35)
	btnMoveToTop.Move(300, 5)
	btnMoveToTop.OnMouseOver = onMouseOver
	btnMoveToTop.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "to Top")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnMoveToTop.OnMouseClick = func() {
		moveTo(true)
		sui.PostUpdate()
	}

	btnMoveToBottom := sui.NewBox(95, 35)
	btnMoveToBottom.Move(400, 5)
	btnMoveToBottom.OnMouseOver = onMouseOver
	btnMoveToBottom.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "to Bottom")
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}
	btnMoveToBottom.OnMouseClick = func() {
		moveTo(false)
		sui.PostUpdate()
	}

	infoBps := sui.NewBox(150, 35)
	infoBps.Move(500, 5)
	infoBps.SetClearColor(sui.Palette.Passive)
	infoBps.OnMouseOver = onMouseOverPassive
	infoBps.OnDraw = func() {
		o := sui.Sender()
		bps := 0
		for _, i := range lbFiles.items {
			item := i.Data.(*ftpItem)
			if item.working {
				bps += item.currentBps
			}
		}
		o.Clear()
		o.WriteText(sui.NewPoint(5, 5), "Speed: "+bpsToStr(bps))
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}

	lbFiles = NewListBox(790, 350)
	lbFiles.Move(5, 45)
	lbFiles.OnMouseOver = onMouseOver

	fInfo := sui.NewBox(790, 195)
	fInfo.Move(5, 400)
	fInfo.SetClearColor(sui.Palette.Passive)
	fInfo.OnMouseOver = onMouseOverPassive
	fInfo.OnDraw = func() {
		o := sui.Sender()
		o.Clear()
		y := 5
		dy := itemHeight
		if lbFiles.itemIndex > -1 && lbFiles.items[lbFiles.itemIndex].Data != nil {
			item := lbFiles.items[lbFiles.itemIndex].Data.(*ftpItem)
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Filename: %s", item.filename))
			y += dy
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("File: %v", item.file))
			y += dy
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Size: %v", item.fileSize))
			y += dy
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Bytes sent: %v", item.bytesSent))
			y += dy
			s := "waiting"
			if item.working {
				s = "transfering"
			}
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Action: %s", s))
			y += dy
			s = "active"
			if item.stopped {
				s = "stopped"
			}
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Status: %s", s))
			y += dy
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("is done: %v", item.done))
			y += dy
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Last error: %v", item.err))
			y += dy
			o.WriteText(sui.NewPoint(10, y), fmt.Sprintf("Prev error: %v", item.oldErr))
		}
		o.Rect(sui.NewRect(sui.NewPoint(0, 0), o.Size()))
	}

	root.OnResize = func() {
		size := root.Size()
		w := size.X - 10
		h := size.Y - 250
		lbFiles.Resize(w, h)
		x := 5
		y := size.Y - 200
		w = size.X - 10
		h = 195
		fInfo.Move(x, y)
		fInfo.Resize(w, h)
	}

	root.AddChild(btnInc)
	root.AddChild(btnDec)
	root.AddChild(lblNumThreads)
	root.AddChild(btnStop)
	root.AddChild(btnPlay)
	root.AddChild(btnMoveToTop)
	root.AddChild(btnMoveToBottom)
	root.AddChild(infoBps)
	root.AddChild(lbFiles)
	root.AddChild(fInfo)

	readFtpItems()

	sui.OnLoop = loop

	sui.Run()

	writeFtpItems()

	root.Close()
}
