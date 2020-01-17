package gel

import (
	"errors"
	"fmt"
	"github.com/xkortex/vprint"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var ignoreDirs = map[string]bool{
	".git":         true,
	".dvc":         true,
	".gel":         true,
	"node_modules": true,
}

type pathInfo struct {
	path string
	info os.FileInfo
}

type walkResults struct {
	nodes []BasicFileNode
}

type parseResults struct {
	node Rmap
	err  error
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

func walkFiles(done <-chan struct{}, root string) (<-chan pathInfo, <-chan error) {
	paths := make(chan pathInfo)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				ok := ignoreDirs[info.Name()]
				if ok {
					fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
					return filepath.SkipDir
				}
			}
			select {
			case paths <- pathInfo{path, info}:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}

// Performs preliminary "parse" based on basic file attributes
func basicFileParser(done <-chan struct{}, infos <-chan pathInfo, out chan<- BasicFileNode) {
	for pi := range infos {
		//fmt.Println(pi.path)
		select {
		case out <- File2basicNode(pi.path, pi.info):
		case <-done:
			vprint.Print(">done<>basicParser\n")
			return
		}
	}
}

var jobber int = 0
var count int64 = 0

// Do something interesting with our file nodes
func nodeParser(done <-chan struct{}, nodes <-chan BasicFileNode, out chan<- parseResults) {
	myjobber := jobber
	jobber++
	for node := range nodes {
		select {
		//case out <- BasicNode2Rmap(&node):
		case out <- func() (parseResults) {
			res := parseResults{
				node: BasicNode2Rmap(&node),
				err: nil,
			}
			fmt.Printf("%4d> %s \n", myjobber, res.node)
			count++
			return res
		}():
		case <-done:
			vprint.Print(">done<>nodeParser\n")
			return
		}
	}
	//vprint.Print("<><>nodeParser loop broken")
}

// MD5All reads all the files in the file tree rooted at root and returns a map
// from file path to the MD5 sum of the file's contents.  If the directory walk
// fails or any read operation fails, MD5All returns an error.  In that case,
// MD5All does not wait for inflight read operations to complete.
func ProcessFileTree(root string) (int, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	done2 := make(chan struct{})
	defer close(done)
	defer close(done2)

	paths, errc := walkFiles(done, root)

	// Start a fixed number of goroutines to read and digest files.
	chanBasicFileNode := make(chan BasicFileNode) // HLc
	chanRmapNode := make(chan parseResults)
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	const num_basicParsers = 20
	const num_postParsers = 20
	wg.Add(num_basicParsers)
	wg2.Add(num_postParsers)
	for i := 0; i < num_basicParsers; i++ {
		go func() {
			basicFileParser(done, paths, chanBasicFileNode) // HLc
			wg.Done()
		}()
	}
	for i := 0; i < num_postParsers; i++ {
		go func() {
			nodeParser(done2, chanBasicFileNode, chanRmapNode) // HLc
			wg2.Done()
		}()
	}
	go func() {
		wg.Wait()
		vprint.Print("Closing basic node chan\n")
		close(chanBasicFileNode) // HLc
	}()
	go func() {
		vprint.Printf("Waiting on postProc")
		wg2.Wait()
		vprint.Print("Closing postProc chan\n")
		close(chanRmapNode) // HLc
	}()
	// End of pipeline. OMIT
	fmt.Println("End of pipeline")

	// aggregate results here
	count := 0
	var nodes []Rmap
	for r := range chanRmapNode {
		if r.err != nil {
			return -1, r.err
		}
		nodes = append(nodes, r.node)
		count++
	}
	// Check whether the Walk failed.
	if err := <-errc; err != nil { // HLerrc
		return -1, err
	}
	vprint.Print("End of ProcessFileTree \n")
	fmt.Print("%d Results: \n%s\n", len(nodes), nodes[:5])
	return count, nil
}
