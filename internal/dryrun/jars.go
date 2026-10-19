package dryrun

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// manifest reads META-INF/MANIFEST.MF from a plugin file, joining the
// continuation lines the format wraps at 72 bytes.
func manifest(path string) (map[string]string, error) {
	z, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer z.Close()
	for _, f := range z.File {
		if f.Name != "META-INF/MANIFEST.MF" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return parseManifest(rc), nil
	}
	return map[string]string{}, nil
}

func parseManifest(r io.Reader) map[string]string {
	out := map[string]string{}
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 4*1024*1024)
	var key string
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		if strings.HasPrefix(line, " ") && key != "" {
			out[key] += line[1:]
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			key = ""
			continue
		}
		key = strings.TrimSpace(k)
		out[key] = strings.TrimSpace(v)
	}
	return out
}

// minimumJava is the highest Minimum-Java-Version among the plugin files in
// dir, 0 if none declares one.
func minimumJava(dir string) (int, string) {
	entries, _ := os.ReadDir(dir)
	best, who := 0, ""
	for _, e := range entries {
		m, err := manifest(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		v := m["Minimum-Java-Version"]
		if v == "" {
			continue
		}
		if strings.HasPrefix(v, "1.") {
			v = v[2:]
		}
		n, err := strconv.Atoi(strings.SplitN(v, ".", 2)[0])
		if err == nil && n > best {
			best, who = n, strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		}
	}
	return best, who
}

// classIndex maps a class file path (org/example/Foo.class) to the plugin
// that ships it, looking inside each plugin's WEB-INF/lib jars.
type classIndex map[string]string

func buildClassIndex(dir string) classIndex {
	idx := classIndex{}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		name := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		z, err := zip.OpenReader(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		for _, f := range z.File {
			switch {
			case strings.HasPrefix(f.Name, "WEB-INF/classes/") && strings.HasSuffix(f.Name, ".class"):
				idx[strings.TrimPrefix(f.Name, "WEB-INF/classes/")] = name
			case strings.HasPrefix(f.Name, "WEB-INF/lib/") && strings.HasSuffix(f.Name, ".jar"):
				rc, err := f.Open()
				if err != nil {
					continue
				}
				b, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					continue
				}
				inner, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
				if err != nil {
					continue
				}
				for _, g := range inner.File {
					if strings.HasSuffix(g.Name, ".class") {
						if _, ok := idx[g.Name]; !ok {
							idx[g.Name] = name
						}
					}
				}
			}
		}
		z.Close()
	}
	return idx
}

// owner returns the plugin that ships the class, or "".
func (idx classIndex) owner(class string) string {
	class = strings.TrimSuffix(class, ".class")
	// Inner classes live in their own file next to the outer one.
	return idx[strings.ReplaceAll(class, ".", "/")+".class"]
}
