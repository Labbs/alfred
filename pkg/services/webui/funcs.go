package webui

import (
	"encoding/json"
	"strings"

	"github.com/labbs/alfred/pkg/config"
	b "github.com/labbs/alfred/pkg/services/bookmark"
	dash "github.com/labbs/alfred/pkg/services/dashboard"
)

func currentVersion() string {
	return config.Version
}

func add(i int) int {
	return i + 1
}

func truncate(s string, i int) string {
	if len(s) > i {
		return s[:i] + "..."
	}
	return s
}

func truncateByWord(s string, i int) string {
	ss := strings.Split(s, " ")
	if len(s) > i {
		return strings.Join(ss[:i], " ") + "..."
	}
	return strings.Join(ss, " ")
}

func joinTags(tags []b.Tag) string {
	var t []string
	for _, tag := range tags {
		t = append(t, tag.Name)
	}
	return strings.Join(t, ",")
}

func toJson(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func widgetConfigTransform(widget dash.Widget) map[string]interface{} {
	return map[string]interface{}{
		"name": widget.Name,
		"h":    widget.H,
		"w":    widget.W,
		"html": widget.HTML,
		"css":  widget.CSS,
		"js":   widget.JS,
	}
}
