/*
Copyright © 2019 MICHAEL McDERMOTT

*/
package cmd

import (
	"fmt"
	"github.com/Wessie/appdirs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xkortex/magellan/gel"
	"github.com/xkortex/vprint"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

var (
	cfgFile       string
	developer     string
	defaultCfgDir string
)

const defaultCfgName = "gel.yml"

var ignoreDirs = map[string]bool {
	".git": true,
	".dvc": true,
	".gel": true,
	"node_modules": true,
}

func elapsed(what string) func() {
    start := time.Now()
    return func() {
        fmt.Printf("%s took %v\n", what, time.Since(start))
    }
}

func do_walk(root string) int {
	defer elapsed("\n<>walk<>\n")()
	var count int = 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			ok := ignoreDirs[info.Name()]
			if ok {
				fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
				return filepath.SkipDir
			}
		}
		count += 1

		fnode := gel.File2basicNode(path, info)
		fmt.Println(gel.BasicNode2Rmap(&fnode))
		//fmt.Printf("v: %q\n", filepath.Join(root, path))
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", root, err)
		return -1
	}
	return count
}

func dumb_tests(target string) {
	hostname, _ := os.Hostname()
	gopath := os.Getenv("GOPATH")
	hosturi, _ := url.Parse("file://" + hostname)
	gouri, _ := url.Parse("file://" + gopath)
	here, _ := filepath.Abs(target)
	hosturi.Path = path.Join(hosturi.Path, here)
	relpath, _ := filepath.Rel(gopath, here)
	//reluri, _ := filepath.Rel(gouri.Path, here)

	fmt.Printf("Abs path: %s\n", here)
	fmt.Printf("Fqn uri : %s\n", hosturi)
	fmt.Println(hosturi.Scheme, hosturi.Host, hosturi.Port(), hosturi.Path)
	hosturi.Host = ""
	fmt.Printf("Go  path: %s\n", relpath)
	fmt.Printf("Go  uri : %s\n", gouri)
	fmt.Printf("Here uri: %s\n", hosturi)

	fmt.Printf("uri path: %s\n", hosturi.Path)
	num_files := do_walk(target)
	fmt.Println("Scanned %d files\n", num_files)
}

// RootCmd represents the root command
var RootCmd = &cobra.Command{
	Use:   "gel",
	Short: "Walk and analyze a file tree",
	Long: `Gel will walk a file tree and parse files into an ontology 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(os.Stdout)

		vprint.Print("root called")
		vprint.Print(args)
		target := args[0]
		timeout, _ := cmd.PersistentFlags().GetFloat64("timeout")
		vprint.Print(target)
		vprint.Print(timeout)
		//if err := cmd.Usage(); err != nil {
		//	log.Fatalf("Error executing root command: %v", err)
		//}
		//log.Fatal("<dbg> silence/usage: ", cmd.SilenceErrors, cmd.SilenceUsage)
		//out := do_walk(target)
		//fmt.Printf("%d files/dirs\n", out)
		dumb_tests(target)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	defaultCfgDir = appdirs.UserConfigDir("gel", "", "", false)
	defaultCfgFile := filepath.Join(defaultCfgDir, "config.yml")
	//RootCmd.AddCommand(RootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// RootCmd.PersistentFlags().String("foo", "", "A help for foo")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c",
		defaultCfgFile,
		"config file, based in UserConfigDir", )

	RootCmd.PersistentFlags().Float64P("timeout", "t", 0.1, "Timeout in seconds")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	RootCmd.PersistentFlags().BoolP("silent", "s", false, "Suppress errors")
	RootCmd.PersistentFlags().BoolP("stdin", "-", false, "Read from standard in")
	RootCmd.Flags().BoolP("verbose", "v", false, "Verbose tracing (in progress)")
	RootCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")

}

func initConfig() {

}
