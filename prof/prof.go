package prof

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
)

type ProfileType string

const (
	GoroutineProfile  ProfileType = "goroutine"
	ThreadcreateProfile  ProfileType = "threadcreate"
	HeapProfile  ProfileType = "heap"
	BlockProfile  ProfileType = "block"
)

var (
	memProfile = flag.String("memprofile", "", "write a memory profile to the named file after execution")
	memProfileRate = flag.Int("memprofilerate", 0, "if > 0, sets runtime.MemProfileRate")
	cpuProfile = flag.String("cpuprofile", "", "write a cpu profile to the named file during execution")
	blockProfile = flag.String("blockprofile", "", "write a goroutine blocking profile to the named file after execution")
	blockProfileRate = flag.Int("blockprofilerate", 1, "if > 0 , calls runtime.SetBlockProfileRate()")
)

var running bool
var lock *sync.Mutex = new(sync.Mutex)

func parseProfFlags() {
	if !flag.Parsed() {
		flag.Parse()
	}
	*cpuProfile = getAbsFilePath(*cpuProfile)
	*blockProfile = getAbsFilePath(*blockProfile)
	*memProfile = getAbsFilePath(*memProfile)
}

func getAbsFilePath(path string) string {
	if path == "" {
		return ""
	}
	path = filepath.FromSlash(path)
	if !filepath.IsAbs(path) {
		baseDir , err := os.Getwd()
		if err != nil {
			panic(errors.New(fmt.Sprintf("Can not get current work dir : %s\n", err)))
		}
		path = filepath.Join(baseDir, path)
	}
	return path
}

func Start() {
	lock.Lock()
	defer lock.Unlock()
	parseProfFlags()
	startBlockProfile()
	startCpuProfile()
	startMemProfile()
	running = true
}

func startBlockProfile() {
	if *blockProfile != "" && *blockProfileRate > 0 {
		runtime.SetBlockProfileRate(*blockProfileRate)
	}
}

func startCpuProfile() {
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not create cpu profile outout file: %s\n", err)
			return
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Can not start cpu profile: %s\n", err)
			f.Close()
			return
		}
	}
}

func startMemProfile() {
	if *memProfile != "" && *memProfileRate > 0 {
		runtime.MemProfileRate = *memProfileRate
	}
}

func Stop() {
	lock.Lock()
	defer lock.Unlock()
	stopBlockProfile()
	stopCpuProfile()
	stopMemProfile()
	running = false
}

func stopBlockProfile() {
	if *blockProfile != "" && *blockProfileRate >= 0 {
		f, err := os.Create(*blockProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not create block profile output file: %s\n", err)
			return
		}
		if err = pprof.Lookup("block").WriteTo(f, 0); err != nil {
			fmt.Fprintf(os.Stderr, "Can not write %s: %s\n", *blockProfile, err)
		}
		f.Close()
	}
}

func stopCpuProfile() {
	if *cpuProfile != "" {
		pprof.StopCPUProfile()
	}
}

func stopMemProfile() {
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not create mem profile output file : %s\n", err)
			return
		}
		if err = pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Can not write %s: %s\n", *memProfile, err)
		}
		f.Close()
	}
}

func SaveProfile(workDir string, profileName string, ptype ProfileType, debug int) {
	absWorkDir := getAbsFilePath(workDir)
	if profileName == "" {
		profileName = string(ptype)
	}
	profileName += ".out"
	profilePath := filepath.Join(absWorkDir, profileName)
	f, err := os.Create(profilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create profile output file: %s\n", err)
		return
	}
	if err = pprof.Lookup(string(ptype)).WriteTo(f, debug); err != nil {
		fmt.Fprintf(os.Stderr, "Can not write %s: %s\n", profilePath, err)
	}
	f.Close()
}
