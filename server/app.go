package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"blitzarx1/wisdom-fort/server/service"
)

type App struct {
	logger  *log.Logger
	service *service.Service
}

func New() (*App, error) {
	logger := service.NewLogger(nil, "server")
	logger.Println("initializing server")

	service, err := service.New(service.NewLogger(logger, "service"))
	if err != nil {
		return nil, err
	}

	return &App{
		logger:  logger,
		service: service,
	}, nil
}

func (a *App) Run() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		a.logError(err)
		return err
	}
	defer ln.Close()

	return a.serve(ln)
}

func (a *App) serve(ln net.Listener) error {
	a.logger.Println("server is listening on port 8080")

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
		a.handleError(conn, nil, service.NewError(service.ErrGeneric, err))
		return
	}

	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		a.handleError(conn, nil, service.NewError(service.ErrInvalidMsgFormat, err))
		return
	}

	ip := conn.RemoteAddr().String()
	token := a.token(ip, &req)

	var respPayload []byte
	var handleErr *service.Error
	switch req.Action {
	case CHALLENGE.String():
		respPayload, handleErr = a.service.HandleChallenge(ip, token)
	case SOLUTION.String():
		respPayload, handleErr = a.service.HandleSolution(ip, token, req.Payload)
	default:
		handleErr = service.NewError(service.ErrInvalidAction, fmt.Errorf("unknown action: %s", req.Action))
	}
	if handleErr != nil {
		a.handleError(conn, &token, handleErr)
		return
	}

	a.write(conn, a.successResponse(token, respPayload))
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

func (a *App) write(conn net.Conn, data []byte) {
	a.logger.Println("writing data: ", string(data))

	_, err := conn.Write(data)
	if err != nil {
		a.logError(err)
	}
}

func (a *App) token(ip string, req *request) service.Token {
	if req.Token != nil && *req.Token != "" {
		t := service.Token(*req.Token)
		return t
	}

	return a.service.GenerateToken(ip)
}

func (a *App) handleError(conn net.Conn, t *service.Token, err *service.Error) {
	a.logError(err)
	a.write(conn, a.errorResponse(t, err))
}

func (a *App) successResponse(token service.Token, payload []byte) []byte {
	resp := response{
		Token:   string(token),
		Payload: payload,
	}

	data, _ := json.Marshal(resp)
	return data
}

func (a *App) errorResponse(token *service.Token, err *service.Error) []byte {
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
