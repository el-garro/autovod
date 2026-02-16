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
  <style>
    .size-col {
      white-space: nowrap;
    }
  </style>
</head>
<body>
  <main class="container">
    <h1>Available VODs</h1>

    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th class="size-col">Size</th>
          <th>Date</th>
        </tr>
      </thead>
      <tbody>
        {{- range .Files }}
        <tr>
          <td>
            📄 <a href="/download/{{ .Name }}" download>{{ .Name }}</a>
          </td>
          <td class="size-col">{{ .Size }}</td>
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

var startTime time.Time
var indexTemplate *template.Template

type pageData struct {
	Files       []FileInfo
	ElapsedTime time.Duration
}

type userCredentials struct {
	Username string
	Password string
}

func authMiddleware(handler http.HandlerFunc, creds userCredentials) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if creds.Username == "" && creds.Password == "" {
			handler(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || user != creds.Username || pass != creds.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := GetFileList(Config.DownloadDir)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime.After(files[j].ModTime)
	})

	page := pageData{
		Files:       files,
		ElapsedTime: time.Since(startTime).Truncate(time.Second),
	}

	if err := indexTemplate.Execute(w, page); err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/download/"):]
	fp := filepath.Join(Config.DownloadDir, file)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", file))
	http.ServeFile(w, r, fp)
}

func WebService() {
	logger := log.NewWithOptions(
		os.Stderr,
		log.Options{
			Level:           log.InfoLevel,
			Prefix:          "WEB",
			ReportTimestamp: true,
		},
	)

	startTime = time.Now()
	indexTemplate = template.Must(template.New("index").Parse(indexHTML))
	os.Mkdir(Config.DownloadDir, os.ModePerm)

	creds := userCredentials{
		Username: Config.HttpUser,
		Password: Config.HttpPassword,
	}

	http.HandleFunc("/", authMiddleware(indexHandler, creds))
	http.HandleFunc("/download/", downloadHandler)

	logger.Info("Service started", "url", fmt.Sprintf("http://localhost:%d", Config.WebPort))
	logger.Fatal("Crashed", "err", http.ListenAndServe(fmt.Sprintf(":%d", Config.WebPort), nil))
}
