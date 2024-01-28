package main

import "flag"

var (
	days     = flag.Int("d", 0, "Add files older than days")
	delete   = flag.Bool("r", false, "Delete files older than days")
	compress = flag.Bool("c", false, "Compress files")
	name     = flag.String("n", "", "Name of the archive")
	source   = flag.String("s", "", "Source directory")
)

func main() {
	flag.Parse()
	if *source == "" {
		flag.Usage()
		return
	}
	if *days == 0 {
		flag.Usage()
		return
	}
	if *name == "" {
		flag.Usage()
		return
	}
	archiver := NewArchiver(*source, *days, *delete, *compress, *name)
	err := archiver.Archive()
	if err != nil {
		println(err.Error())
	}
}
