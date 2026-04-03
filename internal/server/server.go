package server
import ("encoding/json";"log";"net/http";"strings";"github.com/stockyard-dev/stockyard-crossroads/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/links",s.list);s.mux.HandleFunc("POST /api/links",s.create);s.mux.HandleFunc("GET /api/links/{id}",s.get);s.mux.HandleFunc("DELETE /api/links/{id}",s.del)
s.mux.HandleFunc("GET /api/links/{id}/clicks",s.clicks)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard)
s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){
path:=r.URL.Path
if path!="/"&&!strings.HasPrefix(path,"/api/")&&!strings.HasPrefix(path,"/ui"){
slug:=strings.TrimPrefix(path,"/")
if slug!=""&&!strings.Contains(slug,"/"){l:=s.db.GetBySlug(slug);if l!=nil{s.db.RecordClick(slug,r.RemoteAddr,r.UserAgent(),r.Referer());http.Redirect(w,r,l.URL,302);return}}}
s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)list(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"links":oe(s.db.List())})}
func(s *Server)create(w http.ResponseWriter,r *http.Request){var l store.Link;json.NewDecoder(r.Body).Decode(&l);if l.URL==""{we(w,400,"url required");return};s.db.Create(&l);wj(w,201,s.db.GetByID(l.ID))}
func(s *Server)get(w http.ResponseWriter,r *http.Request){l:=s.db.GetByID(r.PathValue("id"));if l==nil{we(w,404,"not found");return};wj(w,200,l)}
func(s *Server)del(w http.ResponseWriter,r *http.Request){s.db.Delete(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)clicks(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"clicks":oe(s.db.ClickHistory(r.PathValue("id"),50))})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"crossroads","links":st.Links,"clicks":st.TotalClicks})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
