package client

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (c *Client) CopyBackup() error {
	c.logf("Copying files from everquest_rof2...")
	// copy all files in everquest_rof2 to current path
	err := filepath.Walk("everquest_rof2", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		dst := strings.TrimPrefix(path, "everquest_rof2/")

		fi, err := os.Stat(dst)
		if err == nil {
			// check if file mod date is newer and file size is around same
			if fi.ModTime().After(info.ModTime()) && fi.Size() > info.Size()-100 && fi.Size() < info.Size()+100 {
				return nil
			}

		}

		r, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("open %s: %w", path, err)
		}
		defer r.Close()

		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
		}

		w, err := os.Create(dst)
		if err != nil {
			return fmt.Errorf("create %s: %w", dst, err)
		}
		defer w.Close()

		_, err = io.Copy(w, r)
		if err != nil {
			return fmt.Errorf("copy %s: %w", dst, err)
		}

		err = w.Sync()
		if err != nil {
			return fmt.Errorf("sync %s: %w", dst, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk: %w", err)
	}

	return nil
}