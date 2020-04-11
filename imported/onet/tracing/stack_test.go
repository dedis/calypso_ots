package tracing

//func TestCallers(t *testing.T) {
//	wg := sync.WaitGroup{}
//	pm := sync.Mutex{}
//	for i := 0; i < 10; i++ {
//		wg.Add(1)
//		go func() {
//			var strs []string
//			strs = append(strs, printCallers()...)
//			strs = append(strs, printCallers()...)
//			strs = append(strs, printCallers()...)
//			pm.Lock()
//			fmt.Println(strs)
//			pm.Unlock()
//			wg.Done()
//		}()
//	}
//	wg.Wait()
//}
//
//func TestGor(t *testing.T) {
//	log.SetDebugVisible(3)
//	go tg()
//	go tg()
//	time.Sleep(time.Second)
//}
//
//var hm = sync.Mutex{}
//
//func tg() {
//	hm.Lock()
//	log.Lvl1("gor:\n", printstacks())
//	hm.Unlock()
//	wg := sync.WaitGroup{}
//	wg.Add(2)
//	go func() {
//		hm.Lock()
//		log.Lvl2("gor:\n", printstacks())
//		hm.Unlock()
//		go func() {
//			hm.Lock()
//			log.Lvl3("gor:\n", printstacks())
//			hm.Unlock()
//			wg.Done()
//		}()
//		//go func() {
//		//	log.Lvl3(log.Stack())
//		//	wg.Done()
//		//}()
//		log.Lvl1("gor")
//	}()
//	wg.Wait()
//}
//
//func printhi() string {
//	return "hi"
//}
//
//func printstacks() string {
//	return string(Stack())
//}
//
//func printcallers() string {
//	ptr := make([]uintptr, 100)
//	ptr = ptr[0:runtime.Callers(0, ptr)]
//	frames := runtime.CallersFrames(ptr)
//	var str string
//	for {
//		fr, ok := frames.Next()
//		if !ok {
//			break
//		}
//		str += fmt.Sprintf("%+v\n", fr)
//	}
//	return str
//}
//
//func printgor() string {
//	buf := bytes.Buffer{}
//	pprof.Lookup("goroutine").WriteTo(&buf, 1)
//	var str string
//	for _, l := range strings.Split(buf.String(), "\n") {
//		str += fmt.Sprintf("   %s\n", l)
//	}
//	return str
//}
//
//func printheap() string {
//	buf := bytes.Buffer{}
//	pprof.Lookup("heap").WriteTo(&buf, 2)
//	var str string
//	for _, l := range strings.Split(buf.String(), "\n") {
//		str += fmt.Sprintf("   %s\n", l)
//	}
//	return str
//}
//
//func Stack() []byte {
//	buf := make([]byte, 1024)
//	for {
//		n := runtime.Stack(buf, true)
//		if n < len(buf) {
//			return buf[:n]
//		}
//		buf = make([]byte, 2*len(buf))
//	}
//}
//
//func printCallers() (ret []string) {
//	//cs := make([]uintptr, 100)
//	//cs = cs[:runtime.Callers(0, cs)]
//	//for _, c := range cs {
//	//	fmt.Printf("%x - ", c)
//	//}
//	gp := make([]runtime.StackRecord, 100)
//	n, _ := runtime.GoroutineProfile(gp)
//	gp = gp[:n]
//	for _, g := range gp {
//		ret = append(ret, fmt.Sprintln(g.Stack0))
//	}
//	ret = append(ret, "\n")
//	return
//}
//
//func TestOldGoroutines(t *testing.T) {
//	sc, tr := newSimulLogger()
//	//tr.PrintSingleSpans = 10
//	tr.AddEntryPoints("go.dedis.ch/onet/v3/tracing.gor")
//	tr.AddDoneMsgs("done sub")
//	wg := sync.WaitGroup{}
//	wg.Add(2)
//	goroutines := 2
//	traces = make([][][]string, goroutines)
//	for i := 0; i < goroutines; i++ {
//		log.Lvl1("launching new go-routine")
//		go func(i int) {
//			gor(i)
//			wg.Done()
//		}(i)
//	}
//	log.Lvl1("waiting on gors")
//	time.Sleep(time.Second * 2)
//	for i, ti := range traces {
//		for j, tj := range ti {
//			for _, k := range tj {
//				fmt.Println(i, j, k)
//			}
//			fmt.Println()
//		}
//	}
//	wg.Wait()
//	sc.waitAndPrint()
//}
//
//func gor(i int) {
//	//log.Lvl1("new go-routine", i)
//	traces[i] = make([][]string, 5)
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	getTrace(i, 0)
//	go func() {
//		//log.Lvl1("new sub-go-routine", i)
//		getTrace(i, 2)
//		time.Sleep(time.Millisecond * 10)
//		getTrace(i, 3)
//		time.Sleep(time.Millisecond * 20)
//		getTrace(i, 4)
//		//pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
//		//fmt.Println(log.Stack())
//		wg.Done()
//	}()
//	getTrace(i, 1)
//	//log.Lvl1("waiting for sub-goroutine")
//	wg.Wait()
//	//log.Lvl1("done sub-goroutine")
//}
//
//var traces [][][]string
//
//var regStackTrace = regexp.MustCompile("^\\d+")
//var regTracing = regexp.MustCompile("tracing\\.(gor|getTrace)")
//var gtMutex = sync.Mutex{}
//
//func getTrace(i, j int) {
//	gtMutex.Lock()
//	defer gtMutex.Unlock()
//	buf := bytes.Buffer{}
//	pprof.Lookup("goroutine").WriteTo(&buf, 1)
//	fmt.Println("*********", i, j)
//	fmt.Println(buf.String())
//	//fmt.Println(log.Stack())
//	cs := make([]uintptr, 100)
//	cs = cs[:runtime.Callers(0, cs)]
//	for _, c := range cs {
//		fmt.Printf("%x - ", c)
//	}
//	fmt.Println()
//	return
//	stackTrace := ""
//	for _, line := range strings.Split(buf.String(), "\n") {
//		if regStackTrace.MatchString(line) {
//			stackTrace = line
//		}
//		if regTracing.MatchString(line) {
//			if stackTrace != "" {
//				traces[i][j] = append(traces[i][j], stackTrace)
//				//fmt.Println(i, j, stackTrace)
//				stackTrace = ""
//			}
//			//fmt.Println(i, j, line)
//			traces[i][j] = append(traces[i][j], line)
//		}
//	}
//	return
//}
//
