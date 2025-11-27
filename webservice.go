package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/log"
)

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Twitch Auto VOD</title>
  <link rel="stylesheet" href="https://unpkg.com/@picocss/pico@latest/css/pico.min.css">
</head>
<body>
  <main class="container">
    <h1>Available VODs</h1>

    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
          <th>Date</th>
        </tr>
      </thead>
      <tbody>
        {{- range .Files }}
        <tr>
          <td>
            📄 <a href="/download/{{ .Name }}" download>{{ .Name }}</a>
          </td>
          <td>{{ .Size }}</td>
          <td>{{ .ModTime.Format "2006-01-02 15:04:05" }}</td>
        </tr>
        {{- else }}
        <tr><td colspan="3">No files found</td></tr>
        {{- end }}
      </tbody>
    </table>

    <p><strong>Running for:</strong> {{ .ElapsedTime }}</p>
  </main>
</body>
</html>`

func WebService() {
	logger := log.NewWithOptions(
		os.Stderr,
		log.Options{
			Level:           log.InfoLevel,
			Prefix:          "WEB",
			ReportTimestamp: true,
		},
	)

	startTime := time.Now()

	type PageData struct {
		Files       []FileInfo
		ElapsedTime time.Duration
	}

	tmpl := template.Must(template.New("index").Parse(indexHTML))
	os.Mkdir(DOWNLOAD_DIR, os.ModePerm)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files, err := GetFileList(DOWNLOAD_DIR)
		if err != nil {
			http.Error(w, "Cannot read directory", http.StatusInternalServerError)
			return
		}

		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime.After(files[j].ModTime)
		})

		page := PageData{
			Files:       files,
			ElapsedTime: time.Since(startTime).Truncate(time.Second),
		}

		if err := tmpl.Execute(w, page); err != nil {
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		file := r.URL.Path[len("/download/"):]
		fp := filepath.Join(DOWNLOAD_DIR, file)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file))
		http.ServeFile(w, r, fp)
	})

	logger.Info("Service started", "url", fmt.Sprintf("http://localhost:%d", Config.WebPort))
	logger.Fatal("Crashed", "err", http.ListenAndServe(fmt.Sprintf(":%d", Config.WebPort), nil))
}
