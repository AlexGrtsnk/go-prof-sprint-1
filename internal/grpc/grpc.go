package grpc

import (
	// импортируем пакет со сгенерированными protobuf-файлами
	"context"
	"fmt"
	pb "go-prof-sprint-1/internal/grpc/proto"
	"log"
	"net"
	"sort"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// URLServer поддерживает все необходимые методы сервера.
type URLServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedURLsServer

	// используем sync.Map для хранения пользователей
	URLs sync.Map
}

// AddURL реализует интерфейс добавления url.
func (s *URLServer) AddURL(ctx context.Context, in *pb.AddURLRequest) (*pb.AddURLResponse, error) {
	var response pb.AddURLResponse

	if _, ok := s.URLs.Load(in.Url.ShortURL); ok {
		response.Error = fmt.Sprintf("Пользователь уже загружал данный url для сокращения %s ", in.Url.LongURL)
	} else {
		s.URLs.Store(in.Url.ShortURL, in.Url)
	}
	return &response, nil
}

// ListURLs реализует интерфейс получения списка пользователей.
func (s *URLServer) ListURLs(ctx context.Context, in *pb.ListURLsRequest) (*pb.ListURLsResponse, error) {
	var list []string

	s.URLs.Range(func(key, _ interface{}) bool {
		list = append(list, key.(string))
		return true
	})
	// сортируем слайс из сокращенных url
	sort.Strings(list)

	offset := int(in.Offset)
	end := int(in.Offset + in.Limit)
	if end > len(list) {
		end = len(list)
	}
	if offset >= end {
		offset = 0
		end = 0
	}
	response := pb.ListURLsResponse{
		Count: int32(len(list)),
		URLs:  list[offset:end],
	}
	return &response, nil
}

// GetURL реализует интерфейс получения информации о URL.
func (s *URLServer) GetUser(ctx context.Context, in *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	var response pb.GetURLResponse

	if url, ok := s.URLs.Load(in.URL); ok {
		response.Url = url.(*pb.URL)
	} else {
		return nil, status.Errorf(codes.NotFound, `полный url для сокращенного %s не найден`, in.URL)
	}
	return &response, nil
}

// DelURL реализует интерфейс удаления URL.
func (s *URLServer) DelURL(ctx context.Context, in *pb.DelURLRequest) (*pb.DelURLResponse, error) {
	var response pb.DelURLResponse

	if _, ok := s.URLs.LoadAndDelete(in.URL); !ok {
		response.Error = fmt.Sprintf("полный url для сокращенного %s не найден, удаление невозможно", in.URL)
	}
	return &response, nil
}

func RunGRPCServer() {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterURLsServer(s, &URLServer{})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
