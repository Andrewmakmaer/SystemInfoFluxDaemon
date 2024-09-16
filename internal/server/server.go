package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/diskinfo"
	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/server/pb"
	"github.com/go-kit/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/kit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type app interface {
	CPUValueAverage(m int) (Sys, Usr, Idle, Iowait float32, err error)
	LAverage(m int) (la1, la5, la15 float32, err error)
	DiskAverage(m int) (resultMap map[string]diskinfo.DiskInfo, err error)
}

type Server struct {
	pb.UnimplementedDaemonServer
	listener net.Listener
	serv     *grpc.Server
	logg     log.Logger
	app      app
}

func NewServer(app app, port string) Server {
	logger := log.NewJSONLogger(os.Stdout)

	lsn, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	serv := grpc.NewServer(
		grpc.ChainStreamInterceptor(grpc_middleware.ChainStreamServer(kit.StreamServerInterceptor(logger))),
	)
	newServer := &Server{listener: lsn, serv: serv, app: app, logg: logger}
	return *newServer
}

func (s *Server) Start(ctx context.Context) error {
	pb.RegisterDaemonServer(s.serv, s)
	reflection.Register(s.serv) // for postman

	if err := s.serv.Serve(s.listener); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.serv.GracefulStop()
	s.listener.Close()
	return nil
}

func (s *Server) EnableStatStream(req *pb.StreamRequest, stream pb.Daemon_EnableStatStreamServer) error {
	streamType := req.GetStatsType()
	streamDelay := int(req.GetSecondsDelay())
	streamRange := req.GetSecondsRange()

	s.stream(streamType, streamDelay, int(streamRange), stream)
	return nil
}

func (s *Server) stream(streamType string, streamDelay, streamRange int, stream pb.Daemon_EnableStatStreamServer) {
	var msg *pb.StatsResponce
	stop := false
	for !stop {
		switch streamType {
		case "cpu":
			msg = s.cpu(streamRange)
		case "la":
			msg = s.la(streamRange)
		case "disk":
			msg = s.disk(streamRange)
		}
		err := stream.Send(msg)
		if err != nil {
			e := s.logg.Log("message", "not send message", "status", "send error")
			if e != nil {
				fmt.Println("Message not send")
			}
			continue
		}
		time.Sleep(time.Duration(streamDelay) * time.Second)
	}
}

func (s *Server) cpu(streamRange int) (msg *pb.StatsResponce) {
	sys, usr, idle, iowait, _ := s.app.CPUValueAverage(streamRange)
	msg = &pb.StatsResponce{Stat: &pb.StatsResponce_CpuStats{CpuStats: &pb.CPUStat{
		Sys:    sys,
		Usr:    usr,
		Idle:   idle,
		Iowait: iowait,
	}}}
	return
}

func (s *Server) la(streamRange int) (msg *pb.StatsResponce) {
	la1, la5, la15, _ := s.app.LAverage(streamRange)
	msg = &pb.StatsResponce{Stat: &pb.StatsResponce_LaStats{LaStats: &pb.LAStat{La1: la1, La5: la5, La15: la15}}}
	return
}

func (s *Server) disk(streamRange int) (msg *pb.StatsResponce) {
	result, _ := s.app.DiskAverage(streamRange)
	diskResult := make(map[string]*pb.DiskInfoStat)
	for k, v := range result {
		diskResult[k] = &pb.DiskInfoStat{Tps: v.Tps, Kbs: v.Kbs}
	}
	msg = &pb.StatsResponce{Stat: &pb.StatsResponce_DiskInfo{DiskInfo: &pb.DiskStat{DiskStat: diskResult}}}
	return
}
