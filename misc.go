/*
Package misc implements a differents trivial functions
*/
package misc

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Appliction exit codes
const (
	// ExPanic --
	ExPanic = 70
	// ExStopped --
	ExStopped = 1
	// ExVersion --
	ExVersion = 64
	// ExMissingConfigFile --
	ExMissingConfigFile = 66
	// ExIncorrectConfigFile --
	ExIncorrectConfigFile = 78
	// ExConfigIncorrect --
	ExConfigIncorrect = 78
	// ExConfigErrors --
	ExConfigErrors = 78
	// ExCreateListenerError --
	ExCreateListenerError = 71
	// ExStartListenerError --
	ExStartListenerError = 71
	// ExServiceInitializationError --
	ExServiceInitializationError = 71
	// ExServiceError --
	ExServiceError = 71
	// ExAccessDenied --
	ExAccessDenied = 77
	// ExProgrammerError --
	ExProgrammerError = 70
)

//----------------------------------------------------------------------------------------------------------------------------//

var (
	// TEST -- test mode
	TEST = false
)

//----------------------------------------------------------------------------------------------------------------------------//

// go build --ldflags "-X github.com/alrusov/misc.appVersion=${VERSION} -X github.com/alrusov/misc.appTags=${TAGS} -X github.com/alrusov/misc.buildTime=`date +'%Y-%m-%d_%H:%M:%S'` -X github.com/alrusov/misc.copyright=${COPYRIGHT}"

var (
	appVersion  string
	appTags     string
	copyright   string
	buildTime   string
	buildTimeTS time.Time
)

//----------------------------------------------------------------------------------------------------------------------------//

var (
	startTime time.Time

	appFullName string
	appExecPath string
	appExecName string
	appName     string
	appWorkDir  string

	appStarted   = int32(1)
	exitLaunched = int32(0)

	exitCode = 0

	sleepInterrupt = make(chan bool, 1)

	cond      *sync.Cond
	exitChain = make([]exitElement, 0)

	// Logger --
	Logger loggerFunc
)

type (
	// ExitFunc --
	ExitFunc func(code int, param interface{})

	exitElement struct {
		name  string
		f     ExitFunc
		param interface{}
	}

	loggerFunc func(facility string, level string, message string, params ...interface{})
)

//----------------------------------------------------------------------------------------------------------------------------//

func killer() {
	time.Sleep(5000 * time.Millisecond)
	Logger("", "CR", "Application shutdown timeout. Forced completion.")
	go Exit()

	time.Sleep(5000 * time.Millisecond)
	Logger("", "CR", "Application forced completion timeout. Forced Killing.")
	os.Exit(exitCode)
}

// StopApp -- set exit code and raise application stop
func StopApp(code int) {
	if atomic.AddInt32(&appStarted, -1) == 0 {
		Logger("", "DE", "Set application exit code %d", code)

		exitCode = code

		cond.Broadcast()
		time.Sleep(100 * time.Millisecond)

		ex := false
		for !ex {
			select {
			case sleepInterrupt <- true:
			default:
				ex = true
			}
		}

		time.Sleep(100 * time.Millisecond)

		go killer()
	}
}

// WaitingForStop --
func WaitingForStop() {
	cond.L.Lock()
	cond.Wait()
	cond.L.Unlock()
}

// Exit -- exit application
func Exit() {
	if atomic.AddInt32(&exitLaunched, 1) == 1 {
		if AppStarted() {
			StopApp(0)
		}

		Logger("", "IN", "Try to finish application with code %d", exitCode)

		time.Sleep(1000 * time.Millisecond)

		for i := len(exitChain) - 1; i >= 0; i-- {
			Logger("", "DE", "Call finalizer \"%s\"", exitChain[i].name)
			exitChain[i].f(exitCode, exitChain[i].param)
		}

		Logger("", "IN", "Application finished with code %d", exitCode)
		os.Exit(exitCode)
	}
}

// AddExitFunc --
func AddExitFunc(name string, f ExitFunc, param interface{}) {
	DelExitFunc(name)
	exitChain = append(exitChain, exitElement{name: name, f: f, param: param})
}

// DelExitFunc --
func DelExitFunc(name string) {
	chain := make([]exitElement, 0)
	for i := 0; i < len(exitChain); i++ {
		if exitChain[i].name != name {
			chain = append(chain, exitChain[i])
		}
	}
	exitChain = chain
}

//----------------------------------------------------------------------------------------------------------------------------//

// SimpleLogger --
func simpleLogger(facility string, level string, message string, params ...interface{}) {
	fmt.Printf(level+": "+message+EOS, params...)
}

func init() {
	startTime = NowUTC()

	Logger = simpleLogger

	cond = sync.NewCond(new(sync.Mutex))

	go signalHandler()

	if appVersion == "" {
		appVersion = "debug"
	}
	copyright = strings.Replace(copyright, "_", " ", -1)
	buildTime = strings.Replace(buildTime, "_", " ", -1)
	appVersion = strings.Replace(appVersion, "_", " ", -1)
	appTags = strings.Replace(appTags, "_", " ", -1)

	buildTimeTS, _ = time.Parse(DateTimeFormatRev, buildTime)

	p, _ := os.Executable()
	appFullName, _ = filepath.Abs(p)
	appExecPath = filepath.Dir(appFullName)
	appExecName = filepath.Base(appFullName)
	appName = strings.TrimSuffix(appExecName, filepath.Ext(appExecName))
	appWorkDir, _ = os.Getwd()
}

//----------------------------------------------------------------------------------------------------------------------------//

// AbsPathEx --
func AbsPathEx(name string, base string) (string, error) {
	prefix := ""
	if name != "" {
		prefix = name[0:1]
	}

	switch prefix {
	case "@":
		name = appWorkDir + "/" + name[1:]
	case "$":
		d, _ := os.Getwd()
		name = d + "/" + name[1:]
	case "^":
		name = base + "/" + name[1:]
	default:
		if !filepath.IsAbs(name) {
			name = AppExecPath() + "/" + name
		}
	}

	return filepath.Abs(name)
}

// AbsPath --
func AbsPath(name string) (string, error) {
	return AbsPathEx(name, AppWorkDir())
}

//----------------------------------------------------------------------------------------------------------------------------//

// IsDebug --
func IsDebug() bool {
	return appExecName == "__debug_bin" // simple workaround for the VS Code
}

//----------------------------------------------------------------------------------------------------------------------------//

// AppStartTime -- time of the apptication start
func AppStartTime() time.Time {
	return startTime
}

//----------------------------------------------------------------------------------------------------------------------------//

// AppVersion -- application version
func AppVersion() string {
	return appVersion
}

// AppTags -- application tags
func AppTags() string {
	return appTags
}

// Copyright --
func Copyright() string {
	return copyright
}

// BuildTime --
func BuildTime() string {
	return buildTime
}

// BuildTimeTS --
func BuildTimeTS() time.Time {
	return buildTimeTS
}

// AppName -- name of the application executable file without last extension
func AppName() string {
	return appName
}

// AppFullName -- application name with full path
func AppFullName() string {
	return appFullName
}

// AppExecPath -- full path of the application executable file
func AppExecPath() string {
	return appExecPath
}

// AppExecName -- name of the application executable file
func AppExecName() string {
	return appExecName
}

// AppWorkDir -- directory where application started from
func AppWorkDir() string {
	return appWorkDir
}

// AppStarted -- is application started?
func AppStarted() bool {
	return atomic.LoadInt32(&appStarted) > 0
}

// ExitCode -- get current exit code
func ExitCode() int {
	return exitCode
}

//----------------------------------------------------------------------------------------------------------------------------//

// GetFuncName -- name of the function from call stack
func GetFuncName(shift int, shortName bool) string {
	ret := ""

	stack := GetCallStack(shift + 1)
	n := len(stack)

	if shortName {
		n = 1
	}

	for i := n - 1; i >= 0; i-- {
		df := stack[i]

		if i != n-1 {
			ret += "->"
		}
		ret += filepath.Base(df.FuncName)
	}

	return ret
}

// CallStackFrame -- call stack element
type CallStackFrame struct {
	FuncName string
	FileName string
	Line     int
}

// GetCallStack -- get call stack
func GetCallStack(shift int) []CallStackFrame {
	var ret []CallStackFrame

	pc := make([]uintptr, 500)
	n := runtime.Callers(2+shift, pc)
	for i := 0; i < n; i++ {
		frame := CallStackFrame{}
		df := pc[i]

		fn := runtime.FuncForPC(df)
		if fn != nil {
			frame.FuncName = fn.Name()
		} else {
			frame.FuncName = "?"
		}

		frame.FileName, frame.Line = fn.FileLine(df)

		ret = append(ret, frame)
	}

	return ret
}

//----------------------------------------------------------------------------------------------------------------------------//

// Sleep --
func Sleep(duration time.Duration) bool {
	if !AppStarted() {
		return false
	}

	select {
	case <-sleepInterrupt:
		return false
	case <-time.After(duration):
		return true
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

// TrimStringAsFloat --
func TrimStringAsFloat(s string) string {
	sp := strings.Split(s, ".")
	v := sp[0]

	if len(sp) > 1 {
		sp[1] = strings.TrimRight(sp[1], "0")
		if sp[1] != "" {
			v += "." + sp[1]
		}
	}

	return v
}

//----------------------------------------------------------------------------------------------------------------------------//

// GetMyIPs --
func GetMyIPs() (map[string]bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	list := make(map[string]bool)

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ip := ""
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP.String()
			case *net.IPAddr:
				ip = v.IP.String()
			}

			if ip != "" {
				list[ip] = true
			}
		}
	}

	return list, nil
}

// IsMyIP --
func IsMyIP(ip string) (bool, error) {
	list, err := GetMyIPs()
	if err != nil {
		return false, err
	}

	_, exists := list[ip]
	return exists, nil
}

//----------------------------------------------------------------------------------------------------------------------------//

// Messages --
type Messages struct {
	mutex *sync.RWMutex
	s     []string
}

// NewMessages --
func NewMessages() *Messages {
	return &Messages{
		mutex: new(sync.RWMutex),
	}
}

// Len --
func (m *Messages) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.s)
}

// Add --
func (m *Messages) Add(msg string, params ...interface{}) {
	if msg != "" {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.s = append(m.s, fmt.Sprintf(msg, params...))
	}
}

// AddError --
func (m *Messages) AddError(err error) {
	if err != nil {
		m.Add(err.Error())
	}
}

// String --
func (m *Messages) String(separators ...string) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.s) == 0 {
		return ""
	}

	sep := "; "
	if len(separators) != 0 {
		sep = strings.Join(separators, "")
	}

	return strings.Join(m.s, sep)
}

// Error --
func (m *Messages) Error(separators ...string) error {
	s := m.String(separators...)
	if s == "" {
		return nil
	}

	return errors.New(s)
}

//----------------------------------------------------------------------------------------------------------------------------//

// LogProcessingTime  --
func LogProcessingTime(facility string, level string, id uint64, module string, message string, t0 int64) int64 {
	if level == "" {
		level = "TM"
	}

	if message == "" {
		message = "Elapsed time"
	} else {
		message += ", elapsed time"
	}

	prefix := ""
	if id != 0 {
		prefix = "[" + strconv.FormatUint(id, 10) + "] "
	}
	if module != "" {
		if prefix == "" {
			prefix = "[" + module + "] "
		} else {
			prefix += module + ": "
		}
	}

	now := NowUnixNano()
	duration := now - t0
	Logger(facility, level, "%s%s %d.%03d ms", prefix, message, duration/int64(time.Millisecond), (duration%int64(time.Millisecond))/1000)
	return now
}

//----------------------------------------------------------------------------------------------------------------------------//

var (
	reSlashes = regexp.MustCompile(`(^|[^:])/{2,}`)
)

// NormalizeSlashes --
func NormalizeSlashes(u string) string {
	u = strings.TrimRight(u, "/")
	return reSlashes.ReplaceAllString(u, `$1/`)
}

//----------------------------------------------------------------------------------------------------------------------------//

// Sha512Hash --
func Sha512Hash(p []byte) []byte {
	h := sha512.Sum512(p)
	s := make([]byte, len(h)*2)
	hex.Encode(s, h[:])
	return s
}

//----------------------------------------------------------------------------------------------------------------------------//

// UnsafeByteSlice2String -- fast convert []byte to string without memory allocs
// Don't forget to use runtime.KeepAlive(b) in the caller if necessary!
func UnsafeByteSlice2String(b []byte) string {
	//return *((*string)(unsafe.Pointer(&b)))
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	s := reflect.StringHeader{
		Data: h.Data,
		Len:  h.Len,
	}

	return *(*string)(unsafe.Pointer(&s))
}

// UnsafeString2ByteSlice -- fast convert string to []byte without memory allocs
// Don't forget to use runtime.KeepAlive(s) in the caller if necessary!
// Don't try to change result without thinking hard before that!
func UnsafeString2ByteSlice(s string) []byte {
	//return *((*[]byte)(unsafe.Pointer(&s)))
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	b := reflect.SliceHeader{
		Data: h.Data,
		Len:  h.Len,
		Cap:  h.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&b))
}

//----------------------------------------------------------------------------------------------------------------------------//

func JoinByteSlices(prefix []byte, suffix []byte, sep []byte, in [][]byte) (out []byte) {
	inCount := len(in)
	ln := len(prefix) + len(suffix)

	if inCount == 0 {
		if ln == 0 {
			return []byte{}
		}

		out = make([]byte, ln)
		pos := copy(out, prefix)
		copy(out[pos:], suffix)
		return
	}

	ln += len(sep) * (inCount - 1)

	for _, v := range in {
		ln += len(v)
	}

	out = make([]byte, ln)
	pos := copy(out, prefix)
	pos += copy(out[pos:], in[0])

	if inCount > 1 {
		for _, v := range in[1:] {
			pos += copy(out[pos:], sep)
			pos += copy(out[pos:], v)
		}
	}

	copy(out[pos:], suffix)

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func JoinStrings(prefix string, suffix string, sep string, in []string) (out string) {
	inCount := len(in)
	ln := len(prefix) + len(suffix)

	if inCount == 0 && ln == 0 {
		return ""
	}

	if inCount > 0 {
		ln += len(sep) * (inCount - 1)

		for _, v := range in {
			ln += len(v)
		}
	}

	var b strings.Builder
	b.Grow(ln)

	b.WriteString(prefix)

	if inCount > 0 {
		b.WriteString(in[0])
		for _, v := range in[1:] {
			b.WriteString(sep)
			b.WriteString(v)
		}
	}

	b.WriteString(suffix)

	return b.String()
}

//----------------------------------------------------------------------------------------------------------------------------//
