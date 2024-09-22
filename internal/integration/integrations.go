package integration

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/Andrewmakmaer/SystemInfoFluxDaemon/internal/server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var daemonAddr = os.Getenv("INTEGRATION_DAEMON_URL")

func RunIntegrationTests() (err error) {
	err = waitForServices()
	if err != nil {
		return fmt.Errorf("ошибка при ожидании запуска сервисов: %w", err)
	}
	fmt.Println("system info flux запущен")

	result, err := getCpuStatFromServ()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	fmt.Println("Test 1: Проверка коректности возвращаеммых данных CPU")
	err = nil
	for i, item := range result {
		s := item.Idle + item.Sys + item.Usr
		s = s - 1
		if s > 100 || s <= 0 {
			err = fmt.Errorf("некоректное значение показателей CPU - sum: %f, idle: %f, sys: %f, usr: %f",
				s, item.Idle, item.Sys, item.Usr)
		}
		if i+1 < len(result) {
			if item.Idle == result[i+1].Idle {
				err = fmt.Errorf("нет динамики в изменении idle статуса cpu %f -> %f | %w",
					item.Idle, result[i+1].Idle, err)
			}
		}
	}

	resultLa, err := getLaStatFromServ()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	fmt.Println("Test 2: Проверка коректности возвращаеммых данных LA")
	for i, item := range resultLa {
		if item.La1 <= 0 || item.La5 <= 0 || item.La15 <= 0 {
			return fmt.Errorf("некоректное значение показателей LA - la1: %f, la5: %f, la15: %f",
				item.La1, item.La5, item.La15)
		}
		if i+1 < len(resultLa) {
			if item.La1 == resultLa[i+1].La1 {
				return fmt.Errorf("нет динамики в изменении idle статуса cpu %f -> %f",
					item.La1, resultLa[i+1].La1)
			}
		}
	}

	return err
}

func waitForServices() error {

	for i := 0; i < 30; i++ {
		conn, err := grpc.NewClient(daemonAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			conn.Close()
			return nil
		}
		fmt.Println(err)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("сервисы не запустились в отведенное время")
}

func getCpuStatFromServ() ([]*pb.CPUStat, error) {
	var result []*pb.CPUStat
	conn, err := grpc.NewClient(daemonAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer conn.Close()

	client := pb.NewDaemonClient(conn)

	ctx, close := context.WithCancel(context.Background())
	defer close()

	req := pb.StreamRequest{StatsType: "cpu", SecondsRange: 5, SecondsDelay: 2}
	stream, err := client.EnableStatStream(ctx, &req)
	if err != nil {
		return result, err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return result, err
		}
		if err != nil {
			slog.Error(err.Error())
			return result, err
		}
		stat := resp.GetCpuStats()
		result = append(result, stat)

		if len(result) > 3 {
			break
		}
	}
	return result, err
}

func getLaStatFromServ() ([]*pb.LAStat, error) {
	var result []*pb.LAStat
	conn, err := grpc.NewClient(daemonAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer conn.Close()

	client := pb.NewDaemonClient(conn)

	ctx, close := context.WithCancel(context.Background())
	defer close()

	req := pb.StreamRequest{StatsType: "la", SecondsRange: 5, SecondsDelay: 2}
	stream, err := client.EnableStatStream(ctx, &req)
	if err != nil {
		return result, err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return result, err
		}
		if err != nil {
			slog.Error(err.Error())
			return result, err
		}
		stat := resp.GetLaStats()
		result = append(result, stat)

		if len(result) > 3 {
			break
		}
	}
	return result, err
}
