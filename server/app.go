package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"blitzarx1/wisdom-fort/pkg/api"
	wfErrors "blitzarx1/wisdom-fort/server/internal/errors"
	"blitzarx1/wisdom-fort/server/internal/logger"
	"blitzarx1/wisdom-fort/server/internal/service/challenges"
	"blitzarx1/wisdom-fort/server/internal/service/handlers"
	"blitzarx1/wisdom-fort/server/internal/service/quotes"
	"blitzarx1/wisdom-fort/server/internal/service/rps"
	"blitzarx1/wisdom-fort/server/internal/service/storage"
	"blitzarx1/wisdom-fort/server/internal/token"
)

const (
	port = 8080

	quotesFilePath = "server/quotes.json"
	rpsUnauthLimit  = 1
)

type App struct {
	logger *log.Logger

	handlersService *handlers.Service
	rpsService      *rps.Service
}

func New() (*App, error) {
	l := logger.NewLogger(nil, "server")
	l.Println("initializing server")

	a := &App{logger: l}
	if err := a.initServices(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	portStr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", portStr)
	if err != nil {
		a.logError(err)
		return err
	}
	defer ln.Close()

	return a.serve(ln)
}

func (a *App) initServices() error {
	a.logger.Println("initializing services")

	var err error

	storageService := storage.New(logger.NewLogger(a.logger, "storage"))
	a.rpsService = rps.New(logger.NewLogger(a.logger, "rps"), storageService)
	quotesService, err := quotes.New(logger.NewLogger(a.logger, "quotes"), quotesFilePath)
	if err != nil {
		return err
	}
	challengesService := challenges.New(logger.NewLogger(a.logger, "challenges"), storageService, a.rpsService)
	a.handlersService, err = handlers.New(
		logger.NewLogger(a.logger, "service"),
		a.rpsService,
		storageService,
		quotesService,
		challengesService,
	)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) serve(ln net.Listener) error {
	a.logger.Println("server is listening on ", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go a.handleConnection(conn)
	}
}

func (a *App) handleConnection(conn net.Conn) {
	defer conn.Close()

	data, err := a.read(conn)
	if err != nil {
		a.handleError(conn, nil, wfErrors.NewError(wfErrors.ErrGeneric, err))
		return
	}

	var req api.Request
	if err := json.Unmarshal(data, &req); err != nil {
		a.handleError(conn, nil, wfErrors.NewError(wfErrors.ErrInvalidMsgFormat, err))
		return
	}

	clientAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		a.handleError(conn, nil, wfErrors.NewError(wfErrors.ErrGeneric, fmt.Errorf("failed to get client addr: %s", clientAddr)))
	}

	ip := clientAddr.IP.String()
	token, authErr := a.auth(ip, &req)
	if authErr != nil {
		a.handleError(conn, nil, authErr)
		return
	}

	var respPayload []byte
	var handleErr *wfErrors.Error
	switch req.Action {
	case api.ActionChallenge:
		respPayload, handleErr = a.handlersService.GenerateChallenge(token)
	case api.ActionSolution:
		respPayload, handleErr = a.handlersService.CheckSolution(token, req.Payload)
	default:
		handleErr = wfErrors.NewError(wfErrors.ErrInvalidAction, fmt.Errorf("unknown action: %s", req.Action))
	}
	if handleErr != nil {
		a.handleError(conn, &token, handleErr)
		return
	}

	if err := a.write(conn, a.successResponse(token, respPayload)); err != nil {
		a.logError(err)
		return
	}
}

func (a *App) logError(err error) {
	a.logger.SetOutput(os.Stderr)
	defer a.logger.SetOutput(os.Stdout)

	a.logger.Println(err.Error())
}

func (a *App) read(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	a.logger.Println("read data: ", string(buffer[:n]))

	return buffer[:n], nil
}

func (a *App) write(conn net.Conn, data []byte) error {
	a.logger.Println("writing data: ", string(data))

	_, err := conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) rpsUnauthorizedGuard(ip string) *wfErrors.Error {
	rps := a.rpsService.Get(ip)
	if rps > rpsUnauthLimit {
		return wfErrors.NewError(wfErrors.ErrTooManyRequests, errors.New("too many requests"))
	}

	return nil
}

// auth gates all calls and returns a token for the given request.
// Returns error if unauthorized request requires a token but none is provided or if it exceedes
// rps limit for unauthorized requests.
func (a *App) auth(ip string, req *api.Request) (token.Token, *wfErrors.Error) {
	a.rpsService.Inc(ip)

	if req.Token != nil && *req.Token != "" {
		t := token.Token(*req.Token)
		return t, nil
	}

	if req.Action == api.ActionSolution {
		return token.Token(""), wfErrors.NewError(wfErrors.ErrMissingToken, errors.New("action requires a token"))
	}

	if err := a.rpsUnauthorizedGuard(ip); err != nil {
		return token.Token(""), err
	}

	return a.handlersService.GenerateToken(ip), nil
}

func (a *App) handleError(conn net.Conn, t *token.Token, err *wfErrors.Error) {
	a.logError(err)
	a.write(conn, a.errorResponse(t, err))
}

func (a *App) successResponse(token token.Token, payload []byte) []byte {
	resp := api.Response{
		Token:   string(token),
		Payload: payload,
	}

	data, _ := json.Marshal(resp)
	return data
}

func (a *App) errorResponse(token *token.Token, err *wfErrors.Error) []byte {
	var t string
	if token != nil {
		t = string(*token)
	}

	errStr := err.Error()
	codeStr := err.Code().String()
	resp := api.Response{
		Token:     t,
		Error:     &errStr,
		ErrorCode: &codeStr,
	}

	data, _ := json.Marshal(resp)
	return data
}
