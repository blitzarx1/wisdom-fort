package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	wfErrors "blitzarx1/wisdom-fort/server/errors"
	"blitzarx1/wisdom-fort/server/logger"
	"blitzarx1/wisdom-fort/server/service/api"
	"blitzarx1/wisdom-fort/server/service/challenges"
	"blitzarx1/wisdom-fort/server/service/quotes"
	"blitzarx1/wisdom-fort/server/service/rps"
	"blitzarx1/wisdom-fort/server/service/storage"
	"blitzarx1/wisdom-fort/server/token"
)

const (
	port = 8080

	quotesFilePath = "server/quotes.json"
)

type App struct {
	logger *log.Logger

	apiService *api.Service
	rpsService *rps.Service
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
	a.apiService, err = api.New(
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

	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		a.handleError(conn, nil, wfErrors.NewError(wfErrors.ErrInvalidMsgFormat, err))
		return
	}

	clientAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		a.handleError(conn, nil, wfErrors.NewError(wfErrors.ErrGeneric, fmt.Errorf("failed to get client addr: %s", clientAddr)))
	}

	ip := clientAddr.IP.String()
	token, err := a.auth(ip, &req)
	if err != nil {
		a.handleError(conn, nil, wfErrors.NewError(wfErrors.ErrMissingToken, err))
		return
	}

	var respPayload []byte
	var handleErr *wfErrors.Error
	switch req.Action {
	case CHALLENGE.String():
		respPayload, handleErr = a.apiService.GenerateChallenge(token)
	case SOLUTION.String():
		respPayload, handleErr = a.apiService.CheckSolution(token, req.Payload)
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

// auth gates all calls and returns a token for the given request.
// Returns error if request requires a token but none is provided.
func (a *App) auth(ip string, req *request) (token.Token, error) {
	a.rpsService.Inc(ip)

	if req.Token != nil && *req.Token != "" {
		t := token.Token(*req.Token)
		return t, nil
	}

	if req.Action == SOLUTION.String() {
		return token.Token(""), errors.New("missing token")
	}

	return a.apiService.GenerateToken(ip), nil
}

func (a *App) handleError(conn net.Conn, t *token.Token, err *wfErrors.Error) {
	a.logError(err)
	a.write(conn, a.errorResponse(t, err))
}

func (a *App) successResponse(token token.Token, payload []byte) []byte {
	resp := response{
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
	resp := response{
		Token:     t,
		Error:     &errStr,
		ErrorCode: &codeStr,
	}

	data, _ := json.Marshal(resp)
	return data
}
