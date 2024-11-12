package main

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
    "google.golang.org/grpc"
)

const (
    port            = "8080"
    defaultCurrency = "USD"
    cookieMaxAge    = 60 * 60 * 48

    cookiePrefix    = "shop_"
    cookieSessionID = cookiePrefix + "session-id"
    cookieCurrency  = cookiePrefix + "currency"
)

var (
    whitelistedCurrencies = map[string]bool{
        "USD": true,
        "EUR": true,
        "CAD": true,
        "JPY": true,
        "GBP": true,
        "TRY": true,
    }

    baseUrl = ""
)

type ctxKeySessionID struct{}

type frontendServer struct {
    productCatalogSvcAddr string
    productCatalogSvcConn *grpc.ClientConn

    currencySvcAddr string
    currencySvcConn *grpc.ClientConn

    cartSvcAddr string
    cartSvcConn *grpc.ClientConn

    recommendationSvcAddr string
    recommendationSvcConn *grpc.ClientConn

    checkoutSvcAddr string
    checkoutSvcConn *grpc.ClientConn

    shippingSvcAddr string
    shippingSvcConn *grpc.ClientConn

    adSvcAddr string
    adSvcConn *grpc.ClientConn

    collectorAddr string
    collectorConn *grpc.ClientConn

    shoppingAssistantSvcAddr string
}

func main() {
    log := logrus.New()
    log.Level = logrus.DebugLevel
    log.Formatter = &logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "severity",
            logrus.FieldKeyMsg:   "message",
        },
        TimestampFormat: time.RFC3339Nano,
    }
    log.Out = os.Stdout

    svc := new(frontendServer)

    otel.SetTextMapPropagator(
        propagation.NewCompositeTextMapPropagator(
            propagation.TraceContext{}, propagation.Baggage{}))

    baseUrl = os.Getenv("BASE_URL")

    srvPort := port
    if os.Getenv("PORT") != "" {
        srvPort = os.Getenv("PORT")
    }
    addr := os.Getenv("LISTEN_ADDR")

    r := mux.NewRouter()
    r.HandleFunc(baseUrl+"/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
    r.HandleFunc(baseUrl+"/product/{id}", svc.productHandler).Methods(http.MethodGet, http.MethodHead)
    r.HandleFunc(baseUrl+"/cart", svc.viewCartHandler).Methods(http.MethodGet, http.MethodHead)
    r.HandleFunc(baseUrl+"/cart", svc.addToCartHandler).Methods(http.MethodPost)
    r.HandleFunc(baseUrl+"/cart/empty", svc.emptyCartHandler).Methods(http.MethodPost)
    r.HandleFunc(baseUrl+"/setCurrency", svc.setCurrencyHandler).Methods(http.MethodPost)
    r.HandleFunc(baseUrl+"/logout", svc.logoutHandler).Methods(http.MethodGet)
    r.HandleFunc(baseUrl+"/cart/checkout", svc.placeOrderHandler).Methods(http.MethodPost)
    r.HandleFunc(baseUrl+"/assistant", svc.assistantHandler).Methods(http.MethodGet)
    r.PathPrefix(baseUrl + "/static/").Handler(http.StripPrefix(baseUrl+"/static/", http.FileServer(http.Dir("./static/"))))
    r.HandleFunc(baseUrl+"/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
    r.HandleFunc(baseUrl+"/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })
    r.HandleFunc(baseUrl+"/product-meta/{ids}", svc.getProductByID).Methods(http.MethodGet)
    r.HandleFunc(baseUrl+"/bot", svc.chatBotHandler).Methods(http.MethodPost)

    var handler http.Handler = r
    handler = &logHandler{log: log, next: handler}     // add logging
    handler = ensureSessionID(handler)                 // add session ID
    handler = otelhttp.NewHandler(handler, "frontend") // add OTel tracing

    log.Infof("starting server on " + addr + ":" + srvPort)
    log.Fatal(http.ListenAndServe(addr+":"+srvPort, handler))
}
