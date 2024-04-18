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

	sleepInterrupted = make(chan bool, 1)

	exitCond  = sync.NewCond(new(sync.Mutex))
	exitChain = make([]exitElement, 0)

	// Logger --
	Logger loggerFunc
)

type (
	// ExitFunc --
	ExitFunc func(code int, param any)

	exitElement struct {
		name  string
		f     ExitFunc
		param any
	}

	loggerFunc func(facility string, level string, message string, params ...any)
)

type (
	CtxKey string
)

//----------------------------------------------------------------------------------------------------------------------------//

var (
	terminationTimeout = 5 * time.Second
	killingTimeout     = 5 * time.Second
)

func SetExitTimeouts(newTerminationTimeout time.Duration, newKillingTimeout time.Duration) (prevTerminationTimeout time.Duration, prevKillingTimeout time.Duration) {
	prevTerminationTimeout, prevKillingTimeout = terminationTimeout, killingTimeout

	if newTerminationTimeout > 0 {
		terminationTimeout = newTerminationTimeout
	}

	if newKillingTimeout > 0 {
		killingTimeout = newKillingTimeout
	}

	return
}

func killer() {
	time.Sleep(terminationTimeout)
	Logger("", "CR", "Application shutdown timeout. Force termination.")
	go Exit()

	time.Sleep(killingTimeout)
	Logger("", "CR", "Application termination timeout. Force killing.")
	os.Exit(exitCode)
}

// StopApp -- set exit code and raise application stop
func StopApp(code int) {
	if atomic.AddInt32(&appStarted, -1) == 0 {
		Logger("", "DE", "Set application exit code %d", code)

		exitCode = code

		exitCond.Broadcast()

		ex := false
		for !ex {
			select {
			case sleepInterrupted <- true:
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
	exitCond.L.Lock()
	exitCond.Wait()
	exitCond.L.Unlock()
}

// WaitingForStopChan --
func WaitingForStopChan() <-chan time.Time {
	c := make(chan time.Time)

	go func() {
		WaitingForStop()
		c <- NowUTC()
	}()

	return c
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
func AddExitFunc(name string, f ExitFunc, param any) {
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
func simpleLogger(facility string, level string, message string, params ...any) {
	fmt.Printf(level+": "+message+EOS, params...)
}

func init() {
	startTime = NowUTC()

	Logger = simpleLogger

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
	return strings.HasPrefix(appExecName, "__debug") // simple workaround for the VS Code
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
	case <-sleepInterrupted:
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
func (m *Messages) Add(msg string, params ...any) {
	if msg != "" {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.s = append(m.s, fmt.Sprintf(msg, params...))
	}
}

// AddError --
func (m *Messages) AddError(err error) {
	if err != nil {
		m.Add("%s", err.Error())
	}
}

func (m *Messages) Content() []string {
	return m.s
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
func UnsafeByteSlice2String(b []byte) (s string) {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// UnsafeString2ByteSlice -- fast convert string to []byte without memory allocs
// Don't forget to use runtime.KeepAlive(s) in the caller if necessary!
// Don't try to change result without thinking hard before that!
func UnsafeString2ByteSlice(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
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

func StructTags(s any, fields []string, tag string) (names []string, err error) {
	return structTags("", s, fields, tag)
}

func structTags(prefix string, s any, fields []string, tag string) (names []string, err error) {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		err = fmt.Errorf("%T is not a struct or pointer to a struct", s)
		return
	}

	if prefix != "" {
		prefix += "."
	}

	names = make([]string, 0, len(fields))
	msgs := NewMessages()

	for _, f := range fields {
		ff := strings.Split(f, ".")

		sf, exists := t.FieldByName(ff[0])
		if !exists {
			msgs.Add(`unknown field "%s"`, f)
			continue
		}

		if !sf.IsExported() {
			msgs.Add(`field "%s" is not exported`, f)
			continue
		}

		name := StructTagName(&sf, tag)

		if len(ff) == 1 {
			// simple type
			if name == "" {
				msgs.Add(`field "%s" has empty value of "%s" tag`, f, tag)
				continue
			}

			names = append(names, prefix+name)
			continue
		}

		// struct expected

		fn := strings.Join(ff[1:], ".")

		fTp := sf.Type
		if fTp.Kind() == reflect.Pointer {
			fTp = fTp.Elem()
		}

		if fTp.Kind() != reflect.Struct {
			msgs.Add("field %s is not a struct or pointer to a struct", fn)
			continue
		}

		var subNames []string
		subNames, err = structTags(
			prefix+name,
			reflect.New(fTp).Interface(),
			[]string{fn},
			tag,
		)
		if err != nil {
			msgs.AddError(err)
			continue
		}

		names = append(names, subNames...)
	}

	err = msgs.Error()
	return
}

func StructTagName(f *reflect.StructField, tag string) (name string) {
	name, ok := f.Tag.Lookup(tag)
	if !ok {
		name = f.Name
		return
	}

	name = strings.TrimSpace(strings.Split(name, ",")[0])
	return
}

func StructTagOpts(f *reflect.StructField, tag string) (opts StringMap) {
	opts = make(StringMap, 8)

	tags, ok := f.Tag.Lookup(tag)
	if !ok {
		return
	}

	list := strings.Split(tags, ",")

	for i := 0; i < len(list); i++ {
		opt := strings.TrimSpace(list[i])
		v := ""
		if i == 0 {
			opt, v = "", opt
		} else {
			sp := strings.Split(opt, "=")
			if len(sp) > 1 {
				opt = strings.TrimSpace(sp[0])
				v = strings.TrimSpace(sp[1])
			}
		}
		opts[opt] = v
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func SplitAndTrim(src string, delimiter string) (dst []string) {
	src = strings.TrimSpace(src)
	if len(src) == 0 {
		return
	}

	dst = strings.Split(src, delimiter)
	dstI := 0

	for srcI, s := range dst {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		if dstI != srcI {
			dst[dstI] = s
		}
		dstI++
	}

	if dstI != len(dst) {
		dst = dst[:dstI]
	}

	return
}

//----------------------------------------------------------------------------------------------------------------------------//

func IsNil(obj any) bool {
	return obj == nil ||
		(reflect.ValueOf(obj).Kind() == reflect.Pointer && reflect.ValueOf(obj).IsNil())
}

//----------------------------------------------------------------------------------------------------------------------------//
