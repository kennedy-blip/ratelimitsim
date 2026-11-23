package handler
import("math";"sync";"time")
type Policy struct{ID string `json:"id"`;Name string `json:"name"`;RoutePattern string `json:"routePattern"`;Capacity float64 `json:"capacity"`;RefillPerSec float64 `json:"refillPerSec"`;CostPerReq float64 `json:"costPerReq"`}
type Bucket struct{Tokens float64 `json:"tokens"`;Capacity float64 `json:"capacity"`;LastRefillAt time.Time `json:"lastRefillAt"`}
type BucketView struct{ClientKey string `json:"clientKey"`;Tokens float64 `json:"tokens"`;Capacity float64 `json:"capacity"`;FillPct float64 `json:"fillPct"`}
type PolicyView struct{Policy;Buckets []BucketView `json:"buckets"`}
type Store struct{mu sync.Mutex;policies map[string]Policy;buckets map[string]*Bucket}
var appStore=newStore()
func newStore()*Store{s:=&Store{policies:map[string]Policy{},buckets:map[string]*Bucket{}};for _,p:=range []Policy{{"payments-stk-api","Payments STK API","/api/payments/stk/*",20,4,1},{"public-search-v1","Public Search v1","/api/search/*",50,10,1},{"heavy-report-export","Heavy Report Export","/api/reports/export",5,.5,3}}{s.policies[p.ID]=p};return s}
func bucketID(id,key string)string{return id+"::"+key}
func(s *Store)bucketLocked(p Policy,key string,now time.Time)*Bucket{id:=bucketID(p.ID,key);b:=s.buckets[id];if b==nil{b=&Bucket{Tokens:p.Capacity,Capacity:p.Capacity,LastRefillAt:now};s.buckets[id]=b};return b}
func refill(b *Bucket,p Policy,now time.Time){if now.After(b.LastRefillAt){b.Tokens=math.Min(p.Capacity,b.Tokens+now.Sub(b.LastRefillAt).Seconds()*p.RefillPerSec);b.LastRefillAt=now};b.Capacity=p.Capacity;if b.Tokens>p.Capacity{b.Tokens=p.Capacity}}
func(s *Store)viewsLocked(p Policy)[]BucketView{out:=[]BucketView{};prefix:=p.ID+"::";for k,b:=range s.buckets{if len(k)>=len(prefix)&&k[:len(prefix)]==prefix{pct:=0.0;if p.Capacity>0{pct=b.Tokens/p.Capacity*100};out=append(out,BucketView{k[len(prefix):],math.Round(b.Tokens*100)/100,p.Capacity,math.Round(pct*10)/10})}};return out}
