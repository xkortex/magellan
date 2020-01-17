package gel

import (
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Halfway between FileInfo and RDF node, can use this to pass around from
// walk to other parsing operations
// todo: maybe a mime type enum?
type BasicFileNode struct {
	A_id         string // @id
	A_type       string // @type
	Label        string // rdfs:Label
	MimeType     string //
	Extension    string // file extension
	BelongsTo    string // parent container/folder
	FileSize     int64  // file size in bytes,
	IsDir        bool
	TextOnly     bool // document is printable-text-only
	LinkedData   bool // document can/does contain parseable linked data structures
	LastModified time.Time
	err          error
}

// todo: need onto for markups
// These are likely going to be in flux for some time
var KnownExtTypes = map[string]string{
	".csv":   "nfo:Spreadsheet",
	".xls":   "nfo:Spreadsheet",
	".xlsx":  "nfo:Spreadsheet",
	".txt":   "nfo:PlainTextDocument",
	".md":    "nfo:PlainTextDocument",
	".html":  "nfo:HtmlDocument",
	".go":    "nfo:SourceCode",
	".py":    "nfo:SourceCode",
	".js":    "nfo:SourceCode",
	".c":     "nfo:SourceCode",
	".c++":   "nfo:SourceCode",
	".cpp":   "nfo:SourceCode",
	".h":     "nfo:SourceCode",
	".hpp":   "nfo:SourceCode",
	".bash":  "nfo:SourceCode",
	".sh":    "nfo:SourceCode",
	".cmake": "nfo:SourceCode",
	".proto": "nfo:SourceCode",
}

// maps mime type (sans subtype) to ontology type
var KnownTypes = map[string]string{
	"image": "nfo:Image",
	"video": "nfo:Video",
	"audio": "nfo:Audio",
}

// Determines if file is text-only
func isText(fnode *BasicFileNode) {
	mt0 := strings.Split(fnode.MimeType, "/")[0]
	if mt0 == "text" {
		fnode.TextOnly = true
	}

}

func parseType(fnode *BasicFileNode) bool {
	atype, ok := KnownExtTypes[fnode.Extension]
	if ok {
		fnode.A_type = atype
		return true
	}
	mt0 := strings.Split(fnode.MimeType, "/")[0]

	atype, ok = KnownTypes[mt0]
	if ok {
		fnode.A_type = atype
		return true
	}

	fnode.A_type = "nfo:FileDataObject"
	return false

}

func File2basicNode(fpath string, info os.FileInfo) (BasicFileNode, ) {
	fnode := &BasicFileNode{}
	fnode.A_id = fpath
	fnode.Label = filepath.Base(fpath)
	fnode.Extension = filepath.Ext(fpath)
	fnode.MimeType = mime.TypeByExtension(fnode.Extension)
	fnode.FileSize = info.Size()

	if info.IsDir() {
		fnode.A_type = "nfo:Folder" // might not be inode/directory on windows, idk
	} else {
		parseType(fnode)
	}

	return *fnode
}

func BasicNode2Rmap(fnode *BasicFileNode) Rmap {
	return Rmap{
		"@id":          fnode.A_id,
		"@type":        fnode.A_type,
		"rdfs:Label":   fnode.Label,
		"extension":    fnode.Extension,
		"mimetype":     fnode.MimeType,
		"nfo:fileSize": strconv.FormatInt(fnode.FileSize, 10),
	}
}
