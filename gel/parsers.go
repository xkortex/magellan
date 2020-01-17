package gel

import (
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// todo: need onto for markups
// These are likely going to be in flux for some time
var KnownExtTypes = map[string]string {
	".csv": "nfo:Spreadsheet",
	".xls": "nfo:Spreadsheet",
	".xlsx": "nfo:Spreadsheet",
	".txt": "nfo:PlainTextDocument",
	".md": "nfo:PlainTextDocument",
	".html": "nfo:HtmlDocument",
	".go": "nfo:SourceCode",
	".py": "nfo:SourceCode",
	".js": "nfo:SourceCode",
	".c": "nfo:SourceCode",
	".c++": "nfo:SourceCode",
	".cpp": "nfo:SourceCode",
	".h": "nfo:SourceCode",
	".hpp": "nfo:SourceCode",
	".bash": "nfo:SourceCode",
	".sh": "nfo:SourceCode",
	".cmake": "nfo:SourceCode",
	".proto": "nfo:SourceCode",
}

var KnownTypes = map[string]string {
	"image": "nfo:Image",
	"video": "nfo:Video",
}

func parseType(ext string, mimetype string) string {
	atype, ok := KnownExtTypes[ext]
	if ok {
		return atype
	}
	mt0 := strings.Split(mimetype, "/")[0]

	atype, ok = KnownTypes[mt0]
	if ok {
		return atype
	}

	return "nfo:FileDataObject"
}

func File2ontology (fpath string, info os.FileInfo) Rmap {
	ext := filepath.Ext(fpath)
	mt := mime.TypeByExtension(filepath.Ext(fpath))

	var atype string
	if info.IsDir() {
		atype = "nfo:Folder" // might not be inode/directory on windows, idk
	} else {
		atype = parseType(ext, mt)
	}

	doc := Rmap {
		"@id": fpath,
		"@type": atype,
		"extension": ext,
		"mimetype": mt,
		"nfo:fileSize": strconv.FormatInt(info.Size(), 10),
	}

	return doc
}
