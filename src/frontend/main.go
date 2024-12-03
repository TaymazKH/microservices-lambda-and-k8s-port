package main

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)

const (
    defaultPort     = "8080"
    defaultCurrency = "USD"
    cookieMaxAge    = 60 * 60 * 48

    cookiePrefix    = "shop_"
    cookieSessionID = cookiePrefix + "session-id"
    cookieCurrency  = cookiePrefix + "currency"
)

var (
    runningInLambda = os.Getenv("RUN_LAMBDA") == "1"
    baseUrl         = os.Getenv("BASE_URL") // must begin with a slash if non-empty
    httpHandler     http.Handler

    whitelistedCurrencies = map[string]bool{
        "USD": true,
        "EUR": true,
        "CAD": true,
        "JPY": true,
        "GBP": true,
        "TRY": true,
    }
)

type ctxKeySessionID struct{}

type frontendServer struct {
    shoppingAssistantSvcAddr string
}

func init() {
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

    homeUrl := "/"
    if baseUrl != "" {
        homeUrl = baseUrl
    }

    r := mux.NewRouter()
    r.HandleFunc(homeUrl, svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
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
    handler = &logHandler{log: log, next: handler} // add logging
    handler = ensureSessionID(handler)             // add session ID

    httpHandler = handler
}

// RequestData represents the structure of the incoming JSON string.
type RequestData struct {
    RequestContext struct {
        HTTP struct {
            Method string `json:"method"`
        } `json:"http"`
    } `json:"requestContext"`
    RawPath         string            `json:"rawPath"`
    RawQueryString  string            `json:"rawQueryString"`
    Headers         map[string]string `json:"headers"`
    Cookies         []string          `json:"cookies"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
    Body            string            `json:"body"`
}

// ResponseData represents the structure of the outgoing JSON string.
type ResponseData struct {
    StatusCode      int               `json:"statusCode"`
    Headers         map[string]string `json:"headers"`
    Cookies         []string          `json:"cookies"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
    Body            string            `json:"body"`
}

// nonSplitHeaders is the set of header keys that should not be split based on a comma.
var nonSplitHeaders = map[string]bool{
    // Authentication
    "authorization":       true,
    "proxy-authorization": true,

    // Cookies
    "cookie":     true,
    "set-cookie": true,

    // User Agent
    "user-agent": true,
    "referer":    true, // May contain query strings with commas

    // Caching
    "if-match":            true,
    "if-none-match":       true,
    "if-unmodified-since": true,
    "if-modified-since":   true,
    "last-modified":       true,

    // Content Headers
    "content-disposition": true,
    "content-type":        true,

    // Range Requests
    "range": true,

    // Miscellaneous
    "location":        true,
    "link":            true, // May contain URIs with commas
    "x-forwarded-for": true, // Often contains IP lists with commas
}

// reconstructHTTPRequest reconstructs the incoming HTTP request.
func reconstructHTTPRequest(reqData *RequestData) (*http.Request, error) {
    rawURL := reqData.RawPath
    if reqData.RawQueryString != "" {
        rawURL += "?" + reqData.RawQueryString
    }
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return nil, err
    }

    var body io.Reader
    if reqData.IsBase64Encoded {
        decodedBody, err := base64.StdEncoding.DecodeString(reqData.Body)
        if err != nil {
            return nil, err
        }
        body = bytes.NewReader(decodedBody)
    } else {
        body = strings.NewReader(reqData.Body)
    }

    req, err := http.NewRequest(reqData.RequestContext.HTTP.Method, parsedURL.String(), body)
    if err != nil {
        return nil, err
    }

    for key, value := range reqData.Headers {
        if nonSplitHeaders[strings.ToLower(key)] {
            req.Header.Add(key, strings.TrimSpace(value))
        } else {
            for _, s := range strings.Split(value, ",") {
                req.Header.Add(key, strings.TrimSpace(s))
            }
        }
    }

    for _, cookieStr := range reqData.Cookies {
        parts := strings.Split(cookieStr, "; ")
        if len(parts) == 0 {
            continue
        }

        nameValue := strings.SplitN(parts[0], "=", 2)
        if len(nameValue) != 2 {
            continue
        }
        cookie := &http.Cookie{
            Name:  nameValue[0],
            Value: nameValue[1],
        }

        for _, attr := range parts[1:] {
            attrParts := strings.SplitN(attr, "=", 2)
            key := strings.ToLower(strings.TrimSpace(attrParts[0]))
            var value string
            if len(attrParts) > 1 {
                value = strings.TrimSpace(attrParts[1])
            }

            switch key {
            case "path":
                cookie.Path = value
            case "domain":
                cookie.Domain = value
            case "expires":
                if t, err := time.Parse(time.RFC1123, value); err == nil {
                    cookie.Expires = t
                }
            case "max-age":
                if maxAge, err := strconv.Atoi(value); err == nil {
                    cookie.MaxAge = maxAge
                }
            case "secure":
                cookie.Secure = true
            case "httponly":
                cookie.HttpOnly = true
            case "samesite":
                switch strings.ToLower(value) {
                case "lax":
                    cookie.SameSite = http.SameSiteLaxMode
                case "strict":
                    cookie.SameSite = http.SameSiteStrictMode
                case "none":
                    cookie.SameSite = http.SameSiteNoneMode
                }
            }
        }

        req.AddCookie(cookie)
    }

    return req, nil
}

// convertToResponseData converts an HTTP response to ResponseData.
func convertToResponseData(resp *http.Response) (*ResponseData, error) {
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    //resp.Header.Set("Content-Type", http.DetectContentType(body))

    headers := make(map[string]string)
    var cookies []string
    for key, values := range resp.Header {
        if key == "Set-Cookie" {
            cookies = append(cookies, values...)
        } else {
            headers[key] = strings.Join(values, ",")
        }
    }

    return &ResponseData{
        StatusCode:      resp.StatusCode,
        Headers:         headers,
        Body:            base64.StdEncoding.EncodeToString(body),
        IsBase64Encoded: true,
        Cookies:         cookies,
    }, nil
}

func runLambda(reqData *RequestData) (*ResponseData, error) {
    log.Infof("Handler started. Event data: %v", reqData)

    httpReq, err := reconstructHTTPRequest(reqData)
    if err != nil {
        return nil, fmt.Errorf("failed to reconstruct HTTP request: %w", err)
    }

    respWriter := httptest.NewRecorder()
    httpHandler.ServeHTTP(respWriter, httpReq)
    httpResp := respWriter.Result()

    respData, err := convertToResponseData(httpResp)
    if err != nil {
        return nil, fmt.Errorf("failed to convert response data: %w", err)
    }

    log.Infof("Handler finished. Response: %v", respData)
    return respData, nil
}

func runHTTPServer() error {
    port := defaultPort
    if p, ok := os.LookupEnv("PORT"); ok {
        port = p
    }
    addr := os.Getenv("LISTEN_ADDR")

    log.Info("Starting HTTP server on " + addr + ":" + port)
    return http.ListenAndServe(addr+":"+port, httpHandler)
}

func main() {
    if runningInLambda {
        lambda.Start(runLambda)
    } else {
        if err := runHTTPServer(); err != nil {
            log.Fatalf("HTTP server ended with error: %v", err)
        }
    }
}
