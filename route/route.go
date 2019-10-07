package route

import (
    "github.com/gorilla/mux"
    gmail "github.com/heaptracetechnology/gmail/service"
    "log"
    "net/http"
)

//Route struct
type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

//Routes list
type Routes []Route

var routes = Routes{
    Route{
        "HealthCheck",
        "GET",
        "/health",
        gmail.HealthCheck,
    },
    Route{
        "Authorization",
        "POST",
        "/authorization",
        gmail.Authorization,
    },
    Route{
        "AccessToken",
        "POST",
        "/accessToken",
        gmail.AccessToken,
    },
    Route{
        "SendMail",
        "POST",
        "/sendMail",
        gmail.SendMail,
    },
    Route{
        "ReceiveEmail",
        "POST",
        "/receive",
        gmail.ReceiveMail,
    },
    Route{
        "RefreshToken",
        "POST",
        "/refreshToken",
        gmail.RefreshToken,
    },
}

//NewRouter func
func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        log.Println(route.Name)
        handler = route.HandlerFunc

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
    return router
}
