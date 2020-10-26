package favicon

import (
	"net/http"
)

// have your image in /opt/images/fav.png
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/opt/images/fav.png")
}
