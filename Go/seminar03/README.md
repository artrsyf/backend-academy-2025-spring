### 1. Введение
GRPC и WebSockets позволяют создавать высокопроизводительные и масштабируемые системы. В этом семинаре мы рассмотрим их принципы работы, использование и реализацию на языке Go.

### 2. Основы GRPC
GRPC — это фреймворк для удалённого вызова процедур (RPC), который использует HTTP/2 и бинарный формат сериализации Protocol Buffers (Protobuf).

**Пример установки GRPC в Go:**
```sh
go install google.golang.org/grpc@latest
```

### 3. Protocol Buffers (Protobuf)
Protobuf используется для сериализации данных в GRPC. Он позволяет определять структуру сообщений и автоматически генерировать код для работы с ними.

**Основные возможности Protobuf:**
- Определение структуры данных через `.proto` файлы
- Автоматическая генерация кода для различных языков
- Эффективное бинарное представление данных

**Пример .proto-файла:**
```proto
syntax = "proto3";

package example;

service ExampleService {
  rpc SayHello (HelloRequest) returns (HelloResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}
```

**Команда для генерации Go-кода:**
```sh
protoc --go_out=. --go-grpc_out=. example.proto
```

### 4. Плагины для Protobuf
Для генерации кода из `.proto` используются плагины, такие как `protoc-gen-go` и `protoc-gen-go-grpc`.

### 5. Инструменты для работы с GRPC
#### 5.1 protoc-gen-go
Этот плагин генерирует Go-структуры из Protobuf.

#### 5.2 protoc-gen-go-grpc
Этот плагин создаёт серверные и клиентские реализации.

**Пример сервера на Go:**
```go
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "path/to/generated/proto"
)

type server struct {
	pb.UnimplementedExampleServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterExampleServiceServer(grpcServer, &server{})
	log.Println("Starting server on :50051")
	grpcServer.Serve(lis)
}
```

### 6. Использование grpc-gateway
GRPC-Gateway позволяет конвертировать GRPC-сервисы в REST API.

**Пример настройки grpc-gateway:**

**Добавьте в `.proto` файл:**
```proto
import "google/api/annotations.proto";

service ExampleService {
  rpc SayHello (HelloRequest) returns (HelloResponse) {
    option (google.api.http) = {
      get: "/v1/hello/{name}"
    };
  };
}
```

**Команды для генерации кода:**
```sh
protoc -I . --go_out=. --go-grpc_out=. --grpc-gateway_out=. example.proto
```

**Пример запуска REST-шлюза:**
```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	pb "path/to/generated/proto"
)

func main() {
	conn, err := grpc.DialContext(context.Background(), "localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	mux := runtime.NewServeMux()
	err = pb.RegisterExampleServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalf("failed to register handler: %v", err)
	}
	log.Println("Starting HTTP server on :8080")
	http.ListenAndServe(":8080", mux)
}
```

### 7. WebSockets
WebSockets обеспечивают двустороннюю связь между клиентом и сервером.

**Пример WebSocket-сервера на Go:**
```go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)
		conn.WriteMessage(messageType, msg)
	}
}

func main() {
	http.HandleFunc("/ws", handler)
	log.Println("WebSocket server started on :8080")
	http.ListenAndServe(":8080", nil)
}
```

### 8. AsyncAPI
AsyncAPI используется для документирования событийных API.

### 9. Практическая часть
1. Разверните сервер GRPC и подключите к нему клиента.
2. Запустите WebSocket-сервер и протестируйте двустороннюю связь.
3. Интегрируйте WebSockets и GRPC в одном сервисе.
4. Настройте grpc-gateway для конвертации GRPC в REST API.

### 10. Заключение
GRPC и WebSockets позволяют создавать эффективные и масштабируемые системы. Мы рассмотрели их принципы работы и примеры реализации на Go.

