package main

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/gorilla/mux"
)

var (
    runningInLambda = os.Getenv("RUN_LAMBDA") == "1"
    router          *mux.Router
)

const (
    defaultPort = "8080"
    baseUrl     = ""
)

// init initializes the router.
func init() {
    router = mux.NewRouter()
    router.HandleFunc(baseUrl+"/", homeHandler).Methods(http.MethodGet, http.MethodHead)
    router.HandleFunc(baseUrl+"/about", aboutHandler).Methods(http.MethodGet, http.MethodHead)
}

// RequestData represents the structure of the incoming JSON string.
type RequestData struct {
    RawPath         string            `json:"rawPath"`
    RawQueryString  string            `json:"rawQueryString"`
    Body            string            `json:"body"`
    Headers         map[string]string `json:"headers"`
    Cookies         []string          `json:"cookies"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
    RequestContext  struct {
        HTTP struct {
            Method string `json:"method"`
        } `json:"http"`
    } `json:"requestContext"`
    BinBody []byte
}

// ResponseData represents the structure of the outgoing JSON string.
type ResponseData struct {
    StatusCode      int               `json:"statusCode"`
    Headers         map[string]string `json:"headers"`
    Body            string            `json:"body"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
    Cookies         []string          `json:"cookies"`
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
        for _, s := range strings.Split(value, ",") {
            req.Header.Add(key, strings.TrimSpace(s)) // fixme
        }
    }

    for _, cookieStr := range reqData.Cookies {
        parts := strings.Split("; ", cookieStr)
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

func runLambda() error {
    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("failed to read from stdin: %w", err)
    }
    request = strings.TrimSpace(request)

    var reqData *RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        return fmt.Errorf("failed to parse request JSON: %w", err)
    }

    httpReq, err := reconstructHTTPRequest(reqData)
    if err != nil {
        return fmt.Errorf("failed to reconstruct HTTP request: %w", err)
    }

    respWriter := httptest.NewRecorder()
    router.ServeHTTP(respWriter, httpReq)
    httpResp := respWriter.Result()

    respData, err := convertToResponseData(httpResp)
    if err != nil {
        return fmt.Errorf("failed to convert response data: %w", err)
    }

    jsonResponse, err := json.Marshal(respData)
    if err != nil {
        return fmt.Errorf("failed to marshal JSON response: %w", err)
    }

    fmt.Println(string(jsonResponse))
    return nil
}

func runHTTPServer() error {
    port := defaultPort
    if p, ok := os.LookupEnv("PORT"); ok {
        port = p
    }
    log.Println("Port:", port)

    return http.ListenAndServe(":"+port, router)
}

func main() {
    httptest.NewRecorder()
    if runningInLambda {
        log.Println("Running Lambda handler.")
        if err := runLambda(); err != nil {
            log.Fatalf("Error running lambda handler: %v", err)
        }
    } else {
        log.Println("Running HTTP server.")
        if err := runHTTPServer(); err != nil {
            log.Fatalf("HTTP server ended with error: %v", err)
        }
    }
}
