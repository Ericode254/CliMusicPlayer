package songs

import (
	"MusicPlayer/logger"
	"io/fs"
	"os"
	"path"
)

func ListFiles(dir string) []string {
	root := os.DirFS(dir)

	mpFiles, err := fs.Glob(root, "*.mp3")
	if err != nil {
		logger.Logger(err)
	}

	var files []string
	for _, v := range mpFiles {
		files = append(files, path.Join(dir, v))
	}

	return files
}
