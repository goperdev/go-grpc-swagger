package main

import (
	"Myswagger/config"
	gw "Myswagger/protos"
	"bytes"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var (
	configPath   = flag.String("c", "./config/config.conf", "Application config file")
	endpoint     = flag.String("endpoint", "", "endpoint of Service")
	swaggerJson  = flag.String("swagger_json", "./protos/test.swagger.json", "swagger json file")
	swaggerUiDir = flag.String("swagger_ui_dir", "./static/swagger-ui", "path to the directory which contains swagger ui statics files")
)

func serveSwaggerJson(w http.ResponseWriter, r *http.Request) {
	log.Println("%v", r)
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		glog.Errorf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	glog.Infof("Serving %s", r.URL.Path)

	contentBytes, err := ioutil.ReadFile(*swaggerJson)
	if err != nil {
		fmt.Print(err)
	}
	js, _ := simplejson.NewJson(contentBytes)
	pathsMap, err := js.Get("paths").Map()
	for url, _ := range pathsMap {
		methodsMap, _ := js.Get("paths").Get(url).Map()
		for method, _ := range methodsMap {
			parametersArr, _ := js.Get("paths").Get(url).Get(method).Get("parameters").Array()
			authMap := []interface{}{

				//这里是在每个请求前加上Authorization鉴权认证

				//map[string]interface{}{
				//	"name":     "Authorization",
				//	"in":       "header",
				//	"required": false,
				//	"type":     "string",
				//},
			}
			js.SetPath([]string{"paths", url, method, "parameters"}, config.Insert(parametersArr, authMap, 0))
		}
	}
	contentBytes, _ = js.MarshalJSON()
	http.ServeContent(w, r, "swagger.json", time.Now(), bytes.NewReader(contentBytes))
}

func serveSwaggerUi(w http.ResponseWriter, r *http.Request) {
	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(*swaggerUiDir, p)
	http.ServeFile(w, r, p)
}

func Run(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/swagger/swagger.json", serveSwaggerJson)
	mux.HandleFunc("/swagger/", serveSwaggerUi)
	gateway, err := newGateway(ctx, opts...)
	if err != nil {
		return nil
	}
	mux.Handle("/", gateway)
	return http.ListenAndServe(address, mux)
}


func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	//注册 hello service
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, *endpoint, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func main() {
	defer glog.Flush()
	flag.Parse()
	if configPath == nil {
		os.Exit(1)
	}
	config.InitConfig(*configPath, config.Settings)
	s := fmt.Sprintf("%v:%v", config.Settings.Server.GrpcHost, config.Settings.Server.GrpcPort)
	endpoint = &s

	if err := Run(fmt.Sprintf(":%v", config.Settings.Server.GatewayPort)); err != nil {
		glog.Fatal(err)
	}
	//fmt.Println("成功")
	//grpcServer := grpc.NewServer()
	//gw.RegisterGreeterServer(grpcServer, &services.ConfigService{})
}
