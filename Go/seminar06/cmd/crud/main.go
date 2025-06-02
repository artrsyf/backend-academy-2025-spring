package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Domain Layer
type Course struct {
	ID    string
	Name  string
	Price float64
}

// Repository Layer
type CourseRepository interface {
	Create(course *Course) error
	GetByID(id string) (*Course, error)
	Update(course *Course) error
	Delete(id string) error
}

type PGXCourseRepository struct {
	db *pgx.Conn
}

func (r *PGXCourseRepository) Create(course *Course) error {
	_, err := r.db.Exec(context.Background(), "INSERT INTO courses (id, name, price) VALUES ($1, $2, $3)", course.ID, course.Name, course.Price)
	return err
}

func (r *PGXCourseRepository) GetByID(id string) (*Course, error) {
	row := r.db.QueryRow(context.Background(), "SELECT id, name, price FROM courses WHERE id=$1", id)
	var course Course
	if err := row.Scan(&course.ID, &course.Name, &course.Price); err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *PGXCourseRepository) Update(course *Course) error {
	_, err := r.db.Exec(context.Background(), "UPDATE courses SET name=$1, price=$2 WHERE id=$3", course.Name, course.Price, course.ID)
	return err
}

func (r *PGXCourseRepository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM courses WHERE id=$1", id)
	return err
}

// Application Layer (CQRS)
type CreateCourseCommand struct {
	Name  string
	Price float64
}

type GetCourseQuery struct {
	ID string
}

// gRPC & OpenAPI setup
func generateGRPCServer(repo CourseRepository) *grpc.Server {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	return grpcServer
}

func generateOpenAPIServer(repo CourseRepository) {
	// Load OpenAPI spec and generate handlers
	// Example: using `oapi-codegen` generated server
}

func main() {
	// Database connection
	conn, err := pgx.Connect(context.Background(), "postgres://user:password@localhost:5432/dbname")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	repo := &PGXCourseRepository{db: conn}
	grpcServer := generateGRPCServer(repo)
	generateOpenAPIServer(repo)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen on port 50051:", err)
	}
	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve gRPC server:", err)
	}
}

const limit = 10

func paginate(page int) string {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

var db = &pgx.Conn{}

func getCourse(id int) (*Course, error) {
	course := &Course{}

	var redis interface {
		Get(key string) ([]byte, error)
		Set(key string, value []byte) error
	}

	res, err := redis.Get(fmt.Sprintf("course:%d", id))
	if err == nil {
		err := json.Unmarshal(res, course)
		if err != nil {
			return nil, err
		}
		return course, nil
	}

	db.
		QueryRow(
			context.Background(),
			"SELECT * FROM courses WHERE id = $1",
			id).
		Scan(&course)

	cacheData, err := json.Marshal(course)
	if err != nil {
		return nil, err
	}

	redis.Set(fmt.Sprintf("course:%d", id), cacheData)

	return course, nil
}

// POST /links
// {
// "url": "https://www.google.com",  //not null
// "tags": ["Google"]
// }

// POST /links
// {
// "url": "https://www.google.com",
// }
// -> ID: 1

// PATCH /links/1
// {
// "tags" : ["Google", "Search"]
// }

// - to main menu
// - add tags

// 2025-03-12-add-initial-tables-up.sql

// BEGIN;

// CREATE TABLE links (
// 	id UUID PRIMARY KEY,
// 	url TEXT NOT NULL
// );

// COMMIT;

// 2025-03-12-add-initial-tables-down.sql

// BEGIN;
// DROP TABLE links;
// COMMIT;
