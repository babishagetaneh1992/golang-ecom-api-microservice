package main

import (
	"context"
	"ecom-api/pkg/auth"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cart-microservice/adaptors/grpc/pb/cart-microservice/services/cart-ms/adaptors/grpc/pb"
	"cart-microservice/internal/adaptors/db"
	grpcAdapter "cart-microservice/internal/adaptors/grpc"
	httpAdapter "cart-microservice/internal/adaptors/http"
	"cart-microservice/internal/application"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "cart-microservice/docs"
)

// @title           Cart Microservice API
// @version         1.0
// @description     This is the Cart service for the e-commerce system.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8083
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.



func main() {
	// --- Load env ---
	if err := godotenv.Load("../../../.env"); err != nil {
    log.Println("Warning: No .env file found, falling back to system environment")
}

auth.InitJWT()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	httpPort := os.Getenv("CART_HTTP_PORT")
	grpcPort := os.Getenv("CART_GRPC_PORT")
	productMsAddr := os.Getenv("PRODUCT_MS_GRPC_ADDR")

	if mongoURI == "" || dbName == "" || httpPort == "" || grpcPort == "" || productMsAddr == "" {
		log.Fatal("❌ Missing required environment variables")
	}

	// --- MongoDB ---
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	dbConn := client.Database(dbName)

	// --- product-ms connect ---
	productConn, err := grpc.Dial(productMsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to product-ms at %s: %v", productMsAddr, err)
	}
	defer productConn.Close()

	productClient := grpcAdapter.NewProductClient(productConn)

	// --- wiring ---
	repo := db.NewMongoCartRepo(dbConn)
	service := application.NewCartService(repo, productClient)

	// HTTP server
	handler := httpAdapter.NewCartHandler(service)
	httpServer := &http.Server{
		Addr:    httpPort,
		Handler: httpAdapter.NewRouter(handler),
	}

	// gRPC server
	grpcServer := grpc.NewServer()
	cartGrpc := grpcAdapter.NewCartGrpcServer(service)
	pb.RegisterCartServiceServer(grpcServer, cartGrpc)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	g := new(errgroup.Group)

	// HTTP
	g.Go(func() error {
		fmt.Println("✅ Cart HTTP server running on", httpPort)
		return httpServer.ListenAndServe()
	})

	// gRPC
	g.Go(func() error {
		fmt.Println("✅ Cart gRPC server running on", grpcPort)
		return grpcServer.Serve(lis)
	})

	// --- Graceful shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	go func() {
		<-stop
		fmt.Println("\n🛑 Shutting down Cart service...")

		// shutdown HTTP
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctxShutdown); err != nil {
			log.Printf("HTTP shutdown error: %v\n", err)
		}

		// shutdown gRPC
		grpcServer.GracefulStop()
	}()

	// wait for error or shutdown
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
