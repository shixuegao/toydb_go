package service

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sxg/toydb_go/config"
	_grpc "sxg/toydb_go/grpc"
	"sxg/toydb_go/grpc/proto"
	logging "sxg/toydb_go/logging"
	"sxg/toydb_go/raft"
	"sxg/toydb_go/raft/node"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	_grpc.UnimplementedGrpcServer
	//common
	conf   config.Raft
	logger logging.Logger
	//node
	node     *node.Node
	caseRecv <-chan *raft.Case
	//grpc server
	grpcServer *grpc.Server
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	//holder
	holder *raft.CaseHolder
}

func NewServer(ctx context.Context, conf config.Config,
	log *raft.Log, state raft.State, logger logging.Logger) (*Server, error) {
	if conf.ID == "" {
		return nil, errors.New("id cannot be empty")
	}
	arr := []string{}
	temp := map[string]struct{}{}
	for id := range conf.Peers {
		if id == "" {
			return nil, errors.New("peer id cannot be empty")
		}
		if _, ok := temp[id]; ok {
			return nil, fmt.Errorf("peer id %s repeated", id)
		}
		arr = append(arr, id)
	}
	caseC := make(chan *raft.Case)
	ctxn, cancel := context.WithCancel(ctx)
	node, err := node.NewNode(ctxn, conf.ID, arr, caseC, log, state, logger)
	if err != nil {
		cancel()
		return nil, err
	}
	server := &Server{
		conf:     conf.Raft,
		logger:   logger,
		node:     node,
		caseRecv: caseC,
		cancel:   cancel,
		wg:       new(sync.WaitGroup),
	}
	server.holder = raft.NewCaseKeeper(ctxn, server.wg)
	return server, nil
}

func (ser *Server) Serve(addr string) (err error) {
	var lis net.Listener
	if lis, err = net.Listen("tcp", addr); err != nil {
		return
	}
	ser.grpcServer = grpc.NewServer()
	_grpc.RegisterGrpcServer(ser.grpcServer, ser)
	if err = ser.grpcServer.Serve(lis); err != nil {
		return
	}
	//holder
	ser.holder.Loop(ser.caseRecv, ser.conf.RpcRequestTimeout)
	return nil
}

func (ser *Server) String() string {
	peers := ""
	if len(ser.conf.Peers) > 0 {
		count := 0
		var builder strings.Builder
		for id, addr := range ser.conf.Peers {
			if count != 0 {
				builder.WriteString(" ,")
			}
			builder.WriteString(id)
			builder.WriteString(":")
			builder.WriteString(addr)
			count++
			peers = builder.String()
		}
	}
	return fmt.Sprintf("Raft: ID -> %v, Peers -> %v", ser.conf.ID, peers)
}

func (ser *Server) Exchange(_ context.Context, msg *proto.Message) (*proto.Message, error) {
	cas := ser.holder.Await(msg, func(cas *raft.Case) {
		ser.node.Step(cas)
	})
	if cas.Timeout {
		ser.logger.Warnf("请求[%v]执行超时!", msg.String())
	}
	res := &proto.Message{}
	return res, nil
}
