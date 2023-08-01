//go:build enable_launchpad
// +build enable_launchpad

package launchpad

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/database"
)

// type responseCaptureWriter struct {
// 	http.ResponseWriter
// 	ContentType string
// 	Body        *bytes.Buffer
// }

// func (w *responseCaptureWriter) Write(p []byte) (int, error) {
// 	if w.Body == nil {
// 		w.Body = bytes.NewBuffer([]byte{})
// 	}
// 	return w.Body.Write(p)
// }

// func (w *responseCaptureWriter) WriteHeader(statusCode int) {
// 	w.ResponseWriter.WriteHeader(statusCode)
// }

// func (w *responseCaptureWriter) Header() http.Header {
// 	return w.ResponseWriter.Header()
// }

// func LaunchpadRequestMiddleware(next http.Handler) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		cookie, err := r.Cookie("LaunchpadUser")

// 		if err != nil {
// 			fmt.Println("not a launchpad request")
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		launchpadPersona := cookie.Value

// 		// Wrap the original ResponseWriter to capture the rendered template
// 		captureWriter := &responseCaptureWriter{ResponseWriter: w}

// 		// Call the next handler in the chain
// 		next.ServeHTTP(captureWriter, r)

// 		if w.Header().Get("Content-Type") == "text/html" || w.Header().Get("Content-Type") == "text/html; charset=utf-8" {
// 			customHTML := fmt.Sprintf(`<div id="launchpad-persona-key" style="padding: 0.5rem 1rem 0.5rem 1rem; border-radius: 0.5rem; position: fixed; top: 0; right: 0; margin: 1rem; display: flex; gap: 1rem; background-color: #93b1ed; align-items: center;box-shadow: 0 1px 2px 0 rgb(0 0 0 / 0.25);">
// 						<p style="margin: 0; color: #04163b;">%s</p>
// 						<a href="/launchpad" style="color: #04163b">Switch Personas</a>
// 						<button style="background-color: transparent; border: 0; font-size: 1rem; cursor: pointer;" onclick="(() => { document.querySelector('#launchpad-persona-key').remove() })()">&times;</button>
// 					</div>`, launchpadPersona)

// 			tmpl := template.Must(template.New("custom").Parse(`{{.}}`))

// 			tmpl.Execute(w, template.HTML(append(captureWriter.Body.Bytes(), []byte(customHTML)...)))
// 		} else {
// 			// Set the status code and copy the captured headers
// 			w.WriteHeader(captureWriter.ResponseWriter.)
// 			for key, values := range captureWriter.Header() {
// 				for _, value := range values {
// 					w.Header().Add(key, value)
// 				}
// 			}

// 			// In case the content type is not "text/html", simply write the captured response back to the original ResponseWriter (w)
// 			w.Write(captureWriter.Body.Bytes())
// 		}
// 	})
// }

func combineMaps(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range map1 {
		result[key] = value
	}

	for key, value := range map2 {
		result[key] = value
	}

	return result
}

func BootstrapUsers(db database.DatabaseAccessor, accessToken string, importUsers map[string]LocksmithLaunchpadUserOptions) {
	password, _ := authentication.CompileLocksmithPassword(accessToken)

	for username, opts := range importUsers {
		_, found := db.FindOne("users", map[string]interface{}{
			"username": username,
		})

		if found {
			fmt.Printf("Launchpad user %s already registered.\n", username)
			continue
		}

		insert := map[string]interface{}{
			"id":          uuid.New().String(),
			"username":    username,
			"password":    password,
			"email":       opts.Email,
			"sessions":    []interface{}{},
			"websessions": []interface{}{},
			"role":        opts.Role,
		}

		if opts.Custom != nil {
			insert = combineMaps(insert, opts.Custom)
		}
		_, err := db.InsertOne("users", insert)

		if err != nil {
			fmt.Println(err)
		}
	}

}
