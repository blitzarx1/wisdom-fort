package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	"blitzarx1/wisdom-fort/pkg/scheme"
	wfErrors "blitzarx1/wisdom-fort/server/internal/errors"
	"blitzarx1/wisdom-fort/server/internal/logger"
	"blitzarx1/wisdom-fort/server/internal/service/challenges"
	"blitzarx1/wisdom-fort/server/internal/service/handlers"
	"blitzarx1/wisdom-fort/server/internal/service/quotes"
	"blitzarx1/wisdom-fort/server/internal/service/rps"
	"blitzarx1/wisdom-fort/server/internal/service/storage"
	"blitzarx1/wisdom-fort/server/internal/token"
)

const quotesFilePath = "server/quotes.json"

type App struct {
	cfg             *Config
	handlersService *handlers.Service
	rpsService      *rps.Service
}

func New(ctx context.Context, cfg *Config) (*App, error) {
	l := logger.New(nil, "server")
	l.Println("initializing server")

	a := &App{cfg: cfg}
	if err := a.initServices(logger.WithCtx(ctx, l, "initServices")); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	l := logger.New(nil, "serverRun")
	l.Println("running server")
	portStr := fmt.Sprintf(":%d", a.cfg.Port)
	ln, err := net.Listen("tcp", portStr)
	if err != nil {
		a.logError(logger.WithCtx(ctx, l, ""), err)
		return err
	}
	defer ln.Close()

	return a.serve(logger.WithCtx(ctx, l, "serve"), ln)
}

func (a *App) initServices(ctx context.Context) error {
	l := logger.MustFromCtx(ctx)
	l.Println("initializing services")

	var err error

	storageService := storage.New(logger.WithCtx(ctx, l, "storage"))
	a.rpsService = rps.New(logger.WithCtx(ctx, l, "rps"), storageService)
	quotesService, err := quotes.New(logger.WithCtx(ctx, l, "quotes"), quotesFilePath)
	if err != nil {
		return err
	}
	challengesService := challenges.New(
		logger.WithCtx(ctx, l, "challenges"),
		a.cfg.DiffMult,
		a.cfg.ChallengeTTLSeconds,
		storageService,
		a.rpsService,
	)
	a.handlersService, err = handlers.New(
		logger.WithCtx(ctx, l, "service"),
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

func (a *App) serve(ctx context.Context, ln net.Listener) error {
	l := logger.MustFromCtx(ctx)
	l.Println("server is listening on ", a.cfg.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go a.handleConnection(logger.WithCtx(ctx, l, "handleConnection"), conn)
	}
}

func (a *App) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	l := logger.MustFromCtx(ctx)

	data, err := a.read(logger.WithCtx(ctx, l, "read"), conn)
	if err != nil {
		a.handleError(logger.WithCtx(ctx, l, ""), conn, nil, wfErrors.NewError(wfErrors.ErrGeneric, err))
		return
	}

	var req scheme.Request
	if err := json.Unmarshal(data, &req); err != nil {
		a.handleError(logger.WithCtx(ctx, l, ""), conn, nil, wfErrors.NewError(wfErrors.ErrInvalidMsgFormat, err))
		return
	}

	clientAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		a.handleError(
			logger.WithCtx(ctx, l, ""),
			conn,
			nil,
			wfErrors.NewError(wfErrors.ErrGeneric, fmt.Errorf("failed to get client addr: %s", clientAddr)),
		)
	}

	ip := clientAddr.IP.String()
	token, authErr := a.auth(ip, &req)
	if authErr != nil {
		a.handleError(logger.WithCtx(ctx, l, ""), conn, nil, authErr)
		return
	}

	tokenLogger := logger.New(l, string(token))

	var respPayload []byte
	var handleErr *wfErrors.Error
	switch req.Action {
	case scheme.ActionChallenge:
		respPayload, handleErr = a.handlersService.GenerateChallenge(
			logger.WithCtx(ctx, tokenLogger, "genChallenge"),
			token,
		)
	case scheme.ActionSolution:
		respPayload, handleErr = a.handlersService.CheckSolution(
			logger.WithCtx(ctx, tokenLogger, "checkSolution"),
			token,
			req.Payload,
		)
	default:
		handleErr = wfErrors.NewError(wfErrors.ErrInvalidAction, fmt.Errorf("unknown action: %s", req.Action))
	}
	if handleErr != nil {
		a.handleError(logger.WithCtx(ctx, tokenLogger, ""), conn, &token, handleErr)
		return
	}

	if err := a.write(logger.WithCtx(ctx, l, "write"), conn, a.successResponse(token, respPayload)); err != nil {
		a.logError(logger.WithCtx(ctx, tokenLogger, ""), err)
		return
	}
}

func (a *App) logError(ctx context.Context, err error) {
	l := logger.MustFromCtx(ctx)
	l.SetOutput(os.Stderr)
	defer l.SetOutput(os.Stdout)

	l.Println(err.Error())
}

func (a *App) read(ctx context.Context, conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	logger.MustFromCtx(ctx).Println("read data: ", string(buffer[:n]))

	return buffer[:n], nil
}

func (a *App) write(ctx context.Context, conn net.Conn, data []byte) error {
	logger.MustFromCtx(ctx).Println("writing data: ", string(data))

	_, err := conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) rpsUnauthorizedGuard(ip string) *wfErrors.Error {
	rps := a.rpsService.Get(ip)
	if rps > a.cfg.RPSLimitUnauth {
		return wfErrors.NewError(wfErrors.ErrTooManyRequests, errors.New("too many requests"))
	}

	return nil
}

// auth gates all calls and returns a token for the given request.
// Returns error if unauthorized request requires a token but none is provided or if it exceedes
// rps limit for unauthorized requests.
func (a *App) auth(ip string, req *scheme.Request) (token.Token, *wfErrors.Error) {
	a.rpsService.Inc(ip)

	if req.Token != nil && *req.Token != "" {
		t := token.Token(*req.Token)
		return t, nil
	}

	if req.Action == scheme.ActionSolution {
		return token.Token(""), wfErrors.NewError(wfErrors.ErrMissingToken, errors.New("action requires a token"))
	}

	if err := a.rpsUnauthorizedGuard(ip); err != nil {
		return token.Token(""), err
	}

	return a.handlersService.GenerateToken(ip), nil
}

func (a *App) handleError(ctx context.Context, conn net.Conn, t *token.Token, err *wfErrors.Error) {
	l := logger.MustFromCtx(ctx)
	a.logError(logger.WithCtx(ctx, l, ""), err)
	a.write(logger.WithCtx(ctx, l, "write"), conn, a.errorResponse(t, err))
}

func (a *App) successResponse(token token.Token, payload []byte) []byte {
	resp := scheme.Response{
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
	resp := scheme.Response{
		Token:     t,
		Error:     &errStr,
		ErrorCode: &codeStr,
	}

	data, _ := json.Marshal(resp)
	return data
}
