package urlshort

import (
	"net/http"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	var m = func (w http.ResponseWriter,r *http.Request) {
		path := r.URL.Path
		if v,ok := pathsToUrls[path]; ok {
			http.Redirect(w,r,v,301)
		}else{
			fallback.ServeHTTP(w,r)
		}
	}
	return http.HandlerFunc(m)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// 下面的代码不完全对，思路没问题，没去仔细看yaml库
	m := make(map[string]string)
	err := yaml.Unmarshal(yml,&m)
	if err!=nil {
		return nil,err
	}
	return MapHandler(m,fallback), nil
}
