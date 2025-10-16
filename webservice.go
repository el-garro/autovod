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
<html>
<head>
  <meta charset="utf-8">
  <title>Twitch Auto VOD</title>
  <style>
    body { font-family: Arial, sans-serif; margin: 2em; }
    h1 { color: #333; }
    table { width: 100%; border-collapse: collapse; }
    table, th, td { border: 1px solid #ccc; }
    colgroup col:first-child { width: 60%; }
    colgroup col:nth-child(2),
    colgroup col:nth-child(3) { width: 20%; }
    th, td { padding: 0.5em 1em; }
    th { background: #f5f5f5; text-align: left; }
    td.size, td.date { text-align: right; }
    td.name a { text-decoration: none; color: #1a0dab; }
    td.name a:hover { text-decoration: underline; }
    .emoji { margin-right: 0.5em; }
  </style>
</head>
<body>
  <h1>Available VODs</h1>
  <table>
    <colgroup>
      <col>
      <col>
      <col>
    </colgroup>
    <thead>
      <tr>
        <th>Name</th>
        <th class="size">Size</th>
        <th class="date">Date</th>
      </tr>
    </thead>
    <tbody>
      {{- range .Files }}
      <tr>
        <td class="name">
          <span class="emoji">📄</span>
          <a href="/download/{{ .Name }}" download>{{ .Name }}</a>
        </td>
        <td class="size">{{ .Size }}</td>
        <td class="date">{{ .ModTime.Format "2006-01-02 15:04:05" }}</td>
      </tr>
      {{- else }}
      <tr><td colspan="3">No files found</td></tr>
      {{- end }}
    </tbody>
  </table>
  <p><strong>Running for:</strong> {{ .ElapsedTime }}</p>
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
