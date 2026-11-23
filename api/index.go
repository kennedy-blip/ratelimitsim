package handler
import "net/http"
func Handler(w http.ResponseWriter,r *http.Request){
 switch r.URL.Path {
 case "/api/policies": policiesHandler(w,r)
 case "/api/simulate": simulateHandler(w,r)
 default: http.NotFound(w,r)
 }
}
